package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/submit_ujian"
	"github.com/stretchr/testify/assert"
)

type fakeSubmitAttemptRepo struct {
	submitErr    error
	submitCalled bool
	gotAttemptID ujian.ID
}

func (*fakeSubmitAttemptRepo) GetAttemptUjianById(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func (*fakeSubmitAttemptRepo) CreateAttemptUjian(context.Context, ujian.AttemptUjian) error {
	return nil
}

func (*fakeSubmitAttemptRepo) UpdateAttemptUjian(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
	return nil
}

func (*fakeSubmitAttemptRepo) DeleteAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func (f *fakeSubmitAttemptRepo) SubmitAttemptUjian(_ context.Context, idAttempt ujian.ID) error {
	f.submitCalled = true
	f.gotAttemptID = idAttempt
	return f.submitErr
}

func (*fakeSubmitAttemptRepo) ListPesertaUjianAttemptSubmittedByIdJadwalUjian(context.Context, ujian.ID) ([]ujian.PesertaUjianSubmitted, error) {
	return nil, nil
}

type fakeSubmitChecker struct {
	ownedRet     bool
	ownedErr     error
	checkCalled  bool
	gotSiswaID   int
	gotAttemptID ujian.ID
}

func (f *fakeSubmitChecker) CheckAttemptOwnershipBySiswa(_ context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	f.checkCalled = true
	f.gotSiswaID = idSiswa
	f.gotAttemptID = idAttempt
	if f.ownedErr != nil {
		return false, f.ownedErr
	}
	return f.ownedRet, nil
}

func (*fakeSubmitChecker) CheckValidSiswaInPesertaUjianById(context.Context, int, int) (bool, int, error) {
	return false, 0, nil
}

func (*fakeSubmitChecker) CheckTokenUjian(context.Context, string, int) (bool, error) {
	return false, nil
}

func (*fakeSubmitChecker) GetDeadlineUjian(context.Context, int) (time.Time, error) {
	return time.Time{}, nil
}

func runSubmitUjianCases(t *testing.T, prefix string) {
	t.Helper()

	ctx := context.Background()
	checkErr := errors.New("checker error")
	submitErr := errors.New("submit error")

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		idSiswa    int
		repo       *fakeSubmitAttemptRepo
		checker    *fakeSubmitChecker
		wantErr    error
		wantCheck  bool
		wantSubmit bool
	}{
		{
			name:       prefix + "1 -> id attempt tidak valid",
			idAttempt:  0,
			idSiswa:    9,
			repo:       &fakeSubmitAttemptRepo{},
			checker:    &fakeSubmitChecker{},
			wantErr:    coreerror.ErrMissingId,
			wantCheck:  false,
			wantSubmit: false,
		},
		{
			name:      prefix + "2 -> checker ownership gagal",
			idAttempt: 11,
			idSiswa:   9,
			repo:      &fakeSubmitAttemptRepo{},
			checker:   &fakeSubmitChecker{ownedErr: checkErr},
			wantErr:   checkErr,
			wantCheck: true,
		},
		{
			name:      prefix + "3 -> attempt bukan milik siswa",
			idAttempt: 11,
			idSiswa:   9,
			repo:      &fakeSubmitAttemptRepo{},
			checker:   &fakeSubmitChecker{ownedRet: false},
			wantErr:   coreerror.ErrNotFound,
			wantCheck: true,
		},
		{
			name:       prefix + "4 -> repo submit gagal",
			idAttempt:  11,
			idSiswa:    9,
			repo:       &fakeSubmitAttemptRepo{submitErr: submitErr},
			checker:    &fakeSubmitChecker{ownedRet: true},
			wantErr:    submitErr,
			wantCheck:  true,
			wantSubmit: true,
		},
		{
			name:       prefix + "5 -> berhasil submit ujian",
			idAttempt:  11,
			idSiswa:    9,
			repo:       &fakeSubmitAttemptRepo{},
			checker:    &fakeSubmitChecker{ownedRet: true},
			wantCheck:  true,
			wantSubmit: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewSubmitUjianService(tc.repo, tc.checker)
			err := svc.SubmitUjian(ctx, tc.idAttempt, tc.idSiswa)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCheck, tc.checker.checkCalled)
			if tc.wantCheck {
				assert.Equal(t, tc.idSiswa, tc.checker.gotSiswaID)
				assert.Equal(t, tc.idAttempt, tc.checker.gotAttemptID)
			}

			assert.Equal(t, tc.wantSubmit, tc.repo.submitCalled)
			if tc.wantSubmit {
				assert.Equal(t, tc.idAttempt, tc.repo.gotAttemptID)
			}
		})
	}
}

func TestSubmitUjianService_BasisPath(t *testing.T) {
	t.Parallel()
	runSubmitUjianCases(t, "Path ")
}
