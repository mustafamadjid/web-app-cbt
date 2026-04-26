package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseListSoalUjianParams(r *http.Request, params httprouter.Params) (int, bool, error) {
	rawIDBankSoal := strings.TrimSpace(params.ByName("idBankSoal"))
	idBankSoal, err := strconv.Atoi(rawIDBankSoal)
	if err != nil || idBankSoal <= 0 {
		return 0, false, errors.New("idBankSoal must be a positive number")
	}

	rawAcakSoal := strings.TrimSpace(r.URL.Query().Get("acak_soal"))
	if rawAcakSoal == "" {
		return idBankSoal, false, nil
	}

	acakSoal, err := strconv.ParseBool(rawAcakSoal)
	if err != nil {
		return 0, false, errors.New("acak_soal must be a boolean")
	}

	return idBankSoal, acakSoal, nil
}

func toListSoalUjianResponses(items []ujian.SoalUjianSiswa) []ListSoalUjianResponse {
	response := make([]ListSoalUjianResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListSoalUjianResponse(item))
	}

	return response
}

func toListSoalUjianResponse(item ujian.SoalUjianSiswa) ListSoalUjianResponse {
	opsiJawaban := make([]ListSoalUjianOpsiJawabanResponse, 0, len(item.OpsiJawaban))
	for _, opsi := range item.OpsiJawaban {
		var opsiContent *content.RichContent
		if !opsi.IsiPilihanContent.Empty() {
			value := opsi.IsiPilihanContent
			opsiContent = &value
		}
		opsiJawaban = append(opsiJawaban, ListSoalUjianOpsiJawabanResponse{
			IDPilihanGanda:    int(opsi.IdPilihanGanda),
			IDSoal:            int(opsi.IdSoal),
			IsiPilihan:        opsi.IsiPilihan,
			IsiPilihanContent: opsiContent,
			IsBenar:           opsi.IsBenar,
		})
	}

	var pertanyaanContent *content.RichContent
	if !item.PertanyaanContent.Empty() {
		value := item.PertanyaanContent
		pertanyaanContent = &value
	}

	return ListSoalUjianResponse{
		IDSoal:            int(item.IdSoal),
		IDBankSoalVersion: int(item.IdBankSoalVersion),
		TipeSoal:          item.TipeSoal,
		Pertanyaan:        item.Pertanyaan,
		PertanyaanContent: pertanyaanContent,
		Gambar:            item.Gambar,
		BobotSoal:         item.BobotSoal,
		NoUrutSoal:        item.NoUrutSoal,
		OpsiJawaban:       opsiJawaban,
	}
}
