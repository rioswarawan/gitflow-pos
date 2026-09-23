package model

import "database/sql"

// InitDB membuka (atau membuat) file SQLite dan menyiapkan skema.
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// SQLite praktis single-writer; batasi koneksi agar aman.
	db.SetMaxOpenConns(1)

	schema := `
CREATE TABLE IF NOT EXISTS products (
	id    INTEGER PRIMARY KEY AUTOINCREMENT,
	sku   TEXT NOT NULL UNIQUE,
	name  TEXT NOT NULL,
	price REAL NOT NULL CHECK (price >= 0),
	stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0)
);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// SeedProducts mengisi data awal jika tabel masih kosong.
func SeedProducts(db *sql.DB) error {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	_, err := db.Exec(`
INSERT INTO products (sku, name, price, stock) VALUES
	('SKU-001', 'Indomie Goreng',      3500.00, 100),
	('SKU-002', 'Teh Kotak 250ml',     5000.00,  80),
	('SKU-003', 'Kopi Sachet',         2000.00, 150),
	('SKU-004', 'Beras 5kg',          68000.00,  20),
	('SKU-005', 'Minyak Goreng 1L',   17500.00,  40);`)
	return err
}
