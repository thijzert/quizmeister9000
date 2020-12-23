package handlers

import "net/http"

type emptyRequest struct{}

func (emptyRequest) FlaggedAsRequest() {}

// HomeDecoder decodes a request for the home page
func HomeDecoder(*http.Request) (Request, error) {
	return emptyRequest{}, nil
}

type homeResponse struct {
}

func (homeResponse) FlaggedAsResponse() {}

// HomeHandler handles requests for the home page
func HomeHandler(s State, _ Request) (State, Response, error) {
	if s.User.Nick == "" {
		return s, homeResponse{}, errRedirect{"profile?new=1"}
	}
	return s, homeResponse{}, nil
}
