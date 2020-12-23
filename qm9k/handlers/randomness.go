package handlers

import (
	"crypto/rand"
)

// randomString generates a random string composed of characters from a fixed
// alphabet. If alphabetAtEnds is not nil, the first and last character will
// come from that separate alphabet instead of the main one
func randomString(length int, alphabet, alphabetAtEnds []byte) string {
	if alphabetAtEnds == nil {
		alphabetAtEnds = alphabet
	}

	buf := make([]byte, length*2)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	var alpha []byte

	rv := make([]byte, 0, length)
	for _, c := range buf {
		alpha = alphabet
		if len(rv) == 0 || len(rv) == length-1 {
			alpha = alphabetAtEnds
		}

		max := len(alphabet) * (256 / len(alpha))
		cc := int(c)
		if cc >= max {
			continue
		}
		rv = append(rv, alpha[cc%len(alpha)])
		if len(rv) == length {
			return string(rv)
		}
	}

	// Try again
	return randomString(length, alphabet, alphabetAtEnds)
}

// NewQuizKey creates a new random QuizKey, with at least 56 bits of entropy.
// QuizKeys should have no ambiguous characters, and shouldn't start or end in a '.'.
func NewQuizKey() QuizKey {
	alphabetA := []byte("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXY01234567890-_.")
	alphabetB := []byte("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXY01234567890_")

	return QuizKey(randomString(11, alphabetA, alphabetB))
}

// NewAccessCode generates a new access code, with roughly 32 bits of entropy
func NewAccessCode() string {
	return randomString(7, []byte("abcdefghijkmnopqrstuvwxyz01234567890-"), []byte("abcdefghijkmnopqrstuvwxyz01234567890"))
}
