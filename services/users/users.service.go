package services

import (
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

func (ts UserService) GetTasks() ([]entities.Users, error) {

	tasks, err := ts.UserRepository.GetTasks()
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
