package httpx

import (
	"errors"
	"net/http"
	"strconv"
	
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/julienschmidt/httprouter"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/delete"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
)

type DeleteKelasHandler struct {
	svc *kelas_service.DeleteKelasService
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}


func NewDeleteKelasHandler(svc *kelas_service.DeleteKelasService, aktivitasUser *aktivitas_user_service.AktivitasUserService) *DeleteKelasHandler {
	return &DeleteKelasHandler{svc: svc, aktivitasUser: aktivitasUser}
}

func(h *DeleteKelasHandler)DeleteKelas(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}


	idNamaKelas, err := strconv.Atoi(ps.ByName("idNamaKelas"))
	if err != nil || idNamaKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id nama kelas")
		return
	}

	if err := h.svc.DeleteNamaKelas(r.Context(),idNamaKelas); err != nil {
		logger.Error(r.Context(), "failed to delete kelas", "layer", "adapter.http.handler", "op", "kelas.delete", "err", err)

		switch {
		case errors.Is(err,coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "not found")
			return
		case errors.Is(err,coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete kelas")
			return
		}
	}

	actor,ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "kelas.delete", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna: actor.IdPengguna,
			Action: aktivitas_user.DELETE,
			Description: "Menghapus kelas",
			IpAddress: httphelper.GetClientIP(r),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "kelas.delete.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(w,http.StatusOK,"Success delete kelas")
}

	


