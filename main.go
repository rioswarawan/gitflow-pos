package main

import (
	"log"
	"net/http"
	"os"

	"pos-module/handler"
	"pos-module/model"
	"pos-module/repo"
	"pos-module/service"

	_ "modernc.org/sqlite" // driver SQLite (pure Go, tanpa CGO)
)

func main() {
	dsn := os.Getenv("POS_DB")
	if dsn == "" {
		dsn = "pos.db"
	}
	addr := os.Getenv("POS_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	db, err := model.InitDB(dsn)
	if err != nil {
		log.Fatalf("gagal inisialisasi db: %v", err)
	}
	defer db.Close()

	if err := model.SeedProducts(db); err != nil {
		log.Fatalf("gagal seed data: %v", err)
	}

	products := repo.NewProductRepo(db)
	checkoutSvc := service.NewCheckoutService(products)
	checkoutH := handler.NewCheckoutHandler(checkoutSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/checkout", checkoutH.Checkout)
	mux.HandleFunc("GET /api/products", checkoutH.ListProducts)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("pos-module listening on %s (db=%s)", addr, dsn)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server berhenti: %v", err)
	}
}
