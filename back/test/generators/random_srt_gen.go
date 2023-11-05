package generators

import (
	"languago/internal/pkg/logger"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringRunes(n int) string {
	if n > 30 {
		n = 30
	}
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func NewPairs() logger.LogFields {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var pairs logger.LogFields = make(logger.LogFields, 30)

	for i := 0; i < rng.Intn(30); i++ {
		pairs[string(i)] = RandStringRunes(rng.Int())
	}
	return pairs
}
