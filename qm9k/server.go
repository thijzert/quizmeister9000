package qm9k

import (
	"net/http"
	"sync"
)

// A Config represents configuration parameters for the Server
type Config struct {
}

// A Server represents a HTTP handler that results in the quizmeister9000 interface
type Server struct {
	config Config

	sessionLock  sync.Mutex
	sessionStore map[SessionID]*Session

	// handle HTTP stuff
	mux *http.ServeMux
}

// NewServer instantiates a new Server instance based on the configuration
func NewServer(c Config) (*Server, error) {
	s := &Server{
		config:       c,
		sessionStore: make(map[SessionID]*Session),
		mux:          http.NewServeMux(),
	}

	s.mux.HandleFunc("/party/", s.serveChat)
	s.mux.HandleFunc("/assets/", serveAsset)
	s.mux.HandleFunc("/", s.homeHandler)

	if !assetsEmbedded {
		// FIXME: find a nicer way of detecting a development version
		s.mux.HandleFunc("/ui-showcase", s.uishowcaseHandler)
	}

	return s, nil
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	var homeData struct {
	}

	s.executeTemplate("home", homeData, w, r)
}

func (s *Server) uishowcaseHandler(w http.ResponseWriter, r *http.Request) {
	s.executeTemplate("ui-showcase", nil, w, r)
}

func (s *Server) serveChat(w http.ResponseWriter, r *http.Request) {
	s.executeTemplate("chat", struct{}{}, w, r)
}
