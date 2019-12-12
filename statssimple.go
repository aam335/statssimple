package statssimple

import (
	"math"
	"sync"
	"time"
)

// StatsSimple ...
type StatsSimple struct {
	unixNsStart,
	max, min, sum, cnt, avg int64
	logchan   chan logchanStruct
	wg        sync.WaitGroup
	startTime time.Time
}

type logchanStruct struct{ min, max, sum, cnt int64 }

// NewStatsSimple initializator for min & max values
func NewStatsSimple() *StatsSimple {
	ch := make(chan logchanStruct)
	s := &StatsSimple{max: math.MinInt64, min: math.MaxInt64, logchan: ch, startTime: time.Now()}
	go func() {
		for v := range ch {
			if s.max < v.max {
				s.max = v.max
			}
			if s.min > v.min {
				s.min = v.min
			}
			s.sum += v.sum
			s.cnt += v.cnt
			s.wg.Done()
		}
	}()
	return s
}

// Append Adds an intermediate set to the main set
func (s *StatsSimple) Append(v *StatsSimple) {
	s.wg.Add(1)
	s.logchan <- logchanStruct{v.min, v.max, v.sum, v.cnt}
}

// PushVal Adds single value to the main set (not thread-safe!)
func (s *StatsSimple) PushVal(v int64) {
	if s.max < v {
		s.max = v
	}
	if s.min > v {
		s.min = v
	}
	s.cnt++
	s.sum += v
}

// StartOne starts one measurement.
func (s *StatsSimple) StartOne() {
	s.unixNsStart = time.Now().UnixNano()
}

// DoneOne one
func (s *StatsSimple) DoneOne() {
	doneTime := time.Now().UnixNano() - s.unixNsStart
	s.PushVal(doneTime)
}

// RunOne exec f & store exec time (not thread-safe!)
func (s *StatsSimple) RunOne(f func()) {
	s.StartOne()
	f()
	s.DoneOne()
}

// GetStatsNs calcs avg & return values
func (s *StatsSimple) GetStatsNs() (min, max, avg, count int64, d time.Duration) {
	if s.cnt > 0 {
		s.avg = s.sum / s.cnt
	}
	d = time.Now().Sub(s.startTime)
	return s.min, s.max, s.avg, s.cnt, d
}

// Wait .... waits for queue is empty
func (s *StatsSimple) Wait() {
	s.wg.Wait()
}

// Shutdown - close chan, goproc
func (s *StatsSimple) Shutdown() {
	close(s.logchan)
}
