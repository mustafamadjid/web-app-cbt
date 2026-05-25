package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
	"github.com/stretchr/testify/assert"
)

type fakeJawabanUjianRepo struct {
	saveFn     func(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error
	saveCalled bool
	gotAttempt ujian.ID
	gotJawaban []ujian.JawabanUjian
}

func (f *fakeJawabanUjianRepo) GetJawabanUjianByAttemptId(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
	return nil, nil
}

func (f *fakeJawabanUjianRepo) SaveJawabanUjian(ctx context.Context, idAttempt ujian.ID, jawaban []ujian.JawabanUjian) error {
	f.saveCalled = true
	f.gotAttempt = idAttempt
	f.gotJawaban = jawaban
	if f.saveFn != nil {
		return f.saveFn(ctx, idAttempt, jawaban)
	}
	return nil
}

func (*fakeJawabanUjianRepo) ListHasilJawabanUjianByIdAttempt(context.Context, ujian.ID) ([]ujian.HasilJawabanUjian, error) {
	return nil, nil
}

func toUjianIDPointer(v ujian.ID) *ujian.ID {
	return &v
}

func toStringPointer(v string) *string {
	return &v
}

func TestSaveJawabanUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Date(2026, time.March, 15, 10, 0, 0, 0, time.UTC)
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		jawaban    []ujian.JawabanUjian
		repo       *fakeJawabanUjianRepo
		wantErr    error
		wantSave   bool
		assertData func(t *testing.T, repo *fakeJawabanUjianRepo)
	}{
		{
			name:      "path 1 -> id attempt tidak valid",
			idAttempt: 0,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10, IdPilihan: toUjianIDPointer(20)},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrMissingId,
			wantSave: false,
		},
		{
			name:      "path 2 -> id soal tidak valid",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 0, IdPilihan: toUjianIDPointer(20)},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrMissingId,
			wantSave: false,
		},
		{
			name:      "path 3 -> id pilihan tidak valid",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10, IdPilihan: toUjianIDPointer(0)},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrMissingId,
			wantSave: false,
		},
		{
			name:      "path 4 -> pilihan dan essay sama-sama terisi",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10, IdPilihan: toUjianIDPointer(20), JawabanEssay: toStringPointer("essay")},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrInvalidInput,
			wantSave: false,
		},
		{
			name:      "path 5 -> jawaban dikosongkan tetap boleh disimpan",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantSave: true,
			assertData: func(t *testing.T, repo *fakeJawabanUjianRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(7), repo.gotAttempt)
				assert.Len(t, repo.gotJawaban, 1)
				assert.Equal(t, ujian.ID(10), repo.gotJawaban[0].IdSoal)
				assert.Nil(t, repo.gotJawaban[0].IdPilihan)
				assert.Nil(t, repo.gotJawaban[0].JawabanEssay)
			},
		},
		{
			name:      "path 6 -> repo gagal menyimpan jawaban",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10, IdPilihan: toUjianIDPointer(20)},
			},
			repo: &fakeJawabanUjianRepo{
				saveFn: func(context.Context, ujian.ID, []ujian.JawabanUjian) error {
					return repoErr
				},
			},
			wantErr:  repoErr,
			wantSave: true,
		},
		{
			name:      "path 7 -> berhasil menyimpan jawaban pilihan",
			idAttempt: 9,
			jawaban: []ujian.JawabanUjian{
				{
					IdSoal:     11,
					IdPilihan:  toUjianIDPointer(22),
					WaktuJawab: &now,
				},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantSave: true,
			assertData: func(t *testing.T, repo *fakeJawabanUjianRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(9), repo.gotAttempt)
				assert.Len(t, repo.gotJawaban, 1)
				assert.Equal(t, ujian.ID(0), repo.gotJawaban[0].IdJawaban)
				assert.Equal(t, ujian.ID(11), repo.gotJawaban[0].IdSoal)
				assert.Equal(t, ujian.ID(22), *repo.gotJawaban[0].IdPilihan)
				assert.Equal(t, now, *repo.gotJawaban[0].WaktuJawab)
			},
		},
		{
			name:      "path 8 -> berhasil menyimpan jawaban essay",
			idAttempt: 9,
			jawaban: []ujian.JawabanUjian{
				{
					IdSoal:       12,
					JawabanEssay: toStringPointer("jawaban essay"),
					WaktuJawab:   &now,
				},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantSave: true,
			assertData: func(t *testing.T, repo *fakeJawabanUjianRepo) {
				t.Helper()
				assert.Equal(t, ujian.ID(9), repo.gotAttempt)
				assert.Len(t, repo.gotJawaban, 1)
				assert.Equal(t, ujian.ID(12), repo.gotJawaban[0].IdSoal)
				assert.Nil(t, repo.gotJawaban[0].IdPilihan)
				assert.Equal(t, "jawaban essay", *repo.gotJawaban[0].JawabanEssay)
				assert.Equal(t, now, *repo.gotJawaban[0].WaktuJawab)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := siswaujian_service.NewJawabanUjianService(tc.repo)
			err := svc.SaveJawabanUjian(ctx, tc.idAttempt, tc.jawaban)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantSave, tc.repo.saveCalled)

			if tc.assertData != nil {
				tc.assertData(t, tc.repo)
			}
		})
	}
}
