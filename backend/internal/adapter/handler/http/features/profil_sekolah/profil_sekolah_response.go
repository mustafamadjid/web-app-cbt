package httpx

type ProfilSekolahResponse struct {
	IDProfil      int     `json:"id_profil"`
	EmailSekolah  string  `json:"email_sekolah"`
	NoTelpSekolah string  `json:"no_telp_sekolah"`
	KepalaSekolah string  `json:"kepala_sekolah"`
	WakaSekolah   string  `json:"waka_sekolah"`
	NamaSekolah   string  `json:"nama_sekolah"`
	AlamatSekolah string  `json:"alamat_sekolah"`
	LogoSekolah   *string `json:"logo_sekolah"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}