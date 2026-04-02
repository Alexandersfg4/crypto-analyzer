package models

import (
	"errors"
)

type Config struct {
	Tokens    []string `json:"tokens"`
	Protocols []string `json:"protocols"`
}

func (c *Config) Validate() error {
	if len(c.Tokens) == 0 {
		return errors.New("tokens must be provided")
	}
	if len(c.Protocols) == 0 {
		return errors.New("protocols must be provided")
	}

	return nil
}
