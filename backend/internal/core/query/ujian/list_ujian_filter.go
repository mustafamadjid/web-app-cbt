package query

type ListUjianFilter struct {
	Search string
	Limit  int
	Offset int

	TanggalUjian   *string
	Tahun          *string
	TingkatKelasID *int
	TingkatKelas   *int
	RuangUjian     *int
}
