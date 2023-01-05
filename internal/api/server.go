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
		M:      chi.NewMux(),
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
	//   None

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

		s.light = l

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
func (s *Server) newPostLightHandlerFunc() http.HandlerFunc {
	// Request
	type Request struct {
		On  bool   `json:"on"`  // Is the Light on or off?
		Bri uint8  `json:"bri"` // Brightness
		Hue uint16 `json:"hue"` // Hue
		Sat uint8  `json:"sat"` // Saturation
	}

	// Response
	//   None

	// Handler
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse input request
		var in Request

		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			logrus.WithError(err).Error("Unable to decode JSON request")
			w.WriteHeader(http.StatusBadRequest)
		}

		// Set light state
		err := s.light.SetState(huego.State{
			On:  in.On,
			Bri: in.Bri,
			Hue: in.Hue,
			Sat: in.Sat,
		})

		if err != nil {
			logrus.WithError(err).WithField("in", in).Error("Unable to set light state")
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
