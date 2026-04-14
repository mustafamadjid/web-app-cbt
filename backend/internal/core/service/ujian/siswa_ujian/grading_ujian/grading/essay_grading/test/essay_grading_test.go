package gradingujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	gradingujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading"
	"github.com/stretchr/testify/assert"
)

type fakeEssayGradingRepo struct {
	updateErr    error
	updateCalled bool
	gotJawaban   []ujian.JawabanUjian
	gotGradedBy  ujian.ID
}

func (*fakeEssayGradingRepo) UpsertNilaiToHasilUjian(context.Context, float64, ujian.HasilUjian) error {
	return nil
}

func (*fakeEssayGradingRepo) UpsertJawabanBenarToStatistikSoal(context.Context, []ujian.StatistikSoal) error {
	return nil
}

func (*fakeEssayGradingRepo) UpsertJawabanSalahToStatistikSoal(context.Context, []ujian.StatistikSoal) error {
	return nil
}

func (f *fakeEssayGradingRepo) UpdateAndGradingEssayUjian(_ context.Context, jawabanSiswa []ujian.JawabanUjian, gradedBy ujian.ID) error {
	f.updateCalled = true
	f.gotJawaban = jawabanSiswa
	f.gotGradedBy = gradedBy
	return f.updateErr
}

func (*fakeEssayGradingRepo) UpsertToStatistikUjian(context.Context, ujian.ID) error {
	return nil
}

func boolPtr(v bool) *bool {
	return &v
}

func runEssayGradingCases(t *testing.T, prefix string) {
	t.Helper()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	validJawaban := []ujian.JawabanUjian{{IdJawaban: 11, EssayIsBenar: boolPtr(true)}}

	tests := []struct {
		name       string
		jawaban    []ujian.JawabanUjian
		gradedBy   ujian.ID
		repo       *fakeEssayGradingRepo
		wantErr    error
		wantCalled bool
	}{
		{
			name:       prefix + "1 -> graded by tidak valid",
			jawaban:    validJawaban,
			gradedBy:   0,
			repo:       &fakeEssayGradingRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: prefix + "2 -> id jawaban tidak valid",
			jawaban: []ujian.JawabanUjian{
				{IdJawaban: 0, EssayIsBenar: boolPtr(true)},
			},
			gradedBy:   7,
			repo:       &fakeEssayGradingRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: prefix + "3 -> status essay benar kosong",
			jawaban: []ujian.JawabanUjian{
				{IdJawaban: 11, EssayIsBenar: nil},
			},
			gradedBy:   7,
			repo:       &fakeEssayGradingRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       prefix + "4 -> repo update essay grading gagal",
			jawaban:    validJawaban,
			gradedBy:   7,
			repo:       &fakeEssayGradingRepo{updateErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       prefix + "5 -> berhasil grading essay",
			jawaban:    validJawaban,
			gradedBy:   7,
			repo:       &fakeEssayGradingRepo{},
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := gradingujian_service.NewEssayGradingUjianService(tc.repo)
			err := svc.EssayGrading(ctx, tc.jawaban, tc.gradedBy)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.updateCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.gradedBy, tc.repo.gotGradedBy)
				assert.Equal(t, tc.jawaban, tc.repo.gotJawaban)
			}
		})
	}
}

func TestEssayGradingUjianService_BranchCoverage(t *testing.T) {
	t.Parallel()
	runEssayGradingCases(t, "branch ")
}

func TestEssayGradingUjianService_BasisPath(t *testing.T) {
	t.Parallel()
	runEssayGradingCases(t, "path ")
}
