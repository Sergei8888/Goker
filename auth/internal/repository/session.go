package repository

import (
	"auth-goker/internal/service"
	"github.com/jmoiron/sqlx"
	"time"
)

type SessionEntity struct {
	Token     string    `db:"token"`
	UserId    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type SessionRepo struct {
	DB *sqlx.DB
}

func (r SessionRepo) CreateSession(userId int64, token string) (service.Session, error) {
	var sessionEntity SessionEntity

	err := r.DB.QueryRowx(
		"INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING *",
		token, userId, time.Now(), time.Now().Add(time.Hour*24)).
		StructScan(&sessionEntity)

	return r.sessionEntityToSession(sessionEntity), err
}

func (r SessionRepo) GetSession(token string) (service.Session, error) {
	var sessionEntity SessionEntity

	err := r.DB.QueryRowx("SELECT * FROM sessions WHERE token = $1", token).StructScan(&sessionEntity)
	if err != nil {
		return service.Session{}, err
	}

	return r.sessionEntityToSession(sessionEntity), nil
}

func (r SessionRepo) UpdateSession(session service.Session) (service.Session, error) {
	var sessionEntity SessionEntity

	err := r.DB.QueryRowx("UPDATE sessions SET expires_at = $1, created_at = $2 WHERE token = $3 RETURNING *", session.ExpiresAt, session.CreatedAt, session.Token).StructScan(&sessionEntity)
	if err != nil {
		return service.Session{}, err
	}

	return r.sessionEntityToSession(sessionEntity), nil
}

func (r SessionRepo) sessionEntityToSession(session SessionEntity) service.Session {
	return service.Session{Token: session.Token, UserId: session.UserId, CreatedAt: session.CreatedAt, ExpiresAt: session.ExpiresAt}
}
