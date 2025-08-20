package game

import (
	"testing"
	"testing/synctest"
	"time"
)

// This uses a fake timer, so the Sleep() calls don't actually cost time

func TestClockStartStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := NewClock(60 * time.Second)

		c.Start()
		time.Sleep(10 * time.Second)
		synctest.Wait()
		c.Stop()

		actual := c.RemainingTime()
		want := 50 * time.Second

		if actual != want {
			t.Errorf("ElapsedTime() initially got %s, want %s", actual, want)
		}

		// Restart
		c.Start()
		time.Sleep(10 * time.Second)
		synctest.Wait()
		c.Stop()

		actual = c.RemainingTime()
		want = 40 * time.Second

		if actual != want {
			t.Errorf("ElapsedTime() after restart got %s, want %s", actual, want)
		}

		// Remaining time live
		c.Start()
		time.Sleep(10 * time.Second)
		synctest.Wait()

		actual = c.RemainingTime()
		want = 30 * time.Second

		if actual != want {
			t.Errorf("ElapsedTime() when live got %s, want %s", actual, want)
		}

		// Increment
		c.Stop()
		c.Increment(5 * time.Second)

		actual = c.RemainingTime()
		want = 35 * time.Second

		if actual != want {
			t.Errorf("ElapsedTime() after increment got %s, want %s", actual, want)
		}

		// Restart
		c.Start()
		time.Sleep(10 * time.Second)
        synctest.Wait()
		c.Stop()

		actual = c.RemainingTime()
		want = 25 * time.Second

		if actual != want {
			t.Errorf("ElapsedTime() after increment & restart got %s, want %s", actual, want)
		}
	})
}
