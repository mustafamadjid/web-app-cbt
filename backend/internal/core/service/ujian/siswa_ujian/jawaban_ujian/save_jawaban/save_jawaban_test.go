package siswaujian_service

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
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

func toUjianIDPointer(v ujian.ID) *ujian.ID {
	return &v
}

func toStringPointer(v string) *string {
	return &v
}

func TestSaveJawabanUjianService(t *testing.T) {
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
			name:      "invalid attempt id",
			idAttempt: 0,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10, IdPilihan: toUjianIDPointer(20)},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrMissingId,
			wantSave: false,
		},
		{
			name:      "missing soal id",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 0, IdPilihan: toUjianIDPointer(20)},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrMissingId,
			wantSave: false,
		},
		{
			name:      "missing pilihan and essay",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrMissingJawabanEssayAndPilgan,
			wantSave: false,
		},
		{
			name:      "pilihan and essay both filled",
			idAttempt: 7,
			jawaban: []ujian.JawabanUjian{
				{IdSoal: 10, IdPilihan: toUjianIDPointer(20), JawabanEssay: toStringPointer("essay")},
			},
			repo:     &fakeJawabanUjianRepo{},
			wantErr:  coreerror.ErrInvalidInput,
			wantSave: false,
		},
		{
			name:      "repo error",
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
			name:      "success without id_jawaban",
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
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewJawabanUjianService(tc.repo)
			err := svc.SaveJawabanUjian(ctx, tc.idAttempt, tc.jawaban)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantSave, tc.repo.saveCalled)

			if tc.assertData != nil {
				tc.assertData(t, tc.repo)
			}
		})
	}
}
