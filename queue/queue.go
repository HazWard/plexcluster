package queue

import "github.com/hazward/plexcluster/types"


// JobQueuer captures the  behavior of a job queuing
// interface where it is able to submit jobs and receive
// notifications when jobs are completed
type JobQueuer interface {
	Submit(job types.Job) error
	WaitForCompletion(jobID string, found chan int) error
}