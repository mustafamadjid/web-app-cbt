package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/get"
)

type GetUjianHandler struct {
	svc *ujian_service.GetUjianService
}

func NewGetUjianHandler(svc *ujian_service.GetUjianService) *GetUjianHandler {
	return &GetUjianHandler{svc: svc}
}

func (h *GetUjianHandler) GetUjianById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idUjianStr := params.ByName("idUjian")
	idUjian, err := strconv.ParseInt(idUjianStr, 10, 64)
	if err != nil || idUjian <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id ujian tidak valid")
		return
	}

	item, err := h.svc.GetUjianByIdService(r.Context(), ujian.ID(idUjian))
	if err != nil {
		logger.Error(r.Context(), "failed get ujian by id", "layer", "adapter.http.handler", "op", "ujian.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id ujian tidak valid")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "ujian not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get ujian by id")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, ToUjianResponse(item), "Success")
}

func (h *GetUjianHandler) GetAllPesertaUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetPesertaUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid peserta ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_all_peserta.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetAllPesertaUjianService(r.Context(), ToPesertaUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get peserta ujian", "layer", "adapter.http.handler", "op", "ujian.get_all_peserta", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get peserta ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toPesertaUjianResponses(items), "Success")
}

func (h *GetUjianHandler) GetPesertaUjianBySiswa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDSiswa := strings.TrimSpace(params.ByName("idSiswa"))
	idSiswa, err := strconv.Atoi(rawIDSiswa)
	if err != nil || idSiswa <= 0 {
		logger.Info(r.Context(), "invalid id siswa", "layer", "adapter.http.handler", "op", "ujian.get_peserta_by_siswa", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id siswa")
		return
	}

	req, err := parseGetPesertaUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid peserta ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_peserta_by_siswa.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	item, err := h.svc.GetPesertaUjianBySiswaService(r.Context(), ujian.ID(idSiswa), ToPesertaUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get peserta ujian by siswa", "layer", "adapter.http.handler", "op", "ujian.get_peserta_by_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get peserta ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, ToPesertaUjianResponse(item), "Success")
}

func (h *GetUjianHandler) GetAllJawabanUjianSiswa(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetJawabanUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid jawaban ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_all_jawaban.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetAllJawabanUjianSiswaService(r.Context(), ToJawabanUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get jawaban ujian siswa", "layer", "adapter.http.handler", "op", "ujian.get_all_jawaban", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jawaban ujian siswa")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toJawabanUjianSiswaResponses(items), "Success")
}

func (h *GetUjianHandler) GetJawabanBySiswa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDSiswa := strings.TrimSpace(params.ByName("idSiswa"))
	idSiswa, err := strconv.Atoi(rawIDSiswa)
	if err != nil || idSiswa <= 0 {
		logger.Info(r.Context(), "invalid id siswa", "layer", "adapter.http.handler", "op", "ujian.get_jawaban_by_siswa", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id siswa")
		return
	}

	req, err := parseGetJawabanUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid jawaban ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_jawaban_by_siswa.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	item, err := h.svc.GetJawabanBySiswaService(r.Context(), ujian.ID(idSiswa), ToJawabanUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get jawaban by siswa", "layer", "adapter.http.handler", "op", "ujian.get_jawaban_by_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jawaban ujian siswa")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, ToJawabanUjianSiswaResponse(item), "Success")
}
