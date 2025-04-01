package util

import (
	"math/rand"
	"strings"
	"time"
)

const allowedChars = "abcdefghijklmopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between max and min
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(allowedChars)

	for range n {
		c := allowedChars[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(5) + ".com"
}

// RandomMoney generates a random ammount
func RandomMoney() int64 {
	return RandomInt(0, 1000000)
}

// RandomCurrency provides a random currency name
func RandomCurrency() string {
	currencies := []string{INR, USD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
