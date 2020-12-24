package handlers

import (
	"net/http"
	"time"

	"github.com/thijzert/speeldoos/lib/properrandom"
)

// The State struct represents the current state of the world
type State struct {
	User      User
	QuizDirty bool
	Quiz      Quiz
}

// A UserID kind-of-uniquely identifies a user
type UserID uint64

// Empty tests whether or not a UserID is set
func (u UserID) Empty() bool {
	return u == 0
}

// NewUserID generates a new user ID
func NewUserID() UserID {
	return UserID(properrandom.Uint64())
}

// A User represents a quiz-taker
type User struct {
	UserID UserID
	Nick   string
	Quest  string `json:",omitempty"`
	Colour string `json:",omitempty"`
	Admin  bool   `json:",omitempty"`
}

// Empty tests whether or not a User has any fields set
func (u User) Empty() bool {
	return u.Nick == "" && u.Quest == ""
}

// A QuizKey is a unique identifier for a video
type QuizKey string

// Empty tests whether or not a Quiz key is set
func (q QuizKey) Empty() bool {
	return string(q) == ""
}

type round struct {
	Quizmaster      UserID
	CurrentQuestion int
	Questions       []question
}

type question struct {
	Question string
	Answers  []answer
}

type answer struct {
	UserID    UserID
	Answer    string
	Timestamp time.Time
	Scored    bool
	Score     int
}

// A Quiz wraps the state of one quiz
type Quiz struct {
	QuizKey      QuizKey
	AccessCode   string
	Started      bool
	Finished     bool
	CurrentRound int
	Contestants  []User
	Votes        []UserID `json:",omitempty"`
	Rounds       []round
}

// Equals tests if two quizzes are the same
func (q Quiz) Equals(b Quiz) bool {
	if q.QuizKey != b.QuizKey || q.AccessCode != b.AccessCode {
		return false
	}
	if q.Started != b.Started || q.Finished != b.Finished || q.CurrentRound != b.CurrentRound {
		return false
	}

	if q.Contestants == nil || b.Contestants == nil || len(q.Contestants) != len(b.Contestants) {
		return false
	}
	for i, u := range q.Contestants {
		if u.UserID != b.Contestants[i].UserID {
			return false
		}
	}

	if (q.Votes == nil && b.Votes != nil) || (q.Votes != nil && b.Votes == nil) || len(q.Votes) != len(b.Votes) {
		return false
	}
	for i, quid := range q.Votes {
		if quid != b.Votes[i] {
			return false
		}
	}

	if (q.Rounds == nil && b.Rounds != nil) || (q.Rounds != nil && b.Rounds == nil) || len(q.Rounds) != len(b.Rounds) {
		return false
	}
	for i, round := range q.Rounds {
		bround := b.Rounds[i]
		if round.Quizmaster != bround.Quizmaster {
			return false
		}

		if (round.Questions == nil && bround.Questions != nil) || (round.Questions != nil && bround.Questions == nil) || len(round.Questions) != len(bround.Questions) {
			return false
		}

		for j, qq := range round.Questions {
			bq := bround.Questions[j]

			if qq.Question != bq.Question {
				return false
			}
			if (qq.Answers == nil && bq.Answers != nil) || (qq.Answers != nil && bq.Answers == nil) || len(qq.Answers) != len(bq.Answers) {
				return false
			}

			for k, qqa := range qq.Answers {
				if qqa != bq.Answers[k] {
					return false
				}
			}
		}
	}

	return true
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
