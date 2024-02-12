package services

import (
	"urlShortenerBack/entities"
	"urlShortenerBack/repositories"
)

type UserService struct {
	UserRepository repositories.UserRepository
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

func (ts UserService) CreateTask(newTask entities.Users) (entities.Users, error) {
	// Llama al método CreateTask del repositorio y pasa la nueva tarea
	createdUser, err := ts.UserRepository.CreateTask(newTask)
	if err != nil {
		return entities.Users{}, err
	}

	return createdUser, nil
}
