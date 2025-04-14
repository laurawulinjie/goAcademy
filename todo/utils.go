package main

import (
	"errors"

	"github.com/google/uuid"
)

func GenerateTraceID() string {
	return uuid.New().String()
}

var ErrInvalidStatus = errors.New("status must be one of: not started, started, completed")

func ValidateStatus(status string) (string, error) {
	switch status {
	case "not started":
		return NotStarted, nil
	case "started":
		return Started, nil
	case "completed":
		return Completed, nil
	case "":
		return "", nil
	default:
		return "", ErrInvalidStatus

	}
}
