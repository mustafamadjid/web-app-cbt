package httpx

type KoreksiEssayRequest struct {
	Jawaban []KoreksiEssayItemRequest `json:"jawaban"`
}

type KoreksiEssayItemRequest struct {
	IDJawaban    int   `json:"id_jawaban"`
	EssayIsBenar *bool `json:"essay_is_benar"`
}
