package utils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/laurawulinjie/goAcademy/pkg/todo"
)

func GenerateTraceID() string {
	return uuid.New().String()
}

var ErrInvalidStatus = errors.New("status must be one of: not started, started, completed")

func ValidateStatus(status string) (string, error) {
	switch status {
	case "not started":
		return todo.NotStarted, nil
	case "started":
		return todo.Started, nil
	case "completed":
		return todo.Completed, nil
	case "":
		return "", nil
	default:
		return "", ErrInvalidStatus
	}
}
