package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/waktu_selesai"
	"github.com/stretchr/testify/assert"
)

type fakeWaktuSelesaiRepo struct {
	getRet      time.Time
	getErr      error
	getCalled   bool
	gotJadwalID int
}

func (*fakeWaktuSelesaiRepo) ListUjianSiswa(context.Context, int, query.ListUjianFilter) ([]ujian.ListUjian, error) {
	return nil, nil
}

func (f *fakeWaktuSelesaiRepo) GetWaktuSelesaiUjian(_ context.Context, idJadwalUjian int) (time.Time, error) {
	f.getCalled = true
	f.gotJadwalID = idJadwalUjian
	if f.getErr != nil {
		return time.Time{}, f.getErr
	}
	return f.getRet, nil
}

func (*fakeWaktuSelesaiRepo) GetActiveUjianAttemptBySiswa(context.Context, int, int) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func TestGetWaktuSelesaiService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		idJadwal   int
		repo       *fakeWaktuSelesaiRepo
		wantErr    error
		wantCalled bool
		wantTime   time.Time
	}{
		{
			name:       "Path 1 -> id jadwal ujian tidak valid",
			idJadwal:   0,
			repo:       &fakeWaktuSelesaiRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Path 2 -> repo get waktu selesai gagal",
			idJadwal:   17,
			repo:       &fakeWaktuSelesaiRepo{getErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "Path 3 -> berhasil get waktu selesai",
			idJadwal:   17,
			repo:       &fakeWaktuSelesaiRepo{getRet: expected},
			wantCalled: true,
			wantTime:   expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := siswaujian_service.NewGetWaktuSelesaiService(tc.repo)
			got, err := svc.GetWaktuSelesai(ctx, tc.idJadwal)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.getCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.idJadwal, tc.repo.gotJadwalID)
			}
			assert.Equal(t, tc.wantTime, got)
		})
	}
}
