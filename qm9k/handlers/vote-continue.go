package handlers

import (
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
	MyVote int
}

func (voteResponse) FlaggedAsResponse() {}

func (s State) hasVoted(u UserID) bool {
	if s.Quiz.Votes == nil {
		return false
	}

	for _, uid := range s.Quiz.Votes {
		if uid == u {
			return true
		}
	}

	return false
}

// VoteHandler handles requests for the vote page
func VoteHandler(s State, r Request) (State, Response, error) {
	var rv voteResponse
	req, ok := r.(voteRequest)
	if !ok {
		return s, rv, errWrongRequestType{}
	}
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	if req.Voted {
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

	return s, rv, nil
}
