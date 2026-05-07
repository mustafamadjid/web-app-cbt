package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/active_attempt"
	"github.com/stretchr/testify/assert"
)

type fakeActiveAttemptRepo struct {
	getRet      ujian.AttemptUjian
	getErr      error
	getCalled   bool
	gotSiswaID  int
	gotJadwalID int
}

func (f *fakeActiveAttemptRepo) ListUjianSiswa(context.Context, int, query.ListUjianFilter) ([]ujian.ListUjian, error) {
	return nil, nil
}

func (f *fakeActiveAttemptRepo) GetWaktuSelesaiUjian(context.Context, int) (time.Time, error) {
	return time.Time{}, nil
}

func (f *fakeActiveAttemptRepo) GetActiveUjianAttemptBySiswa(_ context.Context, idSiswa int, idJadwalUjian int) (ujian.AttemptUjian, error) {
	f.getCalled = true
	f.gotSiswaID = idSiswa
	f.gotJadwalID = idJadwalUjian
	if f.getErr != nil {
		return ujian.AttemptUjian{}, f.getErr
	}
	return f.getRet, nil
}

func runGetActiveAttemptCases(t *testing.T, prefix string) {
	t.Helper()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := ujian.AttemptUjian{IdAttempt: 31, IdPesertaUjian: 17, StatusAttempt: ujian.ATTEMPT_IN_PROGRESS}

	tests := []struct {
		name        string
		idSiswa     int
		idJadwal    int
		repo        *fakeActiveAttemptRepo
		wantErr     error
		wantCalled  bool
		wantAttempt ujian.AttemptUjian
	}{
		{
			name:       prefix + "1 -> id jadwal ujian tidak valid",
			idSiswa:    9,
			idJadwal:   0,
			repo:       &fakeActiveAttemptRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       prefix + "2 -> id siswa tidak valid",
			idSiswa:    0,
			idJadwal:   21,
			repo:       &fakeActiveAttemptRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       prefix + "3 -> repo get active attempt gagal",
			idSiswa:    9,
			idJadwal:   21,
			repo:       &fakeActiveAttemptRepo{getErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:        prefix + "4 -> berhasil get active attempt",
			idSiswa:     9,
			idJadwal:    21,
			repo:        &fakeActiveAttemptRepo{getRet: expected},
			wantCalled:  true,
			wantAttempt: expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewGetActiveAttemptUjianService(tc.repo)
			got, err := svc.GetActiveAttemptUjian(ctx, tc.idSiswa, tc.idJadwal)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.getCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.idSiswa, tc.repo.gotSiswaID)
				assert.Equal(t, tc.idJadwal, tc.repo.gotJadwalID)
			}
			assert.Equal(t, tc.wantAttempt, got)
		})
	}
}

func TestGetActiveAttemptUjianService_BasisPath(t *testing.T) {
	t.Parallel()
	runGetActiveAttemptCases(t, "Path ")
}
