package kelas

type ID int

type TingkatKelas struct {
	IdTingkatKelas ID
	TingkatKelas int
}

type NamaKelas struct {
	IdNamaKelas ID
	IdTingkatKelas ID
	NamaKelas string
}

type FullKelasData struct {
	DataTingkatKelas []TingkatKelas
	DataNamaKelas []NamaKelas
} 