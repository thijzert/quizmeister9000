package handlers

import "net/http"

type joinRequest struct {
	CodePosted bool
}

func (joinRequest) FlaggedAsRequest() {}

// JoinDecoder decodes a request for the join page
func JoinDecoder(r *http.Request) (Request, error) {
	var rv joinRequest

	rv.CodePosted = r.PostFormValue("quiz-join-code") != ""

	return rv, nil
}

type joinResponse struct {
	Placeholder string
	Error       string
}

func (joinResponse) FlaggedAsResponse() {}

// JoinHandler handles requests for the join page
func JoinHandler(s State, r Request) (State, Response, error) {
	var rv joinResponse
	req, ok := r.(joinRequest)
	if !ok {
		return s, rv, errWrongRequestType{}
	}

	if !s.Quiz.QuizKey.Empty() {
		err := s.checkContestant()
		if err != nil {
			// The user isn't a contestant yet - check if we're still okay to join
			if s.Quiz.Started {
				rv.Error = "The quiz you're trying to join has already started"
				return s, rv, nil
			}

			s.Quiz.Contestants = append(s.Quiz.Contestants, s.User)

			// FIXME: Everybody's a quiz master; make configurable
			s.Quiz.Rounds = append(s.Quiz.Rounds, round{
				Quizmaster: s.User.UserID,
			})
		}

		return s, rv, errRedirect{"quiz/" + string(s.Quiz.QuizKey)}
	} else if req.CodePosted {
		rv.Error = "The code you entered was not recognised"
	}

	rv.Placeholder = NewAccessCode()
	return s, rv, nil
}
