package qm9k

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"unicode/utf8"
)

type httpError struct {
	StatusCode int
	Cause      error
}

func errorWithStatus(e error, c int) error {
	if e == nil {
		return nil
	}

	return httpError{
		StatusCode: c,
		Cause:      e,
	}
}

func (e httpError) Error() string {
	return e.Cause.Error()
}

func (e httpError) Unwrap() error {
	return e.Cause
}

// A UserError is an error that may be shown to end users
type UserError interface {
	error
	Headline() string
	Message() string
}

type userError struct {
	headline string
	message  string
	cause    error
}

func errorWithMessage(e error, headline, message string) error {
	if e == nil {
		return nil
	}

	return userError{
		headline: headline,
		message:  message,
		cause:    e,
	}
}

func (e userError) Error() string {
	return e.headline + ": " + e.message
}

func (e userError) Unwrap() error {
	return e.cause
}

func (e userError) Headline() string {
	return e.headline
}

func (e userError) Message() string {
	return e.message
}

// UserHeadline returns the user-visible headline for an error, if it exists
func UserHeadline(e error) string {
	if e == nil {
		return ""
	} else if ue, ok := e.(UserError); ok {
		return ue.Headline()
	} else if c := errors.Unwrap(e); c != nil {
		return UserHeadline(c)
	} else {
		return "An error occurred"
	}
}

// UserMessage returns the user-visible message for an error, if it exists
func UserMessage(e error) string {
	if e == nil {
		return ""
	} else if ue, ok := e.(UserError); ok {
		return ue.Message()
	} else if c := errors.Unwrap(e); c != nil {
		return UserMessage(c)
	} else {
		return "An error occurred"
	}
}

func (s *Server) Error(w http.ResponseWriter, r *http.Request, e error) {
	if e == nil {
		w.Write([]byte("No errors were found."))
		return
	}

	httpcode := httpError{
		StatusCode: 400,
		Cause:      e,
	}
	if !errors.As(e, &httpcode) {
		httpcode.StatusCode = 400
		httpcode.Cause = e
	}

	log.Print(e)
	w.WriteHeader(400)
	w.Header()["Content-Type"] = []string{"text/plain"}
	w.Write([]byte("TODO: better error handling\n"))

	if !assetsEmbedded {
		// FIXME: find a nicer way of detecting a development version
		w.Write([]byte(e.Error()))
	}
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request, url string, statuscode int) {
	log.Printf("redirecting to '%s'", url)

	if len(url) == 0 {
		url = "."
	} else if url[0] != '/' && url[0] != '.' {
		url = "./" + url
	}

	// TODO: some more filtering

	h := w.Header()

	// RFC 7231 notes that a short HTML body is usually included in
	// the response because older user agents may not understand 301/307.
	// Do it only if the request didn't already have a Content-Type header.
	_, hadCT := h["Content-Type"]

	h.Set("Location", hexEscapeNonASCII(url))
	if !hadCT && (r.Method == "GET" || r.Method == "HEAD") {
		h.Set("Content-Type", "text/html; charset=utf-8")
	}
	w.WriteHeader(statuscode)

	// Shouldn't send the body for POST or HEAD; that leaves GET.
	if !hadCT && r.Method == "GET" {
		body := "<a href=\"" + html.EscapeString(url) + "\">We're trying to move you to a new location</a>.\n"
		fmt.Fprintln(w, body)
	}
}

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}
func (s *Server) continueChain(w http.ResponseWriter, r *http.Request) {
	cont := "."
	if cc := r.FormValue("continue"); cc != "" {
		u, err := url.Parse(cc)
		if err == nil && u.Scheme == "" && u.User == nil && u.Host == "" {
			cont = cc
		}
	}
	s.redirect(w, r, cont, http.StatusFound)
}
