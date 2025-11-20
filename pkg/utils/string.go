package utils

import (
	"math/rand/v2"
	"regexp"
)

func GenerateRandomString(num int) string {
	const charset = "abcdef0123456789"
	randomBytes := make([]byte, num)
	for i := 0; i < num; i++ {
		randomBytes[i] = charset[rand.IntN(16)]
	}
	randomString := string(randomBytes)
	return randomString
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}
