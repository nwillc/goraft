package util

import (
	"time"
)

type StopChannel chan bool

// RepeatUntilStopped creates a goroutine that loops, waiting the delay provided and
// then calling the function provided. The loop can be stopped by sending a value to
// returned StopChannel.
func RepeatUntilStopped(delay time.Duration, f func()) StopChannel {
	control := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-control:
				close(control)
				return
			case <-time.After(delay):
			}
			f()
		}
	}()
	return control
}
