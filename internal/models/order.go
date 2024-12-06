package models

import "github.com/guregu/null/zero"

type OrderRequest struct {
	CouponCode zero.String `json:"coupon_code"`
	Items      []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID int     `json:"product_id"`
	OrderID   int     `json:"order_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID         int         `json:"id"`
	CouponCode string      `json:"coupon_code"`
	Items      []OrderItem `json:"items"`
	FinalPrice float64     `json:"final_price"`
}
