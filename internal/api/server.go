package api

import (
	"encoding/json"
	"net/http"

	"github.com/amimof/huego"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// Server wraps all endpoints of the API.
type Server struct {
	// Base muxer
	M *chi.Mux

	// Hue
	bridge *huego.Bridge
	light  *huego.Light
}

// NewServer will create a new Server object.
func NewServer(bridge *huego.Bridge, light *huego.Light) (*Server, error) {
	// Create Router object
	s := &Server{
		M:      chi.NewRouter(),
		bridge: bridge,
		light:  light,
	}

	// Set endpoints
	s.M.Get("/light", s.newGetLightHandlerFunc())
	s.M.Post("/light", s.newPostLightHandlerFunc())

	return s, nil
}

// newGetLightHandlerFunc will return a new http.HandlerFunc to handle "GET /light".
func (s *Server) newGetLightHandlerFunc() http.HandlerFunc {
	// Request

	// Response
	type Response struct {
		On  bool   `json:"on"`  // Is the Light on or off?
		Bri uint8  `json:"bri"` // Brightness
		Hue uint16 `json:"hue"` // Hue
		Sat uint8  `json:"sat"` // Saturation
	}

	// HTTP handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Get updated light state
		l, err := s.bridge.GetLight(s.light.ID)
		if err != nil {
			logrus.WithError(err).Error("Unable to encode JSON response")
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Set headers
		w.Header().Set("Content-Type", "application/json")

		// Encode light state
		err = json.NewEncoder(w).Encode(&Response{
			On:  l.State.On,
			Bri: l.State.Bri,
			Hue: l.State.Hue,
			Sat: l.State.Sat,
		})

		if err != nil {
			logrus.WithError(err).Error("Unable to encode JSON response")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// newPostLightHandlerFunc will return a new http.HandlerFunc to handle "POST /light".
func (*Server) newPostLightHandlerFunc() http.HandlerFunc {
	// Request

	// Response

	// Handler
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("post light")) //nolint:errcheck
	}
}
