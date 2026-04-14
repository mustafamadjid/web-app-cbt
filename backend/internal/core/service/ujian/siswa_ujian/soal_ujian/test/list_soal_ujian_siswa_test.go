package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/soal_ujian"
	"github.com/stretchr/testify/assert"
)

type fakeSoalUjianSiswaRepo struct {
	listRet     []ujian.SoalUjianSiswa
	acakRet     bool
	listErr     error
	listCalled  bool
	gotJadwalID ujian.ID
}

func (*fakeSoalUjianSiswaRepo) GetSoalUjianByBankSoal(context.Context, ujian.ID) ([]ujian.SoalUjianSiswa, error) {
	return nil, nil
}

func (f *fakeSoalUjianSiswaRepo) GetSoalUjianByBankSoalForSiswa(_ context.Context, idJadwalUjian ujian.ID) ([]ujian.SoalUjianSiswa, bool, error) {
	f.listCalled = true
	f.gotJadwalID = idJadwalUjian
	if f.listErr != nil {
		return nil, false, f.listErr
	}
	return f.listRet, f.acakRet, nil
}

func TestListSoalUjianSiswaService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	soal := []ujian.SoalUjianSiswa{
		{IdSoal: 1, Pertanyaan: "Soal 1"},
		{IdSoal: 2, Pertanyaan: "Soal 2"},
	}

	tests := []struct {
		name       string
		idJadwal   ujian.ID
		repo       *fakeSoalUjianSiswaRepo
		wantErr    error
		wantCalled bool
		assertGot  func(t *testing.T, got []ujian.SoalUjianSiswa)
	}{
		{
			name:       "branch 1 -> id jadwal ujian tidak valid",
			idJadwal:   0,
			repo:       &fakeSoalUjianSiswaRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Nil(t, got)
			},
		},
		{
			name:       "branch 2 -> repo get soal siswa gagal",
			idJadwal:   9,
			repo:       &fakeSoalUjianSiswaRepo{listErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Nil(t, got)
			},
		},
		{
			name:       "branch 3 -> berhasil list soal tanpa acak",
			idJadwal:   9,
			repo:       &fakeSoalUjianSiswaRepo{listRet: soal},
			wantCalled: true,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Equal(t, soal, got)
			},
		},
		{
			name:       "branch 4 -> berhasil list soal dengan acak",
			idJadwal:   9,
			repo:       &fakeSoalUjianSiswaRepo{listRet: soal, acakRet: true},
			wantCalled: true,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Len(t, got, 2)
				assert.ElementsMatch(t, soal, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := siswaujian_service.NewListSoalUjianSiswaService(tc.repo)
			got, err := svc.ListSoalUjianSiswa(ctx, tc.idJadwal)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.listCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.idJadwal, tc.repo.gotJadwalID)
			}
			tc.assertGot(t, got)
		})
	}
}
