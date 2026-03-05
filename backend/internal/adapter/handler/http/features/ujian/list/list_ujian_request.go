package httpx

type ListUjianRequest struct {
	Search         string
	Limit          int
	Offset         int
	Tanggal        *string
	Tahun          *string
	TingkatKelasID *int
	TingkatKelas   *int
	RuangUjianID   *int
}
