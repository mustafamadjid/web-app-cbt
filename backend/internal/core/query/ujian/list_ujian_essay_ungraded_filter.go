package query

type ListUjianEssayUngradedFilter struct {
	Search string
	Limit  int
	Offset int

	TanggalUjian   *string
	Tahun          *string
	Bulan          *string
	TingkatKelasID 	*int
	NamaKelasID 	*int
	MapelID 		*int
	SesiID 			*int
}