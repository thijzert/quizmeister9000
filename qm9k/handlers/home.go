package handlers

import "net/http"

type emptyRequest struct{}

func (emptyRequest) FlaggedAsRequest() {}

type homeRequest struct {
	Path string
}

func (homeRequest) FlaggedAsRequest() {}

// HomeDecoder decodes a request for the home page
func HomeDecoder(r *http.Request) (Request, error) {
	return homeRequest{
		Path: r.URL.Path,
	}, nil
}

type homeResponse struct {
}

func (homeResponse) FlaggedAsResponse() {}

// HomeHandler handles requests for the home page
func HomeHandler(s State, r Request) (State, Response, error) {
	var rv homeResponse

	req, ok := r.(homeRequest)
	if !ok {
		return s, rv, errWrongRequestType{}
	}

	if req.Path != "/" {
		return s, rv, errNotFound("", "")
	}

	if s.User.Nick == "" {
		return s, rv, errRedirect{"profile?new=1"}
	}

	return s, rv, nil
}
