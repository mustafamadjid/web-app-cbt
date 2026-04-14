package statisktikujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	statistikrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/statistik_ujian"
	statisktikujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/statistik_ujian"
	"github.com/stretchr/testify/assert"
)

type fakeStatistikUjianRepo struct {
	getRet      ujian.StatistikUjian
	getErr      error
	getCalled   bool
	gotJadwalID ujian.ID
}

func (f *fakeStatistikUjianRepo) GetStatistikUjianByIdJadwal(_ context.Context, idJadwalUjian ujian.ID) (ujian.StatistikUjian, error) {
	f.getCalled = true
	f.gotJadwalID = idJadwalUjian
	if f.getErr != nil {
		return ujian.StatistikUjian{}, f.getErr
	}
	return f.getRet, nil
}

var _ statistikrepo.StatistikUjianRepository = (*fakeStatistikUjianRepo)(nil)

func TestStatistikUjianService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := ujian.StatistikUjian{
		IDStatistikUjian:  3,
		IDJadwalUjian:     15,
		NilaiTertinggi:    98,
		NilaiTerendah:     60,
		NilaiRataRata:     80,
		TotalPesertaUjian: 25,
	}

	tests := []struct {
		name       string
		idJadwal   int
		repo       *fakeStatistikUjianRepo
		wantErr    error
		wantCalled bool
		wantItem   ujian.StatistikUjian
	}{
		{
			name:       "branch 1 -> id jadwal tidak valid",
			idJadwal:   0,
			repo:       &fakeStatistikUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "branch 2 -> repo get statistik ujian gagal",
			idJadwal:   15,
			repo:       &fakeStatistikUjianRepo{getErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "branch 3 -> berhasil get statistik ujian",
			idJadwal:   15,
			repo:       &fakeStatistikUjianRepo{getRet: expected},
			wantCalled: true,
			wantItem:   expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := statisktikujian_service.NewStatistikUjianService(tc.repo)
			got, err := svc.GetStatistikUjian(ctx, tc.idJadwal)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.getCalled)
			if tc.wantCalled {
				assert.Equal(t, ujian.ID(tc.idJadwal), tc.repo.gotJadwalID)
			}
			assert.Equal(t, tc.wantItem, got)
		})
	}
}
