package httpx

import (
	"strings"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func toSaveJawabanPayload(req SaveJawabanRequest) (ujian.ID, []ujian.JawabanUjian) {
	jawaban := make([]ujian.JawabanUjian, 0, len(req.Jawaban))
	for _, item := range req.Jawaban {
		jawaban = append(jawaban, ujian.JawabanUjian{
			IdSoal:       ujian.ID(item.IDSoal),
			IdPilihan:    toSaveJawabanIDPointer(item.IDPilihan),
			JawabanEssay: normalizeJawabanEssay(item.JawabanEssay),
			WaktuJawab:   item.WaktuJawab,
		})
	}

	return ujian.ID(req.IDAttempt), jawaban
}

func toSaveJawabanIDPointer(v *int) *ujian.ID {
	if v == nil {
		return nil
	}

	id := ujian.ID(*v)
	return &id
}

func normalizeJawabanEssay(v *string) *string {
	if v == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
