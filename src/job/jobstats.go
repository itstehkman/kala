package job

import (
	"time"
)

type KalaStats struct {
	ActiveJobs   int
	DisabledJobs int
	Jobs         int

	ErrorCount   uint
	SuccessCount uint

	NextRunAt        time.Time
	LastAttemptedRun time.Time

	CreatedAt time.Time
}

func NewKalaStats() *KalaStats {
	ks := &KalaStats{
		CreatedAt: time.Now(),
	}
	jobs := AllJobs.GetAll()

	ks.Jobs = len(jobs)
	if len(jobs) == 0 {
		return ks
	}

	nextRun := time.Time{}
	lastRun := time.Time{}
	for _, job := range jobs {
		if job.Disabled {
			ks.DisabledJobs += 1
		} else {
			ks.ActiveJobs += 1
		}

		if nextRun.IsZero() {
			nextRun = job.NextRunAt
		} else if (nextRun.UnixNano() - job.NextRunAt.UnixNano()) > 0 {
			nextRun = job.NextRunAt
		}

		if lastRun.IsZero() {
			if !job.LastAttemptedRun.IsZero() {
				lastRun = job.LastAttemptedRun
			}
		} else if (lastRun.UnixNano() - job.LastAttemptedRun.UnixNano()) > 0 {
			lastRun = job.LastAttemptedRun
		}

		ks.ErrorCount += job.ErrorCount
		ks.SuccessCount += job.SuccessCount
	}
	ks.NextRunAt = nextRun
	ks.LastAttemptedRun = lastRun

	return ks
}

type JobStat struct {
	JobId             string
	RanAt             time.Time
	NumberOfRetries   uint
	Success           bool
	ExecutionDuration time.Duration
}

func NewJobStat(id string) *JobStat {
	return &JobStat{
		JobId: id,
		RanAt: time.Now(),
	}
}