package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"pos-module/model"
	"pos-module/service"
)

type CheckoutHandler struct {
	svc *service.CheckoutService
}

func NewCheckoutHandler(svc *service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{svc: svc}
}

// Checkout menangani POST /api/checkout
// Body: {"items":[{"sku":"SKU-001","qty":2}]}
func (h *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, modelError("method not allowed"))
		return
	}

	var req model.CheckoutRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, modelError("invalid json body: "+err.Error()))
		return
	}

	resp, err := h.svc.Calculate(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyCart),
			errors.Is(err, service.ErrInvalidQty),
			errors.Is(err, service.ErrUnknownItem),
			errors.Is(err, service.ErrOutOfStock),
			errors.Is(err, service.ErrNoCashier):
			writeJSON(w, http.StatusBadRequest, modelError(err.Error()))
		default:
			writeJSON(w, http.StatusInternalServerError, modelError("internal error"))
		}
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// ListProducts menangani GET /api/products
func (h *CheckoutHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, modelError("internal error"))
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func modelError(msg string) map[string]string {
	return map[string]string{"error": msg}
}
