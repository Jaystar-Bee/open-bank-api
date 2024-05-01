package utils

import (
	"math/rand"
	"net/mail"
	"net/url"
	"strconv"
	"time"

	"github.com/Jaystar-Bee/open-bank-api/db"
	"github.com/Jaystar-Bee/open-bank-api/emails"
)

func GenerateUniqueNumbers(n, randVal int) []int {
	// Seed the random number generator
	rand.NewSource(time.Now().UnixNano())
	// uniqueNumbers := make([]int, 0)

	// // Generate unique numbers
	// for len(uniqueNumbers) < n {
	// 	num := rand.Intn(9000000000) + 1000000000
	// 	uniqueNumbers = append(uniqueNumbers, num)

	// }
	var numbers []int
	for len(numbers) < n {
		num := rand.Intn(randVal)
		numbers = append(numbers, num)
	}
	return numbers

}

func JoinIntSlice(slice []int) string {
	var str string
	for _, num := range slice {
		str += strconv.Itoa(num) + ""
	}
	return str
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsConvertibleToNumber(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}

func SetOTP(key, name string) {

	arrayOfNumbers := GenerateUniqueNumbers(1, 99999)
	otp := JoinIntSlice(arrayOfNumbers)

	db.RDB.Set(db.Ctx, key, otp, time.Minute+10)

	emails.SendOTPToMail(key, name, otp)
}

func ConfirmURL(link string) (*url.URL, error) {
	return url.ParseRequestURI(link)
}
