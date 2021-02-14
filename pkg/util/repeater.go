package util

import (
	"time"
)

type StopChannel chan bool

func RepeatUntilStopped(delay time.Duration, f func()) StopChannel {
	control := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-control:
				return
			case <-time.After(delay):
			}
			f()
		}
	}()
	return control
}
