package ujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCreateUjianRepo struct {
	createErr    error
	createCalled bool
	gotCreate    ujian.PenjadwalanUjian
}

func (*fakeCreateUjianRepo) GetIdUjianByAttempt(context.Context, ujian.ID) (ujian.ID, error) {
	return 0, nil
}

func (f *fakeCreateUjianRepo) CreateUjian(_ context.Context, data ujian.PenjadwalanUjian) error {
	f.createCalled = true
	f.gotCreate = data
	return f.createErr
}

func (*fakeCreateUjianRepo) UpdateUjian(context.Context, ujian.ID, updatepatch.UpdatePenjadwalanUjian) error {
	return nil
}

func (*fakeCreateUjianRepo) DeleteUjian(context.Context, ujian.ID) error {
	return nil
}

func baseCreateUjianData() ujian.PenjadwalanUjian {
	desc := "  Deskripsi Ujian  "
	idNamaKelas := ujian.ID(15)
	tanggal := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	mulai := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)
	selesai := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)

	return ujian.PenjadwalanUjian{
		Ujian: ujian.Ujian{
			IdBankSoal:     3,
			IdKelas:        4,
			IdNamaKelas:    &idNamaKelas,
			IdGuru:         5,
			NamaUjian:      "  UTS Matematika  ",
			DeskripsiUjian: &desc,
		},
		JadwalUjian: ujian.JadwalUjian{
			IdSesi:       6,
			IdRuangan:    7,
			IdPengawas:   8,
			TanggalUjian: tanggal,
			WaktuMulai:   mulai,
			WaktuSelesai: selesai,
			Token:        "  abc123  ",
			StatusUjian:  ujian.StatusUjian(" mulai "),
		},
	}
}

func TestCreateUjianService_BasisPathValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		data       ujian.PenjadwalanUjian
		repo       *fakeCreateUjianRepo
		wantErr    error
		wantCalled bool
		assertData func(t *testing.T, repo *fakeCreateUjianRepo)
	}{
		{
			name: "Path 1 -> id utama ujian tidak valid",
			data: func() ujian.PenjadwalanUjian {
				item := baseCreateUjianData()
				item.Ujian.IdBankSoal = 0
				return item
			}(),
			repo:       &fakeCreateUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "Path 2 -> id nama kelas opsional tidak valid",
			data: func() ujian.PenjadwalanUjian {
				item := baseCreateUjianData()
				invalid := ujian.ID(0)
				item.Ujian.IdNamaKelas = &invalid
				return item
			}(),
			repo:       &fakeCreateUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "Path 3 -> nama ujian kosong setelah trim",
			data: func() ujian.PenjadwalanUjian {
				item := baseCreateUjianData()
				item.Ujian.NamaUjian = "   "
				return item
			}(),
			repo:       &fakeCreateUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name: "Path 4 -> waktu mulai tidak sebelum waktu selesai",
			data: func() ujian.PenjadwalanUjian {
				item := baseCreateUjianData()
				item.JadwalUjian.WaktuMulai = item.JadwalUjian.WaktuSelesai
				return item
			}(),
			repo:       &fakeCreateUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 5 -> token kosong setelah trim",
			data: func() ujian.PenjadwalanUjian {
				item := baseCreateUjianData()
				item.JadwalUjian.Token = "   "
				return item
			}(),
			repo:       &fakeCreateUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name: "Path 6 -> status ujian tidak valid",
			data: func() ujian.PenjadwalanUjian {
				item := baseCreateUjianData()
				item.JadwalUjian.StatusUjian = ujian.StatusUjian("selain-valid")
				return item
			}(),
			repo:       &fakeCreateUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 7 -> repo create ujian gagal",
			data: baseCreateUjianData(),
			repo: &fakeCreateUjianRepo{
				createErr: repoErr,
			},
			wantErr:    repoErr,
			wantCalled: true,
			assertData: func(t *testing.T, repo *fakeCreateUjianRepo) {
				t.Helper()
				assert.Equal(t, "UTS Matematika", repo.gotCreate.Ujian.NamaUjian)
				require.NotNil(t, repo.gotCreate.Ujian.DeskripsiUjian)
				assert.Equal(t, "Deskripsi Ujian", *repo.gotCreate.Ujian.DeskripsiUjian)
				assert.Equal(t, "ABC123", repo.gotCreate.JadwalUjian.Token)
				assert.Equal(t, ujian.MULAI, repo.gotCreate.JadwalUjian.StatusUjian)
			},
		},
		{
			name:       "Path 8 -> berhasil create ujian",
			data:       baseCreateUjianData(),
			repo:       &fakeCreateUjianRepo{},
			wantCalled: true,
			assertData: func(t *testing.T, repo *fakeCreateUjianRepo) {
				t.Helper()
				assert.Equal(t, "UTS Matematika", repo.gotCreate.Ujian.NamaUjian)
				require.NotNil(t, repo.gotCreate.Ujian.DeskripsiUjian)
				assert.Equal(t, "Deskripsi Ujian", *repo.gotCreate.Ujian.DeskripsiUjian)
				assert.Equal(t, "ABC123", repo.gotCreate.JadwalUjian.Token)
				assert.Equal(t, ujian.MULAI, repo.gotCreate.JadwalUjian.StatusUjian)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewCreateUjianService(tc.repo)
			err := svc.CreateUjianService(ctx, tc.data)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.createCalled)
			if tc.assertData != nil {
				tc.assertData(t, tc.repo)
			}
		})
	}
}
