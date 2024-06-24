package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"urlShortenerBack/utils"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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
		"userID":   userID,
		"username": username,
	})

	tokenString, err := token.SignedString([]byte(SecretWord))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Middleware de autenticación
type contextKey string

const userIDKey = contextKey("userID")

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

		// Depuración: imprimir el token
		fmt.Printf("Token recibido: %s\n", jwtToken)

		// Verificar la validez del token
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretWord), nil
		})
		if err != nil {
			fmt.Printf("Error al parsear el token: %v\n", err)
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			fmt.Println("Token no es válido")
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Extraer los datos del token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("No se pudieron obtener los datos del token")
			http.Error(w, "No se pudieron obtener los datos del token", http.StatusUnauthorized)
			return
		}

		// Depuración: imprimir los claims del token
		fmt.Printf("Claims del token: %+v\n", claims)

		// Verificar si el claim "id" está presente
		userID, ok := claims["userID"].(float64) // JWT claims are typically float64
		if !ok {
			fmt.Println("Claim 'userID' no está presente en el token")
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Depuración: imprimir el userID
		fmt.Printf("ID de usuario extraído del token: %v\n", userID)

		// Agregar el userID al contexto de la solicitud
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		r = r.WithContext(ctx)

		// Continuar con la solicitud
		next.ServeHTTP(w, r)
	})
}

// metodo para validadr el id token de google
func (as AuthService) ValidateGoogleID(ctx context.Context, googleIDToken string) (bool, error) {
	// Tu CLIENT_ID de Google
	clientID := "MyGoogleClientID"

	// Verificar el ID token
	payload, err := idtoken.Validate(ctx, googleIDToken, clientID)
	if err != nil {
		return false, err
	}
	if payload == nil {
		return false, errors.New("token de Google no válido")
	}

	return true, nil
}
