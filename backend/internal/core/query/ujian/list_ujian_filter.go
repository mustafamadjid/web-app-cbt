package query


type kategoriUjian string

const (
	MENDATANG kategoriUjian = "mendatang"
	BERLANGSUNG kategoriUjian = "berlangsung"
	SELESAI kategoriUjian = "selesai"
)

type ListUjianFilter struct {
	Search string
	Limit  int
	Offset int

	TanggalUjian   *string
	Tahun          *string
	Bulan          *string
	TingkatKelasID *int
	TingkatKelas   *int
	RuangUjian     *int
	IDMapel        *int

	KategoriUjian kategoriUjian 
}
