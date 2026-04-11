package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseGetStatistikUjianRequest(params httprouter.Params) (GetStatistikUjianRequest, error) {
	rawIDJadwalUjian := strings.TrimSpace(params.ByName("idJadwalUjian"))
	idJadwalUjian, err := strconv.Atoi(rawIDJadwalUjian)
	if err != nil || idJadwalUjian <= 0 {
		return GetStatistikUjianRequest{}, errors.New("idJadwalUjian must be a positive number")
	}

	return GetStatistikUjianRequest{IDJadwalUjian: idJadwalUjian}, nil
}

func toGetStatistikUjianResponse(item ujian.StatistikUjian) GetStatistikUjianResponse {
	return GetStatistikUjianResponse{
		IDStatistikUjian: int(item.IDStatistikUjian),
		IDJadwalUjian:    int(item.IDJadwalUjian),
		NilaiTertinggi:   item.NilaiTertinggi,
		NilaiTerendah:    item.NilaiTerendah,
		RataRata:         item.NilaiRataRata,
		JumlahPeserta:    item.TotalPesertaUjian,
	}
}
