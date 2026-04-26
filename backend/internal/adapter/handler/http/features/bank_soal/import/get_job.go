package importhandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/get_job"
)

type GetJobHandler struct {
	svc *get_job.GetJobService
}

func NewGetJobHandler(svc *get_job.GetJobService) *GetJobHandler {
	return &GetJobHandler{svc: svc}
}

type jobResponse struct {
	IDJob      int64  `json:"id_job"`
	IDBankSoal int64  `json:"id_bank_soal"`
	Status     string `json:"status"`
	ErrorMsg   string `json:"error_msg,omitempty"`
	WarningMsg string `json:"warning_msg,omitempty"`
	TotalSoal  int    `json:"total_soal"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (h *GetJobHandler) GetJobByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idJobStr := ps.ByName("idJob")
	idJob, err := strconv.ParseInt(idJobStr, 10, 64)
	if err != nil || idJob <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id job tidak valid")
		return
	}

	job, err := h.svc.GetByID(r.Context(), idJob)
	if err != nil {
		logger.Error(r.Context(), "failed getting import job", "layer", "adapter.http.handler", "op", "import_soal.get_job", "err", err)
		if errors.Is(err, coreerror.ErrImportJobNotFound) {
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "import job tidak ditemukan")
			return
		}
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toJobResponse(job), "success")
}

func (h *GetJobHandler) GetJobsByBankSoal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idBankSoalStr := ps.ByName("idBankSoal")
	idBankSoal, err := strconv.ParseInt(idBankSoalStr, 10, 64)
	if err != nil || idBankSoal <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id bank soal tidak valid")
		return
	}

	jobs, err := h.svc.GetByBankSoal(r.Context(), idBankSoal)
	if err != nil {
		logger.Error(r.Context(), "failed getting import jobs", "layer", "adapter.http.handler", "op", "import_soal.get_jobs_by_bank_soal", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toJobResponses(jobs), "success")
}
