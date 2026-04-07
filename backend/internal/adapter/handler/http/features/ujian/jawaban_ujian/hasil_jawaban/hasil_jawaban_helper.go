package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseHasilJawabanUjianIDAttempt(params httprouter.Params) (int, error) {
	rawIDAttempt := strings.TrimSpace(params.ByName("idAttempt"))
	if rawIDAttempt == "" {
		return 0, errors.New("id attempt is required")
	}

	idAttempt, err := strconv.Atoi(rawIDAttempt)
	if err != nil || idAttempt <= 0 {
		return 0, errors.New("id attempt must be a positive number")
	}

	return idAttempt, nil
}

func toHasilJawabanUjianResponse(idAttempt int, items []ujian.HasilJawabanUjian) HasilJawabanUjianResponse {
	responseItems := make([]HasilJawabanUjianItemResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, toHasilJawabanUjianItemResponse(item))
	}

	return HasilJawabanUjianResponse{
		IDAttempt:    idAttempt,
		HasilJawaban: responseItems,
	}
}

func toHasilJawabanUjianItemResponse(item ujian.HasilJawabanUjian) HasilJawabanUjianItemResponse {
	opsiJawaban := make([]HasilJawabanUjianOpsiJawabanResponse, 0, len(item.SoalUjianSiswa.OpsiJawaban))
	for _, opsi := range item.SoalUjianSiswa.OpsiJawaban {
		opsiJawaban = append(opsiJawaban, HasilJawabanUjianOpsiJawabanResponse{
			IDPilihanGanda: int(opsi.IdPilihanGanda),
			IDSoal:         int(opsi.IdSoal),
			IsiPilihan:     opsi.IsiPilihan,
			IsBenar:        opsi.IsBenar,
		})
	}

	return HasilJawabanUjianItemResponse{
		IDSoal:            int(item.SoalUjianSiswa.IdSoal),
		IDBankSoalVersion: int(item.SoalUjianSiswa.IdBankSoalVersion),
		TipeSoal:          item.SoalUjianSiswa.TipeSoal,
		Pertanyaan:        item.SoalUjianSiswa.Pertanyaan,
		Gambar:            item.SoalUjianSiswa.Gambar,
		BobotSoal:         item.SoalUjianSiswa.BobotSoal,
		NoUrutSoal:        item.SoalUjianSiswa.NoUrutSoal,
		OpsiJawaban:       opsiJawaban,
		JawabanSiswa: HasilJawabanUjianJawabanSiswaResponse{
			IDJawaban:    toHasilJawabanNullableInt(item.JawabanSiswa.IdJawaban),
			IDPilihan:    toHasilJawabanIntPointer(item.JawabanSiswa.IdPilihan),
			JawabanEssay: item.JawabanSiswa.JawabanEssay,
			WaktuJawab:   httphelper.FormatRFC3339Ptr(item.JawabanSiswa.WaktuJawab),
			EssayIsBenar: item.JawabanSiswa.EssayIsBenar,
		},
	}
}

func toHasilJawabanIntPointer(value *ujian.ID) *int {
	if value == nil {
		return nil
	}

	result := int(*value)
	return &result
}

func toHasilJawabanNullableInt(value ujian.ID) *int {
	if value <= 0 {
		return nil
	}

	result := int(value)
	return &result
}
