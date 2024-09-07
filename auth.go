package ledger

import (
	"context"
	"errors"
	"github.com/timhugh/ctxlogger"
)

var ErrInvalidPassword = errors.New("invalid password")

const pepper = "pepper" // TODO: real pepper in config

type User struct {
	UserUUID     string `json:"user_uuid" db:"user_uuid"`
	Login        string `json:"user_login" db:"user_login"`
	PasswordHash string `json:"user_password_hash" db:"user_password_hash"`
	PasswordSalt string `json:"user_password_salt" db:"user_password_salt"`
}

type UserGetter interface {
	GetUser(ctx context.Context, uuid string) (*User, error)
}

type Session struct {
	SessionUUID string `json:"session_uuid" db:"session_uuid"`
	UserUUID    string `json:"user_uuid" db:"session_user_uuid"`

	User *User `json:"user"`
}

type SessionGetter interface {
	GetSession(ctx context.Context, uuid string) (*Session, error)
}

type HashFunction func(password string, salt string, pepper string) string

func AuthenticateUser(ctx context.Context, repo UserGetter, hashFunc HashFunction, login string, password string) (*User, error) {
	user, err := repo.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// Use empty user to avoid timing attacks
			user = &User{}
		}
	}

	hash := hashFunc(password, user.PasswordSalt, pepper)
	if err != nil {
		ctxlogger.Error(ctx, "failed to create password hash: %s", err.Error())
		return nil, ErrInvalidPassword
	}
	if hash != user.PasswordHash {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
