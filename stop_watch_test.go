package main

import (
	"fmt"
	"testing"
	"time"
)

func TestStopWatch(t *testing.T) {
	t.Run("test stopWatch for 3 seconds", func(t *testing.T) {
		var start, now time.Time
		stopWatch := NewStopWatch()
		stopWatch.Reset()
		go stopWatch.Run()
		start = time.Now()
		stopWatch.Start()
		for now.Sub(start).Nanoseconds()/1000000 < 3000 {
			now = time.Now()
			fmt.Println(now.Second()-start.Second(), now.Sub(start).Nanoseconds()/1000000, stopWatch)
			time.Sleep(10 * time.Millisecond)
		}
		stopWatch.Stop()

		_, got, _, _ := stopWatch.GetStopWatch()
		want := int(now.Sub(start).Seconds())
		if got != want {
			t.Errorf("got:%v,want:%v", got, want)
		}
	})
	t.Run("test stopWatch for 1 minute 30 seconds", func(t *testing.T) {
		var start, now time.Time
		stopWatch := NewStopWatch()
		stopWatch.Reset()
		start = time.Now()
		go stopWatch.Run()
		stopWatch.Start()
		for now.Sub(start).Seconds() < 90 {
			now = time.Now()
			fmt.Println(now.Second()-start.Second(), now.Sub(start).Nanoseconds()/1000000, stopWatch)
			time.Sleep(100 * time.Millisecond)
		}
		stopWatch.Stop()

		mSec, sec, minutes, hour := stopWatch.GetStopWatch()
		gotSeconds := sec + minutes*60
		wantSeconds := int(now.Sub(start).Seconds())
		if gotSeconds != wantSeconds {
			t.Errorf("got:%v,want:%v, stopWatch:%v:%v:%v:%v", gotSeconds, wantSeconds, hour, minutes, sec, mSec)
		}
		got := minutes
		want := int(now.Sub(start).Minutes())
		if got != want {
			t.Errorf("got:%v,want:%v, stopWatch:%v:%v:%v:%v", minutes, want, hour, minutes, sec, mSec)
		}
	})
}
