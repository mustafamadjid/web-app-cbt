package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
)

type GetBankSoalHandler struct {
	svc *bank_soal_service.GetBankSoalService
}

func NewGetBankSoalHandler(svc *bank_soal_service.GetBankSoalService) *GetBankSoalHandler {
	return &GetBankSoalHandler{svc: svc}
}

func (h *GetBankSoalHandler) GetBankSoal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filter, err := parseListBankSoalRequest(r)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	items, err := h.svc.GetBankSoalService(r.Context(), query.BankSoalFilter{
		Search:       filter.Search,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
		TingkatKelas: filter.TingkatKelas,
		Mapel:        filter.Mapel,
	})
	if err != nil {
		logger.Error(r.Context(), "failed getting bank soal", "layer", "adapter.http.handler", "op", "bank_soal.get", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toBankSoalResponses(items), "Success")
}

func (h *GetBankSoalHandler) GetBankSoalUploaded(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filter, err := parseListBankSoalRequest(r)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	items, err := h.svc.GetBankSoalUploadedService(r.Context(), query.BankSoalFilter{
		Search:       filter.Search,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
		TingkatKelas: filter.TingkatKelas,
		Mapel:        filter.Mapel,
	})
	if err != nil {
		logger.Error(r.Context(), "failed getting uploaded bank soal", "layer", "adapter.http.handler", "op", "bank_soal.get_uploaded", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toBankSoalResponses(items), "Success")
}

func (h *GetBankSoalHandler) GetBankSoalByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idBankSoal, err := strconv.Atoi(ps.ByName("idBankSoal"))
	if err != nil || idBankSoal <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id bank soal")
		return
	}

	item, err := h.svc.GetBankSoalByIdService(r.Context(), bank_soal.ID(idBankSoal))
	if err != nil {
		logger.Error(r.Context(), "failed getting bank soal by id", "layer", "adapter.http.handler", "op", "bank_soal.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id bank soal")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toBankSoalResponse(item), "Success")
}

func (h *GetBankSoalHandler) GetBankSoalByGuru(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idPengguna, err := strconv.Atoi(ps.ByName("idPengguna"))
	if err != nil || idPengguna <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id pengguna")
		return
	}

	items, err := h.svc.GetBankSoalByGuruService(r.Context(), bank_soal.ID(idPengguna))
	if err != nil {
		logger.Error(r.Context(), "failed getting bank soal by guru", "layer", "adapter.http.handler", "op", "bank_soal.get_by_guru", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id pengguna")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toBankSoalResponses(items), "Success")
}
