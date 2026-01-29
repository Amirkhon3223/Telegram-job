package bot

import (
	"sync"

	"telegram-job/internal/domain"
)

type State int

const (
	StateNone State = iota
	StateWaitLanguage
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
)

type Language string

const (
	LangRU Language = "ru"
	LangEN Language = "en"
)

type JobDraft struct {
	Company     string
	Contact     string
	Title       string
	Level       domain.JobLevel
	Type        domain.JobType
	Category    domain.JobCategory
	SalaryFrom  *int
	SalaryTo    *int
	Description string
	ApplyLink   string
	Language    string
}

type UserState struct {
	State    State
	Language Language
	Draft    JobDraft
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
	return LangRU // default
}

func (f *FSM) UpdateDraft(userID int64, updater func(*JobDraft)) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.states[userID]; !ok {
		f.states[userID] = &UserState{}
	}
	updater(&f.states[userID].Draft)
}

func (f *FSM) GetDraft(userID int64) *JobDraft {
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

func (d *JobDraft) ToCreateRequest() *domain.CreateJobRequest {
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
