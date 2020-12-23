package qm9k

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
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
		Secure:   s.config.SecureCookies,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	r = r.WithContext(ctx)

	s.mux.ServeHTTP(w, r)
}
