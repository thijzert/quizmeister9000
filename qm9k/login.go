package qm9k

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
)

type SessionID string

func newSessionID() SessionID {
	b := make([]byte, 16)
	io.ReadFull(rand.Reader, b)
	return SessionID(hex.EncodeToString(b))
}

type contextKey int

const (
	cKeySession contextKey = iota
)

// Session contains the longer-lived variables pertaining to one browser session
type Session struct {
	ID SessionID

	UserID handlers.UserID

	CreatedAt time.Time
	DeletedAt *time.Time
}

func newSession() *Session {
	rv := &Session{
		ID:     newSessionID(),
		UserID: handlers.NewUserID(),
	}

	return rv
}

// MaybeSession gets the session object for this request, if present.
func (s *Server) MaybeSession(r *http.Request) *Session {
	ctx := r.Context()
	if ses, ok := ctx.Value(cKeySession).(*Session); ok {
		return ses
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var session *Session

	existingSession := true

	c, _ := r.Cookie("TFSESSIONID")
	if c == nil {
		log.Printf("No session ID sent. Initiating new session.")
		existingSession = false
	} else {
		session = s.getSession(SessionID(c.Value))
		if session == nil {
			log.Printf("Session ID '%s' does not exist. Initiating new session.", c.Value)
			existingSession = false
		} else {
			if session.DeletedAt != nil {
				log.Printf("Session '%s' was deleted at %s. Initiating new session.", c.Value, session.DeletedAt)
				existingSession = false
			}
		}
	}

	if !existingSession {
		session = newSession()
		s.saveSession(session)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, cKeySession, session)

	c = &http.Cookie{
		Name:     "TFSESSIONID",
		Path:     "/",
		Value:    string(session.ID),
		Expires:  time.Now().Add(168 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	r = r.WithContext(ctx)

	s.mux.ServeHTTP(w, r)
}

func (s *Server) getSession(id SessionID) *Session {
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	rv, ok := s.sessionStore[id]
	if !ok {
		return nil
	}

	if rv.UserID.Empty() {
		rv.UserID = handlers.NewUserID()
		s.sessionStore[id] = rv
	}

	return rv
}

func (s *Server) initialiseSessionStore() error {
	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	s.sessionStore = make(map[SessionID]*Session)
	f, err := os.Open(path.Join(s.config.BDFATFJF, "session-store.json"))
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&s.sessionStore)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) saveSession(ses *Session) {
	if ses == nil {
		return
	}

	s.sessionLock.Lock()
	defer s.sessionLock.Unlock()

	s.sessionStore[ses.ID] = ses

	// TODO: swap out for proper database
	f, err := os.Create(path.Join(s.config.BDFATFJF, "session-store.json"))
	defer f.Close()
	if err != nil {
		log.Printf("error saving session store: %s", err)
		return
	}
	enc := json.NewEncoder(f)
	enc.Encode(s.sessionStore)
}

func (s *Server) getUser(id handlers.UserID) (handlers.User, bool) {
	if id.Empty() {
		return handlers.User{}, false
	}

	s.userLock.RLock()
	defer s.userLock.RUnlock()

	u, ok := s.userStore[id]
	u.UserID = id
	return u, ok
}

func (s *Server) initialiseUserStore() error {
	s.userLock.Lock()
	defer s.userLock.Unlock()

	s.userStore = make(map[handlers.UserID]handlers.User)
	f, err := os.Open(path.Join(s.config.BDFATFJF, "user-store.json"))
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&s.userStore)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) saveUser(u handlers.User) {
	if u.UserID.Empty() {
		return
	}

	s.userLock.Lock()
	defer s.userLock.Unlock()

	currentU, ok := s.userStore[u.UserID]
	if ok {
		if currentU == u {
			return
		}
	}

	log.Printf("saving user %+v", u)
	s.userStore[u.UserID] = u

	// TODO: swap out for proper database
	f, err := os.Create(path.Join(s.config.BDFATFJF, "user-store.json"))
	defer f.Close()
	if err != nil {
		log.Printf("error saving session store: %s", err)
		return
	}
	enc := json.NewEncoder(f)
	enc.Encode(s.userStore)
}
