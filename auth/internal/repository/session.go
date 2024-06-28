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
	var session SessionEntity
	err := r.DB.QueryRowx(
		"INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING *",
		token, userId, time.Now(), time.Now().Add(time.Hour*24)).
		StructScan(&session)

	return r.sessionEntityToSession(session), err
}

func (r SessionRepo) GetSessionsByUser(userId int64) ([]service.Session, error) {
	var sessions []SessionEntity
	err := r.DB.Select(&sessions, "SELECT * FROM sessions WHERE user_id = $1", userId)

	return func() []service.Session {
		var result []service.Session
		for _, session := range sessions {
			result = append(result, r.sessionEntityToSession(session))
		}
		return result
	}(), err
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
	return service.Session{Token: session.Token, CreatedAt: session.CreatedAt, ExpiresAt: session.ExpiresAt}
}
