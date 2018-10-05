package botbase

import "time"
import "log"
import "sort"

// Cron interface declares interfaces for communication with some cron daemon
type Cron interface {
	AddJob(t time.Time, job CronJob)
}

// CronJob provides a piece of work which should be done once its time has come
type CronJob interface {
	Do(now time.Time, cron Cron)
}

type cronJobDesc struct {
	execTime time.Time
	job      CronJob
}

type cron struct {
	newJobCh chan cronJobDesc
	timer    *time.Timer

	jobs           map[time.Time][]CronJob
	sortedJobTimes []time.Time
}

func (c *cron) AddJob(t time.Time, job CronJob) {
	c.newJobCh <- cronJobDesc{
		execTime: t,
		job:      job}
}

func (c *cron) executeJobs(jobsToExecute map[time.Time][]CronJob, now time.Time) {
	for scheduledTime, jobs := range jobsToExecute {
		log.Printf("Executing %d jobs at time %s (scheduled %s; diff %s)", len(jobsToExecute), now, scheduledTime, now.Sub(scheduledTime))
		for _, j := range jobs {
			go j.Do(scheduledTime, c)
		}
	}
}

func (c *cron) processNewJob(execTime time.Time, job CronJob) {
	if _, found := c.jobs[execTime]; found {
		log.Printf("New job with known time %s has arrived", execTime)
		c.jobs[execTime] = append(c.jobs[execTime], job)
	} else {
		log.Printf("New job with not yet known time %s has arrived", execTime)
		c.jobs[execTime] = []CronJob{job}
		c.sortedJobTimes = append(c.sortedJobTimes, execTime)
		sort.Slice(c.sortedJobTimes, func(i int, j int) bool {
			return c.sortedJobTimes[i].Before(c.sortedJobTimes[j])
		})
		c.resetTimer(time.Now())
	}
}

func (c *cron) resetTimer(now time.Time) {
	nextTimer := c.sortedJobTimes[0].Sub(now)
	if !c.timer.Stop() {
		<-c.timer.C
	}
	c.timer.Reset(nextTimer)
}

func (c *cron) run() {
	isRunning := true
	for isRunning {
		select {
		case j := <-c.newJobCh:
			c.processNewJob(j.execTime, j.job)
		case now := <-c.timer.C:
			pos := sort.Search(len(c.sortedJobTimes), func(i int) bool {
				return c.sortedJobTimes[i].Before(now)
			})
			if pos == len(c.sortedJobTimes) {
				panic("cron scheduling inconsistency")
			}
			// preparing list of jobs which should be executed, removing them from internal structures
			jobsToExecute := make(map[time.Time][]CronJob, pos+1)
			for i := 0; i <= pos; i++ {
				t := c.sortedJobTimes[i]
				jobsToExecute[t] = c.jobs[t]
				delete(c.jobs, t)
			}
			c.sortedJobTimes = c.sortedJobTimes[pos+1:]
			if len(jobsToExecute) == 0 {
				panic("cron time-to-jobs inconsistency")
			}
			c.executeJobs(jobsToExecute, now)
			c.resetTimer(now)
		}
	}
}

// NewCron creates an instance of cron
func NewCron() Cron {
	c := cron{
		newJobCh:       make(chan cronJobDesc, 0),
		jobs:           make(map[time.Time][]CronJob, 0),
		sortedJobTimes: make([]time.Time, 0)}

	go c.run()

	return &c
}
