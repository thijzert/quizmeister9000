package qm9k

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
)

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
