package httpx

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	siswaujianlist_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list"
)

type ListUjianSiswaHandler struct {
	svc *siswaujianlist_service.ListUjianSiswaService
}

func NewListUjianSiswaHandler(svc *siswaujianlist_service.ListUjianSiswaService) *ListUjianSiswaHandler {
	return &ListUjianSiswaHandler{svc: svc}
}

func (h *ListUjianSiswaHandler) ListUjianSiswa(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "siswa_ujian.list", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	req, err := parseListUjianSiswaRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid ujian siswa filters", "layer", "adapter.http.handler", "op", "siswa_ujian.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		return
	}

	items, err := h.svc.ListUjianSiswa(r.Context(), int(actor.IdPengguna), toListUjianSiswaFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed listing ujian siswa", "layer", "adapter.http.handler", "op", "siswa_ujian.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list ujian siswa")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListUjianSiswaResponses(items), "Success")
}

func toListUjianSiswaFilter(req ListUjianSiswaRequest) query.ListUjianFilter {
	filter := query.ListUjianFilter{
		Search:         req.Search,
		Limit:          req.Limit,
		Offset:         req.Offset,
		TanggalUjian:   req.Tanggal,
		Tahun:          req.Tahun,
		Bulan:          req.Bulan,
		TingkatKelasID: req.TingkatKelasID,
		TingkatKelas:   req.TingkatKelas,
		RuangUjian:     req.RuangUjianID,
		IDMapel:        req.IDMapel,
	}

	switch strings.ToLower(strings.TrimSpace(req.KategoriUjian)) {
	case string(query.MENDATANG):
		filter.KategoriUjian = query.MENDATANG
	case string(query.BERLANGSUNG):
		filter.KategoriUjian = query.BERLANGSUNG
	case string(query.SELESAI):
		filter.KategoriUjian = query.SELESAI
	case string(query.DIBATALKAN):
		filter.KategoriUjian = query.DIBATALKAN
	}

	return filter
}
