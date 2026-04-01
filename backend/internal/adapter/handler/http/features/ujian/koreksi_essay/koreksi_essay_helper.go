package httpx

import ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"

func toKoreksiEssayPayload(req KoreksiEssayRequest) []ujian.JawabanUjian {
	jawaban := make([]ujian.JawabanUjian, 0, len(req.Jawaban))

	for _, item := range req.Jawaban {
		jawaban = append(jawaban, ujian.JawabanUjian{
			IdJawaban:    ujian.ID(item.IDJawaban),
			EssayIsBenar: toKoreksiEssayBoolPointer(item.EssayIsBenar),
		})
	}

	return jawaban
}

func toKoreksiEssayBoolPointer(value *bool) *bool {
	if value == nil {
		return nil
	}

	result := *value
	return &result
}
