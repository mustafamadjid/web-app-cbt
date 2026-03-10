package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
)

func parseListBankSoalRequest(r *http.Request) (ListBankSoalRequest, error) {
	values := r.URL.Query()
	req := ListBankSoalRequest{
		Search: strings.TrimSpace(values.Get("q")),
	}
	if req.Search == "" {
		req.Search = strings.TrimSpace(values.Get("search"))
	}

	search, err := validator.ValidateOptionalPrintableText(req.Search, "search")
	if err != nil {
		return ListBankSoalRequest{}, err
	}
	req.Search = search

	if raw := strings.TrimSpace(values.Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListBankSoalRequest{}, errors.New("limit must be a number")
		}
		req.Limit = parsed
	}

	if raw := strings.TrimSpace(values.Get("offset")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListBankSoalRequest{}, errors.New("offset must be a number")
		}
		req.Offset = parsed
	}

	tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas"))
	if tingkatKelasRaw == "" {
		tingkatKelasRaw = strings.TrimSpace(values.Get("id_kelas"))
	}
	if tingkatKelasRaw != "" {
		parsed, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return ListBankSoalRequest{}, errors.New("tingkat_kelas must be a number")
		}
		req.TingkatKelas = &parsed
	}

	mapelRaw := strings.TrimSpace(values.Get("mapel"))
	if mapelRaw == "" {
		mapelRaw = strings.TrimSpace(values.Get("id_mapel"))
	}
	if mapelRaw != "" {
		parsed, err := strconv.Atoi(mapelRaw)
		if err != nil {
			return ListBankSoalRequest{}, errors.New("mapel must be a number")
		}
		req.Mapel = &parsed
	}

	return req, nil
}

func toBankSoalResponse(item bank_soal.BankSoal) BankSoalResponse {
	kelasLabel := "-"
	if item.TingkatKelas > 0 {
		kelasLabel = fmt.Sprintf("Kelas %d", item.TingkatKelas)
	}

	return BankSoalResponse{
		IDBankSoal:    int(item.IdBankSoal),
		IDMapel:       int(item.IdMapel),
		IDKelas:       int(item.IdKelas),
		IDPengguna:    int(item.IdPengguna),
		Mapel:         item.Mapel,
		GuruPembuat:   item.GuruPembuat,
		Kelas:         kelasLabel,
		NamaBankSoal:  item.NamaBankSoal,
		Deskripsi:     item.Deskripsi,
		Materi:        item.Materi,
		TanggalDibuat: item.TanggalDibuat,
		SoalUploaded:  item.SoalUploaded,
	}
}

func toBankSoalResponses(items []bank_soal.BankSoal) []BankSoalResponse {
	response := make([]BankSoalResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toBankSoalResponse(item))
	}

	return response
}
