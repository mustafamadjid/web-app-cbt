package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseListSoalUjianSiswaRequest(params httprouter.Params) (ListSoalUjianSiswaRequest, error) {
	rawIDJadwalUjian := strings.TrimSpace(params.ByName("idJadwalUjian"))
	idJadwalUjian, err := strconv.Atoi(rawIDJadwalUjian)
	if err != nil || idJadwalUjian <= 0 {
		return ListSoalUjianSiswaRequest{}, errors.New("idJadwalUjian must be a positive number")
	}

	return ListSoalUjianSiswaRequest{IDJadwalUjian: idJadwalUjian}, nil
}

func toListSoalUjianSiswaResponses(items []ujian.SoalUjianSiswa) []ListSoalUjianSiswaResponse {
	response := make([]ListSoalUjianSiswaResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListSoalUjianSiswaResponse(item))
	}

	return response
}

func toListSoalUjianSiswaResponse(item ujian.SoalUjianSiswa) ListSoalUjianSiswaResponse {
	opsiJawaban := make([]ListSoalUjianSiswaOpsiJawabanResponse, 0, len(item.OpsiJawaban))
	for _, opsi := range item.OpsiJawaban {
		opsiJawaban = append(opsiJawaban, ListSoalUjianSiswaOpsiJawabanResponse{
			IDPilihanGanda: int(opsi.IdPilihanGanda),
			IsiPilihan:     opsi.IsiPilihan,
		})
	}

	return ListSoalUjianSiswaResponse{
		IDSoal:      int(item.IdSoal),
		TipeSoal:    item.TipeSoal,
		Pertanyaan:  item.Pertanyaan,
		Gambar:      item.Gambar,
		BobotSoal:   item.BobotSoal,
		NoUrutSoal:  item.NoUrutSoal,
		OpsiJawaban: opsiJawaban,
	}
}
