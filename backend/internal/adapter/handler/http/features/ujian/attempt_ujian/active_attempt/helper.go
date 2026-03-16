package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseGetActiveAttemptUjianRequest(r *http.Request) (GetActiveAttemptUjianRequest, error) {
	rawIDJadwalUjian := strings.TrimSpace(r.URL.Query().Get("id_jadwal_ujian"))
	if rawIDJadwalUjian == "" {
		return GetActiveAttemptUjianRequest{}, errors.New("id_jadwal_ujian is required")
	}

	idJadwalUjian, err := strconv.Atoi(rawIDJadwalUjian)
	if err != nil || idJadwalUjian <= 0 {
		return GetActiveAttemptUjianRequest{}, errors.New("id_jadwal_ujian must be a positive number")
	}

	return GetActiveAttemptUjianRequest{IDJadwalUjian: idJadwalUjian}, nil
}

func toGetActiveAttemptUjianResponse(item ujian.AttemptUjian) GetActiveAttemptUjianResponse {
	return GetActiveAttemptUjianResponse{
		IDAttempt:      int(item.IdAttempt),
		IDPesertaUjian: int(item.IdPesertaUjian),
		StatusAttempt:  string(item.StatusAttempt),
		WaktuMulai:     httphelper.FormatRFC3339Ptr(item.WaktuMulai),
		WaktuSubmit:    httphelper.FormatRFC3339Ptr(item.WaktuSubmit),
		DeadlineAt:     httphelper.FormatRFC3339Ptr(item.DeadlineAt),
	}
}
