package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDb() (*sql.DB, error) {
	// Cambiar la cadena de conexión y el controlador para PostgreSQL
	dataSourceName := "user=sebas_cruds password=root dbname=url_short_db sslmode=disable"
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err // Devolvemos el error en lugar de usar log.Fatal
	}

	err = db.Ping()
	if err != nil {
		db.Close() // Cerramos la conexión antes de devolver el error
		return nil, err
	}

	fmt.Println("Conexión a la base de datos PostgreSQL establecida exitosamente.")
	return db, nil
}
