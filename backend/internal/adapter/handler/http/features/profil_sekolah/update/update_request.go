package httpx

type updateProfilSekolahRequest struct {
	EmailSekolah  *string
	NoTelpSekolah *string
	KepalaSekolah *string
	WakaSekolah   *string
	NamaSekolah   *string
	AlamatSekolah *string
}

func (u updateProfilSekolahRequest) IsEmpty() bool {
	return u.EmailSekolah == nil &&
		u.NoTelpSekolah == nil &&
		u.KepalaSekolah == nil &&
		u.WakaSekolah == nil &&
		u.NamaSekolah == nil &&
		u.AlamatSekolah == nil
}
