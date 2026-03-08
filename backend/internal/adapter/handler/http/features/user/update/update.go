package httpx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type UpdateHandler struct {
	svc           *user_service.UpdateTx
	storeImage    httphelper.ImageStore
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}

func NewUpdateUserHandler(svc *user_service.UpdateTx, storeImage httphelper.ImageStore, aktivitasUser *aktivitas_user_service.AktivitasUserService) *UpdateHandler {
	return &UpdateHandler{svc: svc, storeImage: storeImage, aktivitasUser: aktivitasUser}
}

func (h *UpdateHandler) UpdateGuru(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	cmd, err := h.parseUpdateGuruMultipart(w, r, ps)
	if err != nil {
		logger.Error(r.Context(), "failed parsing update guru request", "layer", "adapter.http.handler", "op", "user.update_guru", "err", err)
		h.writeRequestError(w, err)
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.update_guru", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.UpdateGuru(r.Context(), cmd, actor); err != nil {
		logger.Error(r.Context(), "failed updating guru", "layer", "adapter.http.handler", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		h.writeUpdateError(w, err, "guru")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.UPDATE,
			Description: "Memperbarui data guru",
			IpAddress:   httphelper.GetClientIP(r),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "user.update_guru.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}

func (h *UpdateHandler) UpdateSiswa(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	cmd, err := h.parseUpdateSiswaMultipart(w, r, ps)
	if err != nil {
		logger.Error(r.Context(), "failed parsing update siswa request", "layer", "adapter.http.handler", "op", "user.update_siswa", "err", err)
		h.writeRequestError(w, err)
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.update_siswa", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.UpdateSiswa(r.Context(), cmd, actor); err != nil {
		logger.Error(r.Context(), "failed updating siswa", "layer", "adapter.http.handler", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
		h.writeUpdateError(w, err, "siswa")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.UPDATE,
			Description: "Memperbarui data siswa",
			IpAddress:   httphelper.GetClientIP(r),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "user.update_siswa.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
