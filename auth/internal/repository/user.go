package repository

import (
	"auth-goker/internal/service"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	DB *sqlx.DB
}

type UserEntity struct {
	Id       int64  `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (r UserRepo) CreateUser(login string, password string) (service.User, error) {
	var user UserEntity

	err := r.DB.QueryRowx("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING *", login, password).StructScan(&user)
	if err != nil {
		return service.User{}, err
	}

	return service.User{Id: user.Id, Login: user.Login, Password: user.Password}, nil
}
