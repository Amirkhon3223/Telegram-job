package bot

import (
	"sync"

	"telegram-job/internal/domain"
)

type State int

const (
	StateNone State = iota
	StateWaitPostType    // NEW: Choose vacancy or resume
	StateWaitLanguage    // Choose post language

	// Vacancy states
	StateWaitCompany
	StateWaitContact
	StateWaitTitle
	StateWaitLevel
	StateWaitType
	StateWaitCategory
	StateWaitDescription
	StateWaitSalaryFrom
	StateWaitSalaryTo
	StateWaitApplyLink
	StatePreview

	// Resume states
	StateResumeWaitTitle
	StateResumeWaitLevel
	StateResumeWaitExperience
	StateResumeWaitType
	StateResumeWaitEmployment
	StateResumeWaitSalaryFrom
	StateResumeWaitSalaryTo
	StateResumeWaitAbout
	StateResumeWaitContact
	StateResumeWaitLink
	StateResumePreview
)

type Language string

const (
	LangRU Language = "ru"
	LangEN Language = "en"
)

// PostDraft holds data for both vacancy and resume
type PostDraft struct {
	PostType domain.PostType

	// Vacancy fields
	Company     string
	Contact     string // Author contact (for admins)
	Title       string
	Level       domain.JobLevel
	Type        domain.JobType
	Category    domain.JobCategory
	SalaryFrom  *int
	SalaryTo    *int
	Description string
	ApplyLink   string // For candidates
	Language    string

	// Resume fields
	ExperienceYears *float64
	Employment      domain.EmploymentType
	About           string
	ResumeLink      string
	ResumeContact   string // Candidate contact
}

// JobDraft is alias for backward compatibility
type JobDraft = PostDraft

type UserState struct {
	State    State
	Language Language
	Draft    PostDraft
}

type FSM struct {
	mu     sync.RWMutex
	states map[int64]*UserState
}

func NewFSM() *FSM {
	return &FSM{
		states: make(map[int64]*UserState),
	}
}

func (f *FSM) GetState(userID int64) *UserState {
	f.mu.RLock()
	defer f.mu.RUnlock()

	state, ok := f.states[userID]
	if !ok {
		return &UserState{State: StateNone}
	}
	return state
}

func (f *FSM) SetState(userID int64, state State) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.states[userID]; !ok {
		f.states[userID] = &UserState{}
	}
	f.states[userID].State = state
}

func (f *FSM) SetLanguage(userID int64, lang Language) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.states[userID]; !ok {
		f.states[userID] = &UserState{}
	}
	f.states[userID].Language = lang
}

func (f *FSM) GetLanguage(userID int64) Language {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if state, ok := f.states[userID]; ok {
		return state.Language
	}
	return LangEN // default to English
}

func (f *FSM) SetPostType(userID int64, postType domain.PostType) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.states[userID]; !ok {
		f.states[userID] = &UserState{}
	}
	f.states[userID].Draft.PostType = postType
}

func (f *FSM) GetPostType(userID int64) domain.PostType {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if state, ok := f.states[userID]; ok {
		return state.Draft.PostType
	}
	return domain.PostTypeVacancy // default
}

func (f *FSM) UpdateDraft(userID int64, updater func(*PostDraft)) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.states[userID]; !ok {
		f.states[userID] = &UserState{}
	}
	updater(&f.states[userID].Draft)
}

func (f *FSM) GetDraft(userID int64) *PostDraft {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if state, ok := f.states[userID]; ok {
		return &state.Draft
	}
	return nil
}

func (f *FSM) Reset(userID int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.states, userID)
}

func (d *PostDraft) ToCreateJobRequest() *domain.CreateJobRequest {
	return &domain.CreateJobRequest{
		Company:     d.Company,
		Contact:     d.Contact,
		Title:       d.Title,
		Level:       d.Level,
		Type:        d.Type,
		Category:    d.Category,
		SalaryFrom:  d.SalaryFrom,
		SalaryTo:    d.SalaryTo,
		Description: d.Description,
		ApplyLink:   d.ApplyLink,
		Language:    d.Language,
	}
}

// ToCreateRequest is alias for backward compatibility
func (d *PostDraft) ToCreateRequest() *domain.CreateJobRequest {
	return d.ToCreateJobRequest()
}

func (d *PostDraft) ToCreateResumeRequest() *domain.CreateResumeRequest {
	return &domain.CreateResumeRequest{
		Title:           d.Title,
		Level:           d.Level,
		Type:            d.Type,
		Employment:      d.Employment,
		SalaryFrom:      d.SalaryFrom,
		SalaryTo:        d.SalaryTo,
		ExperienceYears: d.ExperienceYears,
		About:           d.About,
		Contact:         d.ResumeContact,
		ResumeLink:      d.ResumeLink,
		Language:        d.Language,
	}
}
