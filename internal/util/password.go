package util

import "golang.org/x/crypto/bcrypt"

const HASHCONST = 12

func HashPassword(password []byte) ([]byte, error) {
	hashpassword, err := bcrypt.GenerateFromPassword(password, HASHCONST)
	if err != nil {
		return nil, err
	}
	return hashpassword, nil
}

func ComparePassword(password, hashpassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashpassword, password)
	return err
}
