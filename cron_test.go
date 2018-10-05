package botbase

import "testing"
import "time"
import "math/rand"

type testCronCountingJob struct {
	count  int
	repeat *time.Duration
}

func (j *testCronCountingJob) Do(t time.Time, c Cron) {
	j.count++
	if j.repeat != nil {
		c.AddJob(t.Add(*j.repeat), j)
	}
}

func TestCallOnce(t *testing.T) {
	c := NewCron()
	j := &testCronCountingJob{}
	c.AddJob(time.Now(), j)
	time.Sleep(time.Second)
	if j.count != 1 {
		t.Fatal(j.count)
	}
}

func TestCallXTimes(t *testing.T) {
	c := NewCron()
	j := &testCronCountingJob{}

	now := time.Now()
	n := 5 + rand.Intn(5)
	for i := 0; i < n; i++ {
		c.AddJob(now, j)
	}

	time.Sleep(time.Second)
	if j.count != n {
		t.Fatal(j.count, n)
	}
}
