package ujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUpdateUjianRepo struct {
	updateErr    error
	updateCalled bool
	gotID        ujian.ID
	gotPayload   updatepatch.UpdatePenjadwalanUjian
}

func (*fakeUpdateUjianRepo) GetIdUjianByAttempt(context.Context, ujian.ID) (ujian.ID, error) {
	return 0, nil
}

func (*fakeUpdateUjianRepo) CreateUjian(context.Context, ujian.PenjadwalanUjian) error {
	return nil
}

func (f *fakeUpdateUjianRepo) UpdateUjian(_ context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	f.updateCalled = true
	f.gotID = id
	f.gotPayload = payload
	return f.updateErr
}

func (*fakeUpdateUjianRepo) DeleteUjian(context.Context, ujian.ID) error {
	return nil
}

func idPtr(v ujian.ID) *ujian.ID {
	return &v
}

func stringPatchPtr(v string) *string {
	return &v
}

func statusPtr(v ujian.StatusUjian) *ujian.StatusUjian {
	return &v
}

func timePtr(v time.Time) *time.Time {
	return &v
}

func validUpdatePayload() updatepatch.UpdatePenjadwalanUjian {
	nama := "  UTS Final  "
	deskripsi := "  Deskripsi Baru  "
	token := "  xyz123  "
	status := ujian.StatusUjian(" selesai ")
	tanggal := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	mulai := time.Date(2026, 6, 15, 9, 0, 0, 0, time.UTC)
	selesai := time.Date(2026, 6, 15, 11, 0, 0, 0, time.UTC)

	return updatepatch.UpdatePenjadwalanUjian{
		Ujian: updatepatch.UpdateUjianPatch{
			NamaUjian:      &nama,
			DeskripsiUjian: &deskripsi,
		},
		JadwalUjian: updatepatch.UpdateJadwalUjianPatch{
			Token:        &token,
			StatusUjian:  &status,
			TanggalUjian: &tanggal,
			WaktuMulai:   &mulai,
			WaktuSelesai: &selesai,
		},
	}
}

func TestUpdateUjianService_BasisPathValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name        string
		id          ujian.ID
		payload     updatepatch.UpdatePenjadwalanUjian
		repo        *fakeUpdateUjianRepo
		wantErr     error
		wantCalled  bool
		assertPatch func(t *testing.T, repo *fakeUpdateUjianRepo)
	}{
		{
			name:       "Path 1 -> tidak ada field untuk diupdate",
			id:         10,
			payload:    updatepatch.UpdatePenjadwalanUjian{},
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrNoFieldToUpdate,
			wantCalled: false,
		},
		{
			name:       "Path 2 -> id ujian service tidak valid",
			id:         0,
			payload:    validUpdatePayload(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "Path 3 -> id patch ujian tidak valid",
			id:   10,
			payload: func() updatepatch.UpdatePenjadwalanUjian {
				payload := validUpdatePayload()
				payload.Ujian.IdBankSoal = idPtr(0)
				return payload
			}(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "Path 4 -> id patch jadwal tidak valid",
			id:   10,
			payload: func() updatepatch.UpdatePenjadwalanUjian {
				payload := validUpdatePayload()
				payload.JadwalUjian.IdSesi = idPtr(0)
				return payload
			}(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name: "Path 5 -> nama ujian kosong setelah trim",
			id:   10,
			payload: func() updatepatch.UpdatePenjadwalanUjian {
				payload := validUpdatePayload()
				payload.Ujian.NamaUjian = stringPatchPtr("   ")
				return payload
			}(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name: "Path 6 -> token kosong setelah trim",
			id:   10,
			payload: func() updatepatch.UpdatePenjadwalanUjian {
				payload := validUpdatePayload()
				payload.JadwalUjian.Token = stringPatchPtr("   ")
				return payload
			}(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCalled: false,
		},
		{
			name: "Path 7 -> status ujian tidak valid",
			id:   10,
			payload: func() updatepatch.UpdatePenjadwalanUjian {
				payload := validUpdatePayload()
				payload.JadwalUjian.StatusUjian = statusPtr(ujian.StatusUjian("unknown"))
				return payload
			}(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 8 -> waktu jadwal tidak valid",
			id:   10,
			payload: func() updatepatch.UpdatePenjadwalanUjian {
				payload := validUpdatePayload()
				mulai := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
				selesai := time.Date(2026, 6, 15, 11, 0, 0, 0, time.UTC)
				payload.JadwalUjian.WaktuMulai = timePtr(mulai)
				payload.JadwalUjian.WaktuSelesai = timePtr(selesai)
				return payload
			}(),
			repo:       &fakeUpdateUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name:       "Path 9 -> repo update ujian gagal",
			id:         10,
			payload:    validUpdatePayload(),
			repo:       &fakeUpdateUjianRepo{updateErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
			assertPatch: func(t *testing.T, repo *fakeUpdateUjianRepo) {
				t.Helper()
				require.NotNil(t, repo.gotPayload.Ujian.NamaUjian)
				assert.Equal(t, "UTS Final", *repo.gotPayload.Ujian.NamaUjian)
				require.NotNil(t, repo.gotPayload.Ujian.DeskripsiUjian)
				assert.Equal(t, "Deskripsi Baru", *repo.gotPayload.Ujian.DeskripsiUjian)
				require.NotNil(t, repo.gotPayload.JadwalUjian.Token)
				assert.Equal(t, "XYZ123", *repo.gotPayload.JadwalUjian.Token)
				require.NotNil(t, repo.gotPayload.JadwalUjian.StatusUjian)
				assert.Equal(t, ujian.SELESAI, *repo.gotPayload.JadwalUjian.StatusUjian)
			},
		},
		{
			name:       "Path 10 -> berhasil update ujian",
			id:         10,
			payload:    validUpdatePayload(),
			repo:       &fakeUpdateUjianRepo{},
			wantCalled: true,
			assertPatch: func(t *testing.T, repo *fakeUpdateUjianRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(10), repo.gotID)
				require.NotNil(t, repo.gotPayload.Ujian.NamaUjian)
				assert.Equal(t, "UTS Final", *repo.gotPayload.Ujian.NamaUjian)
				require.NotNil(t, repo.gotPayload.JadwalUjian.Token)
				assert.Equal(t, "XYZ123", *repo.gotPayload.JadwalUjian.Token)
				require.NotNil(t, repo.gotPayload.JadwalUjian.StatusUjian)
				assert.Equal(t, ujian.SELESAI, *repo.gotPayload.JadwalUjian.StatusUjian)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewUpdateUjianService(tc.repo)
			err := svc.UpdateUjianService(ctx, tc.id, tc.payload)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.updateCalled)
			if tc.assertPatch != nil {
				tc.assertPatch(t, tc.repo)
			}
		})
	}
}
