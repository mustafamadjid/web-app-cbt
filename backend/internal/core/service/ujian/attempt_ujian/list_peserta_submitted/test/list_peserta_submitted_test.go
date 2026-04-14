package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/list_peserta_submitted"
	"github.com/stretchr/testify/assert"
)

type fakeListPesertaSubmittedRepo struct {
	listRet     []ujian.PesertaUjianSubmitted
	listErr     error
	listCalled  bool
	gotJadwalID ujian.ID
}

func (*fakeListPesertaSubmittedRepo) GetAttemptUjianById(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func (*fakeListPesertaSubmittedRepo) CreateAttemptUjian(context.Context, ujian.AttemptUjian) error {
	return nil
}

func (*fakeListPesertaSubmittedRepo) UpdateAttemptUjian(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
	return nil
}

func (*fakeListPesertaSubmittedRepo) DeleteAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func (*fakeListPesertaSubmittedRepo) SubmitAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func (f *fakeListPesertaSubmittedRepo) ListPesertaUjianAttemptSubmittedByIdJadwalUjian(_ context.Context, idJadwalUjian ujian.ID) ([]ujian.PesertaUjianSubmitted, error) {
	f.listCalled = true
	f.gotJadwalID = idJadwalUjian
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listRet, nil
}

func runListPesertaSubmittedCases(t *testing.T, prefix string) {
	t.Helper()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := []ujian.PesertaUjianSubmitted{
		{IdPesertaUjian: 1, IdAttempt: 11, IdSiswa: 101, NamaLengkap: "Siswa A"},
		{IdPesertaUjian: 2, IdAttempt: 12, IdSiswa: 102, NamaLengkap: "Siswa B"},
	}

	tests := []struct {
		name       string
		idJadwal   int
		repo       *fakeListPesertaSubmittedRepo
		wantErr    error
		wantCalled bool
		wantItems  []ujian.PesertaUjianSubmitted
	}{
		{
			name:       prefix + "1 -> id jadwal ujian tidak valid",
			idJadwal:   0,
			repo:       &fakeListPesertaSubmittedRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
			wantItems:  []ujian.PesertaUjianSubmitted{},
		},
		{
			name:       prefix + "2 -> repo list peserta submitted gagal",
			idJadwal:   13,
			repo:       &fakeListPesertaSubmittedRepo{listErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
			wantItems:  []ujian.PesertaUjianSubmitted{},
		},
		{
			name:       prefix + "3 -> berhasil list peserta submitted",
			idJadwal:   13,
			repo:       &fakeListPesertaSubmittedRepo{listRet: expected},
			wantCalled: true,
			wantItems:  expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewPesertaUjianSubmittedService(tc.repo)
			got, err := svc.ListPesertaUjianSubmitted(ctx, tc.idJadwal)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.listCalled)
			if tc.wantCalled {
				assert.Equal(t, ujian.ID(tc.idJadwal), tc.repo.gotJadwalID)
			}
			assert.Equal(t, tc.wantItems, got)
		})
	}
}

func TestPesertaUjianSubmittedService_BranchCoverage(t *testing.T) {
	t.Parallel()
	runListPesertaSubmittedCases(t, "branch ")
}

func TestPesertaUjianSubmittedService_BasisPath(t *testing.T) {
	t.Parallel()
	runListPesertaSubmittedCases(t, "path ")
}
