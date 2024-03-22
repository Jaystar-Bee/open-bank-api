package utils

import (
	"math/rand"
	"net/mail"
	"strconv"
	"time"
)

func GenerateUniqueNumbers(n int) []int {
	// Seed the random number generator
	rand.NewSource(time.Now().UnixNano())
	uniqueNumbers := make([]int, 0)

	// Generate unique numbers
	for len(uniqueNumbers) < n {
		num := rand.Intn(9000000000) + 1000000000
		uniqueNumbers = append(uniqueNumbers, num)

	}

	return uniqueNumbers
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsConvertibleToNumber(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}
