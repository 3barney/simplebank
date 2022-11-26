package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Called automatically when package is first used
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt Generates random integer between min and max
func RandomInt(min, max int64) int64 {
	return max - rand.Int63n(max-min+1)
}

// RandomString Generates a string of length n
func RandomString(n int) string {
	var stringBuilder strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		stringBuilder.WriteByte(c)
	}

	return stringBuilder.String()
}

// RandomName Generates a random name using
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney Generates a random money amount
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency Generates a random currency
func RandomCurrency() string {
	currencies := []string{"EUR", "KSH", "USD"}
	n := len(currencies)
	return currencies[rand.Intn(n)] // Get random data between 0 and n-1

}
