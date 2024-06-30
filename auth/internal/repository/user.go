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

	return r.userEntityToUser(user), nil
}

func (r UserRepo) GetUserByLogin(login string) (service.User, error) {
	var user UserEntity

	err := r.DB.QueryRowx("SELECT * FROM users WHERE login = $1", login).StructScan(&user)
	if err != nil {
		return service.User{}, err
	}

	return r.userEntityToUser(user), nil

}

func (r UserRepo) GetUserById(id int64) (service.User, error) {
	var user UserEntity

	err := r.DB.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(&user)
	if err != nil {
		return service.User{}, err
	}

	return r.userEntityToUser(user), nil
}

func (r UserRepo) userEntityToUser(entity UserEntity) service.User {
	return service.User{Id: entity.Id, Login: entity.Login, Password: entity.Password}
}
