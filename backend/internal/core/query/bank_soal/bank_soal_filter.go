package query

type BankSoalFilter struct {
	Search string
	Limit  int
	Offset int

	TingkatKelas *int
	Mapel *int
}