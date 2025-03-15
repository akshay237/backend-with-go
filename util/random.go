package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Random Int generates a random integer b/w min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Random String generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Random owner genrates a owner with given length
func RandomOwner() string {
	return RandomString(6)
}

// Random Balance generates balance with in a limit
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

// Random Currency from the currency list
func RandomCurrency() string {
	currencies := []string{"INR", "USD", "EUR"}
	return currencies[rand.Intn(len(currencies))]
}
