# POS Module

Aplikasi Point of Sale (POS) sederhana berbahasa Go dengan arsitektur **monolith**, database **SQLite** (file lokal, tanpa server), dan **REST API**.

Fitur utama: menghitung total harga barang yang dibeli (`jumlah * harga` per item, lalu dijumlahkan).

## Struktur Project

```
pos-module/
├── main.go                  # Entry point: wiring dependency + HTTP server
├── model/
│   ├── product.go           # Struct: Product, CartItem, CheckoutRequest/Response
│   └── db.go                # Inisialisasi SQLite + skema + seed data
├── repo/
│   └── product.go           # Akses data produk (query SQLite)
├── service/
│   └── checkout.go          # Business logic kalkulasi total belanja
├── handler/
│   └── checkout.go          # HTTP handler (JSON request/response)
└── go.mod
```

## Prasyarat

- Go **1.25+** — cek dengan perintah:

```bash
go version
```

## Instalasi

1. Clone / masuk ke direktori project:

```bash
cd pos-module
```

2. Download dependency:

```bash
go mod tidy
```

3. (Opsional) Build binary:

```bash
go build -o bin/pos-module .
```

> Tidak perlu instalasi SQLite terpisah. Project memakai driver `modernc.org/sqlite` (pure Go, tanpa CGO). File database `pos.db` dibuat otomatis saat pertama kali dijalankan, beserta 5 data produk contoh.

## Cara Run

Jalankan langsung dengan `go run`:

```bash
go run .
```

Atau jalankan binary hasil build:

```bash
./bin/pos-module
```

Output:

```
pos-module listening on :8080 (db=pos.db)
```

### Konfigurasi (environment variable)

| Variable  | Default  | Keterangan                     |
|-----------|----------|--------------------------------|
| `POS_ADDR` | `:8080`  | Alamat & port HTTP server      |
| `POS_DB`   | `pos.db` | Lokasi file database SQLite    |

Contoh dengan konfigurasi kustom:

```bash
POS_ADDR=:3000 POS_DB=/tmp/pos.db go run .
```

---

## API

### `POST /api/checkout`

Menghitung total harga barang yang dibeli. Server mengambil harga dari database berdasarkan SKU, lalu menghitung subtotal tiap item (`qty * harga`) dan total keseluruhan.

#### Request

- Header: `Content-Type: application/json`
- Body:

| Field        | Tipe   | Wajib | Keterangan                        |
|--------------|--------|-------|-----------------------------------|
| `items`      | array  | Ya    | Daftar barang yang dibeli         |
| `items[].sku`| string | Ya    | SKU produk (lihat `GET /api/products`) |
| `items[].qty`| number | Ya    | Jumlah barang, harus > 0          |

Contoh:

```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -d '{"items":[{"sku":"SKU-001","qty":3},{"sku":"SKU-004","qty":1}]}'
```

#### Response sukses — `200 OK`

| Field               | Tipe    | Keterangan                              |
|---------------------|---------|-----------------------------------------|
| `items`             | array   | Detail tiap item (diisi oleh server)    |
| `items[].product`   | string  | Nama produk                             |
| `items[].price`     | number  | Harga satuan                            |
| `items[].subtotal`  | number  | `qty * price`                           |
| `total_qty`         | number  | Total jumlah barang                     |
| `total_price`       | number  | Total harga seluruh item                |

```json
{
  "items": [
    {"sku": "SKU-001", "product": "Indomie Goreng",  "qty": 3, "price": 3500,  "subtotal": 10500},
    {"sku": "SKU-004", "product": "Beras 5kg",       "qty": 1, "price": 68000, "subtotal": 68000}
  ],
  "total_qty": 4,
  "total_price": 78500
}
```

#### Response error — `400 Bad Request`

Format: `{"error": "<pesan>"}`

| Kondisi                     | Contoh pesan                            |
|-----------------------------|-----------------------------------------|
| Body bukan JSON valid       | `invalid json body: ...`                |
| `items` kosong              | `items must not be empty`               |
| `qty` <= 0                  | `qty must be greater than 0: SKU-001`   |
| SKU tidak ditemukan         | `unknown product sku: SKU-999`          |
| Stok tidak cukup            | `insufficient stock: SKU-004 (stok 20)` |

Contoh:

```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -d '{"items":[{"sku":"SKU-999","qty":1}]}'
# {"error":"unknown product sku: SKU-999"}
```

---

### `GET /api/products`

Menampilkan daftar semua produk yang tersedia beserta stoknya. Tidak butuh request body.

#### Response sukses — `200 OK`

Content-Type: `application/json`

| Field        | Tipe    | Keterangan                     |
|--------------|---------|--------------------------------|
| `id`         | number  | ID produk (urut menaik)        |
| `sku`        | string  | Kode unik produk               |
| `name`       | string  | Nama produk                    |
| `price`      | number  | Harga satuan                   |
| `stock`      | number  | Sisa stok                      |

Contoh:

```bash
curl http://localhost:8080/api/products
```

```json
[
  {"id":1,"sku":"SKU-001","name":"Indomie Goreng","price":3500,"stock":100},
  {"id":2,"sku":"SKU-002","name":"Teh Kotak 250ml","price":5000,"stock":80},
  {"id":3,"sku":"SKU-003","name":"Kopi Sachet","price":2000,"stock":150},
  {"id":4,"sku":"SKU-004","name":"Beras 5kg","price":68000,"stock":20},
  {"id":5,"sku":"SKU-005","name":"Minyak Goreng 1L","price":17500,"stock":40}
]
```

#### Response error — `500 Internal Server Error`

Format: `{"error": "<pesan>"}`

| Kondisi             | Contoh pesan      |
|---------------------|-------------------|
| Gagal query database | `internal error` |

---

### Endpoint lain

| Method | Path      | Fungsi                    |
|--------|-----------|---------------------------|
| GET    | `/health` | Health check (balas `ok`) |
