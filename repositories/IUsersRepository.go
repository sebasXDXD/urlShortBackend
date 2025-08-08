package repositories

import (
	"database/sql"
	"urlShortenerBack/entities"
)

// UserRepository
type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{DB: db}
}

func (tr UserRepository) GetUsers() ([]entities.Users, error) {
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

func (tr UserRepository) CreateUser(newUser entities.Users) (entities.Users, error) {
	// Define la consulta SQL para insertar un nuevo usuario
	query := "INSERT INTO users (first_name, last_name, username, email, google_id, password) VALUES ($1, $2, $3, $4, $5, $6)"

	result, err := tr.DB.Exec(query, newUser.FirstName, newUser.LastName, newUser.Username, newUser.Email, newUser.GoogleID, newUser.Password)
	if err != nil {
		return entities.Users{}, err
	}

	// Obtén el ID del usuario creado
	userID, _ := result.LastInsertId()

	// Asigna el ID al usuario creado
	newUser.ID = int(userID)

	return newUser, nil
}

func (tr UserRepository) GetUserProfileByID(userID int) (*entities.Users, error) {
	query := `
		SELECT id, first_name, last_name, username, email, company, country, phone, plan, created_at, updated_at, url_limit
		FROM users 
		WHERE id = $1
	`

	row := tr.DB.QueryRow(query, userID)

	var user entities.Users
	var createdAtNull, updatedAtNull sql.NullTime
	var company, country, phone, plan sql.NullString
	var urlLimit sql.NullInt64

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&company,
		&country,
		&phone,
		&plan,
		&createdAtNull,
		&updatedAtNull,
		&urlLimit,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Asignar campos nulos si existen
	// user.Company = company.String
	// user.Country = country.String
	// user.Phone = phone.String
	// user.Plan = plan.String
	// user.URLLimit = int(urlLimit.Int64)
	// user.CreatedAt = createdAtNull.Time
	// user.UpdatedAt = updatedAtNull.Time

	return &user, nil
}
func (tr UserRepository) GetByID(userID int) (*entities.Users, error) {
	query := `
		SELECT id, first_name, last_name, username, email, password, google_id, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	row := tr.DB.QueryRow(query, userID)

	var user entities.Users
	var createdAtNull, updatedAtNull sql.NullTime
	var password, googleID sql.NullString

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&password,
		&googleID,
		&createdAtNull,
		&updatedAtNull,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Asignar campos nulos con verificación
	user.Password = password.String
	user.GoogleID = googleID.String

	if createdAtNull.Valid {
		user.CreatedAt = createdAtNull.Time
	}
	if updatedAtNull.Valid {
		user.UpdatedAt = updatedAtNull.Time
	}

	return &user, nil
}
