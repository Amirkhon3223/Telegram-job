package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"telegram-job/internal/domain"
	"telegram-job/internal/service"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	telegramID, err := strconv.ParseInt(r.Header.Get("X-Telegram-ID"), 10, 64)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid telegram id")
		return
	}

	var req domain.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	username := r.Header.Get("X-Telegram-Username")
	job, err := h.jobService.CreateJob(r.Context(), telegramID, username, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":     job.ID,
		"status": job.Status,
	})
}

func (h *JobHandler) GetPendingJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobService.GetPendingJobs(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, jobs)
}

func (h *JobHandler) ApproveJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid job id")
		return
	}

	adminID, err := strconv.ParseInt(r.Header.Get("X-Telegram-ID"), 10, 64)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid telegram id")
		return
	}

	err = h.jobService.ApproveJob(r.Context(), jobID, adminID)
	if err != nil {
		switch err {
		case service.ErrForbidden:
			writeError(w, http.StatusForbidden, "forbidden")
		case service.ErrNotFound:
			writeError(w, http.StatusNotFound, "job not found")
		case service.ErrInvalidTransition:
			writeError(w, http.StatusBadRequest, "invalid status transition")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	job, _ := h.jobService.GetJobWithCompany(r.Context(), jobID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":       job.Status,
		"published_at": job.PublishedAt,
	})
}

func (h *JobHandler) RejectJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid job id")
		return
	}

	adminID, err := strconv.ParseInt(r.Header.Get("X-Telegram-ID"), 10, 64)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid telegram id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err = h.jobService.RejectJob(r.Context(), jobID, adminID, req.Reason)
	if err != nil {
		switch err {
		case service.ErrForbidden:
			writeError(w, http.StatusForbidden, "forbidden")
		case service.ErrNotFound:
			writeError(w, http.StatusNotFound, "job not found")
		case service.ErrInvalidTransition:
			writeError(w, http.StatusBadRequest, "invalid status transition")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": domain.JobStatusRejected,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
