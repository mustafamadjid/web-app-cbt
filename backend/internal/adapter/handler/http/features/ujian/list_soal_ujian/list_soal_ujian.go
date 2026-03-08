package httpx

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
)

type listSoalUjianService interface {
	ListSoalUjian(ctx context.Context, idBankSoal ujian.ID, acakSoal bool) ([]ujian.SoalUjianSiswa, error)
}

type ListSoalUjianHandler struct {
	svc listSoalUjianService
}

func NewListSoalUjianHandler(svc *ujian_service.ListSoalUjianService) *ListSoalUjianHandler {
	return &ListSoalUjianHandler{svc: svc}
}

func (h *ListSoalUjianHandler) ListSoalUjian(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idBankSoal, acakSoal, err := parseListSoalUjianParams(r, params)
	if err != nil {
		logger.Info(r.Context(), "invalid list soal ujian request", "layer", "adapter.http.handler", "op", "ujian.list_soal", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	items, err := h.svc.ListSoalUjian(r.Context(), ujian.ID(idBankSoal), acakSoal)
	if err != nil {
		logger.Error(r.Context(), "failed list soal ujian", "layer", "adapter.http.handler", "op", "ujian.list_soal", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get soal ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListSoalUjianResponses(items), "Success")
}
