package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseGetJawabanUjianRequest(params httprouter.Params) (GetJawabanUjianRequest, error) {
	rawIDAttempt := strings.TrimSpace(params.ByName("idAttempt"))
	if rawIDAttempt == "" {
		return GetJawabanUjianRequest{}, errors.New("id attempt is required")
	}

	idAttempt, err := strconv.Atoi(rawIDAttempt)
	if err != nil || idAttempt <= 0 {
		return GetJawabanUjianRequest{}, errors.New("id attempt must be a positive number")
	}

	return GetJawabanUjianRequest{IDAttempt: idAttempt}, nil
}

func toGetJawabanUjianResponse(idAttempt int, items []ujian.JawabanUjian) GetJawabanUjianResponse {
	jawaban := make([]GetJawabanUjianItemResponse, 0, len(items))
	for _, item := range items {
		jawaban = append(jawaban, GetJawabanUjianItemResponse{
			IDJawaban:    int(item.IdJawaban),
			IDSoal:       int(item.IdSoal),
			IDPilihan:    toGetJawabanIntPointer(item.IdPilihan),
			JawabanEssay: item.JawabanEssay,
			WaktuJawab:   httphelper.FormatRFC3339Ptr(item.WaktuJawab),
		})
	}

	return GetJawabanUjianResponse{
		IDAttempt: idAttempt,
		Jawaban:   jawaban,
	}
}

func toGetJawabanIntPointer(value *ujian.ID) *int {
	if value == nil {
		return nil
	}

	result := int(*value)
	return &result
}
