package random

import (
	"math/rand"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int, src *rand.Rand) string {
	b := RandBytes(n, src)
	return *(*string)(unsafe.Pointer(&b))
}

func RandBytes(n int, src *rand.Rand) []byte {
	int63 := func() int64 {
		if src != nil {
			return src.Int63()
		}
		return rand.Int63()
	}
	b := make([]byte, n)
	// A int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

// Dur returns a pseudo-random Duration in [0, max)
func Dur(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}

// Uniformly jitters the provided duration by +/- 10%.
func Jitter(period time.Duration) time.Duration {
	return JitterFraction(period, .9)
}

// Uniformly jitters the provided duration by +/- the given fraction.  NOTE:
// fraction must be in (0, 1].
func JitterFraction(period time.Duration, fraction float64) time.Duration {
	fixed := time.Duration(float64(period) * (1 - fraction))
	return fixed + Dur(2*(period-fixed))
}
