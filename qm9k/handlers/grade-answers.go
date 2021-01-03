package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type gradeAnswerRequest struct {
	Grading  bool
	Question int
	Answer   string
	Score    int
}

func (gradeAnswerRequest) FlaggedAsRequest() {}

// GradeAnswerDecoder decodes a request for the gradeAnswer page
func GradeAnswerDecoder(r *http.Request) (Request, error) {
	var rv gradeAnswerRequest
	var err error

	if r.Method == "POST" {
		rv.Grading = true

		rv.Question, err = strconv.Atoi(r.PostFormValue("question"))
		if err != nil {
			return rv, err
		}
		if rv.Question < 0 {
			return rv, errors.New("invalid question number")
		}
		rv.Score, err = strconv.Atoi(r.PostFormValue("score"))
		if err != nil {
			return rv, err
		}
		if rv.Score < 0 || rv.Score > 2 {
			return rv, errors.New("score should be 0, 1, or 2")
		}

		rv.Answer = r.PostFormValue("answer")
	}

	return rv, nil
}

type gradedAnswer struct {
	Answer string
	Scored bool
	Score  int
}

type gradeable struct {
	Question string
	Answers  []gradedAnswer
}

type gradeAnswerResponse struct {
	Questions []gradeable
}

func (gradeAnswerResponse) FlaggedAsResponse() {}

// GradeAnswerHandler handles requests for the gradeAnswer page
func GradeAnswerHandler(s State, r Request) (State, Response, error) {
	var rv gradeAnswerResponse
	req, ok := r.(gradeAnswerRequest)
	if !ok {
		return s, rv, errWrongRequestType{}
	}
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	if !s.Quiz.Started || s.Quiz.Finished {
		return s, rv, errors.New("quiz not active")
	}
	if s.Quiz.CurrentRound < len(s.Quiz.Rounds) {
		return s, rv, errors.New("quiz still active")
	}

	myRoundIdx := -1
	for i, round := range s.Quiz.Rounds {
		if round.Quizmaster == s.User.UserID {
			myRoundIdx = i
		}
	}
	if myRoundIdx < 0 {
		return s, rv, errNotFound("Not found", "You don't appear to have been a quiz master")
	}
	questions := s.Quiz.Rounds[myRoundIdx].Questions

	if req.Grading {
		if req.Question >= len(questions) {
			return s, rv, errNotFound("Not found", "question not found")
		}

		for i, ans := range questions[req.Question].Answers {
			if normaliseNewlines(ans.Answer) == normaliseNewlines(req.Answer) {
				questions[req.Question].Answers[i].Score = req.Score
				questions[req.Question].Answers[i].Scored = true
			}
		}

		s.Quiz.Rounds[myRoundIdx].Questions = questions
		s.QuizDirty = true
	}

	for _, q := range questions {
		answers := make(map[string]gradedAnswer)
		for _, ans := range q.Answers {
			answers[normaliseNewlines(ans.Answer)] = gradedAnswer{
				Answer: ans.Answer,
				Scored: ans.Scored,
				Score:  ans.Score,
			}
		}

		grd := gradeable{
			Question: q.Question,
		}
		for _, ans := range answers {
			grd.Answers = append(grd.Answers, ans)
		}
		rv.Questions = append(rv.Questions, grd)
	}

	return s, rv, nil
}

func normaliseNewlines(str string) string {
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\n", " ", -1)
	return str
}
