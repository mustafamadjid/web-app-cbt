package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/update"
)

type UpdatePengumumanHandler struct {
	svc           *pengumuman_service.UpdatePengumumanService
	storeDocument httphelper.DocumentStore
}

func NewUpdatePengumumanHandler(svc *pengumuman_service.UpdatePengumumanService, storeDocument httphelper.DocumentStore) *UpdatePengumumanHandler {
	return &UpdatePengumumanHandler{svc: svc, storeDocument: storeDocument}
}

func (h *UpdatePengumumanHandler) UpdatePengumuman(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idPengumuman, err := strconv.Atoi(ps.ByName("idPengumuman"))
	if err != nil || idPengumuman <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id pengumuman")
		return
	}

	if err := httphelper.MultipartHeaderBodyValidator(w, r, 10<<20); err != nil {
		logger.Error(r.Context(), "failed parsing multipart form", "layer", "adapter.http.handler", "op", "pengumuman.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrContentTypeMustMultipart):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be multipart/form-data")
		default:
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid multipart form")
		}
		return
	}

	mf := r.MultipartForm
	getOptional := func(key string) (string, bool) {
		if mf == nil || mf.Value == nil {
			return "", false
		}

		vals, ok := mf.Value[key]
		if !ok || len(vals) == 0 {
			return "", false
		}

		return vals[0], true
	}

	req := UpdatePengumumanRequest{}

	if judul, ok := getOptional("judul_pengumuman"); ok {
		if judul == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: judul_pengumuman is required")
			return
		}

		if err := validator.ValidateInputSafe(judul, "judul_pengumuman"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}

		req.JudulPengumuman = &judul
	}

	if isi, ok := getOptional("isi_pengumuman"); ok {
		if isi == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: isi_pengumuman is required")
			return
		}

		req.IsiPengumuman = &isi
	}

	if tanggalRilis, ok := getOptional("tanggal_rilis_pengumuman"); ok {
		if tanggalRilis == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: tanggal_rilis_pengumuman is required")
			return
		}

		req.TanggalRilisPengumuman = &tanggalRilis
	}

	if tanggalSelesai, ok := getOptional("tanggal_selesai_pengumuman"); ok {
		if tanggalSelesai == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: tanggal_selesai_pengumuman is required")
			return
		}

		req.TanggalSelesaiPengumuman = &tanggalSelesai
	}

	dokumenPath, err := httphelper.StoreFileToDisk(r, "dokumen_pengumuman", false, h.storeDocument.SaveDocumentRelative)
	if err != nil {
		if errors.Is(err, coreerror.ErrFileTooLarge) {
			httpResponse.WriteErr(w, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
			return
		}

		logger.Error(r.Context(), "failed saving dokumen pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.update", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid dokumen_pengumuman")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "pengumuman.update", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get actor from context")
		return
	}

	patch := updatepatch.PengumumanUpdatePatch{
		IdPengguna:               ptrPengumumanID(pengumuman.ID(actor.IdPengguna)),
		JudulPengumuman:          req.JudulPengumuman,
		IsiPengumuman:            req.IsiPengumuman,
		TanggalRilisPengumuman:   req.TanggalRilisPengumuman,
		TanggalSelesaiPengumuman: req.TanggalSelesaiPengumuman,
		DokumenPengumuman:        dokumenPath,
	}

	if err := h.svc.UpdatePengumumanService(r.Context(), pengumuman.ID(idPengumuman), patch); err != nil {
		logger.Error(r.Context(), "failed updating pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
		case errors.Is(err, coreerror.ErrMissingField), errors.Is(err, coreerror.ErrNoFieldToUpdate):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing fields")
		case errors.Is(err, coreerror.ErrInvalidDateFormat):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid date format")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update pengumuman")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}

func ptrPengumumanID(id pengumuman.ID) *pengumuman.ID {
	return &id
}
