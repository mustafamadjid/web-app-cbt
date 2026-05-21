package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseGetListAnalisisSoalRequest(params httprouter.Params) (GetListAnalisisSoalRequest, error) {
	rawIDJadwalUjian := strings.TrimSpace(params.ByName("idJadwalUjian"))
	idJadwalUjian, err := strconv.Atoi(rawIDJadwalUjian)
	if err != nil || idJadwalUjian <= 0 {
		return GetListAnalisisSoalRequest{}, errors.New("idJadwalUjian must be a positive number")
	}

	return GetListAnalisisSoalRequest{IDJadwalUjian: idJadwalUjian}, nil
}

func toGetListAnalisisSoalResponse(idJadwalUjian int, items []ujian.AnalisisSoal) GetListAnalisisSoalResponse {
	responseItems := make([]AnalisisSoalItemResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, toAnalisisSoalItemResponse(item))
	}

	return GetListAnalisisSoalResponse{
		IDJadwalUjian: idJadwalUjian,
		AnalisisSoal:  responseItems,
	}
}

func toAnalisisSoalItemResponse(item ujian.AnalisisSoal) AnalisisSoalItemResponse {
	opsiJawaban := make([]AnalisisSoalOpsiResponse, 0, len(item.Soal.OpsiJawaban))
	for _, opsi := range item.Soal.OpsiJawaban {
		var opsiContent *content.RichContent
		if !opsi.IsiPilihanContent.Empty() {
			value := opsi.IsiPilihanContent
			opsiContent = &value
		}

		opsiJawaban = append(opsiJawaban, AnalisisSoalOpsiResponse{
			IDPilihanGanda:    int(opsi.IdPilihanGanda),
			IsiPilihan:        opsi.IsiPilihan,
			IsiPilihanContent: opsiContent,
			IsBenar:           opsi.IsBenar,
		})
	}

	var pertanyaanContent *content.RichContent
	if !item.Soal.PertanyaanContent.Empty() {
		value := item.Soal.PertanyaanContent
		pertanyaanContent = &value
	}

	return AnalisisSoalItemResponse{
		IDSoal:             int(item.Soal.IdSoal),
		IDBankSoalVersion:  int(item.Soal.IdBankSoalVersion),
		TipeSoal:           item.Soal.TipeSoal,
		Pertanyaan:         item.Soal.Pertanyaan,
		PertanyaanContent:  pertanyaanContent,
		Gambar:             item.Soal.Gambar,
		BobotSoal:          item.Soal.BobotSoal,
		NoUrutSoal:         item.Soal.NoUrutSoal,
		JumlahJawabanBenar: item.JumlahJawabanBenar,
		JumlahJawabanSalah: item.JumlahJawabanSalah,
		OpsiJawaban:        opsiJawaban,
	}
}
