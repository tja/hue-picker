package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Mux wraps all endpoints of the API.
type Mux struct {
	// Base muxer
	M *chi.Mux
}

// NewMux will create a new Mux object.
func NewMux() (*Mux, error) {
	// Create Router object
	m := &Mux{
		M: chi.NewRouter(),
	}

	// Set endpoints
	m.M.Get("/light", m.newGetLightHandlerFunc())
	m.M.Post("/light", m.newPostLightHandlerFunc())

	return m, nil
}

// newGetLightHandlerFunc will return a new http.HandlerFunc to handle "GET /light".
func (*Mux) newGetLightHandlerFunc() http.HandlerFunc {
	// Request

	// Response

	// Handler
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get light")) //nolint:errcheck
	}
}

// newPostLightHandlerFunc will return a new http.HandlerFunc to handle "POST /light".
func (*Mux) newPostLightHandlerFunc() http.HandlerFunc {
	// Request

	// Response

	// Handler
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("post light")) //nolint:errcheck
	}
}
