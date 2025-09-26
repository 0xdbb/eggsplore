package util

import (
	"fmt"
	"math/rand/v2"
	"time"
)

var (
	seed1 = uint64(time.Now().UnixNano())
	seed2 = rand.Uint64()
	r     = rand.New(rand.NewPCG(seed1, seed2))
)

// RandomString generates a random string of length n
func RandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const l = len(letterBytes)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.IntN(l)]
	}
	return string(b)
}

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min, max int) int {
	if max <= min {
		panic("max must be greater than min")
	}

	// Generate a random number in the range [min, max]
	return r.IntN(max-min+1) + min
}

// RandomPhone generates a random 10-digit phone number as a string.
func RandomPhone() string {
	// Ensure a valid range for a 10-digit number
	phone := RandomInt(1_000_000_000, 9_999_999_999)
	return fmt.Sprintf("%010d", phone)
}

func RandomEmail() string {
	domains := []string{"gmail", "yahoo", "outlook"}
	randomEmail := RandomString(6)
	return fmt.Sprintf("%s@%s.com", randomEmail, domains[r.IntN(len(domains))])
}
func RandomAddress() string {
	return RandomString(10)

}

func RandomName() string {
	return RandomString(5)
}

func RandomPassword() string {
	n := 10
	const letterBytes = "abcdefghijklmnopqrstuvwxyz12345678910ABCDEFGHIJKLMNOPQRSTUVWXYZ{/}*&^%$#@!_-=+"
	const l = len(letterBytes)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.IntN(l)]
	}
	return string(b)
}
