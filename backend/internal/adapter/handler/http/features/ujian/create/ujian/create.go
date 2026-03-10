package httpx

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/create"
)

type CreateRuangUjianHandler struct {
	svc *ujian_service.CreateUjianService
}

func NewCreateUjianHandler(svc *ujian_service.CreateUjianService) *CreateRuangUjianHandler {
	return &CreateRuangUjianHandler{svc: svc}
}

func (s *CreateRuangUjianHandler) CreateUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest CreatePenjadwalanUjianRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		case strings.Contains(err.Error(), "parsing time") || strings.Contains(err.Error(), "cannot parse"):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid time format")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	validatedRequest, err := ValidateCreateUjianRequestFields(dataRequest)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest = validatedRequest

	var idNamaKelas *ujian.ID
	if dataRequest.IdNamaKelas != nil {
		id := ujian.ID(*dataRequest.IdNamaKelas)
		idNamaKelas = &id
	}

	data := ujian.PenjadwalanUjian{
		Ujian: ujian.Ujian{
			IdBankSoal:     ujian.ID(dataRequest.IdBankSoal),
			IdKelas:        ujian.ID(dataRequest.IdKelas),
			IdGuru:         ujian.ID(dataRequest.IdGuru),
			IdNamaKelas:    idNamaKelas,
			NamaUjian:      dataRequest.NamaUjian,
			DeskripsiUjian: dataRequest.DeskripsiUjian,
			AcakSoal:       dataRequest.AcakSoal,
		},
		JadwalUjian: ujian.JadwalUjian{
			IdSesi:       ujian.ID(dataRequest.IdSesi),
			IdRuangan:    ujian.ID(dataRequest.IdRuangan),
			TanggalUjian: dataRequest.TanggalUjian,
			WaktuMulai:   dataRequest.WaktuMulai,
			WaktuSelesai: dataRequest.WaktuSelesai,
			StatusUjian:  ujian.StatusUjian(dataRequest.StatusUjian),
			Token:        dataRequest.Token,
			IdPengawas:   ujian.ID(dataRequest.IdPengawas),
		},
	}
	if err := s.svc.CreateUjianService(r.Context(), data); err != nil {
		logger.Error(r.Context(), "failed create ujian", "layer", "handler.http", "op", "ujian.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
			return
		case errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing field")
			return
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid input")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
			return
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
