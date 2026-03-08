package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
)

type GetKelasHandler struct {
	svc *kelas_service.GetKelasService
}

func NewGetKelasHandler(svc *kelas_service.GetKelasService) *GetKelasHandler {
	return &GetKelasHandler{svc: svc}
}

func (h *GetKelasHandler) ListKelas(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListKelasFilters(r)
	if err != nil {
		logger.Info(r.Context(), "invalid kelas filters", "layer", "adapter.http.handler", "op", "kelas.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetFullKelas(r.Context(), filters)
	if err != nil {
		logger.Error(r.Context(), "failed listing kelas", "layer", "adapter.http.handler", "op", "kelas.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get kelas")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toFullKelasResponse(items), "Success")
}

func (h *GetKelasHandler) GetKelasByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDTingkat := strings.TrimSpace(params.ByName("idTingkatKelas"))
	rawIDNama := strings.TrimSpace(params.ByName("idNamaKelas"))
	if rawIDTingkat == "" || rawIDNama == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id_tingkat_kelas dan id_nama_kelas wajib diisi")
		return
	}

	idTingkatKelas, err := strconv.Atoi(rawIDTingkat)
	if err != nil || idTingkatKelas <= 0 {
		logger.Info(r.Context(), "invalid id tingkat kelas", "layer", "adapter.http.handler", "op", "kelas.get_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id_tingkat_kelas")
		return
	}

	idNamaKelas, err := strconv.Atoi(rawIDNama)
	if err != nil || idNamaKelas <= 0 {
		logger.Info(r.Context(), "invalid id nama kelas", "layer", "adapter.http.handler", "op", "kelas.get_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id_nama_kelas")
		return
	}

	data, err := h.svc.GetKelasById(r.Context(), idTingkatKelas, idNamaKelas)
	if err != nil {
		logger.Error(r.Context(), "failed get kelas by id", "layer", "adapter.http.handler", "op", "kelas.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get kelas")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toDataKelasResponse(data), "Success")
}
