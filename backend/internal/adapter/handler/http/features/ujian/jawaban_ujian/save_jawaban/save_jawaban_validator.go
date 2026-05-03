package httpx

import (
	"errors"
	"strings"
	"unicode/utf8"
)

func sanitizeAndValidateSaveJawabanRequest(req SaveJawabanRequest) (SaveJawabanRequest, error) {
	if req.IDAttempt <= 0 {
		return SaveJawabanRequest{}, errors.New("id_attempt must be a positive number")
	}

	for i := range req.Jawaban {
		if req.Jawaban[i].IDSoal <= 0 {
			return SaveJawabanRequest{}, errors.New("id_soal must be a positive number")
		}

		if req.Jawaban[i].IDPilihan != nil && *req.Jawaban[i].IDPilihan <= 0 {
			return SaveJawabanRequest{}, errors.New("id_pilihan must be a positive number")
		}

		if req.Jawaban[i].JawabanEssay != nil {
			trimmed := strings.TrimSpace(*req.Jawaban[i].JawabanEssay)
			if trimmed == "" {
				req.Jawaban[i].JawabanEssay = nil
			} else {
				if !utf8.ValidString(trimmed) {
					return SaveJawabanRequest{}, errors.New("jawaban_essay must be valid utf-8")
				}
				req.Jawaban[i].JawabanEssay = &trimmed
			}
		}

		hasPilihan := req.Jawaban[i].IDPilihan != nil
		hasEssay := req.Jawaban[i].JawabanEssay != nil
		if hasPilihan == hasEssay {
			return SaveJawabanRequest{}, errors.New("each jawaban must contain exactly one of id_pilihan or jawaban_essay")
		}

		if req.Jawaban[i].WaktuJawab != nil && req.Jawaban[i].WaktuJawab.IsZero() {
			return SaveJawabanRequest{}, errors.New("waktu_jawab is invalid")
		}
	}

	return req, nil
}
