package handlers

import (
	"net/http"
)

type quizViewerRequest struct {
}

func (quizViewerRequest) FlaggedAsRequest() {}

// QuizViewerDecoder decodes a request for the quizViewer page
func QuizViewerDecoder(*http.Request) (Request, error) {
	return emptyRequest{}, nil
}

type quizViewerResponse struct {
	Quiz Quiz
}

func (quizViewerResponse) FlaggedAsResponse() {}

func (s State) checkContestant() error {
	if s.User.UserID.Empty() || s.Quiz.QuizKey.Empty() {
		return errNotFound("Not Found", "This quiz does not exist")
	}

	for _, uu := range s.Quiz.Contestants {
		if uu.UserID == s.User.UserID {
			return nil
		}
	}

	if s.Quiz.Started {
		return errForbidden("Access denied", "You are not a participant in this quiz")
	}

	return errForbidden("Access denied", "You aren't a participant in this quiz, but you can still ask to be let in")
}

// QuizViewerHandler handles requests for the quizViewer page
func QuizViewerHandler(s State, _ Request) (State, Response, error) {
	var rv quizViewerResponse
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	rv.Quiz = s.Quiz
	return s, rv, nil
}
