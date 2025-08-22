package game

import "time"

type Clock struct {
	running     bool
	restartTime time.Time
	remaining   time.Duration
}

func NewClock(duration time.Duration) *Clock {
	return &Clock{remaining: duration}
}

func (c *Clock) Start() {
	if c.running {
		return
	}
	c.restartTime = time.Now()
	c.running = true
}

func (c *Clock) Running() bool {
	return c.running
}

func (c *Clock) Increment(duration time.Duration) {
	c.remaining += duration
}

func (c *Clock) RemainingTime() time.Duration {
	if c.running {
		return max(c.remaining-time.Now().Sub(c.restartTime), 0)
	}
	return max(c.remaining, 0)
}

func (c *Clock) Stop() {
	if !c.running {
		return
	}
	c.remaining -= time.Now().Sub(c.restartTime)
	c.running = false
}
