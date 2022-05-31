package random

import (
	"math/rand"
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
