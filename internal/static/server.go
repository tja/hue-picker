package static

//go:generate go run github.com/tdewolff/minify/v2/cmd/minify --recursive --output=. ../../etc/www

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed www
var www embed.FS

// Server wraps serving static files.
type Server struct {
	// Base handler
	M http.Handler
}

// NewServer will create a new Server object.
func NewServer(path string) (*Server, error) {
	// Detect file system to use
	var fsys http.FileSystem

	if path == "" {
		// Serve from the embedded file system
		sub, err := fs.Sub(www, "www")
		if err != nil {
			return nil, fmt.Errorf("embedded sub-system: %w", err)
		}

		fsys = http.FS(sub)
	} else {
		// Serve from the given path
		fsys = http.Dir(path)
	}

	// Return full server object
	return &Server{
		M: http.FileServer(fsys),
	}, nil
}
