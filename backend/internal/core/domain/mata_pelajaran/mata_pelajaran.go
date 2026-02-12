package matapelajaran

type ID int

type MataPelajaran struct {
	IdMapel ID
	IdKelas ID
	KodeMapel string
	NamaMapel string
	Deskripsi string
}
