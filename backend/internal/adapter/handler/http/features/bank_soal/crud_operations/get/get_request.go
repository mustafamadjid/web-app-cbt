package httpx

type ListBankSoalRequest struct {
	Search       string
	Limit        int
	Offset       int
	TingkatKelas *int
	Mapel        *int
}
