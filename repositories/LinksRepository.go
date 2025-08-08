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
func (lr LinkRepository) GetStatsByUserID(userID int) (entities.LinkStats, error) {
	stats := entities.LinkStats{}

	// Total de links creados por el usuario
	queryTotalLinks := `
		SELECT COUNT(*) FROM links WHERE user_created_id = $1 AND is_deleted = false
	`
	err := lr.DB.QueryRow(queryTotalLinks, userID).Scan(&stats.TotalCreated)
	if err != nil {
		return stats, err
	}

	// Total de clics en todos sus links
	queryTotalClicks := `
		SELECT COUNT(*) 
		FROM clicks 
		WHERE link_id IN (SELECT id FROM links WHERE user_created_id = $1 AND is_deleted = false)
	`
	err = lr.DB.QueryRow(queryTotalClicks, userID).Scan(&stats.TotalClicks)
	if err != nil {
		return stats, err
	}

	// URL más popular
	queryPopular := `
		SELECT l.name, COUNT(*) as clicks 
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_created_id = $1 AND l.is_deleted = false
		GROUP BY l.name
		ORDER BY clicks DESC
		LIMIT 1
	`
	err = lr.DB.QueryRow(queryPopular, userID).Scan(&stats.MostClickedURL, &stats.MostClickedCount)
	if err == sql.ErrNoRows {
		stats.MostClickedURL = ""
		stats.MostClickedCount = 0
	} else if err != nil {
		return stats, err
	}

	// Último acceso (último click)
	queryLastAccess := `
		SELECT MAX(c.created_at)
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.user_created_id = $1 AND l.is_deleted = false
	`
	var lastAccess sql.NullTime
	err = lr.DB.QueryRow(queryLastAccess, userID).Scan(&lastAccess)
	if err != nil {
		return stats, err
	}
	if lastAccess.Valid {
		stats.LastAccess = lastAccess.Time
	}

	return stats, nil
}
