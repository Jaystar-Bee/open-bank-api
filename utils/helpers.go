package utils

import (
	"math/rand"
	"net/mail"
	"net/url"
	"strconv"
	"time"
)

func GenerateUniqueNumbers(low, high int) int {
	rand.NewSource(time.Now().UnixNano())

	return low + rand.Intn(high-low)

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

func ConfirmURL(link string) (*url.URL, error) {
	return url.ParseRequestURI(link)
}
