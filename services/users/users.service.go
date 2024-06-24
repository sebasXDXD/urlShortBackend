package services

import (
	"context"
	"errors"
	"urlShortenerBack/auth"
	"urlShortenerBack/entities"
	"urlShortenerBack/repositories"
)

type UserService struct {
	UserRepository repositories.UserRepository
	AuthService    auth.AuthService
}

func NewUserService(userRepo repositories.UserRepository) UserService {

	return UserService{UserRepository: userRepo}
}

func (ts UserService) GetUsers() ([]entities.Users, error) {

	tasks, err := ts.UserRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (us UserService) CreateUser(newUser entities.Users) (entities.Users, error) {
	// Hashear la contraseña del nuevo usuario
	hashedPassword, err := us.AuthService.HashPassword(newUser.Password)
	if err != nil {
		return entities.Users{}, err
	}

	// Actualizar la contraseña con la versión hasheada
	newUser.Password = hashedPassword

	// Llamar al método CreateUser del repositorio y pasar la nueva tarea
	createdUser, err := us.UserRepository.CreateUser(newUser)
	if err != nil {
		return entities.Users{}, err
	}

	return createdUser, nil
}

func (us UserService) CreateGoogleUser(ctx context.Context, newUser entities.Users) (entities.Users, error) {
	// // Validar el google_id token del nuevo usuario
	// isValid, err := us.AuthService.ValidateGoogleID(ctx, googleIDToken)
	// if err != nil {
	// 	return entities.Users{}, err
	// }

	// if !isValid {
	// 	return entities.Users{}, errors.New("google_id token no es válido")
	// }

	// // Aquí no es necesario hashear la contraseña porque estamos creando un usuario de Google
	// // Si quieres puedes establecer la contraseña como una cadena vacía
	// newUser.Password = ""

	// // Llamar al método CreateUser del repositorio y pasar la nueva tarea
	createdUser, err := us.UserRepository.CreateUser(newUser)
	if err != nil {
		return entities.Users{}, err
	}

	return createdUser, nil
}
func (us UserService) Login(inputUser entities.Users) (*entities.Users, error) {
	// Buscar el usuario por su nombre de usuario en el repositorio
	existingUser, err := us.UserRepository.GetUserByUsername(inputUser.Username)
	if err != nil {
		return nil, err
	}

	// Verificar si el usuario existe
	if existingUser == nil {
		// El usuario no existe, puedes devolver un error o un mensaje adecuado
		return nil, errors.New("El usuario no existe")
	}

	// Hashear la contraseña proporcionada para compararla con la contraseña almacenada
	err = us.AuthService.ComparePasswords(existingUser.Password, inputUser.Password)
	if err != nil {
		// Las contraseñas no coinciden
		return nil, errors.New("Contraseña incorrecta")
	}

	// Continuar con el flujo de inicio de sesión si todo está correcto
	// Puedes devolver el usuario autenticado o la información necesaria
	return existingUser, nil
}
func (us UserService) LoginGoogle(inputUser entities.Users) (*entities.Users, error) {
	// Buscar el usuario por su nombre de usuario en el repositorio
	existingUser, err := us.UserRepository.GetUserByUsername(inputUser.Username)
	if err != nil {
		return nil, err
	}

	// Verificar si el usuario existe
	if existingUser == nil {
		// El usuario no existe, puedes devolver un error o un mensaje adecuado
		return nil, errors.New("El usuario no existe")
	}
	return existingUser, nil
}
