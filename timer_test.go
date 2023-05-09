package go_timer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTimer_Start_with_context(t *testing.T) {
	timeout, cancelFunc := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancelFunc()

	count := 0
	timer, err := New(time.Second, func() (bool, error) {
		if count == 15 {
			return false, nil
		}
		count += 1
		return true, nil
	})
	if err != nil {
		t.Errorf("didn't initialized timer")
		return
	}
	defer func() {
		if timer.isRunning {
			t.Errorf("should be closed")
		}
		timer = nil
	}()
	go timer.StartWithContext(timeout)
	time.Sleep(time.Second * 3)
	if count != 2 {
		fmt.Println(count)
		t.Errorf("count should be equal to 2 cause timeout exeeded after 2 seconds")
	}
}

func TestTimer_Start(t *testing.T) {
	count := 0
	timer, err := New(time.Second, func() (bool, error) {
		if count == 15 {
			return false, nil
		}
		count += 1
		return true, nil
	})
	if err != nil {
		t.Errorf("didn't initialized timer")
		return
	}
	defer func() {
		timer.Stop()
		if timer.isRunning {
			t.Errorf("should be stopped")
		}
		timer = nil
	}()
	go timer.Start()
	time.Sleep(time.Second * 3)
	if count != 3 {
		t.Errorf("count should be equal to 3 ")
	}
}
