package cron

import (
	"context"
	"time"
)

type Cron struct {
	ticker        *time.Ticker
	executionTime time.Time
}

func New(t time.Time) *Cron {
	executionTime, duration := nextExecutionTimeDuration(t)

	return &Cron{
		executionTime: executionTime,
		ticker:        time.NewTicker(duration),
	}
}

func (c *Cron) Run(ctx context.Context, task func()) {
	for {
		select {
		case <-c.ticker.C:
			task()
			executionTime, duration := nextExecutionTimeDuration(c.executionTime)
			c.executionTime = executionTime
			c.ticker.Reset(duration)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Cron) Reset(t time.Time) {
	c.ticker.Stop()
	executionTime, duration := nextExecutionTimeDuration(t)
	c.executionTime = executionTime
	c.ticker.Reset(duration)
}

func (c *Cron) ExecutionTime() time.Time {
	return c.executionTime
}

func nextExecutionTimeDuration(t time.Time) (time.Time, time.Duration) {
	now := time.Now()
	tommorow := now.AddDate(0, 0, 1)
	executionTime := time.Date(tommorow.Year(), tommorow.Month(), tommorow.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	return executionTime, executionTime.Sub(now)
}
