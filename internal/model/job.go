package model

import "time"

type JobStatus string

const (
	StatusPending  JobStatus = "pending"
	StatusRunning  JobStatus = "running"
	StatusSuccess  JobStatus = "success"
	StatusFailed   JobStatus = "failed"
)

type PipelineRun struct {
	ID         int
	Repo       string
	Status     JobStatus
	Logs       string
	StartedAt  time.Time
	FinishedAt time.Time
}
