package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDb() (*sql.DB, error) {
	var dataSourceName string

	// Si existe DATABASE_URL, la usamos (Render lo da así)
	if os.Getenv("DATABASE_URL") != "" {
		dataSourceName = os.Getenv("DATABASE_URL")
	} else {
		// Caso local: usamos las variables separadas del .env
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")

		dataSourceName = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode,
		)
	}

	// Conectar
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Probar conexión
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("✅ Conexión a la base de datos PostgreSQL establecida exitosamente.")
	return db, nil
}
