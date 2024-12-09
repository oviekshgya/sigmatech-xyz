package pkg

import (
	"math/rand"
	"time"
)

func KodeVerify(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letters = []rune("1234567890")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
