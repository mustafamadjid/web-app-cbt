package importhandler_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	importhandler "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/import"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/get_job"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeImportJobRepo struct {
	createFn         func(ctx context.Context, job importsoal.ImportSoalJob) (int64, error)
	getByIDFn        func(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error)
	getByBankSoalFn  func(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error)
	createCalled     bool
	gotCreateJob     importsoal.ImportSoalJob
	gotGetByID       int64
	gotGetByBankSoal int64
}

func (f *fakeImportJobRepo) CreateJob(ctx context.Context, job importsoal.ImportSoalJob) (int64, error) {
	f.createCalled = true
	f.gotCreateJob = job
	if f.createFn != nil {
		return f.createFn(ctx, job)
	}
	return 77, nil
}

func (f *fakeImportJobRepo) GetPendingJobs(context.Context, int) ([]importsoal.ImportSoalJob, error) {
	return nil, nil
}

func (f *fakeImportJobRepo) UpdateJobStatus(context.Context, int64, importsoal.JobStatus, string, string, int) error {
	return nil
}

func (f *fakeImportJobRepo) GetJobByID(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error) {
	f.gotGetByID = jobID
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, jobID)
	}
	return importsoal.ImportSoalJob{IDJob: jobID, IDBankSoal: 5, Status: importsoal.StatusCompleted, TotalSoal: 3}, nil
}

func (f *fakeImportJobRepo) GetJobsByBankSoal(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error) {
	f.gotGetByBankSoal = bankSoalID
	if f.getByBankSoalFn != nil {
		return f.getByBankSoalFn(ctx, bankSoalID)
	}
	return []importsoal.ImportSoalJob{{IDJob: 1, IDBankSoal: bankSoalID, Status: importsoal.StatusPending}}, nil
}

func newImportUploadRequest(t *testing.T, method string, fileName string, withFile bool, withActor bool) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if withFile {
		part, err := writer.CreateFormFile("file", fileName)
		require.NoError(t, err)
		_, err = part.Write([]byte("docx bytes"))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	req := httptest.NewRequest(method, "/bank-soal/5/import", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if withActor {
		req = req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 9, Role: user.GURU}))
	}
	return req
}

func TestImportSoal_Success(t *testing.T) {
	repo := &fakeImportJobRepo{}
	handler := importhandler.NewImportHandler(create_job.NewCreateJobService(repo), t.TempDir())
	req := newImportUploadRequest(t, http.MethodPost, "soal.docx", true, true)
	rec := httptest.NewRecorder()

	handler.ImportSoal(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "5"}})

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id_job":77`)
	assert.True(t, repo.createCalled)
	assert.Equal(t, int64(5), repo.gotCreateJob.IDBankSoal)
	assert.Equal(t, int64(9), repo.gotCreateJob.IDPengguna)
	assert.Equal(t, importsoal.StatusPending, repo.gotCreateJob.Status)
	assert.Equal(t, ".docx", filepath.Ext(repo.gotCreateJob.FilePath))
}

