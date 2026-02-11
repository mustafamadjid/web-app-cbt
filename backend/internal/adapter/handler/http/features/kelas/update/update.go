package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
)

type UpdateKelasHandler struct {
	svc           *kelas_service.UpdateKelasService
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}

func NewUpdateKelasHandler(svc *kelas_service.UpdateKelasService, aktivitasUser *aktivitas_user_service.AktivitasUserService) *UpdateKelasHandler {
	return &UpdateKelasHandler{svc: svc, aktivitasUser: aktivitasUser}
}

func (h *UpdateKelasHandler) UpdateNamaKelas(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	idNamaKelas, err := strconv.Atoi(ps.ByName("idNamaKelas"))
	if err != nil || idNamaKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id nama kelas")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest UpdateKelasRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if dec.Decode(&struct{}{}) != io.EOF {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if dataRequest.IdTingkatKelas == nil && dataRequest.NamaKelas == nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: no fields to update")
		return
	}

	patch := updatepatch.NamaKelasPatch{}
	if dataRequest.IdTingkatKelas != nil {
		if *dataRequest.IdTingkatKelas <= 0 {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id tingkat kelas is required")
			return
		}
		idTingkatKelas := kelas.ID(*dataRequest.IdTingkatKelas)
		patch.IdTingkatKelas = &idTingkatKelas
	}

	if dataRequest.NamaKelas != nil {
		namaKelas := strings.TrimSpace(*dataRequest.NamaKelas)
		if namaKelas == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama kelas is required")
			return
		}
		patch.NamaKelas = &namaKelas
	}

	if err := h.svc.UpdateNamaKelas(r.Context(), idNamaKelas, patch); err != nil {
		logger.Error(r.Context(), "failed updating nama kelas", "layer", "adapter.http.handler", "op", "kelas.update_nama_kelas_handler", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data kelas tidak ditemukan")
			return
		case errors.Is(err, coreerror.ErrMissingId), errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update nama kelas")
			return
		}
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "kelas.update_nama_kelas", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.UPDATE,
			Description: "Memperbarui nama kelas",
			IpAddress:   httphelper.GetClientIP(r),
		}

		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "aktivitas_user.update_nama_kelas", "err", err)
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success update nama kelas")
}
