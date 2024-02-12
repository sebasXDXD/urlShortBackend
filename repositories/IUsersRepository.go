package repositories

import (
	"database/sql"
	"urlShortenerBack/entities"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{DB: db}
}

func (tr UserRepository) GetTasks() ([]entities.Users, error) {
	query := "SELECT id, username, password,email, created_at, updated_at FROM users"
	rows, err := tr.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entities.Users{}

	for rows.Next() {
		user := entities.Users{}
		updatedAtNull := sql.NullTime{}
		createdAtNull := sql.NullTime{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &createdAtNull, &updatedAtNull); err != nil {
			return nil, err
		}
		user.CreatedAt = createdAtNull.Time
		user.UpdatedAt = updatedAtNull.Time
		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (tr UserRepository) CreateTask(newUser entities.Users) (entities.Users, error) {
	// Define la consulta SQL para insertar una nueva tarea con columnas opcionales
	query := "INSERT INTO users (Title, Content) VALUES (?, ?)"
	_, err := tr.DB.Exec(query, newUser.Username, newUser.Password)
	if err != nil {
		return entities.Users{}, err
	}

	return newUser, nil
}
