package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseListPesertaUjianSubmittedRequest(params httprouter.Params) (ListPesertaUjianSubmittedRequest, error) {
	rawIDJadwalUjian := strings.TrimSpace(params.ByName("idJadwalUjian"))
	idJadwalUjian, err := strconv.Atoi(rawIDJadwalUjian)
	if err != nil || idJadwalUjian <= 0 {
		return ListPesertaUjianSubmittedRequest{}, errors.New("idJadwalUjian must be a positive number")
	}

	return ListPesertaUjianSubmittedRequest{IDJadwalUjian: idJadwalUjian}, nil
}

func toListPesertaUjianSubmittedResponses(items []ujian.PesertaUjianSubmitted) []ListPesertaUjianSubmittedResponse {
	response := make([]ListPesertaUjianSubmittedResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListPesertaUjianSubmittedResponse(item))
	}

	return response
}

func toListPesertaUjianSubmittedResponse(item ujian.PesertaUjianSubmitted) ListPesertaUjianSubmittedResponse {
	return ListPesertaUjianSubmittedResponse{
		IDPesertaUjian: int(item.IdPesertaUjian),
		IDAttempt:      int(item.IdAttempt),
		IDSiswa:        int(item.IdSiswa),
		TingkatKelas:   item.TingkatKelas,
		NamaKelas:      item.NamaKelas,
		NamaLengkap:    item.NamaLengkap,
		NoAbsen:        item.NoAbsen,
		NilaiAkhir:     item.NilaiAkhir,
		WaktuMulai:     httphelper.FormatRFC3339Ptr(item.WaktuMulai),
		WaktuSubmit:    httphelper.FormatRFC3339Ptr(item.WaktuSubmit),
	}
}
