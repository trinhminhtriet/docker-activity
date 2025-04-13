package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
	"github.com/trinhminhtriet/docker-activity/internal/docker"
	"github.com/trinhminhtriet/docker-activity/internal/exporter"
	"github.com/trinhminhtriet/docker-activity/internal/model"
)

type Params struct {
	BufferSize int
	Containers string
	Exporter   exporter.Config
}

func parseParams() (Params, error) {
	var p Params
	fs := pflag.NewFlagSet("docker-activity", pflag.ContinueOnError)
	fs.IntVar(&p.BufferSize, "buffer-size", 32, "Size of the channel buffer")
	fs.StringVar(&p.Containers, "containers", "", "Comma-separated container names or IDs to monitor")

	// Subcommands for exporters
	file := &exporter.FileConfig{}
	stdout := &exporter.StdoutConfig{}
	fs.AddFlagSet(file.Flags())
	fs.AddFlagSet(stdout.Flags())

	if err := fs.Parse(os.Args[1:]); err != nil {
		return Params{}, err
	}

	// Determine exporter based on flags
	if file.Path != "" {
		p.Exporter = file
	} else {
		p.Exporter = stdout
	}
	return p, nil
}

func main() {
	params, err := parseParams()
	if err != nil {
		slog.Error("Failed to parse parameters", "error", err)
		os.Exit(1)
	}

	// Initialize logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	ctx := context.Background()

	// Create Docker client
	cli, err := docker.NewClient()
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		os.Exit(1)
	}

	// Create channel for records
	recordChan := make(chan model.Record, params.BufferSize)

	// Start exporter
	go func() {
		exp := params.Exporter.Exporter()
		for record := range recordChan {
			if err := exp.Handle(record); err != nil {
				slog.Warn("Failed to export record", "error", err)
			}
		}
	}()

	// Start orchestrator
	orch, err := docker.NewOrchestrator(cli, params.Containers)
	if err != nil {
		slog.Error("Failed to create orchestrator", "error", err)
		os.Exit(1)
	}
	if err := orch.Run(ctx, recordChan); err != nil {
		slog.Error("Orchestrator failed", "error", err)
		os.Exit(1)
	}
}
