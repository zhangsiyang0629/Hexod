package utils

import (
	"sync"
	"time"
)

var diliverytimer = sync.Pool{
	New: func() interface{} {
		return time.NewTimer(time.Second * 600)
	},
}

type TimerPool sync.Pool

var defaultTimerPool = NewTimerPool()

func NewTimerPool() *TimerPool {
	return (*TimerPool)(&sync.Pool{
		New: func() interface{} {
			return time.NewTimer(time.Second * 600)
		},
	})
}

func (pool *TimerPool) NewTimer(d time.Duration) *time.Timer {
	timer := diliverytimer.Get().(*time.Timer)
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	timer.Reset(d)
	return timer
}

func (pool *TimerPool) FreeTimer(t *time.Timer) {
	t.Stop()
	diliverytimer.Put(t)
}

func NewTimer(d time.Duration) *time.Timer {
	return defaultTimerPool.NewTimer(d)
}

func FreeTimer(t *time.Timer) {
	defaultTimerPool.FreeTimer(t)
}
