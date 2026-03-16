package httpx

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
)

func parseGetWaktuSelesaiUjianRequest(params httprouter.Params) (GetWaktuSelesaiUjianRequest, error) {
	rawIDJadwalUjian := strings.TrimSpace(params.ByName("idJadwalUjian"))
	idJadwalUjian, err := strconv.Atoi(rawIDJadwalUjian)
	if err != nil || idJadwalUjian <= 0 {
		return GetWaktuSelesaiUjianRequest{}, errors.New("idJadwalUjian must be a positive number")
	}

	return GetWaktuSelesaiUjianRequest{IDJadwalUjian: idJadwalUjian}, nil
}

func toGetWaktuSelesaiUjianResponse(idJadwalUjian int, waktuSelesai time.Time) GetWaktuSelesaiUjianResponse {
	return GetWaktuSelesaiUjianResponse{
		IDJadwalUjian: idJadwalUjian,
		WaktuSelesai:  httphelper.FormatRFC3339(waktuSelesai),
	}
}
