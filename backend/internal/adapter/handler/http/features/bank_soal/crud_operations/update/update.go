package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
)

type UpdateBankSoalHandler struct {
	svc *bank_soal_service.UpdateBankSoalService
}

func NewUpdateBankSoalHandler(svc *bank_soal_service.UpdateBankSoalService) *UpdateBankSoalHandler {
	return &UpdateBankSoalHandler{svc: svc}
}

func (h *UpdateBankSoalHandler) UpdateBankSoal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idBankSoal, err := strconv.Atoi(ps.ByName("idBankSoal"))
	if err != nil || idBankSoal <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id bank soal")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be application/json")
		return
	}

	var req UpdateBankSoalRequest
	if err := httphelper.JSONDecoder(r, &req); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if req.NamaBankSoal != nil {
		if err := validator.ValidateInputSafe(*req.NamaBankSoal, "nama_bank_soal"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
	}
	if req.Deskripsi != nil {
		if err := validator.ValidateInputSafe(*req.Deskripsi, "deskripsi"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
	}
	if req.Materi != nil {
		if err := validator.ValidateInputSafe(*req.Materi, "materi"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
	}

	patch := updatepatch.UpdateBankSoalPatch{
		IdMapel:      toBankSoalIDPointer(req.IDMapel),
		IdKelas:      toBankSoalIDPointer(req.IDKelas),
		IdPengguna:   toBankSoalIDPointer(req.IDPengguna),
		NamaBankSoal: req.NamaBankSoal,
		Deskripsi:    req.Deskripsi,
		Materi:       req.Materi,
	}

	if err := h.svc.UpdateBankSoalService(r.Context(), bank_soal.ID(idBankSoal), patch); err != nil {
		logger.Error(r.Context(), "failed updating bank soal", "layer", "adapter.http.handler", "op", "bank_soal.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrMissingField),
			errors.Is(err, coreerror.ErrNoFieldToUpdate):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid update payload")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}

func toBankSoalIDPointer(v *int) *bank_soal.ID {
	if v == nil {
		return nil
	}

	id := bank_soal.ID(*v)
	return &id
}
