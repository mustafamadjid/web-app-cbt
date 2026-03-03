package query

type ListUjianFilter struct {
	Search string
	Limit  int
	Offset int

	TanggalUjian *string
	Tahun        *string
	TingkatKelas  *int
	RuangUjian 	  *int
}