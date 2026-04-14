package gradingujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	gradingujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/statistik_ujian"
	"github.com/stretchr/testify/assert"
)

type fakeGradingStatistikRepo struct {
	upsertErr    error
	upsertCalled bool
	gotAttemptID ujian.ID
}

func (*fakeGradingStatistikRepo) UpsertNilaiToHasilUjian(context.Context, float64, ujian.HasilUjian) error {
	return nil
}

func (*fakeGradingStatistikRepo) UpsertJawabanBenarToStatistikSoal(context.Context, []ujian.StatistikSoal) error {
	return nil
}

func (*fakeGradingStatistikRepo) UpsertJawabanSalahToStatistikSoal(context.Context, []ujian.StatistikSoal) error {
	return nil
}

func (*fakeGradingStatistikRepo) UpdateAndGradingEssayUjian(context.Context, []ujian.JawabanUjian, ujian.ID) error {
	return nil
}

func (f *fakeGradingStatistikRepo) UpsertToStatistikUjian(_ context.Context, idAttempt ujian.ID) error {
	f.upsertCalled = true
	f.gotAttemptID = idAttempt
	return f.upsertErr
}

func runGradingStatistikCases(t *testing.T, prefix string) {
	t.Helper()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idAttempt  int
		repo       *fakeGradingStatistikRepo
		wantErr    error
		wantCalled bool
	}{
		{
			name:       prefix + "1 -> id attempt tidak valid",
			idAttempt:  0,
			repo:       &fakeGradingStatistikRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       prefix + "2 -> repo upsert statistik ujian gagal",
			idAttempt:  9,
			repo:       &fakeGradingStatistikRepo{upsertErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       prefix + "3 -> berhasil upsert statistik ujian",
			idAttempt:  9,
			repo:       &fakeGradingStatistikRepo{},
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := gradingujian_service.NewStatistikUjianService(tc.repo)
			err := svc.StatistikUjian(ctx, tc.idAttempt)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.upsertCalled)
			if tc.wantCalled {
				assert.Equal(t, ujian.ID(tc.idAttempt), tc.repo.gotAttemptID)
			}
		})
	}
}

func TestGradingStatistikUjianService_BranchCoverage(t *testing.T) {
	t.Parallel()
	runGradingStatistikCases(t, "branch ")
}

func TestGradingStatistikUjianService_BasisPath(t *testing.T) {
	t.Parallel()
	runGradingStatistikCases(t, "path ")
}
