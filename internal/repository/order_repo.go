package repository

import (
	"database/sql"
	"fmt"
	"order_food_online/internal/cache"
	"order_food_online/internal/models"
	"time"
)

type OrderRepository interface {
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id int) (*models.Order, error)
	PlaceOrder(orderReq models.OrderRequest) (*models.Order, error)
	CheckProductExists(id int) (bool, error)
}

type OrderRepo struct {
	db    *sql.DB
	cache cache.OrderCache
}

func NewOrderRepository(db *sql.DB, cache cache.OrderCache) OrderRepository {
	return &OrderRepo{db: db, cache: cache}
}

// GetAllOrders retrieves all orders, attempting to use cache first.
func (r *OrderRepo) GetAllOrders() ([]models.Order, error) {
	// Try Redis cache first
	cachedOrders, err := r.cache.GetAllOrders()
	if err == nil {
		return cachedOrders, nil
	}

	// Fallback to DB
	orders, err := r.fetchAllOrdersFromDB()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all orders from database: %w", err)
	}

	// Update Redis cache (non-blocking)
	_ = r.cache.SetAllOrders(orders, 10*time.Minute)

	return orders, nil
}

// GetOrderByID retrieves a specific order by ID, attempting to use cache first.
func (r *OrderRepo) GetOrderByID(id int) (*models.Order, error) {
	// Try Redis cache
	cachedOrder, err := r.cache.GetOrderByID(id)
	if err == nil {
		return cachedOrder, nil
	}

	// Fallback to DB
	order, err := r.fetchOrderByIDFromDB(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order by ID %d from database: %w", id, err)
	}

	// Cache result (non-blocking)
	_ = r.cache.SetOrderByID(id, order, 10*time.Minute)

	return order, nil
}

// PlaceOrder inserts a new order into the database and updates the cache.
func (r *OrderRepo) PlaceOrder(orderReq models.OrderRequest) (*models.Order, error) {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Insert the order
	var order models.Order
	err = tx.QueryRow(
		`INSERT INTO orders (coupon_code, final_price) VALUES ($1, 0) RETURNING id, coupon_code`,
		orderReq.CouponCode,
	).Scan(&order.ID, &order.CouponCode)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert the order items and calculate the final price
	var finalPrice float64
	for _, item := range orderReq.Items {
		// Retrieve product price
		var price float64
		err := tx.QueryRow(
			`SELECT price FROM products WHERE id = $1`, item.ProductID,
		).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product price for product ID %d: %w", item.ProductID, err)
		}

		// Insert the item
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`,
			order.ID, item.ProductID, item.Quantity, price,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item for product ID %d: %w", item.ProductID, err)
		}

		// Accumulate final price
		finalPrice += price * float64(item.Quantity)
	}

	// Update the order's final price
	_, err = tx.Exec(
		`UPDATE orders SET final_price = $1 WHERE id = $2`,
		finalPrice, order.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update final price for order ID %d: %w", order.ID, err)
	}

	// Set order details
	order.FinalPrice = finalPrice
	order.Items = orderReq.Items

	// Invalidate and update cache
	_ = r.invalidateAllOrdersCache()
	_ = r.cache.SetOrderByID(order.ID, &order, 10*time.Minute)

	return &order, nil
}

// invalidateAllOrdersCache refreshes the cache for all orders.
func (r *OrderRepo) invalidateAllOrdersCache() error {
	orders, err := r.fetchAllOrdersFromDB()
	if err != nil {
		return fmt.Errorf("failed to refresh cache for all orders: %w", err)
	}

	// Update Redis cache
	return r.cache.SetAllOrders(orders, 10*time.Minute)
}

// fetchAllOrdersFromDB retrieves all orders from the database.
func (r *OrderRepo) fetchAllOrdersFromDB() ([]models.Order, error) {
	rows, err := r.db.Query("SELECT id, coupon_code, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CouponCode, &order.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// fetchOrderByIDFromDB retrieves a specific order by ID from the database.
func (r *OrderRepo) fetchOrderByIDFromDB(id int) (*models.Order, error) {
	var order models.Order
	err := r.db.QueryRow("SELECT id, coupon_code, final_price FROM orders WHERE id = $1", id).
		Scan(&order.ID, &order.CouponCode, &order.FinalPrice)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CheckProductExists check if product with id exists
func (r *OrderRepo) CheckProductExists(productID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`,
		productID,
	).Scan(&exists)
	return exists, err
}
