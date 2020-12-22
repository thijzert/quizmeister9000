package handlers

import "net/http"

// UIShowcaseDecoder decodes a request for the home page
func UIShowcaseDecoder(*http.Request) (Request, error) {
	return emptyRequest{}, nil
}

type uiShowcaseResponse struct {
}

func (uiShowcaseResponse) FlaggedAsResponse() {}

// UIShowcaseHandler handles requests for the home page
func UIShowcaseHandler(s State, _ Request) (State, Response, error) {
	return s, uiShowcaseResponse{}, nil
}
