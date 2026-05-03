package httpx

import "errors"

func sanitizeAndValidateKoreksiEssayRequest(req KoreksiEssayRequest) (KoreksiEssayRequest, error) {
	for _, item := range req.Jawaban {
		if item.IDJawaban <= 0 {
			return KoreksiEssayRequest{}, errors.New("id_jawaban must be a positive number")
		}
		if item.EssayIsBenar == nil {
			return KoreksiEssayRequest{}, errors.New("essay_is_benar is required")
		}
	}

	return req, nil
}
