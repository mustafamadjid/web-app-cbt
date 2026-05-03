package ujian_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	attempt_create "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/attempt_ujian/attempt/create"
	jawaban_save "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/jawaban_ujian/save_jawaban"
	submit_ujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/submit_ujian"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	ujiandomain "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	attempt_create_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/create"
	submit_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/submit_ujian"
	jawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAttemptUjian(t *testing.T) {
	mockRepo := new(MockAttemptUjianRepository)
	mockChecker := new(MockSiswaUjianChecker)
	svc := attempt_create_service.NewAttemptUjianService(mockChecker, mockRepo)
	handler := attempt_create.NewAttemptUjianHandler(svc)

	t.Run("Success Attempt Ujian", func(t *testing.T) {
		reqBody := `{"id_siswa":1, "id_jadwal_ujian":10, "token_ujian":"TOKEN123", "waktu_mulai":"2024-04-14T20:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/attempt-ujian", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockChecker.On("CheckValidSiswaInPesertaUjianById", mock.Anything, 1, 10).Return(true, 1, nil).Once()
		mockChecker.On("GetDeadlineUjian", mock.Anything, 10).Return(time.Now().Add(time.Hour), nil).Once()
		mockChecker.On("CheckTokenUjian", mock.Anything, "TOKEN123", 10).Return(true, nil).Once()
		mockRepo.On("CreateAttemptUjian", mock.Anything, mock.Anything).Return(nil).Once()

		handler.AttemptUjian(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Forbidden - Siswa Not Allowed", func(t *testing.T) {
		reqBody := `{"id_siswa":1, "id_jadwal_ujian":10, "token_ujian":"TOKEN123", "waktu_mulai":"2024-04-14T20:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/attempt-ujian", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockChecker.On("CheckValidSiswaInPesertaUjianById", mock.Anything, 1, 10).Return(false, 0, nil).Once()

		handler.AttemptUjian(w, req, nil)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Sanitizes Token Input", func(t *testing.T) {
		reqBody := `{"id_siswa":1, "id_jadwal_ujian":10, "token_ujian":" token123 ", "waktu_mulai":"2024-04-14T20:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/attempt-ujian", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockChecker.On("CheckValidSiswaInPesertaUjianById", mock.Anything, 1, 10).Return(true, 1, nil).Once()
		mockChecker.On("GetDeadlineUjian", mock.Anything, 10).Return(time.Now().Add(time.Hour), nil).Once()
		mockChecker.On("CheckTokenUjian", mock.Anything, "TOKEN123", 10).Return(true, nil).Once()
		mockRepo.On("CreateAttemptUjian", mock.Anything, mock.Anything).Return(nil).Once()

		handler.AttemptUjian(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSaveJawaban(t *testing.T) {
	mockRepo := new(MockJawabanUjianRepository)
	svc := jawaban_service.NewJawabanUjianService(mockRepo)
	handler := jawaban_save.NewSaveJawabanUjianHandler(svc)

	t.Run("Success Save Jawaban", func(t *testing.T) {
		reqBody := `{"id_attempt":101, "jawaban":[{"id_soal":1, "id_pilihan":5}]}`
		req := httptest.NewRequest(http.MethodPost, "/save-jawaban", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockRepo.On("SaveJawabanUjian", mock.Anything, ujiandomain.ID(101), mock.Anything).Return(nil).Once()

		handler.SaveJawabanUjian(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Rejects Invalid Jawaban Payload", func(t *testing.T) {
		reqBody := `{"id_attempt":101, "jawaban":[{"id_soal":1, "id_pilihan":5, "jawaban_essay":"opsi ganda"}]}`
		req := httptest.NewRequest(http.MethodPost, "/save-jawaban", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.SaveJawabanUjian(w, req, nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "BAD_REQUEST", resp["error"].(map[string]any)["code"])
	})
}

func TestSubmitUjian(t *testing.T) {
	mockRepo := new(MockAttemptUjianRepository)
	mockChecker := new(MockSiswaUjianChecker)
	svc := submit_service.NewSubmitUjianService(mockRepo, mockChecker)
	handler := submit_ujian.NewSubmitUjianHandler(svc)

	t.Run("Success Submit Ujian", func(t *testing.T) {
		actor := user.Actor{IdPengguna: 1, Role: user.SISWA}
		ctx := middleware.WithActor(context.Background(), actor)

		req := httptest.NewRequest(http.MethodPatch, "/submit-ujian/101", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idAttempt", Value: "101"}}

		mockChecker.On("CheckAttemptOwnershipBySiswa", mock.Anything, 1, ujiandomain.ID(101)).Return(true, nil).Once()
		mockRepo.On("SubmitAttemptUjian", mock.Anything, ujiandomain.ID(101)).Return(nil).Once()

		handler.SubmitUjian(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
