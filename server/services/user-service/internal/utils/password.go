package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}

func ComparePassword(pw string, hpw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hpw), []byte(pw))
	return err == nil
}
