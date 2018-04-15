package storage

import (
	"math/rand"
	"time"
)

const Nothing = "I got nothing"

type Storager interface {
	Add(string, string)
	Get(string) string
}

func Random(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
