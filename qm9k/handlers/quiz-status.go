package handlers

import "net/http"

type quizStatusRequest struct {
}

func (quizStatusRequest) FlaggedAsRequest() {}

// QuizStatusDecoder decodes a request for the quizStatus page
func QuizStatusDecoder(*http.Request) (Request, error) {
	return quizStatusRequest{}, nil
}

type quizStatusQ struct {
	Question string
	MyAnswer string
}

type quizStatusResponse struct {
	VotingEnabled bool
	MyVote        int
	QuizStatus    struct {
		Started  bool
		Grading  bool
		Finished bool
	}
	CurrentRound struct {
		RoundNo    int
		QuizMaster User
		ThisIsMe   bool
		Questions  []quizStatusQ
	}
}

func (quizStatusResponse) FlaggedAsResponse() {}

// QuizStatusHandler handles requests for the quizStatus page
func QuizStatusHandler(s State, _ Request) (State, Response, error) {
	var rv quizStatusResponse
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	rv.VotingEnabled = s.votingEnabled()
	if s.hasVoted(s.User.UserID) {
		rv.MyVote = 1
	}

	rv.CurrentRound.RoundNo = -1

	rv.QuizStatus.Started = s.Quiz.Started
	rv.QuizStatus.Finished = s.Quiz.Finished
	if s.Quiz.Started && !s.Quiz.Finished {
		rv.CurrentRound.RoundNo = s.Quiz.CurrentRound
		if s.Quiz.CurrentRound < len(s.Quiz.Rounds) {
			thisRound := s.Quiz.Rounds[s.Quiz.CurrentRound]

			// Set the quiz master for this round
			for _, u := range s.Quiz.Contestants {
				if thisRound.Quizmaster == u.UserID {
					rv.CurrentRound.QuizMaster = u
				}
			}
			rv.CurrentRound.ThisIsMe = (thisRound.Quizmaster == s.User.UserID)

			// Add the questions in this round, together with my answer
			for _, q := range thisRound.Questions {
				nq := quizStatusQ{
					Question: q.Question,
				}
				for _, ans := range q.Answers {
					if ans.UserID == s.User.UserID {
						nq.MyAnswer = ans.Answer
					}
				}
				rv.CurrentRound.Questions = append(rv.CurrentRound.Questions, nq)
			}
		} else {
			rv.QuizStatus.Grading = true
		}
	}

	return s, rv, nil
}
