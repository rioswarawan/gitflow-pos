package service

import (
	"errors"
	"fmt"

	"pos-module/model"
	"pos-module/repo"
)

var (
	ErrEmptyCart   = errors.New("items must not be empty")
	ErrInvalidQty  = errors.New("qty must be greater than 0")
	ErrUnknownItem = errors.New("unknown product sku")
	ErrOutOfStock  = errors.New("insufficient stock")
	ErrNoCashier   = errors.New("cashier must not be empty")
)

// CheckoutService menghitung total belanja dari daftar item.
type CheckoutService struct {
	products *repo.ProductRepo
}

func NewCheckoutService(products *repo.ProductRepo) *CheckoutService {
	return &CheckoutService{products: products}
}

// Calculate memvalidasi tiap item, mengambil harga dari database,
// menghitung subtotal per item (qty * harga) dan total keseluruhan.
func (s *CheckoutService) Calculate(req model.CheckoutRequest) (model.CheckoutResponse, error) {
	if len(req.Items) == 0 {
		return model.CheckoutResponse{}, ErrEmptyCart
	}
	if req.Cashier == "" {
		return model.CheckoutResponse{}, ErrNoCashier
	}
	cashier := req.Cashier

	resp := model.CheckoutResponse{Cashier: cashier, Items: make([]model.CartItem, 0, len(req.Items))}

	for _, it := range req.Items {
		if it.Qty <= 0 {
			return model.CheckoutResponse{}, fmt.Errorf("%w: %s", ErrInvalidQty, it.SKU)
		}
		p, err := s.products.FindBySKU(it.SKU)
		if errors.Is(err, repo.ErrNotFound) {
			return model.CheckoutResponse{}, fmt.Errorf("%w: %s", ErrUnknownItem, it.SKU)
		}
		if err != nil {
			return model.CheckoutResponse{}, err
		}
		if it.Qty > p.Stock {
			return model.CheckoutResponse{}, fmt.Errorf("%w: %s (stok %d)", ErrOutOfStock, p.SKU, p.Stock)
		}

		subtotal := float64(it.Qty) * p.Price
		resp.Items = append(resp.Items, model.CartItem{
			SKU:      p.SKU,
			Product:  p.Name,
			Qty:      it.Qty,
			Price:    p.Price,
			Subtotal: subtotal,
		})
		resp.TotalQty += it.Qty
		resp.TotalPrice += subtotal
	}
	return resp, nil
}

// ListProducts menampilkan semua produk yang tersedia.
func (s *CheckoutService) ListProducts() ([]model.Product, error) {
	return s.products.List()
}
