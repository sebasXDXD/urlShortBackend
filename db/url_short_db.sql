-- Crear la función para actualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear la tabla users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear el trigger para llamar a la función antes de una actualización en users
CREATE TRIGGER users_update_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- Crear la tabla links
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    redirect_to TEXT,
    user_created_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP
);

-- Crear un trigger para actualizar updated_at cuando se realicen cambios en la fila
CREATE OR REPLACE FUNCTION update_links_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER links_updated_at_trigger
BEFORE UPDATE ON links
FOR EACH ROW
EXECUTE FUNCTION update_links_updated_at();

-- Crear un trigger para establecer deleted_at cuando is_deleted se establece como true
CREATE OR REPLACE FUNCTION update_links_deleted_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_deleted = TRUE THEN
        NEW.deleted_at = CURRENT_TIMESTAMP;
    ELSE
        NEW.deleted_at = NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER links_deleted_at_trigger
BEFORE UPDATE ON links
FOR EACH ROW
EXECUTE FUNCTION update_links_deleted_at();
