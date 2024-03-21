package utils

import (
	"math/rand"
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
