package qm9k

import (
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
)

// A Config represents configuration parameters for the Server
type Config struct {
	// Base Directory For All The JSON Files
	BDFATFJF string

	// Set the 'secure' flag on all cookies (if the reverse proxy is unable to do so)
	SecureCookies bool
}

// A Server represents a HTTP handler that results in the quizmeister9000 interface
type Server struct {
	config Config

	sessionLock  sync.Mutex
	sessionStore map[SessionID]*Session

	userLock  sync.RWMutex
	userStore map[handlers.UserID]handlers.User

	quizLock        sync.RWMutex
	quizByAccesskey map[string]handlers.QuizKey
	quizStore       map[handlers.QuizKey]handlers.Quiz

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

	err = s.initialiseQuizStore()
	if err != nil {
		return nil, err
	}

	s.mux.Handle("/profile", s.HTMLFunc(handlers.ProfileHandler, handlers.ProfileDecoder, "full/profile"))
	s.mux.Handle("/new-quiz", s.HTMLFunc(handlers.NewQuizHandler, handlers.NewQuizDecoder, "full/new-quiz"))
	s.mux.Handle("/join-quiz", s.HTMLFunc(handlers.JoinHandler, handlers.JoinDecoder, "full/join-quiz"))

	s.mux.Handle("/quiz/", s.HTMLFunc(handlers.QuizViewerHandler, handlers.QuizViewerDecoder, "full/quiz-viewer"))

	s.mux.Handle("/grade-answers/", s.JSONFunc(handlers.GradeAnswerHandler, handlers.GradeAnswerDecoder))
	s.mux.Handle("/peer-status/", s.JSONFunc(handlers.PeerStatusHandler, handlers.PeerStatusDecoder))
	s.mux.Handle("/quiz-status/", s.JSONFunc(handlers.QuizStatusHandler, handlers.QuizStatusDecoder))
	s.mux.Handle("/set-answer/", s.JSONFunc(handlers.SetAnswerHandler, handlers.SetAnswerDecoder))
	s.mux.Handle("/vote-continue/", s.JSONFunc(handlers.VoteHandler, handlers.VoteDecoder))
	s.mux.Handle("/leaderboard/", s.JSONFunc(handlers.LeaderboardHandler, handlers.LeaderboardDecoder))

	s.mux.HandleFunc("/assets/", s.serveStaticAsset)
	// s.mux.HandleFunc("/", s.homeHandler)
	s.mux.Handle("/", s.HTMLFunc(handlers.HomeHandler, handlers.HomeDecoder, "full/home"))

	if !assetsEmbedded {
		// FIXME: find a nicer way of detecting a development version
		s.mux.Handle("/ui-showcase", s.HTMLFunc(handlers.UIShowcaseHandler, handlers.UIShowcaseDecoder, "full/ui-showcase"))
	}

	return s, nil
}

// appRoot finds the relative path to the application root
func (*Server) appRoot(r *http.Request) string {
	// Find the relative path for the application root by counting the number of slashes in the relative URL
	c := strings.Count(r.URL.Path, "/") - 1
	if c == 0 {
		return "./"
	}
	return strings.Repeat("../", c)
}

func (s *Server) getState(r *http.Request) handlers.State {
	var rv handlers.State

	ses := s.MaybeSession(r)
	if ses != nil && !ses.UserID.Empty() {
		rv.User, _ = s.getUser(ses.UserID)
	}

	var quiz handlers.Quiz
	var ok bool = false

	if r.PostFormValue("quiz-join-code") != "" {
		// A code was posted; add the corresponding quiz to the context state
		quiz, ok = s.getQuizByAccessKey(r.PostFormValue("quiz-join-code"))
	} else {
		spp := strings.Split(r.URL.Path, "/")
		if len(spp) > 2 && len(spp[2]) > 5 {
			// TODO: figure out a less hacky way of doing this
			quiz, ok = s.getQuiz(handlers.QuizKey(spp[2]))
		}
	}

	if ok {
		rv.Quiz = quiz

		for i, u := range rv.Quiz.Contestants {
			uu, uok := s.getUser(u.UserID)
			if uok {
				rv.Quiz.Contestants[i] = uu
			}
		}
	}

	return rv
}

// setState writes back any modified fields to the global state
func (s *Server) setState(st handlers.State) error {
	if !st.User.UserID.Empty() {
		s.saveUser(st.User)
	}
	if st.QuizDirty {
		s.saveQuiz(st.Quiz)
	}
	return nil
}
