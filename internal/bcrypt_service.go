package internal

import (
	"golang.org/x/crypto/bcrypt"
)

type IBCryptService interface {
	HashPassword(password string) (string, error)
	PasswordsMatch(hash, password string) bool
}

type BCryptService struct{}

func NewBCryptService() *BCryptService {
	return &BCryptService{}
}

func (bs *BCryptService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(bytes), err
}

func (bs *BCryptService) PasswordsMatch(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
