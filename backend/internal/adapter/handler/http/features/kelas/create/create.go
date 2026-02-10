package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
)

type CreateKelasHandler struct {
	svc           *kelas_service.CreateKelasService
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}

func NewCreateKelasHandler(svc *kelas_service.CreateKelasService) *CreateKelasHandler {
	return &CreateKelasHandler{svc: svc}
}

func (h *CreateKelasHandler) CreateTingkatKelas(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest CreateTingkatKelasReq
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

	if dataRequest.TingkatKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: tingkat kelas is required")
		return
	}

	cmdData := kelas_service.CreateTingkatKelasCmd{
		TingkatKelas: dataRequest.TingkatKelas,
	}

	if err := h.svc.CreateTingkatKelas(r.Context(), cmdData); err != nil {
		logger.Error(r.Context(), "failed creating kelas", "layer", "adapter.http.handler", "op", "kelas.create_tingkat_kelas_handler", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrTingkatKelasExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: tingkat kelas already exist")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create kelas")
			return
		}
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.create_guru", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.CREATE,
			Description: "Membuat akun guru",
			IpAddress:   httpx.GetClientIP(r),
		}

		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "aktivitas_user.create_tingkat_kelas", "err", err)
			return
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success create tingkat kelas")

}

func (h *CreateKelasHandler) CreateNamaKelas(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest CreateNamaKelasReq
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

	if dataRequest.IdTingkatKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id tingkat kelas is required")
		return
	}

	if strings.TrimSpace(dataRequest.NamaKelas) == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama kelas is required")
		return
	}

	cmdData := kelas_service.CreateNamaKelasCmd{
		IdTingkatKelas: kelas.ID(dataRequest.IdTingkatKelas),
		NamaKelas:      dataRequest.NamaKelas,
	}

	if err := h.svc.CreateNamaKelas(r.Context(), cmdData); err != nil {
		logger.Error(r.Context(), "failed creating kelas", "layer", "adapter.http.handler", "op", "kelas.create_nama_kelas_handler", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNamaKelasExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama kelas already exist")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create kelas")
			return
		}
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "kelas.create_nama_kelas", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.CREATE,
			Description: "Membuat nama kelas",
			IpAddress:   httpx.GetClientIP(r),
		}

		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "aktivitas_user.create_nama_kelas", "err", err)
			return
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success create nama kelas")
}
