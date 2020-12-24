package handlers

import (
	"net/http"
	"sort"
)

type leaderboardRequest struct {
}

func (leaderboardRequest) FlaggedAsRequest() {}

// LeaderboardDecoder decodes a request for the leaderboard page
func LeaderboardDecoder(*http.Request) (Request, error) {
	return leaderboardRequest{}, nil
}

type peerScore struct {
	Nick   string
	Colour string
	Quest  string
	Status string
	Score  int
}

type leaderboardResponse struct {
	Peers []peerScore
}

func (leaderboardResponse) FlaggedAsResponse() {}

// LeaderboardHandler handles requests for the leaderboard page
func LeaderboardHandler(s State, r Request) (State, Response, error) {
	var rv leaderboardResponse
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}
	if !s.Quiz.Finished {
		return s, rv, errForbidden("Forbidden", "The quiz is still going")
	}

	rv.Peers = make([]peerScore, 0, len(s.Quiz.Contestants))
	for _, peer := range s.Quiz.Contestants {
		pst := peerScore{
			Nick:   peer.Nick,
			Colour: peer.Colour,
			Quest:  peer.Quest,
			Status: "neutral",
			Score:  0,
		}

		for _, r := range s.Quiz.Rounds {
			for _, q := range r.Questions {
				for _, a := range q.Answers {
					if a.UserID == peer.UserID {
						pst.Score += a.Score
					}
				}
			}
		}

		rv.Peers = append(rv.Peers, pst)
	}

	sort.Slice(rv.Peers, func(i, j int) bool {
		return rv.Peers[i].Score > rv.Peers[j].Score
	})

	return s, rv, nil
}
