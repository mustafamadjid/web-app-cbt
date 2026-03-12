package httpx

type ListUjianSiswaRequest struct {
	Search         string
	Limit          int
	Offset         int
	Tanggal        *string
	Tahun          *string
	Bulan          *string
	TingkatKelasID *int
	TingkatKelas   *int
	RuangUjianID   *int
	IDMapel        *int
	KategoriUjian  string
}
