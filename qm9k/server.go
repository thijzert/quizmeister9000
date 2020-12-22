package qm9k

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
)

// A Config represents configuration parameters for the Server
type Config struct {
	// Base Directory For All The JSON Files
	BDFATFJF string
}

// A Server represents a HTTP handler that results in the quizmeister9000 interface
type Server struct {
	config Config

	sessionLock  sync.Mutex
	sessionStore map[SessionID]*Session

	userLock  sync.RWMutex
	userStore map[handlers.UserID]handlers.User

	// handle HTTP stuff
	mux *http.ServeMux

	parsedTemplates map[string]*template.Template
}

// NewServer instantiates a new Server instance based on the configuration
func NewServer(c Config) (*Server, error) {
	s := &Server{
		config: c,
		mux:    http.NewServeMux(),
	}

	err := s.initialiseSessionStore()
	if err != nil {
		return nil, err
	}

	err = s.initialiseUserStore()
	if err != nil {
		return nil, err
	}

	s.mux.Handle("/profile", s.HTMLFunc(handlers.ProfileHandler, handlers.ProfileDecoder, "full/profile"))
	s.mux.HandleFunc("/assets/", s.serveStaticAsset)
	// s.mux.HandleFunc("/", s.homeHandler)
	s.mux.Handle("/", s.HTMLFunc(handlers.HomeHandler, handlers.HomeDecoder, "full/home"))

	if !assetsEmbedded {
		// FIXME: find a nicer way of detecting a development version
		s.mux.Handle("/ui-showcase", s.HTMLFunc(handlers.UIShowcaseHandler, handlers.UIShowcaseDecoder, "full/ui-showcase"))
	}

	return s, nil
}

func (s *Server) getState(r *http.Request) handlers.State {
	var rv handlers.State

	ses := s.MaybeSession(r)
	if ses != nil && !ses.UserID.Empty() {
		rv.User, _ = s.getUser(ses.UserID)
	}

	return rv
}

// setState writes back any modified fields to the global state
func (s *Server) setState(st handlers.State) error {
	if st.User.UserID != 0 {
		s.saveUser(st.User)
	}
	return nil
}
