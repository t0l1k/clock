package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	t.Run("test timer", func(t *testing.T) {
		var start, now time.Time
		timer := NewTimer()
		timer.Reset()
		go timer.Run()
		start = time.Now()
		timer.Start()
		for now.Sub(start).Nanoseconds()/1000000 < 3000 {
			now = time.Now()
			fmt.Println(now.Second()-start.Second(), now.Sub(start).Nanoseconds()/1000000, timer)
			time.Sleep(10 * time.Millisecond)
		}
		timer.Stop()

		_, got, _, _ := timer.GetTimer()
		want := int(now.Sub(start).Seconds())
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
	})
}
