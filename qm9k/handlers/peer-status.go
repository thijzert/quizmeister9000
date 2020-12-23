package handlers

import (
	"net/http"
	"time"
)

type peerStatusRequest struct {
}

func (peerStatusRequest) FlaggedAsRequest() {}

// PeerStatusDecoder decodes a request for the peerStatus page
func PeerStatusDecoder(*http.Request) (Request, error) {
	return peerStatusRequest{}, nil
}

type peerStatus struct {
	UserID UserID
	Nick   string
	Colour string
	Quest  string
	Status string
	Voted  bool
}

type peerStatusResponse struct {
	Peers []peerStatus
}

func (peerStatusResponse) FlaggedAsResponse() {}

// PeerStatusHandler handles requests for the peerStatus page
func PeerStatusHandler(s State, r Request) (State, Response, error) {
	var rv peerStatusResponse
	err := s.checkContestant()
	if err != nil {
		return s, rv, err
	}

	someoneHasTyped := false
	var answers []answer
	if len(s.Quiz.Rounds) > s.Quiz.CurrentRound {
		cr := s.Quiz.Rounds[s.Quiz.CurrentRound]
		if len(cr.Questions) > cr.CurrentQuestion {
			answers = cr.Questions[cr.CurrentQuestion].Answers
		}
	}
	if answers != nil {
		for _, ans := range answers {
			if ans.Answer != "" || !ans.Timestamp.IsZero() {
				someoneHasTyped = true
			}
		}
	}

	rv.Peers = make([]peerStatus, 0, len(s.Quiz.Contestants)-1)
	for _, peer := range s.Quiz.Contestants {
		if peer.UserID == s.User.UserID {
			continue
		}
		pst := peerStatus{
			UserID: peer.UserID,
			Nick:   peer.Nick,
			Colour: peer.Colour,
			Quest:  peer.Quest,
			Status: "neutral",
			Voted:  s.hasVoted(peer.UserID),
		}

		if answers != nil {
			for _, ans := range answers {
				if ans.UserID != peer.UserID {
					continue
				}
				if time.Since(ans.Timestamp) < 2500*time.Millisecond {
					pst.Status = "typing"
				} else if ans.Answer != "" {
					pst.Status = "done"
				} else if someoneHasTyped {
					pst.Status = "thinking"
				}
			}
		}

		rv.Peers = append(rv.Peers, pst)
	}

	return s, rv, nil
}
