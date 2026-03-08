package domain

import (
	"CHIPMUNK-T0T/ito_web_app/internal/functional"
	"fmt"
)

type User struct {
	id       uint
	username string
	password functional.Hash
}

func NewUser(username string, password string) (User, error) {
	hash, err := functional.Encrypt(password)
	if err != nil {
		return User{}, err
	}
	return User{
		username: username,
		password: hash,
	}, nil
}

func NewUserWithID(id uint, username string, passwordHash functional.Hash) User {
	return User{
		id:       id,
		username: username,
		password: passwordHash,
	}
}

// Getters
func (u User) ID() uint {
	return u.id
}

func (u *User) SetID(id uint) {
	u.id = id
}

func (u User) Username() string {
	return u.username
}

func (u User) Password() functional.Hash {
	return u.password
}

// Validation
func (u *User) Validate(username string, password string) error {
	if u.username != username {
		return fmt.Errorf("ユーザー名が一致しません")
	}
	if !u.password.Validate(password) {
		return fmt.Errorf("パスワードが一致しません")
	}
	return nil
}
