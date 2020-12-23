package qm9k

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
)

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
	} else if u.Empty() {
		return
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
