package query

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type GuruListItem struct {
	IdPengguna  user.ID
	Role user.Role
	Username    string
	NoHp        string
	Email       user.Email
	NamaLengkap string
	StatusAkun  user.StatusAkun
	Nip         user.NIP
	Jabatan     string
	BidangStudi string
	Foto string
	JenisKelamin string
}

type ListGuruFilter struct {
    Search    string          // cari di nama / username / nip
    Status    *user.StatusAkun // optional: AKTIF/NONAKTIF
    Bidang    *string          // optional: "Matematika"
    Limit     int              // pagination
    Offset    int              // pagination
    SortBy    string           // "nama", "created_at", etc
    SortDesc  bool
}