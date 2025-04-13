package error

import "fmt"

// Error represents a custom error.
type Error struct {
	Kind    string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Kind, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

// NewCustom creates a custom error.
func NewCustom(format string, args ...any) error {
	return &Error{Kind: "Custom", Message: fmt.Sprintf(format, args...)}
}

// WrapIO wraps an IO error.
func WrapIO(err error) error {
	if err == nil {
		return nil
	}
	return &Error{Kind: "IO", Message: "IO error", Cause: err}
}
