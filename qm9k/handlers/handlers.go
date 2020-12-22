package handlers

import (
	"net/http"

	"github.com/thijzert/speeldoos/lib/properrandom"
)

// The State struct represents the current state of the world
type State struct {
	User User
	Quiz Quiz
}

type UserID uint64

func (u UserID) Empty() bool {
	return u == 0
}

func NewUserID() UserID {
	return UserID(properrandom.Uint64())
}

// A User represents a quiz-taker
type User struct {
	UserID UserID
	Nick   string
	Quest  string
	Colour string
}

// A Quiz wraps the state of one quiz
type Quiz struct {
	Started      bool
	Finished     bool
	CurrentRound int
	Contestants  []User
	Rounds       []struct {
		Quizmaster UserID
		Questions  []struct {
			Question string
			Answers  []struct {
				UserID UserID
				Answer string
				Scored bool
				Score  int
			}
		}
	}
}

var (
	// ErrParser is thrown when a request object is of the wrong type
	ErrParser error = errParse{}
)

type errParse struct{}

func (errParse) Error() string {
	return "parse error while decoding request"
}

// A Request flags any request type
type Request interface {
	FlaggedAsRequest()
}

// A Response flags any response type
type Response interface {
	FlaggedAsResponse()
}

// A RequestDecoder turns a HTTP request into a domain-specific request type
type RequestDecoder func(*http.Request) (Request, error)

// A RequestHandler is a monadic definition of a request handler. The inputs are
// the current state of the world, and a handler-specific request type, and the
// output is the new state of the world (which may or may not be the same), a
// handler-specific response type, and/or an error.
type RequestHandler func(State, Request) (State, Response, error)

type errWrongRequestType struct{}

func (errWrongRequestType) Error() string {
	return "wrong request type"
}

func (errWrongRequestType) HTTPCode() int {
	return 400
}
