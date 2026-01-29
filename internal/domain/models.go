package domain

import (
	"time"

	"github.com/google/uuid"
)

// PostType distinguishes between vacancy and resume
type PostType string

const (
	PostTypeVacancy PostType = "vacancy"
	PostTypeResume  PostType = "resume"
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

// EmploymentType for resumes
type EmploymentType string

const (
	EmploymentFullTime  EmploymentType = "full-time"
	EmploymentPartTime  EmploymentType = "part-time"
	EmploymentContract  EmploymentType = "contract"
	EmploymentFreelance EmploymentType = "freelance"
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

// Post represents both vacancy and resume
type Post struct {
	ID               uuid.UUID   `json:"id"`
	PostType         PostType    `json:"post_type"`
	UserID           *uuid.UUID  `json:"user_id,omitempty"`
	CompanyID        *uuid.UUID  `json:"company_id,omitempty"`
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
	// Resume-specific fields
	ExperienceYears *float64       `json:"experience_years,omitempty"`
	Employment      EmploymentType `json:"employment,omitempty"`
	About           string         `json:"about,omitempty"`
	ResumeLink      string         `json:"resume_link,omitempty"`
	Contact         string         `json:"contact,omitempty"`
}

// Job is alias for Post (backward compatibility)
type Job = Post

// PostWithDetails includes company info for display (used for both vacancies and resumes)
type PostWithDetails struct {
	Post
	CompanyName      string `json:"company_name,omitempty"`
	CompanyContact   string `json:"company_contact,omitempty"`
	AuthorTelegramID int64  `json:"author_telegram_id"`
}

// JobWithCompany is alias for backward compatibility
type JobWithCompany = PostWithDetails

// CreateJobRequest is used when creating a new vacancy
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

// CreateResumeRequest is used when creating a new resume
type CreateResumeRequest struct {
	Title           string         `json:"title"`           // Position title
	Level           JobLevel       `json:"level"`           // Experience level
	Type            JobType        `json:"type"`            // Work format (remote/hybrid/onsite)
	Employment      EmploymentType `json:"employment"`      // full-time, part-time, etc.
	SalaryFrom      *int           `json:"salary_from,omitempty"`
	SalaryTo        *int           `json:"salary_to,omitempty"`
	ExperienceYears *float64       `json:"experience_years,omitempty"`
	About           string         `json:"about"`           // About the candidate
	Contact         string         `json:"contact"`         // Contact info
	ResumeLink      string         `json:"resume_link"`     // Link to CV (optional)
	Language        string         `json:"language"`
}

// Stats contains job statistics
type Stats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Published int `json:"published"`
	Rejected  int `json:"rejected"`
	Archived  int `json:"archived"`
}
