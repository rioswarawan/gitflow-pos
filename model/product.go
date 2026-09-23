package model

// Product merepresentasikan barang yang dijual.
type Product struct {
	ID    int64   `json:"id"`
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// CartItem adalah item yang dibeli: produk + jumlah.
type CartItem struct {
	SKU     string `json:"sku"`
	Product string `json:"product,omitempty"` // nama produk, diisi server
	Qty     int    `json:"qty"`
	Price   float64 `json:"price,omitempty"` // harga satuan, diisi server
	Subtotal float64 `json:"subtotal,omitempty"` // qty * harga, diisi server
}

// CheckoutRequest adalah body POST /api/checkout.
type CheckoutRequest struct {
	Items []CartItem `json:"items"`
}

// CheckoutResponse adalah hasil perhitungan total belanja.
type CheckoutResponse struct {
	Items      []CartItem `json:"items"`
	TotalQty   int        `json:"total_qty"`
	TotalPrice float64    `json:"total_price"`
}

// Error adalah format response error standar.
type Error struct {
	Error string `json:"error"`
}
