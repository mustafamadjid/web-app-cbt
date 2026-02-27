package query

type BankSoalFilter struct {
	Search string
	Limit  int
	Offset int

	IdPengguna *int
	TingkatKelas *int
	Mapel *int
}