func TestImportSoal_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		req        func(t *testing.T) *http.Request
		params     httprouter.Params
		wantStatus int
		wantBody   string
	}{
		{
			name: "method not allowed",
			req: func(t *testing.T) *http.Request {
				return newImportUploadRequest(t, http.MethodGet, "soal.docx", true, true)
			},
			params:     httprouter.Params{{Key: "idBankSoal", Value: "5"}},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name: "invalid bank soal id",
			req: func(t *testing.T) *http.Request {
				return newImportUploadRequest(t, http.MethodPost, "soal.docx", true, true)
			},
			params:     httprouter.Params{{Key: "idBankSoal", Value: "abc"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "id bank soal tidak valid",
		},
		{
			name: "missing actor",
			req: func(t *testing.T) *http.Request {
				return newImportUploadRequest(t, http.MethodPost, "soal.docx", true, false)
			},
			params:     httprouter.Params{{Key: "idBankSoal", Value: "5"}},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "unauthorized",
		},
		{
			name: "missing multipart content type",
			req: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/bank-soal/5/import", bytes.NewBufferString(""))
				return req.WithContext(middleware.WithActor(req.Context(), user.Actor{IdPengguna: 9, Role: user.GURU}))
			},
			params:     httprouter.Params{{Key: "idBankSoal", Value: "5"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "content type must be multipart/form-data",
		},
		{
			name: "missing file",
			req: func(t *testing.T) *http.Request {
				return newImportUploadRequest(t, http.MethodPost, "soal.docx", false, true)
			},
			params:     httprouter.Params{{Key: "idBankSoal", Value: "5"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "file tidak ditemukan",
		},
		{
			name: "invalid extension",
			req: func(t *testing.T) *http.Request {
				return newImportUploadRequest(t, http.MethodPost, "soal.pdf", true, true)
			},
			params:     httprouter.Params{{Key: "idBankSoal", Value: "5"}},
			wantStatus: http.StatusBadRequest,
			wantBody:   "format file harus .docx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := importhandler.NewImportHandler(create_job.NewCreateJobService(&fakeImportJobRepo{}), t.TempDir())
			rec := httptest.NewRecorder()

			handler.ImportSoal(rec, tc.req(t), tc.params)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestImportSoal_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "bank soal not found", err: coreerror.ErrBankSoalNotFound, wantStatus: http.StatusNotFound, wantBody: "bank soal tidak ditemukan"},
		{name: "conflict", err: coreerror.ErrConflict, wantStatus: http.StatusConflict, wantBody: "konflik import bank soal"},
		{name: "internal", err: errors.New("db down"), wantStatus: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeImportJobRepo{
				createFn: func(context.Context, importsoal.ImportSoalJob) (int64, error) {
					return 0, tc.err
				},
			}
			handler := importhandler.NewImportHandler(create_job.NewCreateJobService(repo), t.TempDir())
			req := newImportUploadRequest(t, http.MethodPost, "soal.docx", true, true)
			rec := httptest.NewRecorder()

			handler.ImportSoal(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "5"}})

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestGetImportJobByID(t *testing.T) {
	now := time.Date(2026, time.March, 1, 10, 0, 0, 0, time.UTC)
	repo := &fakeImportJobRepo{
		getByIDFn: func(_ context.Context, id int64) (importsoal.ImportSoalJob, error) {
			return importsoal.ImportSoalJob{
				IDJob:      id,
				IDBankSoal: 5,
				Status:     importsoal.StatusCompleted,
				TotalSoal:  12,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
	}
	handler := importhandler.NewGetJobHandler(get_job.NewGetJobService(repo))
	req := httptest.NewRequest(http.MethodGet, "/import-jobs/77", nil)
	rec := httptest.NewRecorder()

	handler.GetJobByID(rec, req, httprouter.Params{{Key: "idJob", Value: "77"}})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id_job":77`)
	assert.Contains(t, rec.Body.String(), `"status":"completed"`)
	assert.Equal(t, int64(77), repo.gotGetByID)
}

func TestGetImportJobByID_Errors(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		idJob      string
		repo       *fakeImportJobRepo
		wantStatus int
		wantBody   string
	}{
		{name: "method not allowed", method: http.MethodPost, idJob: "77", repo: &fakeImportJobRepo{}, wantStatus: http.StatusMethodNotAllowed, wantBody: "method not allowed"},
		{name: "invalid id", method: http.MethodGet, idJob: "abc", repo: &fakeImportJobRepo{}, wantStatus: http.StatusBadRequest, wantBody: "id job tidak valid"},
		{name: "not found", method: http.MethodGet, idJob: "77", repo: &fakeImportJobRepo{getByIDFn: func(context.Context, int64) (importsoal.ImportSoalJob, error) {
			return importsoal.ImportSoalJob{}, coreerror.ErrImportJobNotFound
		}}, wantStatus: http.StatusNotFound, wantBody: "import job tidak ditemukan"},
		{name: "internal", method: http.MethodGet, idJob: "77", repo: &fakeImportJobRepo{getByIDFn: func(context.Context, int64) (importsoal.ImportSoalJob, error) {
			return importsoal.ImportSoalJob{}, errors.New("db down")
		}}, wantStatus: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := importhandler.NewGetJobHandler(get_job.NewGetJobService(tc.repo))
			req := httptest.NewRequest(tc.method, "/import-jobs/"+tc.idJob, nil)
			rec := httptest.NewRecorder()

			handler.GetJobByID(rec, req, httprouter.Params{{Key: "idJob", Value: tc.idJob}})

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}

func TestGetImportJobsByBankSoal(t *testing.T) {
	repo := &fakeImportJobRepo{}
	handler := importhandler.NewGetJobHandler(get_job.NewGetJobService(repo))
	req := httptest.NewRequest(http.MethodGet, "/bank-soal/5/import-jobs", nil)
	rec := httptest.NewRecorder()

	handler.GetJobsByBankSoal(rec, req, httprouter.Params{{Key: "idBankSoal", Value: "5"}})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id_job":1`)
	assert.Equal(t, int64(5), repo.gotGetByBankSoal)
}

func TestGetImportJobsByBankSoal_Errors(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		idBankSoal string
		repo       *fakeImportJobRepo
		wantStatus int
		wantBody   string
	}{
		{name: "method not allowed", method: http.MethodPost, idBankSoal: "5", repo: &fakeImportJobRepo{}, wantStatus: http.StatusMethodNotAllowed, wantBody: "method not allowed"},
		{name: "invalid id", method: http.MethodGet, idBankSoal: "abc", repo: &fakeImportJobRepo{}, wantStatus: http.StatusBadRequest, wantBody: "id bank soal tidak valid"},
		{name: "internal", method: http.MethodGet, idBankSoal: "5", repo: &fakeImportJobRepo{getByBankSoalFn: func(context.Context, int64) ([]importsoal.ImportSoalJob, error) {
			return nil, errors.New("db down")
		}}, wantStatus: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := importhandler.NewGetJobHandler(get_job.NewGetJobService(tc.repo))
			req := httptest.NewRequest(tc.method, "/bank-soal/"+tc.idBankSoal+"/import-jobs", nil)
			rec := httptest.NewRecorder()

			handler.GetJobsByBankSoal(rec, req, httprouter.Params{{Key: "idBankSoal", Value: tc.idBankSoal}})

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBody)
		})
	}
}
