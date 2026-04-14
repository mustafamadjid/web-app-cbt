package ujian_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	koreksi_essay "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/koreksi_essay"
	statistik_ujian "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/ujian/statistik_ujian"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	essaygrading_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading"
	statistikujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/statistik_ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGradingRepository struct {
	mock.Mock
}

func (m *MockGradingRepository) UpsertNilaiToHasilUjian(ctx context.Context, totalNilai float64, hasilUjian ujian.HasilUjian) error {
	return m.Called(ctx, totalNilai, hasilUjian).Error(0)
}

func (m *MockGradingRepository) UpsertJawabanBenarToStatistikSoal(ctx context.Context, s []ujian.StatistikSoal) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockGradingRepository) UpsertJawabanSalahToStatistikSoal(ctx context.Context, s []ujian.StatistikSoal) error {
	return m.Called(ctx, s).Error(0)
}

func (m *MockGradingRepository) UpdateAndGradingEssayUjian(ctx context.Context, j []ujian.JawabanUjian, g ujian.ID) error {
	return m.Called(ctx, j, g).Error(0)
}

func (m *MockGradingRepository) UpsertToStatistikUjian(ctx context.Context, id ujian.ID) error {
	return m.Called(ctx, id).Error(0)
}

type MockStatistikRepository struct {
	mock.Mock
}

func (m *MockStatistikRepository) GetStatistikUjianByIdJadwal(ctx context.Context, id ujian.ID) (ujian.StatistikUjian, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ujian.StatistikUjian), args.Error(1)
}

func TestKoreksiAndStatistikHandlers(t *testing.T) {
	mockGrading := new(MockGradingRepository)
	mockStatistik := new(MockStatistikRepository)

	essaySvc := essaygrading_service.NewEssayGradingUjianService(mockGrading)
	statSvc := statistikujian_service.NewStatistikUjianService(mockStatistik)

	koreksiHandler := koreksi_essay.NewKoreksiEssayHandler(essaySvc)
	statHandler := statistik_ujian.NewGetStatistikUjianHandler(statSvc)

	t.Run("Koreksi Essay Success", func(t *testing.T) {
		reqBody := `{"jawaban":[{"id_jawaban":1, "essay_is_benar":true}]}`
		req := httptest.NewRequest(http.MethodPatch, "/ujian/koreksi-essay", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
		ctx := middleware.WithActor(context.Background(), actor)
		req = req.WithContext(ctx)
		
		w := httptest.NewRecorder()

		mockGrading.On("UpdateAndGradingEssayUjian", mock.Anything, mock.Anything, ujian.ID(1)).Return(nil).Once()

		koreksiHandler.KoreksiEssay(w, req, nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Statistik Ujian Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ujian/statistik/1", nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{{Key: "idJadwalUjian", Value: "1"}}

		mockStatistik.On("GetStatistikUjianByIdJadwal", mock.Anything, ujian.ID(1)).Return(ujian.StatistikUjian{}, nil).Once()

		statHandler.GetStatistikUjian(w, req, params)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
