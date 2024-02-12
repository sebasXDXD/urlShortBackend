package services

import (
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
