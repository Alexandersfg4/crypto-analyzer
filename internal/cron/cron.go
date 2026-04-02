package cron

import "time"

type Cron struct {
	ticker *time.Ticker
}

func New(duration time.Duration) *Cron {
	return &Cron{
		ticker: time.NewTicker(duration),
	}
}

func (c *Cron) Run(task func()) {
	for {
		select {
		case <-c.ticker.C:
			task()
		}
	}
}

func (c *Cron) Reset(duration time.Duration) {
	c.ticker.Reset(duration)
}

func (c *Cron) Stop() {
	c.ticker.Stop()
}
