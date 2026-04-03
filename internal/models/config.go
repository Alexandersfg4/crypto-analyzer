package models

import (
	"errors"
	"time"
)

type Config struct {
	Tokens                []string  `json:"tokens"`
	Protocols             []string  `json:"protocols"`
	CronNextExecutionTime time.Time `json:"cron_next_execution_time"`
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
