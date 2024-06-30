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
	UserId    int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

type UserRepo interface {
	CreateUser(login string, password string) (User, error)
	GetUserByLogin(login string) (User, error)
	GetUserById(id int64) (User, error)
}

type SessionRepo interface {
	CreateSession(userId int64, token string) (Session, error)
	GetSession(token string) (Session, error)
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

func (r AuthService) ValidateSession(token string) (bool, error) {
	session, err := r.Sr.GetSession(token)
	if err != nil {
		return false, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return false, nil
	}

	session.ExpiresAt = time.Now().Add(time.Hour * 24)

	_, err = r.Sr.UpdateSession(session)
	if err != nil {
		return false, err
	}

	return true, nil
}

type UserInfoDto struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
}

func (r AuthService) GetInfoAboutUserBySessionToken(token string) (UserInfoDto, error) {
	session, err := r.Sr.GetSession(token)
	if err != nil {
		return UserInfoDto{}, err
	}

	user, err := r.Ur.GetUserById(session.UserId)
	if err != nil {
		return UserInfoDto{}, err
	}

	return UserInfoDto{Id: user.Id, Login: user.Login}, err
}
