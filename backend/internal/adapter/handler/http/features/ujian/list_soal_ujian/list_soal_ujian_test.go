package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	responsehelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/require"
)

type fakeListSoalUjianService struct {
	items        []ujian.SoalUjianSiswa
	err          error
	called       bool
	capturedID   ujian.ID
	capturedAcak bool
}

func (f *fakeListSoalUjianService) ListSoalUjian(ctx context.Context, idBankSoal ujian.ID, acakSoal bool) ([]ujian.SoalUjianSiswa, error) {
	f.called = true
	f.capturedID = idBankSoal
	f.capturedAcak = acakSoal

	if f.err != nil {
		return nil, f.err
	}

	return f.items, nil
}

func TestListSoalUjianHandler_Success(t *testing.T) {
	svc := &fakeListSoalUjianService{
		items: []ujian.SoalUjianSiswa{
			{
				IdSoal:            101,
				IdBankSoalVersion: 202,
				TipeSoal:          "PILIHAN_GANDA",
				Pertanyaan:        "Apa ibu kota Indonesia?",
				Gambar:            "gambar.png",
				BobotSoal:         5,
				NoUrutSoal:        1,
				OpsiJawaban: []ujian.OpsiPilganUjian{
					{
						IdPilihanGanda: 1,
						IdSoal:         101,
						IsiPilihan:     "Jakarta",
						IsBenar:        true,
					},
				},
			},
		},
	}
	handler := &ListSoalUjianHandler{svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/ujian/soal/bank-soal/10?acak_soal=true", nil)
	rec := httptest.NewRecorder()

	handler.ListSoalUjian(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "10"}})

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, svc.called)
	require.Equal(t, ujian.ID(10), svc.capturedID)
	require.True(t, svc.capturedAcak)

	var response responsehelper.APIResponse[[]ListSoalUjianResponse]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.NotNil(t, response.Data)
	require.Len(t, *response.Data, 1)
	require.Equal(t, "Success", response.Message)
	require.Equal(t, 101, (*response.Data)[0].IDSoal)
	require.Len(t, (*response.Data)[0].OpsiJawaban, 1)
	require.Equal(t, "Jakarta", (*response.Data)[0].OpsiJawaban[0].IsiPilihan)
}

func TestListSoalUjianHandler_InvalidID(t *testing.T) {
	svc := &fakeListSoalUjianService{}
	handler := &ListSoalUjianHandler{svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/ujian/soal/bank-soal/abc", nil)
	rec := httptest.NewRecorder()

	handler.ListSoalUjian(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "abc"}})

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.False(t, svc.called)

	var response struct {
		Error responsehelper.APIError `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.Equal(t, "BAD_REQUEST", response.Error.Code)
}

func TestListSoalUjianHandler_InvalidAcakSoal(t *testing.T) {
	svc := &fakeListSoalUjianService{}
	handler := &ListSoalUjianHandler{svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/ujian/soal/bank-soal/10?acak_soal=bukan-bool", nil)
	rec := httptest.NewRecorder()

	handler.ListSoalUjian(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "10"}})

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.False(t, svc.called)
}

func TestListSoalUjianHandler_MissingIDError(t *testing.T) {
	svc := &fakeListSoalUjianService{err: coreerror.ErrMissingId}
	handler := &ListSoalUjianHandler{svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/ujian/soal/bank-soal/10", nil)
	rec := httptest.NewRecorder()

	handler.ListSoalUjian(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "10"}})

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var response struct {
		Error responsehelper.APIError `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.Equal(t, "MISSING_ID", response.Error.Code)
}

func TestListSoalUjianHandler_InternalError(t *testing.T) {
	svc := &fakeListSoalUjianService{err: errors.New("boom")}
	handler := &ListSoalUjianHandler{svc: svc}

	req := httptest.NewRequest(http.MethodGet, "/ujian/soal/bank-soal/10", nil)
	rec := httptest.NewRecorder()

	handler.ListSoalUjian(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "10"}})

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var response struct {
		Error responsehelper.APIError `json:"error"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.Equal(t, "INTERNAL_SERVER_ERROR", response.Error.Code)
}

func TestToListSoalUjianResponse_EmptyOptions(t *testing.T) {
	response := toListSoalUjianResponse(ujian.SoalUjianSiswa{
		IdSoal:            5,
		IdBankSoalVersion: 7,
		TipeSoal:          "ESSAY",
		Pertanyaan:        "Jelaskan",
		BobotSoal:         10,
		NoUrutSoal:        2,
	})

	require.NotNil(t, response.OpsiJawaban)
	require.Len(t, response.OpsiJawaban, 0)
}
