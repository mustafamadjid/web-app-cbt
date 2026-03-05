package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/update"
)

type UpdateUjianHandler struct {
	svc *ujian_service.UpdateUjianService
}

func NewUpdateUjianHandler(svc *ujian_service.UpdateUjianService) *UpdateUjianHandler {
	return &UpdateUjianHandler{svc: svc}
}

func (h *UpdateUjianHandler) UpdateUjian(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDUjian := strings.TrimSpace(ps.ByName("idUjian"))
	idUjian, err := strconv.Atoi(rawIDUjian)
	if err != nil || idUjian <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ujian")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest UpdatePenjadwalanUjianRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if err := ValidateInputIDRequestUpdateUjian(dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id")
		return
	}

	if err := ValidateInputSafeRequestUpdateUjian(dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid input")
		return
	}

	patch := updatepatch.UpdatePenjadwalanUjian{
		Ujian: updatepatch.UpdateUjianPatch{
			IdBankSoal:     toIDUjianPointer(dataRequest.IdBankSoal),
			IdKelas:        toIDUjianPointer(dataRequest.IdKelas),
			IdNamaKelas:    toIDUjianPointer(dataRequest.IdNamaKelas),
			IdGuru:         toIDUjianPointer(dataRequest.IdGuru),
			NamaUjian:      dataRequest.NamaUjian,
			DeskripsiUjian: dataRequest.DeskripsiUjian,
			AcakSoal:       dataRequest.AcakSoal,
		},
		JadwalUjian: updatepatch.UpdateJadwalUjianPatch{
			IdSesi:       toIDUjianPointer(dataRequest.IdSesi),
			IdRuangan:    toIDUjianPointer(dataRequest.IdRuangan),
			IdPengawas:   toIDUjianPointer(dataRequest.IdPengawas),
			TanggalUjian: dataRequest.TanggalUjian,
			Token:        dataRequest.Token,
			WaktuMulai:   dataRequest.WaktuMulai,
			WaktuSelesai: dataRequest.WaktuSelesai,
			StatusUjian:  toStatusUjianPointer(dataRequest.StatusUjian),
		},
	}

	if err := h.svc.UpdateUjianService(r.Context(), ujian.ID(idUjian), patch); err != nil {
		logger.Error(r.Context(), "failed updating ujian", "layer", "adapter.http.handler", "op", "ujian.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrMissingField),
			errors.Is(err, coreerror.ErrNoFieldToUpdate),
			errors.Is(err, coreerror.ErrInvalidInput),
			errors.Is(err, coreerror.ErrInvalidInputSafe):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update ujian")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}

func toIDUjianPointer(value *int) *ujian.ID {
	if value == nil {
		return nil
	}

	id := ujian.ID(*value)
	return &id
}

func toStatusUjianPointer(value *string) *ujian.StatusUjian {
	if value == nil {
		return nil
	}

	status := ujian.StatusUjian(*value)
	return &status
}
