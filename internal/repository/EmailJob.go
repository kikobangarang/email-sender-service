package repository

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
)

type EmailJob struct {
	ID        int64
	To        string
	Subject   string
	Body      string
	Status    Status
	Attempts  int
	CreatedAt time.Time
}
