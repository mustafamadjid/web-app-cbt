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
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
)

type CreatePengumumanHandler struct {
	svc           *pengumuman_service.CreatePengumumanService
	storeDocument httphelper.DocumentStore
}

func NewCreatePengumumanHandler(svc *pengumuman_service.CreatePengumumanService, storeDocument httphelper.DocumentStore) *CreatePengumumanHandler {
	return &CreatePengumumanHandler{svc: svc, storeDocument: storeDocument}
}

func (h *CreatePengumumanHandler) CreatePengumuman(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.MultipartHeaderBodyValidator(w, r, 10<<20); err != nil {
		logger.Error(r.Context(), "failed parsing multipart form", "layer", "adapter.http.handler", "op", "pengumuman.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrContentTypeMustMultipart):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be multipart/form-data")
		default:
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid multipart form")
		}
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "pengumuman.create", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get actor from context")
		return
	}

	req := CreatePengumumanRequest{
		JudulPengumuman:          r.FormValue("judul_pengumuman"),
		IsiPengumuman:            r.FormValue("isi_pengumuman"),
		TanggalRilisPengumuman:   r.FormValue("tanggal_rilis_pengumuman"),
		TanggalSelesaiPengumuman: r.FormValue("tanggal_selesai_pengumuman"),
	}

	if req.JudulPengumuman == "" || req.IsiPengumuman == "" || req.TanggalRilisPengumuman == "" || req.TanggalSelesaiPengumuman == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing fields")
		return
	}

	judulPengumuman, err := validator.ValidateRequiredPrintableText(req.JudulPengumuman, "judul_pengumuman")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	req.JudulPengumuman = judulPengumuman

	dokumenPath := ""
	relativePathPtr, err := httphelper.StoreFileToDisk(r, "dokumen_pengumuman", false, h.storeDocument.SaveDocumentRelative)
	if err != nil {
		if errors.Is(err, coreerror.ErrFileTooLarge) {
			httpResponse.WriteErr(w, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
			return
		}

		logger.Error(r.Context(), "failed saving dokumen pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.create", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid dokumen_pengumuman")
		return
	}
	if relativePathPtr != nil {
		dokumenPath = *relativePathPtr
	}

	payload := pengumuman.Pengumuman{
		IdPengguna:               pengumuman.ID(actor.IdPengguna),
		JudulPengumuman:          req.JudulPengumuman,
		IsiPengumuman:            req.IsiPengumuman,
		TanggalRilisPengumuman:   req.TanggalRilisPengumuman,
		TanggalSelesaiPengumuman: req.TanggalSelesaiPengumuman,
		DokumenPengumuman:        dokumenPath,
	}

	if err := h.svc.CreatePengumuman(r.Context(), payload); err != nil {
		logger.Error(r.Context(), "failed creating pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
		case errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing fields")
		case errors.Is(err, coreerror.ErrInvalidDateFormat):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid date format")
		case errors.Is(err, coreerror.ErrInvalidInput), errors.Is(err, coreerror.ErrInvalidInputSafe):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create pengumuman")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
