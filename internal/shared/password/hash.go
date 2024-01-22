package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err != nil {
		return false
	}
	return true
}
