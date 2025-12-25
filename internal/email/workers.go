package email

import (
	"context"
	"log"
	"time"

	"github.com/kikobangarang/email-sender-service/internal/repository"
)

type WorkerConfig struct {
	WorkerCount  int
	PollInterval time.Duration
	MaxRetries   int
	BatchSize    int
}

type WorkerPool struct {
	repo   repository.SQLiteRepository
	sender Sender
	cfg    WorkerConfig
}

func NewWorkerPool(
	repo repository.SQLiteRepository,
	sender Sender,
	cfg WorkerConfig,
) *WorkerPool {
	return &WorkerPool{
		repo:   repo,
		sender: sender,
		cfg:    cfg,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	log.Printf("starting %d email workers\n", wp.cfg.WorkerCount)

	for i := 0; i < wp.cfg.WorkerCount; i++ {
		go wp.workerLoop(ctx, i)
	}
}

func (wp *WorkerPool) workerLoop(ctx context.Context, workerID int) {
	log.Printf("worker %d started\n", workerID)

	ticker := time.NewTicker(wp.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d shutting down\n", workerID)
			return

		case <-ticker.C:
			wp.processBatch(workerID)
		}
	}
}

func (wp *WorkerPool) processBatch(workerID int) {
	jobs, err := wp.repo.FetchPending(wp.cfg.BatchSize)
	if err != nil {
		log.Printf("worker %d: failed fetching jobs: %v\n", workerID, err)
		return
	}

	for _, job := range jobs {
		wp.processJob(workerID, job)
	}
}

func (wp *WorkerPool) processJob(workerID int, job repository.EmailJob) {
	log.Printf("worker %d: sending email job %d\n", workerID, job.ID)

	err := wp.sender.Send(job.To, job.Subject, job.Body)
	if err != nil {
		wp.handleFailure(workerID, job, err)
		return
	}

	if err := wp.repo.MarkSent(job.ID); err != nil {
		log.Printf("worker %d: failed marking job %d as sent: %v\n",
			workerID, job.ID, err)
		return
	}

	log.Printf("worker %d: job %d sent successfully\n", workerID, job.ID)
}

func (wp *WorkerPool) handleFailure(workerID int, job repository.EmailJob, sendErr error) {
	job.Attempts++

	log.Printf(
		"worker %d: job %d failed attempt %d: %v\n",
		workerID,
		job.ID,
		job.Attempts,
		sendErr,
	)

	if job.Attempts >= wp.cfg.MaxRetries {
		if err := wp.repo.MarkFailed(job.ID, sendErr); err != nil {
			log.Printf(
				"worker %d: failed marking job %d as failed: %v\n",
				workerID,
				job.ID,
				err,
			)
		}
		return
	}

	if err := wp.repo.IncrementAttempts(job.ID); err != nil {
		log.Printf(
			"worker %d: failed incrementing attempts for job %d: %v\n",
			workerID,
			job.ID,
			err,
		)
	}
}
