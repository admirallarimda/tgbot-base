package botbase

import "testing"
import "time"
import "math/rand"
import "sync/atomic"

type testCronCountingJob struct {
	count  int32
	repeat *time.Duration
}

func (j *testCronCountingJob) Do(t time.Time, c Cron) {
	atomic.AddInt32(&j.count, 1) // atomic to avoid race detector
	if j.repeat != nil {
		c.AddJob(t.Add(*j.repeat), j)
	}
}

func TestCallOnce(t *testing.T) {
	c := NewCron()
	j := &testCronCountingJob{}
	c.AddJob(time.Now(), j)
	time.Sleep(time.Second)
	atomic.LoadInt32(&j.count)
	if j.count != 1 {
		t.Fatal(j.count)
	}
}

func TestCallXTimes(t *testing.T) {
	c := NewCron()
	j := &testCronCountingJob{}

	now := time.Now()
	n := 5 + rand.Int31n(5)
	var i int32
	for ; i < n; i++ {
		c.AddJob(now, j)
	}

	time.Sleep(time.Second)
	atomic.LoadInt32(&j.count)
	if j.count != n {
		t.Fatal(j.count, n)
	}
}
