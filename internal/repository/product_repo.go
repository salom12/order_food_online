package repository

import (
	"database/sql"
	"order_food_online/internal/cache"
	"order_food_online/internal/models"
	"time"
)

type ProductRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id int) (*models.Product, error)
}

type ProductRepo struct {
	db    *sql.DB
	cache cache.ProductCache
}

func NewProductRepository(db *sql.DB, cache cache.ProductCache) ProductRepository {
	return &ProductRepo{db: db, cache: cache}
}

func (r *ProductRepo) GetAllProducts() ([]models.Product, error) {
	// Try Redis cache first
	cachedProducts, err := r.cache.GetAllProducts()
	if err == nil {
		return cachedProducts, nil
	}

	// Fallback to DB
	rows, err := r.db.Query("SELECT id, name, price, category FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Update Redis cache
	r.cache.SetAllProducts(products, 10*time.Minute)
	return products, nil
}

func (r *ProductRepo) GetProductByID(id int) (*models.Product, error) {
	// Try Redis cache
	product, err := r.cache.GetProductByID(id)
	if err == nil {
		return product, nil
	}

	// Fallback to DB
	var p models.Product
	err = r.db.QueryRow("SELECT id, name, price, category FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if err != nil {
		return nil, err
	}

	// Cache result
	r.cache.SetProductByID(id, &p, 10*time.Minute)
	return &p, nil
}
