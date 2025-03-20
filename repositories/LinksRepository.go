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
	query := "INSERT INTO links (name, redirect_to, user_created_id) VALUES ($1, $2, $3)"

	result, err := lr.DB.Exec(query, newLink.Name, newLink.RedirectTo, newLink.UserCreatedID)
	if err != nil {
		return entities.Link{}, err
	}

	linkID, _ := result.LastInsertId()
	newLink.ID = int(linkID)

	return newLink, nil
}

// Método para obtener un enlace por su nombre
func (lr LinkRepository) GetLinkByString(name string) (entities.Link, error) {
	query := "SELECT id, name, redirect_to, user_created_id, created_at, updated_at FROM links WHERE name = $1 AND is_deleted = false"
	row := lr.DB.QueryRow(query, name)

	link := entities.Link{}
	updatedAtNull := sql.NullTime{}
	createdAtNull := sql.NullTime{}
	err := row.Scan(&link.ID, &link.Name, &link.RedirectTo, &link.UserCreatedID, &createdAtNull, &updatedAtNull)
	if err != nil {
		if err == sql.ErrNoRows {
			return link, nil
		}
		return link, err
	}
	link.CreatedAt = createdAtNull.Time
	link.UpdatedAt = updatedAtNull.Time
	return link, nil
}
func (lr LinkRepository) GetLinkByID(id int) (entities.Link, error) {
	query := "SELECT id, name, redirect_to, user_created_id, created_at, updated_at FROM links WHERE id = $1 AND is_deleted = false"
	row := lr.DB.QueryRow(query, id)

	link := entities.Link{}
	updatedAtNull := sql.NullTime{}
	createdAtNull := sql.NullTime{}
	err := row.Scan(&link.ID, &link.Name, &link.RedirectTo, &link.UserCreatedID, &createdAtNull, &updatedAtNull)
	if err != nil {
		if err == sql.ErrNoRows {
			return link, nil
		}
		return link, err
	}
	link.CreatedAt = createdAtNull.Time
	link.UpdatedAt = updatedAtNull.Time
	return link, nil
}

func (lr LinkRepository) UpdateLink(linkID int, name string, redirectTo string) error {
	query := "UPDATE links SET name = $1, redirect_to = $2, updated_at = NOW() WHERE id = $3 AND is_deleted = false"
	_, err := lr.DB.Exec(query, name, redirectTo, linkID)
	return err
}
