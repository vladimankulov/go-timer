package go_timer

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Timer struct {
	duration  time.Duration
	isStarted bool
	isRunning bool
	count     int
	err       error
	f         process
	mu        sync.Mutex
	ticker    *time.Ticker
}

// function that will be called on every tick of ticker
// bool return type will mark if timer will proceed
// true timer will proceed, false timer will stop
// err points the latest error in func
type process func() (bool, error)

// New create object with duration on when the function will be called
func New(d time.Duration, f process) (*Timer, error) {
	if d <= 0 {
		return nil, fmt.Errorf("duration should be grater than 0")
	}
	if f == nil {
		return nil, fmt.Errorf("should specify func to process")
	}

	return &Timer{
		duration: d,
		f:        f,
	}, nil
}

// StartWithContext the timer
// on every ticker of timer the func will be called
// stops on bool mark of func or depends on context
func (timer *Timer) StartWithContext(ctx context.Context) {
	timer.ticker = time.NewTicker(timer.duration)
	timer.isStarted = true
	timer.isRunning = true
	for {
		select {
		case <-timer.ticker.C:
			timer.count++
			isContinue, err := timer.f()
			timer.err = err
			if !isContinue {
				timer.isRunning = false
				return
			}
		case <-ctx.Done():
			fmt.Println("stopping timer cause ctx is done...")
			timer.isRunning = false
			timer.Stop()
			return
		}
	}
}

// Start the timer
// on every ticker of timer the process will be called
// stops on bool mark of func
func (timer *Timer) Start() {
	timer.ticker = time.NewTicker(timer.duration)
	timer.isStarted = true
	timer.isRunning = true
	for range timer.ticker.C {
		timer.count++
		isContinue, err := timer.f()
		timer.err = err
		if !isContinue {
			timer.isRunning = false
			return
		}
	}
}

// Reset the timer with new duration, no need to start ticker again
func (timer *Timer) Reset(d time.Duration) error {
	if d > 0 {
		timer.ticker.Reset(d)
		return nil
	} else {
		return fmt.Errorf("duration has to be greater than 0")
	}

}

// RestartWithContext the timer with the same duration and new context
func (timer *Timer) RestartWithContext(ctx context.Context) {
	timer.err = nil
	timer.Stop()
	timer.StartWithContext(ctx)
}

// Restart the timer without context
func (timer *Timer) Restart() {
	timer.err = nil
	timer.Stop()
	timer.Start()
}

// Stop the timer
func (timer *Timer) Stop() {
	timer.mu.Lock()
	defer timer.mu.Unlock()
	timer.ticker.Stop()
	timer.isRunning = false
	timer.isStarted = false
}

func (timer *Timer) HasStarted() bool {
	return timer.isStarted
}

func (timer *Timer) IsRunning() bool {
	return timer.isStarted && timer.isRunning
}

func (timer *Timer) Error() error {
	return timer.err
}
