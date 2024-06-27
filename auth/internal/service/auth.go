package service

import (
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(login string, password string) (User, error)
}

type User struct {
	Id    int64
	Login string
	// hash of password
	Password string
}

type AuthService struct {
	Ur UserRepo
}

func (r AuthService) Signup(login string, password string) (User, error) {
	hash, err := r.hashPassword(password)
	if err != nil {
		return User{}, err
	}

	return r.Ur.CreateUser(login, hash)
}

func (r AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
