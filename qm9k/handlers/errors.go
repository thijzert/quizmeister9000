package handlers

import "fmt"

type errWrongRequestType struct{}

func (errWrongRequestType) Error() string {
	return "wrong request type"
}

func (errWrongRequestType) HTTPCode() int {
	return 400
}

type errRedirect struct {
	URL string
}

func (errRedirect) Error() string {
	return "you are being redirected to another page"
}

func (errRedirect) Headline() string {
	return "Redirecting..."
}

func (e errRedirect) Message() string {
	return fmt.Sprintf("You are being redirected to the address '%s'", e.URL)
}

func (errRedirect) HTTPCode() int {
	return 302
}

func (e errRedirect) RedirectLocation() string {
	return e.URL
}
