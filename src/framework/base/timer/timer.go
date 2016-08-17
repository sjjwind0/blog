package timer

import (
	"time"
)

type Timer struct {
	delay       time.Duration
	receiver    func()
	isRepeating bool
	baseTimer   *time.Timer
}

func newTimer(isRepeating bool) *Timer {
	return &Timer{isRepeating: isRepeating}
}

func NewRepeatingTimer() *Timer {
	return newTimer(true)
}

func NewOneShotTimer() *Timer {
	return newTimer(false)
}

func (t *Timer) Start(delay time.Duration, receiver func()) {
	t.delay = delay
	t.receiver = receiver
	t.baseTimer = time.NewTimer(delay)
	go func() {
		for {
			<-t.baseTimer.C
			t.receiver()
			if t.isRepeating {
				t.baseTimer.Reset(t.delay)
			} else {
				break
			}
		}
	}()
}

func (t *Timer) Stop() {
	t.baseTimer.Stop()
}
