package kelas

type ID int

type TingkatKelas struct {
	IdTingkatKelas ID
	TingkatKelas   int
}

type NamaKelas struct {
	IdNamaKelas    ID
	IdTingkatKelas ID
	NamaKelas      string
}

type FullKelasData struct {
	ItemsTingkatKelas []TingkatKelas
	ItemsNamaKelas    []NamaKelas
}

type KelasData struct {
	ItemsTingkatKelas TingkatKelas
	ItemsNamaKelas    NamaKelas
}
