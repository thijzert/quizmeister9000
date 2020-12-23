package handlers

import (
	"log"

	"github.com/thijzert/speeldoos/lib/properrandom"
)

func (s State) votingEnabled() bool {
	if !s.Quiz.Started {
		return true
	}

	return true
}

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

func (s State) everyoneHasVoted() bool {
	for _, u := range s.Quiz.Contestants {
		if !s.hasVoted(u.UserID) {
			return false
		}
	}

	return true
}

func advanceQuiz(quiz Quiz) Quiz {
	// If the quiz hadn't started yet, proceed to the first round (that's -1 + 1 = index 0)
	if !quiz.Started {
		quiz.Started = true
		quiz.CurrentRound = 0

		// fall through into next block
	}

	if quiz.CurrentRound < len(quiz.Rounds) {
		quiz.CurrentRound++
		if quiz.CurrentRound < len(quiz.Rounds) {
			properrandom.Shuffle(len(quiz.Rounds[quiz.CurrentRound:]), func(i, j int) {
				quiz.Rounds[quiz.CurrentRound+i], quiz.Rounds[quiz.CurrentRound+j] = quiz.Rounds[quiz.CurrentRound+j], quiz.Rounds[quiz.CurrentRound+i]
			})

			if len(quiz.Rounds[quiz.CurrentRound].Questions) == 0 {
				// FIXME: remove hardcoded 10 questions per round
				for i := 0; i < 2; i++ {
					quiz.Rounds[quiz.CurrentRound].Questions = append(quiz.Rounds[quiz.CurrentRound].Questions, question{})
				}
			}
		}
		return quiz
	}

	log.Printf("next stage ill-defined")
	return quiz
}
