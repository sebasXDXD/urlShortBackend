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
func (tr UserRepository) GetUserByUsername(username string) (*entities.Users, error) {
	query := "SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = $1"
	row := tr.DB.QueryRow(query, username)

	user := entities.Users{}
	updatedAtNull := sql.NullTime{}
	createdAtNull := sql.NullTime{}

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &createdAtNull, &updatedAtNull)
	if err == sql.ErrNoRows {
		// No se encontró ningún usuario con ese nombre de usuario
		return nil, nil
	} else if err != nil {
		// Otro error, devolverlo
		return nil, err
	}

	user.CreatedAt = createdAtNull.Time
	user.UpdatedAt = updatedAtNull.Time

	return &user, nil
}

// UserRepository
func (tr UserRepository) CreateUser(newUser entities.Users) (entities.Users, error) {
	// Define la consulta SQL para insertar un nuevo usuario
	query := "INSERT INTO users (first_name, last_name, username, email, password) VALUES ($1, $2, $3, $4, $5)"

	result, err := tr.DB.Exec(query, newUser.FirstName, newUser.LastName, newUser.Username, newUser.Email, newUser.Password)
	if err != nil {
		return entities.Users{}, err
	}

	// Obtén el ID del usuario creado
	userID, _ := result.LastInsertId()

	// Asigna el ID al usuario creado
	newUser.ID = int(userID)

	return newUser, nil
}
