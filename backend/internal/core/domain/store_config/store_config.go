package store_config

type DocumentStore struct {
	Dir      string // contoh: "./public/uploads"
	BaseURL  string // contoh: "https://example.com"
	Route    string // contoh: "/uploads" (route publik)
	MaxBytes int64  // contoh: 5 << 20
}

type ImageStore struct {
	Dir      string // contoh: "./public/uploads"
	BaseURL  string // contoh: "https://example.com"
	Route    string // contoh: "/uploads" (route publik)
	MaxBytes int64  // contoh: 5 << 20
}
