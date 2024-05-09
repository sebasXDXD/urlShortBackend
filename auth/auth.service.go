package auth

import (
	"context"
	"net/http"
	"strings"

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

// Middleware de autenticación
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token de la solicitud
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Token de autorización faltante", http.StatusUnauthorized)
			return
		}

		// Verificar el formato del token
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}
		jwtToken := parts[1]

		// Verificar la validez del token
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			// Aquí debes usar la misma clave secreta utilizada para firmar el token
			return []byte(SecretWord), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Extraer los datos del token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "No se pudieron obtener los datos del token", http.StatusUnauthorized)
			return
		}
		type contextKey string
		var userIDKey = contextKey("userID")
		ctx := context.WithValue(r.Context(), userIDKey, claims["id"])
		r = r.WithContext(ctx)

		// Si el token es válido, continuamos con la solicitud
		next.ServeHTTP(w, r)
	})
}
