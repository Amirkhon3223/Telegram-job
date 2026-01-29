package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobLevel string

const (
	JobLevelJunior     JobLevel = "junior"
	JobLevelMiddle     JobLevel = "middle"
	JobLevelSenior     JobLevel = "senior"
	JobLevelInternship JobLevel = "internship"
	JobLevelSkip       JobLevel = ""
)

type JobType string

const (
	JobTypeRemote JobType = "remote"
	JobTypeHybrid JobType = "hybrid"
	JobTypeOnsite JobType = "onsite"
)

type JobCategory string

const (
	JobCategoryWeb2 JobCategory = "web2"
	JobCategoryWeb3 JobCategory = "web3"
	JobCategoryDev  JobCategory = "dev"
)

type JobStatus string

const (
	JobStatusDraft     JobStatus = "draft"
	JobStatusPending   JobStatus = "pending"
	JobStatusApproved  JobStatus = "approved"
	JobStatusPublished JobStatus = "published"
	JobStatusRejected  JobStatus = "rejected"
	JobStatusArchived  JobStatus = "archived"
)

type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleRecruiter UserRole = "recruiter"
)

type User struct {
	ID                uuid.UUID `json:"id"`
	TelegramID        int64     `json:"telegram_id"`
	Username          string    `json:"username"`
	Role              UserRole  `json:"role"`
	InterfaceLanguage *string   `json:"interface_language,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type Company struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	Telegram  string    `json:"telegram,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Job struct {
	ID               uuid.UUID   `json:"id"`
	CompanyID        uuid.UUID   `json:"company_id"`
	Title            string      `json:"title"`
	Level            JobLevel    `json:"level"`
	Type             JobType     `json:"type"`
	Category         JobCategory `json:"category"`
	SalaryFrom       *int        `json:"salary_from,omitempty"`
	SalaryTo         *int        `json:"salary_to,omitempty"`
	Description      string      `json:"description"`
	ApplyLink        string      `json:"apply_link"`
	Status           JobStatus   `json:"status"`
	Language         string      `json:"language"`
	ChannelMessageID *int        `json:"channel_message_id,omitempty"`
	PublishedAt      *time.Time  `json:"published_at,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
}

// JobWithCompany includes company info for display
type JobWithCompany struct {
	Job
	CompanyName    string `json:"company_name"`
	CompanyContact string `json:"company_contact"`
	AuthorTelegramID int64 `json:"author_telegram_id"`
}

// CreateJobRequest is used when creating a new job
type CreateJobRequest struct {
	Company     string      `json:"company"`
	Contact     string      `json:"contact"`
	Title       string      `json:"title"`
	Level       JobLevel    `json:"level"`
	Type        JobType     `json:"type"`
	Category    JobCategory `json:"category"`
	SalaryFrom  *int        `json:"salary_from,omitempty"`
	SalaryTo    *int        `json:"salary_to,omitempty"`
	Description string      `json:"description"`
	ApplyLink   string      `json:"apply_link"`
	Language    string      `json:"language"`
}

// Stats contains job statistics
type Stats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Published int `json:"published"`
	Rejected  int `json:"rejected"`
	Archived  int `json:"archived"`
}
