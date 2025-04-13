package docker

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/trinhminhtriet/docker-activity/internal/model"
	// "github.com/trinhminhtriet/docker-activity/pkg/error"
)

type ContainerWatcher struct {
	Docker *Client
	Name   string
}

func (w *ContainerWatcher) IsAlive(ctx context.Context) (bool, error) {
	filters := filters.NewArgs()
	filters.Add("name", w.Name)
	filters.Add("name", strings.TrimPrefix(w.Name, "/"))
	filters.Add("status", "running")

	containers, err := w.Docker.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return false, error.NewCustom("couldn't list containers: %v", err)
	}
	return len(containers) > 0, nil
}

func (w *ContainerWatcher) Run(ctx context.Context, register *sync.Map, tx chan<- model.Record) error {
	slog.Info("Watching container", "name", w.Name)
	for {
		alive, err := w.IsAlive(ctx)
		if err != nil {
			return err
		}
		if !alive {
			break
		}

		stats, err := w.Docker.ContainerStats(ctx, strings.TrimPrefix(w.Name, "/"), true)
		if err != nil {
			slog.Warn("Failed to get stats", "name", w.Name, "error", err)
			continue
		}

		decoder := json.NewDecoder(stats.Body)
		for {
			var stat types.StatsJSON
			if err := decoder.Decode(&stat); err != nil {
				slog.Warn("Lost stats connection", "name", w.Name, "error", err)
				break
			}
			select {
			case tx <- model.NewRecordFromStats(stat):
			default:
				slog.Warn("Channel full, dropping snapshot", "name", w.Name)
			}
		}
		stats.Body.Close()
	}

	register.Delete(w.Name)
	slog.Info("Done watching container", "name", w.Name)
	return nil
}

type Orchestrator struct {
	Docker *Client
	Names  map[string]bool
	Tasks  *sync.Map
}

func NewOrchestrator(cli *Client, containers string) (*Orchestrator, error) {
	names := make(map[string]bool)
	if containers != "" {
		for _, name := range strings.Split(containers, ",") {
			names[name] = true
		}
	}
	return &Orchestrator{
		Docker: cli,
		Names:  names,
		Tasks:  &sync.Map{},
	}, nil
}

func (o *Orchestrator) IsRunning(name string) bool {
	_, exists := o.Tasks.Load(name)
	return exists
}

func (o *Orchestrator) RegisterTask(name string) {
	o.Tasks.Store(name, true)
}

func (o *Orchestrator) HandleStartEvent(ctx context.Context, containerName string, tx chan model.Record) error {
	if o.IsRunning(containerName) {
		slog.Debug("Container already running", "name", containerName)
		return nil
	}
	if len(o.Names) > 0 && !o.Names[containerName] {
		return nil
	}
	o.RegisterTask(containerName)

	watcher := &ContainerWatcher{Docker: o.Docker, Name: containerName}
	go func() {
		if err := watcher.Run(ctx, o.Tasks, tx); err != nil {
			slog.Warn("Container watcher errored", "name", containerName, "error", err)
		}
	}()
	return nil
}

func (o *Orchestrator) ListRunning(ctx context.Context) ([]string, error) {
	filters := filters.NewArgs()
	filters.Add("status", "running")
	containers, err := o.Docker.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return nil, error.NewCustom("couldn't list running containers: %v", err)
	}
	var names []string
	for _, c := range containers {
		if len(c.Names) > 0 {
			names = append(names, c.Names[0])
		}
	}
	return names, nil
}

func (o *Orchestrator) Run(ctx context.Context, tx chan model.Record) error {
	// Start watching existing containers
	names, err := o.ListRunning(ctx)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err := o.HandleStartEvent(ctx, name, tx); err != nil {
			slog.Warn("Failed to handle start event", "name", name, "error", err)
		}
	}

	// Watch Docker events
	filters := filters.NewArgs()
	filters.Add("type", "container")
	events, errs := o.Docker.Events(ctx, types.EventsOptions{Filters: filters})
	for {
		select {
		case event := <-events:
			if event.Action == "start" && event.Actor.Attributes != nil {
				containerName := event.Actor.Attributes["name"]
				if containerName == "" {
					continue
				}
				if err := o.HandleStartEvent(ctx, containerName, tx); err != nil {
					slog.Warn("Failed to handle event", "name", containerName, "error", err)
				}
			}
		case err := <-errs:
			if err != nil {
				return error.WrapIO(err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
