package qm9k

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/thijzert/quizmeister9000/qm9k/handlers"
)

func (s *Server) initialiseQuizStore() error {
	s.quizLock.Lock()
	defer s.quizLock.Unlock()

	s.quizStore = make(map[handlers.QuizKey]handlers.Quiz)
	s.quizByAccesskey = make(map[string]handlers.QuizKey)

	f, err := os.Open(path.Join(s.config.BDFATFJF, "quiz-store.json"))
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&s.quizStore)
	if err != nil {
		return err
	}

	for qkey, quiz := range s.quizStore {
		if !quiz.Started && quiz.AccessCode != "" {
			s.quizByAccesskey[quiz.AccessCode] = qkey
		}
	}

	return nil
}

func (s *Server) getQuiz(id handlers.QuizKey) (handlers.Quiz, bool) {
	if len(id) < 5 {
		return handlers.Quiz{}, false
	}

	s.quizLock.RLock()
	defer s.quizLock.RUnlock()

	q, ok := s.quizStore[id]
	q.QuizKey = id
	return q, ok
}

func (s *Server) getQuizByAccessKey(accessKey string) (handlers.Quiz, bool) {
	s.quizLock.RLock()
	id := s.quizByAccesskey[accessKey]
	s.quizLock.RUnlock()

	return s.getQuiz(id)
}

func (s *Server) saveQuiz(q handlers.Quiz) {
	if len(q.QuizKey) < 5 {
		return
	}

	s.quizLock.Lock()
	defer s.quizLock.Unlock()

	log.Printf("saving quiz %s", q.QuizKey)
	s.quizStore[q.QuizKey] = q

	if q.AccessCode != "" {
		if q.Started {
			delete(s.quizByAccesskey, q.AccessCode)
		} else {
			s.quizByAccesskey[q.AccessCode] = q.QuizKey
		}
	}

	// TODO: swap out for proper database
	f, err := os.Create(path.Join(s.config.BDFATFJF, "quiz-store.json"))
	defer f.Close()
	if err != nil {
		log.Printf("error saving quiz store: %s", err)
		return
	}

	enc := json.NewEncoder(f)
	enc.Encode(s.quizStore)
}
