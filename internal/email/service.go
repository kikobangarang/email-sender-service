package email

import (
	"errors"
	"strings"
	"time"

	"github.com/kikobangarang/email-sender-service/internal/repository"
)

type Service struct {
	repo repository.SQLiteRepository
}

func NewService(repo repository.SQLiteRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendEmail(to, subject, body string) error {
	to = strings.TrimSpace(to)
	subject = strings.TrimSpace(subject)
	body = strings.TrimSpace(body)

	if to == "" || subject == "" || body == "" {
		return errors.New("to, subject and body are required")
	}

	if len(subject) > 255 {
		return errors.New("subject too long")
	}

	if len(body) > 100_000 {
		return errors.New("body too large")
	}

	job := repository.EmailJob{
		To:        to,
		Subject:   subject,
		Body:      body,
		Status:    repository.StatusPending,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(job)
}

func (s *Service) GetJobByID(id int) (*repository.EmailJob, error) {
	job, err := s.repo.FetchJobById(id)
	if (*job == repository.EmailJob{}) {
		return nil, errors.New("job not found")
	}
	return job, err
}
