package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
)

type CreateBankSoalHandler struct {
	svc *bank_soal_service.CreateBankSoalService
}

func NewCreateBankSoalHandler(svc *bank_soal_service.CreateBankSoalService) *CreateBankSoalHandler {
	return &CreateBankSoalHandler{svc: svc}
}

func (h *CreateBankSoalHandler) CreateBankSoal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be application/json")
		return
	}

	var dataRequest CreateBankSoalRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		case errors.As(err, &maxErr):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: request body is too large")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	deskripsi, err := validator.ValidateRequiredPrintableText(dataRequest.Deskripsi, "deskripsi")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.Deskripsi = deskripsi

	materi, err := validator.ValidateRequiredPrintableText(dataRequest.Materi, "materi")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.Materi = materi

	namaBankSoal, err := validator.ValidateRequiredPrintableText(dataRequest.NamaBankSoal, "nama_bank_soal")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.NamaBankSoal = namaBankSoal

	if dataRequest.IdKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id kelas is required")
		return
	}

	if dataRequest.IdMapel <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id mapel is required")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get actor from context")
		return
	}

	data := bank_soal.BankSoal{
		IdMapel:      dataRequest.IdMapel,
		IdKelas:      dataRequest.IdKelas,
		IdPengguna:   bank_soal.ID(actor.IdPengguna),
		NamaBankSoal: dataRequest.NamaBankSoal,
		Deskripsi:    dataRequest.Deskripsi,
		Materi:       dataRequest.Materi,
	}

	if err := h.svc.CreateBankSoalService(r.Context(), data); err != nil {
		logger.Error(r.Context(), "failed creating bank soal", "layer", "adapter.http.handler", "op", "bank_soal.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId), errors.Is(err, coreerror.ErrMissingField), errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
