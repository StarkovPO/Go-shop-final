package service

import (
	"crypto/sha1"
	"fmt"
)

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func comparePasswordHash(hashedPassword, password string) bool {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) == hashedPassword
}
