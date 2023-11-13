package password

import "golang.org/x/crypto/bcrypt"

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (us *PasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (us *PasswordService) ComparePasswords(hashed string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err != nil {
		return false
	}
	return true
}
