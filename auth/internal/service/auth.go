package service

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Id    int64
	Login string
	// hash of password
	Password string
}

type Session struct {
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type UserRepo interface {
	CreateUser(login string, password string) (User, error)
	GetUserByLogin(login string) (User, error)
}

type SessionRepo interface {
	CreateSession(userId int64, token string) (Session, error)
	GetSessionsByUser(userId int64) ([]Session, error)
	UpdateSession(session Session) (Session, error)
}

type AuthService struct {
	Ur UserRepo
	Sr SessionRepo
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

func (r AuthService) Signin(login string, password string) (Session, error) {
	user, err := r.Ur.GetUserByLogin(login)
	if err != nil {
		return Session{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return Session{}, err
	}

	session, err := r.Sr.CreateSession(user.Id, r.generateToken(login))
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (r AuthService) generateToken(login string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(login), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))

	return base64.StdEncoding.EncodeToString(hash)
}
