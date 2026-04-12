package telegram

import (
	"errors"
	"strings"
)

func parseCommandArg(text string) (string, error) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		return "", errors.New("invalid input")
	}
	return parts[1], nil
}

func parseCommandArgs(text string) ([]string, error) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid input")
	}
	return strings.Split(parts[1], ","), nil
}
