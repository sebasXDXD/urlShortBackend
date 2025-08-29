package repositories

import (
	"database/sql"
	"urlShortenerBack/entities"
)

// Tipos auxiliares de retorno
type LinkWithClicks struct {
	Link   entities.Link
	Clicks int
}

type MonthCount struct {
	Month string
	Count int
}

type LinkRepository struct {
	DB *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return LinkRepository{DB: db}
}

// Devuelve todos los links (no borrados)
func (lr LinkRepository) GetLinks() ([]entities.Link, error) {
	query := "SELECT id, name, redirect_to, user_created_id, created_at, updated_at FROM links WHERE is_deleted = false"
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

// CreateLink: usa RETURNING para recuperar id y timestamps (Postgres)
func (lr LinkRepository) CreateLink(newLink entities.Link) (entities.Link, error) {
	query := `
		INSERT INTO links (name, redirect_to, user_created_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	var id int
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	err := lr.DB.QueryRow(query, newLink.Name, newLink.RedirectTo, newLink.UserCreatedID).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return entities.Link{}, err
	}
	newLink.ID = id
	if createdAt.Valid {
		newLink.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		newLink.UpdatedAt = updatedAt.Time
	}
	return newLink, nil
}

// Obtener enlace por short name (string)
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

// UpdateLink: actualiza nombre y redirect
func (lr LinkRepository) UpdateLink(linkID int, name string, redirectTo string) error {
	query := "UPDATE links SET name = $1, redirect_to = $2, updated_at = NOW() WHERE id = $3 AND is_deleted = false"
	_, err := lr.DB.Exec(query, name, redirectTo, linkID)
	return err
}

// Marca un click (insert en clicks)
func (lr LinkRepository) RegisterClick(linkID int, ip, userAgent string) error {
	query := `
		INSERT INTO clicks (link_id, ip_address, user_agent) 
		VALUES ($1, $2, $3)
	`
	_, err := lr.DB.Exec(query, linkID, ip, userAgent)
	return err
}

// Estadísticas generales por usuario (ya implementado con pequeñas mejoras)
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

/*
  Consultas para el panel
*/

// 1) Clicks por mes (últimos `months` meses) para un usuario
// Retorna un slice con MonthCount{Month, Count} ordenado asc por mes (más antiguo -> más reciente)
func (lr LinkRepository) GetClicksByMonth(userID, months int) ([]MonthCount, error) {
	/*
	   Usamos generate_series para generar los meses y hacer LEFT JOIN con la agregación real.
	*/
	query := `
	SELECT to_char(m.month, 'Mon YYYY') AS month_label,
	       COALESCE(c.clicks, 0) AS clicks
	FROM generate_series(
	       date_trunc('month', CURRENT_DATE) - interval '1 month' * ($2 - 1),
	       date_trunc('month', CURRENT_DATE),
	       interval '1 month'
	     ) AS m(month)
	LEFT JOIN (
	    SELECT date_trunc('month', c.created_at) AS mth, COUNT(*) AS clicks
	    FROM clicks c
	    JOIN links l ON c.link_id = l.id
	    WHERE l.user_created_id = $1 AND l.is_deleted = false
	    GROUP BY mth
	) c ON c.mth = m.month
	ORDER BY m.month;
	`

	rows, err := lr.DB.Query(query, userID, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []MonthCount{}
	for rows.Next() {
		var m MonthCount
		if err := rows.Scan(&m.Month, &m.Count); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// 2) Top N enlaces por clicks para un usuario
func (lr LinkRepository) GetTopLinksByUser(userID, limit int) ([]LinkWithClicks, error) {
	query := `
	SELECT l.id, l.name, l.redirect_to, l.user_created_id, l.created_at, l.updated_at, COUNT(c.*) AS clicks
	FROM links l
	LEFT JOIN clicks c ON c.link_id = l.id
	WHERE l.user_created_id = $1 AND l.is_deleted = false
	GROUP BY l.id
	ORDER BY clicks DESC
	LIMIT $2
	`
	rows, err := lr.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []LinkWithClicks{}
	for rows.Next() {
		var lw LinkWithClicks
		createdAtNull := sql.NullTime{}
		updatedAtNull := sql.NullTime{}
		if err := rows.Scan(&lw.Link.ID, &lw.Link.Name, &lw.Link.RedirectTo, &lw.Link.UserCreatedID, &createdAtNull, &updatedAtNull, &lw.Clicks); err != nil {
			return nil, err
		}
		if createdAtNull.Valid {
			lw.Link.CreatedAt = createdAtNull.Time
		}
		if updatedAtNull.Valid {
			lw.Link.UpdatedAt = updatedAtNull.Time
		}
		res = append(res, lw)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// 3) Links creados por mes (últimos `months` meses) para un usuario
func (lr LinkRepository) GetLinksCreatedByMonth(userID, months int) ([]MonthCount, error) {
	query := `
	SELECT to_char(m.month, 'Mon YYYY') AS month_label,
	       COALESCE(lc.count_links, 0) AS count_links
	FROM generate_series(
	       date_trunc('month', CURRENT_DATE) - interval '1 month' * ($2 - 1),
	       date_trunc('month', CURRENT_DATE),
	       interval '1 month'
	     ) AS m(month)
	LEFT JOIN (
	    SELECT date_trunc('month', created_at) AS mth, COUNT(*) AS count_links
	    FROM links
	    WHERE user_created_id = $1 AND is_deleted = false
	    GROUP BY mth
	) lc ON lc.mth = m.month
	ORDER BY m.month;
	`

	rows, err := lr.DB.Query(query, userID, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []MonthCount{}
	for rows.Next() {
		var mc MonthCount
		if err := rows.Scan(&mc.Month, &mc.Count); err != nil {
			return nil, err
		}
		res = append(res, mc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// 4) Links recientes por usuario (limit)
func (lr LinkRepository) GetRecentLinksByUser(userID, limit int) ([]entities.Link, error) {
	query := `
		SELECT id, name, redirect_to, user_created_id, created_at, updated_at
		FROM links
		WHERE user_created_id = $1 AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := lr.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []entities.Link{}
	for rows.Next() {
		link := entities.Link{}
		createdAtNull := sql.NullTime{}
		updatedAtNull := sql.NullTime{}
		if err := rows.Scan(&link.ID, &link.Name, &link.RedirectTo, &link.UserCreatedID, &createdAtNull, &updatedAtNull); err != nil {
			return nil, err
		}
		if createdAtNull.Valid {
			link.CreatedAt = createdAtNull.Time
		}
		if updatedAtNull.Valid {
			link.UpdatedAt = updatedAtNull.Time
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
