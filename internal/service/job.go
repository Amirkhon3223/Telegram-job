package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"telegram-job/internal/config"
	"telegram-job/internal/domain"
	"telegram-job/internal/repository"
)

var (
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrNotFound        = errors.New("job not found")
)

type Publisher interface {
	Publish(ctx context.Context, post *domain.PostWithDetails) (int, error)
	Delete(ctx context.Context, messageID int) error
}

type AdminNotifier interface {
	NotifyNewJob(ctx context.Context, post *domain.PostWithDetails) error
}

type JobService struct {
	cfg         *config.Config
	jobRepo     *repository.JobRepository
	companyRepo *repository.CompanyRepository
	userRepo    *repository.UserRepository
	publisher   Publisher
	notifier    AdminNotifier
}

func NewJobService(
	cfg *config.Config,
	jobRepo *repository.JobRepository,
	companyRepo *repository.CompanyRepository,
	userRepo *repository.UserRepository,
	publisher Publisher,
	notifier AdminNotifier,
) *JobService {
	return &JobService{
		cfg:         cfg,
		jobRepo:     jobRepo,
		companyRepo: companyRepo,
		userRepo:    userRepo,
		publisher:   publisher,
		notifier:    notifier,
	}
}

func (s *JobService) CreateJob(ctx context.Context, telegramID int64, username string, req *domain.CreateJobRequest) (*domain.Job, error) {
	// Get or create user
	user, err := s.userRepo.GetOrCreate(ctx, telegramID, username)
	if err != nil {
		return nil, err
	}

	// Create company
	company := &domain.Company{
		UserID:  user.ID,
		Name:    req.Company,
		Contact: req.Contact,
	}
	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}

	companyID := company.ID
	job := &domain.Post{
		PostType:    domain.PostTypeVacancy,
		UserID:      &user.ID,
		CompanyID:   &companyID,
		Title:       req.Title,
		Level:       req.Level,
		Type:        req.Type,
		Category:    req.Category,
		SalaryFrom:  req.SalaryFrom,
		SalaryTo:    req.SalaryTo,
		Description: req.Description,
		ApplyLink:   req.ApplyLink,
		Status:      domain.JobStatusPending,
		Language:    req.Language,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	if s.notifier != nil {
		jobWithCompany := &domain.PostWithDetails{
			Post:             *job,
			CompanyName:      company.Name,
			CompanyContact:   company.Contact,
			AuthorTelegramID: telegramID,
		}
		_ = s.notifier.NotifyNewJob(ctx, jobWithCompany)
	}

	return job, nil
}

func (s *JobService) CreateResume(ctx context.Context, telegramID int64, username string, req *domain.CreateResumeRequest) (*domain.Post, error) {
	// Get or create user
	user, err := s.userRepo.GetOrCreate(ctx, telegramID, username)
	if err != nil {
		return nil, err
	}

	resume := &domain.Post{
		PostType:        domain.PostTypeResume,
		UserID:          &user.ID,
		Title:           req.Title,
		Level:           req.Level,
		Type:            req.Type,
		SalaryFrom:      req.SalaryFrom,
		SalaryTo:        req.SalaryTo,
		ExperienceYears: req.ExperienceYears,
		Employment:      req.Employment,
		About:           req.About,
		Contact:         req.Contact,
		ResumeLink:      req.ResumeLink,
		Status:          domain.JobStatusPending,
		Language:        req.Language,
	}
	if err := s.jobRepo.Create(ctx, resume); err != nil {
		return nil, err
	}

	if s.notifier != nil {
		resumeWithDetails := &domain.PostWithDetails{
			Post:             *resume,
			AuthorTelegramID: telegramID,
		}
		_ = s.notifier.NotifyNewJob(ctx, resumeWithDetails)
	}

	return resume, nil
}

func (s *JobService) GetPendingJobs(ctx context.Context) ([]domain.JobWithCompany, error) {
	return s.jobRepo.GetByStatus(ctx, domain.JobStatusPending)
}

func (s *JobService) GetJobWithCompany(ctx context.Context, id uuid.UUID) (*domain.JobWithCompany, error) {
	return s.jobRepo.GetWithCompany(ctx, id)
}

func (s *JobService) ApproveJob(ctx context.Context, jobID uuid.UUID, adminTelegramID int64) error {
	// Check admin permission
	if !s.cfg.IsAdmin(adminTelegramID) {
		return ErrForbidden
	}

	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return ErrNotFound
	}

	if job.Status != domain.JobStatusPending {
		return ErrInvalidTransition
	}

	if err := s.jobRepo.UpdateStatus(ctx, jobID, domain.JobStatusApproved); err != nil {
		return err
	}

	jobWithCompany, err := s.jobRepo.GetWithCompany(ctx, jobID)
	if err != nil {
		return err
	}

	var channelMessageID int
	if s.publisher != nil {
		channelMessageID, err = s.publisher.Publish(ctx, jobWithCompany)
		if err != nil {
			return err
		}
	}

	// Set published status with channel message ID
	return s.jobRepo.SetPublished(ctx, jobID, channelMessageID)
}

func (s *JobService) RejectJob(ctx context.Context, jobID uuid.UUID, adminTelegramID int64, reason string) error {
	// Check admin permission
	if !s.cfg.IsAdmin(adminTelegramID) {
		return ErrForbidden
	}

	// Get job
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return ErrNotFound
	}

	// Validate transition: only pending -> rejected allowed
	if job.Status != domain.JobStatusPending {
		return ErrInvalidTransition
	}

	return s.jobRepo.UpdateStatus(ctx, jobID, domain.JobStatusRejected)
}

func (s *JobService) ArchiveJob(ctx context.Context, jobID uuid.UUID, adminTelegramID int64) error {
	// Check admin permission
	if !s.cfg.IsAdmin(adminTelegramID) {
		return ErrForbidden
	}

	// Get job
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return ErrNotFound
	}

	// Can only archive published jobs
	if job.Status != domain.JobStatusPublished {
		return ErrInvalidTransition
	}

	// Delete from channel if we have the message ID
	if s.publisher != nil && job.ChannelMessageID != nil {
		_ = s.publisher.Delete(ctx, *job.ChannelMessageID)
	}

	// Archive in DB
	return s.jobRepo.Archive(ctx, jobID)
}

func (s *JobService) GetJob(ctx context.Context, jobID uuid.UUID) (*domain.Job, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

func (s *JobService) GetUserJobs(ctx context.Context, telegramID int64) ([]domain.Job, error) {
	return s.jobRepo.GetByUserTelegramID(ctx, telegramID)
}

func (s *JobService) GetStats(ctx context.Context) (*domain.Stats, error) {
	return s.jobRepo.GetStats(ctx)
}
