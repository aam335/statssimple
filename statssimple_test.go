package statssimple

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []int64{10, 20, 30, 40, 50}
	stats := NewStatsSimple()
	stats1 := NewStatsSimple()

	for _, v := range tests {
		stats1.PushVal(v)
	}

	min, max, avg, cnt := stats1.GetStatsNs()

	require.Equal(t, int64(10), min, "min")
	require.Equal(t, int64(50), max, "max")
	require.Equal(t, int64(30), avg, "avg")
	require.Equal(t, int64(5), cnt, "cnt")

	stats.Append(stats1)
	stats.Wait()
	min, max, avg, cnt = stats.GetStatsNs()

	require.Equal(t, int64(10), min, "min")
	require.Equal(t, int64(50), max, "max")
	require.Equal(t, int64(30), avg, "avg")
	require.Equal(t, int64(5), cnt, "cnt")
}

func TestRunOne(t *testing.T) {
	testsCnt := 10
	dt := time.Millisecond

	stats := NewStatsSimple()

	wg := sync.WaitGroup{}
	wg.Add(testsCnt)
	for cnt := 0; cnt < testsCnt; cnt++ {
		go func() {
			statsn := NewStatsSimple()
			statsn.RunOne(func() {
				time.Sleep(dt)
			})
			stats.Append(statsn)
			wg.Done()
		}()
	}
	wg.Wait()
	stats.Wait()
	// stats
	min, max, avg, cnt := stats.GetStatsNs()
	ns := int64(dt.Nanoseconds())
	require.LessOrEqual(t, ns, min, "min")
	require.LessOrEqual(t, ns, max, "max<dt")
	require.GreaterOrEqual(t, ns*100, max, "max >100*dt")
	require.LessOrEqual(t, ns, avg, "avg")
	require.Equal(t, int64(testsCnt), cnt, "cnt")
}
