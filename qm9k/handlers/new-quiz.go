package handlers

import "net/http"

type newQuizRequest struct {
}

func (newQuizRequest) FlaggedAsRequest() {}

// NewQuizDecoder decodes a request for the newQuiz page
func NewQuizDecoder(*http.Request) (Request, error) {
	return emptyRequest{}, nil
}

type newQuizResponse struct {
}

func (newQuizResponse) FlaggedAsResponse() {}

// NewQuizHandler handles requests for the newQuiz page
func NewQuizHandler(s State, _ Request) (State, Response, error) {
	if !s.User.Admin {
		return s, newQuizResponse{}, errForbidden("Access Denied", "You don't have permission to start new quizzes")
	}

	s.Quiz = Quiz{
		QuizKey:     NewQuizKey(),
		AccessCode:  NewAccessCode(),
		Contestants: []User{s.User},
		Rounds: []round{{
			Quizmaster: s.User.UserID,
			Questions:  nil,
		}},
	}

	s.QuizDirty = true

	return s, newQuizResponse{}, errRedirect{"quiz/" + string(s.Quiz.QuizKey)}
}
