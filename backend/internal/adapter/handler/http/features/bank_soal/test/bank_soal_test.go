package bank_soal_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	banksoal_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/create"
	banksoal_delete "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/delete"
	banksoal_get "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/get"
	banksoal_update "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/crud_operations/update"
	banksoal_import "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/bank_soal/import"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	bank_soal_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/create"
	bank_soal_delete_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
	bank_soal_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
	bank_soal_update_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/update"
	import_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	import_get_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/get_job"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBankSoal(t *testing.T) {
	mockRepo := new(MockBankSoalRepository)
	svc := bank_soal_create_service.NewCreateBankSoalService(mockRepo)
	handler := banksoal_create.NewCreateBankSoalHandler(svc)

	t.Run("Success Create Bank Soal", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.GURU}
		ctx := middleware.WithActor(context.Background(), actor)

		reqBody := `{"id_mapel":1, "id_kelas":1, "nama_bank_soal":"Bank Soal 1", "deskripsi":"Deskripsi", "materi":"Materi"}`
		req := httptest.NewRequest(http.MethodPost, "/bank-soal", strings.NewReader(reqBody)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("CreateBankSoal", mock.Anything, mock.Anything).Return(nil).Once()

		handler.CreateBankSoal(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetBankSoal(t *testing.T) {
	mockRepo := new(MockBankSoalRepository)
	svc := bank_soal_get_service.NewGetBankSoalService(mockRepo)
	handler := banksoal_get.NewGetBankSoalHandler(svc)

	t.Run("Success Get Bank Soal By Guru", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.GURU}
		ctx := middleware.WithActor(context.Background(), actor)

		req := httptest.NewRequest(http.MethodGet, "/bank-soal/guru/1", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idPengguna", Value: "1"}}

		mockRepo.On("GetBankSoalByGuru", mock.Anything, bank_soal.ID(1)).Return([]bank_soal.BankSoal{}, nil).Once()

		handler.GetBankSoalByGuru(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success Get Bank Soal By ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bank-soal/10", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idBankSoal", Value: "10"}}

		mockRepo.On("GetBankSoalById", mock.Anything, bank_soal.ID(10)).Return(bank_soal.BankSoal{IdBankSoal: 10}, nil).Once()

		handler.GetBankSoalByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bank-soal/99", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idBankSoal", Value: "99"}}

		mockRepo.On("GetBankSoalById", mock.Anything, bank_soal.ID(99)).Return(bank_soal.BankSoal{}, coreerror.ErrNotFound).Once()

		handler.GetBankSoalByID(w, req, params)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateBankSoal(t *testing.T) {
	mockRepo := new(MockBankSoalRepository)
	svc := bank_soal_update_service.NewUpdateBankSoalService(mockRepo)
	handler := banksoal_update.NewUpdateBankSoalHandler(svc)

	t.Run("Success Update Bank Soal", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.GURU}
		ctx := middleware.WithActor(context.Background(), actor)

		reqBody := `{"nama_bank_soal":"Updated name"}`
		req := httptest.NewRequest(http.MethodPatch, "/bank-soal/10", strings.NewReader(reqBody)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idBankSoal", Value: "10"}}

		mockRepo.On("UpdateBankSoal", mock.Anything, bank_soal.ID(10), mock.Anything).Return(nil).Once()

		handler.UpdateBankSoal(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteBankSoal(t *testing.T) {
	mockRepo := new(MockBankSoalRepository)
	svc := bank_soal_delete_service.NewDeleteBankSoalService(mockRepo)
	handler := banksoal_delete.NewDeleteBankSoalHandler(svc)

	t.Run("Success Delete Bank Soal", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.GURU}
		ctx := middleware.WithActor(context.Background(), actor)

		req := httptest.NewRequest(http.MethodDelete, "/bank-soal/10", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idBankSoal", Value: "10"}}

		mockRepo.On("DeleteBankSoal", mock.Anything, bank_soal.ID(10)).Return(nil).Once()

		handler.DeleteBankSoal(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestImportSoal(t *testing.T) {
	mockJobRepo := new(MockImportSoalJobRepo)
	svc := import_create_service.NewCreateJobService(mockJobRepo)
	
	tempDir, err := os.MkdirTemp("", "test_uploads")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	handler := banksoal_import.NewImportHandler(svc, tempDir)

	t.Run("Success Import Soal", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.GURU}
		ctx := middleware.WithActor(context.Background(), actor)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "questions.docx")
		_, _ = io.WriteString(part, "fake docx content")
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/bank-soal/10/import", body).WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idBankSoal", Value: "10"}}

		mockJobRepo.On("CreateJob", mock.Anything, mock.MatchedBy(func(job importsoal.ImportSoalJob) bool {
			return job.IDBankSoal == 10 && job.IDPengguna == 1
		})).Return(int64(101), nil).Once()

		handler.ImportSoal(w, req, params)

		assert.Equal(t, http.StatusAccepted, w.Code)
	})
}

func TestGetImportJob(t *testing.T) {
	mockJobRepo := new(MockImportSoalJobRepo)
	svc := import_get_service.NewGetJobService(mockJobRepo)
	handler := banksoal_import.NewGetJobHandler(svc)

	t.Run("Success Get Job By ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/import-job/101", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idJob", Value: "101"}}

		mockJobRepo.On("GetJobByID", mock.Anything, int64(101)).Return(importsoal.ImportSoalJob{
			IDJob: 101,
			Status: importsoal.StatusPending,
		}, nil).Once()

		handler.GetJobByID(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/import-job/999", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idJob", Value: "999"}}

		mockJobRepo.On("GetJobByID", mock.Anything, int64(999)).Return(importsoal.ImportSoalJob{}, coreerror.ErrImportJobNotFound).Once()

		handler.GetJobByID(w, req, params)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
