package auth

import (
	"github.com/dgrijalva/jwt-go"
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

// ComparePasswords compara la contraseña proporcionada con la contraseña hasheada
func (as *AuthService) ComparePasswords(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func (as *AuthService) AssignToken(userID int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
	})

	tokenString, err := token.SignedString([]byte(SecretWord))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
