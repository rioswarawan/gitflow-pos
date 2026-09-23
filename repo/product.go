package repo

import (
	"database/sql"
	"errors"

	"pos-module/model"
)

// ProductRepo mengakses data produk di SQLite.
type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

var ErrNotFound = errors.New("product not found")

// FindBySKU mengambil satu produk berdasarkan SKU.
func (r *ProductRepo) FindBySKU(sku string) (model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(
		`SELECT id, sku, name, price, stock FROM products WHERE sku = ?`, sku,
	).Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Product{}, ErrNotFound
	}
	return p, err
}

// List mengambil semua produk.
func (r *ProductRepo) List() ([]model.Product, error) {
	rows, err := r.db.Query(`SELECT id, sku, name, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
