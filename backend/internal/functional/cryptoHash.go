package functional

import "golang.org/x/crypto/bcrypt"

type Hash string

func Encrypt(password string) (Hash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return Hash(hash), nil
}

func (h *Hash) Validate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*h), []byte(password))
	return err == nil
}
