package bot

import (
	"context"
	"log"
	"time"

	"telegram-job/internal/repository"
	"telegram-job/internal/service"
)

type CleanupService struct {
	jobRepo   *repository.JobRepository
	publisher service.Publisher
	interval  time.Duration
	maxDays   int
}

func NewCleanupService(jobRepo *repository.JobRepository, publisher service.Publisher, maxDays int) *CleanupService {
	return &CleanupService{
		jobRepo:   jobRepo,
		publisher: publisher,
		interval:  1 * time.Hour, // Проверяем каждый час
		maxDays:   maxDays,
	}
}

func (c *CleanupService) Start(ctx context.Context) {
	log.Printf("Cleanup service started. Will archive jobs older than %d days", c.maxDays)

	// Запускаем сразу при старте
	c.cleanup(ctx)

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup service stopped")
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *CleanupService) cleanup(ctx context.Context) {
	log.Println("Running cleanup check...")

	jobs, err := c.jobRepo.GetExpiredJobs(ctx, c.maxDays)
	if err != nil {
		log.Printf("Error getting expired jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		log.Println("No expired jobs found")
		return
	}

	log.Printf("Found %d expired jobs to archive", len(jobs))

	for _, job := range jobs {
		// Удаляем из канала
		if c.publisher != nil && job.ChannelMessageID != nil {
			err := c.publisher.Delete(ctx, *job.ChannelMessageID)
			if err != nil {
				log.Printf("Error deleting job %s from channel: %v", job.ID, err)
			} else {
				log.Printf("Deleted job %s from channel (message_id: %d)", job.ID, *job.ChannelMessageID)
			}
		}

		// Архивируем в БД
		err := c.jobRepo.Archive(ctx, job.ID)
		if err != nil {
			log.Printf("Error archiving job %s: %v", job.ID, err)
		} else {
			log.Printf("Archived job %s: %s", job.ID, job.Title)
		}
	}

	log.Printf("Cleanup complete. Archived %d jobs", len(jobs))
}
