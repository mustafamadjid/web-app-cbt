package ujian

type StatistikSoal struct {
	IDStatistikSoal ID
	IDSoal ID `json:"id_soal"`
	IDUjian ID	`json:"id_ujian"` 
	JumlahJawabanBenar int
	JumlahJawabanSalah int	
}