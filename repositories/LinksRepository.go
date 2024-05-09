package repositories

import (
	"database/sql"
	"urlShortenerBack/entities"
)

type LinkRepository struct {
	DB *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return LinkRepository{DB: db}
}

func (lr LinkRepository) GetLinks() ([]entities.Link, error) {
	query := "SELECT id, name, redirect_to, user_created_id, created_at, updated_at FROM links"
	rows, err := lr.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []entities.Link{}

	for rows.Next() {
		link := entities.Link{}
		updatedAtNull := sql.NullTime{}
		createdAtNull := sql.NullTime{}
		if err := rows.Scan(&link.ID, &link.Name, &link.RedirectTo, &link.UserCreatedID, &createdAtNull, &updatedAtNull); err != nil {
			return nil, err
		}
		link.CreatedAt = createdAtNull.Time
		link.UpdatedAt = updatedAtNull.Time
		links = append(links, link)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (lr LinkRepository) CreateLink(newLink entities.Link) (entities.Link, error) {
	// Define la consulta SQL para insertar un nuevo enlace
	query := "INSERT INTO links (name, redirect_to, user_created_id) VALUES ($1, $2, $3)"

	result, err := lr.DB.Exec(query, newLink.Name, newLink.RedirectTo, newLink.UserCreatedID)
	if err != nil {
		return entities.Link{}, err
	}

	// Obtén el ID del enlace creado
	linkID, _ := result.LastInsertId()

	// Asigna el ID al enlace creado
	newLink.ID = int(linkID)

	return newLink, nil
}

// Otras funciones de repositorio específicas para la entidad Link, como GetUserByUsername, etc.
