package importhandler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
)

type ImportHandler struct {
	svc       *create_job.CreateJobService
	uploadDir string // directory where uploaded .docx files are stored
}

func NewImportHandler(svc *create_job.CreateJobService, uploadDir string) *ImportHandler {
	return &ImportHandler{svc: svc, uploadDir: uploadDir}
}

type importResponse struct {
	IDJob int64 `json:"id_job"`
}

func (h *ImportHandler) ImportSoal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idBankSoalStr := ps.ByName("idBankSoal")
	idBankSoal, err := strconv.ParseInt(idBankSoalStr, 10, 64)
	if err != nil || idBankSoal <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id bank soal tidak valid")
		return
	}

	// Get actor from auth middleware context
	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		httpResponse.WriteErr(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	if err := httphelper.MultipartHeaderBodyValidator(w, r, 25<<20); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrContentTypeMustMultipart):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "content type must be multipart/form-data")
		default:
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "file terlalu besar atau request tidak valid")
		}
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "file tidak ditemukan dalam request")
		return
	}
	defer file.Close()

	// Validate .docx extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".docx" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "format file harus .docx")
		return
	}

	// Save file to upload directory
	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		logger.Error(r.Context(), "failed creating upload dir", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	dstFilename := strconv.FormatInt(idBankSoal, 10) + "_" + header.Filename
	dstPath := filepath.Join(h.uploadDir, dstFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		logger.Error(r.Context(), "failed creating destination file", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		logger.Error(r.Context(), "failed saving uploaded file", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	// Create import job
	result, err := h.svc.Execute(r.Context(), create_job.CreateJobCmd{
		IDBankSoal: idBankSoal,
		IDPengguna: int64(actor.IdPengguna),
		FilePath:   dstPath,
	})
	if err != nil {
		logger.Error(r.Context(), "failed creating import job", "layer", "adapter.http.handler", "op", "import_soal.import", "err", err)
		if errors.Is(err, coreerror.ErrBankSoalNotFound) {
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "bank soal tidak ditemukan")
			return
		}
		if errors.Is(err, coreerror.ErrConflict) {
			httpResponse.WriteErr(w, http.StatusConflict, "CONFLICT", "konflik import bank soal")
			return
		}
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}

	httpResponse.WriteOK(w, http.StatusAccepted, importResponse{IDJob: result.IDJob}, "import job created")
}
