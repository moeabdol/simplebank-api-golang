package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates random int64 number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner function
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney function
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency function
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "SAR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
