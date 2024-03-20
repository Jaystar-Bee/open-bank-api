package utils

import "golang.org/x/crypto/bcrypt"

func HashText(text string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(text), 13)
	return string(hashByte), err
}

func CompareHash(hash, text string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
}
