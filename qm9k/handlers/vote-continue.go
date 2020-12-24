package handlers

import (
	"log"
	"net/http"
	"strconv"
)

type voteRequest struct {
	Voted bool
	Vote  int
}

func (voteRequest) FlaggedAsRequest() {}

// VoteDecoder decodes a request for the vote page
func VoteDecoder(r *http.Request) (Request, error) {
	var rv voteRequest

	voteStr := r.PostFormValue("vote")
	if voteStr != "" {
		rv.Voted = true
		rv.Vote, _ = strconv.Atoi(voteStr)
	}
	return rv, nil
}

type voteResponse struct {
	MyVote        int
	VotingEnabled bool
}

func (voteResponse) FlaggedAsResponse() {}

// VoteHandler handles requests for the vote page
func VoteHandler(s State, r Request) (State, Response, error) {
	rv := voteResponse{
		VotingEnabled: s.votingEnabled(),
	}
	req, ok := r.(voteRequest)
	if !ok {
		return s, rv, errWrongRequestType{}
	}
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	if req.Voted && rv.VotingEnabled {
		s.QuizDirty = true

		idx := -1
		if s.Quiz.Votes != nil {
			for i, uid := range s.Quiz.Votes {
				if uid == s.User.UserID {
					idx = i
				}
			}
		}

		if req.Vote == 0 {
			if idx >= 0 {
				s.Quiz.Votes = append(s.Quiz.Votes[:idx], s.Quiz.Votes[idx+1:]...)
			}
		} else {
			if idx == -1 {
				s.Quiz.Votes = append(s.Quiz.Votes, s.User.UserID)
			}
		}
	}

	if s.hasVoted(s.User.UserID) {
		rv.MyVote = 1
	}

	// Sneaky hack: since it's nice to have button feedback, calculate global
	// changes to the quiz structure after setting the output

	if s.everyoneHasVoted() {
		// Advance the main game loop

		if !s.Quiz.Started {
			s.Quiz.Started = true
			s.Quiz.CurrentRound = -1
		}

		log.Printf("everyone has voted in quiz '%s'", s.Quiz.QuizKey)
		s.Quiz = advanceQuiz(s.Quiz)

		s.Quiz.Votes = s.Quiz.Votes[:0]
	}

	return s, rv, nil
}
