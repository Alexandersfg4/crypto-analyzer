package models

import "errors"

type Config struct {
	Tokens          []string `json:"tokens"`
	Protocols       []string `json:"protocols"`
	OpenrouterModel string   `json:"openrouter_model"`
}

func (c *Config) Validate() error {
	if len(c.Tokens) == 0 {
		return errors.New("tokens must be provided")
	}
	if len(c.Protocols) == 0 {
		return errors.New("protocols must be provided")
	}

	if c.OpenrouterModel == "" {
		return errors.New("openrouter_model must be provided")
	}

	return nil
}
