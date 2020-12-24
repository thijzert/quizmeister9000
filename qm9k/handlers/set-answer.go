package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type setAnswerRequest struct {
	Round    int
	Question int
	Text     string
}

func (setAnswerRequest) FlaggedAsRequest() {}

// SetAnswerDecoder decodes a request for the setAnswer page
func SetAnswerDecoder(r *http.Request) (Request, error) {
	var rv setAnswerRequest
	var err error

	rv.Round, err = strconv.Atoi(r.PostFormValue("round"))
	if err != nil {
		return rv, err
	}
	rv.Question, err = strconv.Atoi(r.PostFormValue("question"))
	if err != nil {
		return rv, err
	}

	rv.Text = r.PostFormValue("text")

	return rv, nil
}

type setAnswerResponse struct {
}

func (setAnswerResponse) FlaggedAsResponse() {}

// SetAnswerHandler handles requests for the setAnswer page
func SetAnswerHandler(s State, r Request) (State, Response, error) {
	var rv setAnswerResponse
	req, ok := r.(setAnswerRequest)
	if !ok {
		return s, rv, errWrongRequestType{}
	}
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	if !s.Quiz.Started {
		return s, rv, fmt.Errorf("the quiz hasn't started yet")
	}
	if s.Quiz.Finished || s.Quiz.CurrentRound >= len(s.Quiz.Rounds) {
		return s, rv, fmt.Errorf("it's over, man")
	}

	if s.Quiz.CurrentRound != req.Round {
		return s, rv, fmt.Errorf("wrong round")
	}

	thisRound := s.Quiz.Rounds[s.Quiz.CurrentRound]
	if req.Question < 0 || req.Question >= len(thisRound.Questions) {
		return s, rv, fmt.Errorf("invalid question id")
	}

	if thisRound.Quizmaster == s.User.UserID {
		s.Quiz.Rounds[s.Quiz.CurrentRound].Questions[req.Question].Question = req.Text
	} else {
		found := false
		for i, ans := range thisRound.Questions[req.Question].Answers {
			if ans.UserID == s.User.UserID {
				found = true
				s.Quiz.Rounds[s.Quiz.CurrentRound].Questions[req.Question].Answers[i].Answer = req.Text
				s.Quiz.Rounds[s.Quiz.CurrentRound].Questions[req.Question].Answers[i].Timestamp = time.Now()
			}
		}
		if !found {
			ans := answer{
				UserID:    s.User.UserID,
				Answer:    req.Text,
				Timestamp: time.Now(),
			}
			s.Quiz.Rounds[s.Quiz.CurrentRound].Questions[req.Question].Answers = append(s.Quiz.Rounds[s.Quiz.CurrentRound].Questions[req.Question].Answers, ans)
		}
	}

	s.QuizDirty = true

	return s, rv, nil
}
