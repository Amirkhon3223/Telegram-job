package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(jobHandler *JobHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Route("/jobs", func(r chi.Router) {
			r.Post("/", jobHandler.CreateJob)
			r.Get("/", jobHandler.GetPendingJobs)
			r.Post("/{id}/approve", jobHandler.ApproveJob)
			r.Post("/{id}/reject", jobHandler.RejectJob)
		})
	})

	return r
}
