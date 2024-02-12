package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// AuthService maneja la lógica de autenticación
type AuthService struct{}

// HashPassword hashea la contraseña utilizando bcrypt
func (as *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
