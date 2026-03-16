package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	getjawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/get_jawaban"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeJawabanOwnershipChecker struct {
	checkFn      func(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)
	checkCalled  bool
	gotSiswaID   int
	gotAttemptID ujian.ID
}

func (f *fakeJawabanOwnershipChecker) CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	f.checkCalled = true
	f.gotSiswaID = idSiswa
	f.gotAttemptID = idAttempt
	if f.checkFn != nil {
		return f.checkFn(ctx, idSiswa, idAttempt)
	}
	return false, nil
}

type fakeJawabanRepo struct {
	getFn         func(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error)
	getCalled     bool
	gotAttemptID  ujian.ID
	gotSaveCalled bool
}

func (f *fakeJawabanRepo) GetJawabanUjianByAttemptId(ctx context.Context, idAttempt ujian.ID) ([]ujian.JawabanUjian, error) {
	f.getCalled = true
	f.gotAttemptID = idAttempt
	if f.getFn != nil {
		return f.getFn(ctx, idAttempt)
	}
	return nil, nil
}

func (f *fakeJawabanRepo) SaveJawabanUjian(context.Context, ujian.ID, []ujian.JawabanUjian) error {
	f.gotSaveCalled = true
	return nil
}

func TestSiswaGetJawabanUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	waktuJawab := time.Date(2026, time.March, 16, 11, 45, 0, 0, time.UTC)

	tests := []struct {
		name        string
		idSiswa     int
		idAttempt   ujian.ID
		checker     *fakeJawabanOwnershipChecker
		repo        *fakeJawabanRepo
		wantErr     error
		wantCheck   bool
		wantGet     bool
		assertItems func(t *testing.T, items []ujian.JawabanUjian)
	}{
		{
			name:      "invalid siswa id",
			idAttempt: 7,
			checker:   &fakeJawabanOwnershipChecker{},
			repo:      &fakeJawabanRepo{},
			wantErr:   coreerror.ErrMissingId,
		},
		{
			name:    "invalid attempt id",
			idSiswa: 9,
			checker: &fakeJawabanOwnershipChecker{},
			repo:    &fakeJawabanRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:      "ownership checker error",
			idSiswa:   9,
			idAttempt: 7,
			checker: &fakeJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return false, repoErr
				},
			},
			repo:      &fakeJawabanRepo{},
			wantErr:   repoErr,
			wantCheck: true,
		},
		{
			name:      "attempt not owned by siswa",
			idSiswa:   9,
			idAttempt: 7,
			checker: &fakeJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return false, nil
				},
			},
			repo:      &fakeJawabanRepo{},
			wantErr:   coreerror.ErrNotFound,
			wantCheck: true,
		},
		{
			name:      "getter error bubbles up",
			idSiswa:   9,
			idAttempt: 7,
			checker: &fakeJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo: &fakeJawabanRepo{
				getFn: func(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
					return nil, repoErr
				},
			},
			wantErr:   repoErr,
			wantCheck: true,
			wantGet:   true,
		},
		{
			name:      "success",
			idSiswa:   9,
			idAttempt: 7,
			checker: &fakeJawabanOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo: &fakeJawabanRepo{
				getFn: func(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
					essay := "jawaban siswa"
					idPilihan := ujian.ID(12)
					return []ujian.JawabanUjian{
						{
							IdJawaban:    3,
							IdSoal:       11,
							IdPilihan:    &idPilihan,
							JawabanEssay: &essay,
							WaktuJawab:   &waktuJawab,
						},
					}, nil
				},
			},
			wantCheck: true,
			wantGet:   true,
			assertItems: func(t *testing.T, items []ujian.JawabanUjian) {
				t.Helper()
				require.Len(t, items, 1)
				assert.Equal(t, ujian.ID(3), items[0].IdJawaban)
				assert.Equal(t, ujian.ID(11), items[0].IdSoal)
				require.NotNil(t, items[0].IdPilihan)
				assert.Equal(t, ujian.ID(12), *items[0].IdPilihan)
				require.NotNil(t, items[0].JawabanEssay)
				assert.Equal(t, "jawaban siswa", *items[0].JawabanEssay)
				require.NotNil(t, items[0].WaktuJawab)
				assert.Equal(t, waktuJawab, *items[0].WaktuJawab)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			getter := getjawaban_service.NewGetJawabanUjianService(tc.repo)
			svc := getjawaban_service.NewSiswaGetJawabanUjianService(tc.checker, getter)
			items, err := svc.GetJawabanUjianByAttemptId(ctx, tc.idSiswa, tc.idAttempt)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCheck, tc.checker.checkCalled)
			assert.Equal(t, tc.wantGet, tc.repo.getCalled)
			assert.False(t, tc.repo.gotSaveCalled)

			if tc.wantCheck {
				assert.Equal(t, tc.idSiswa, tc.checker.gotSiswaID)
				assert.Equal(t, tc.idAttempt, tc.checker.gotAttemptID)
			}
			if tc.wantGet {
				assert.Equal(t, tc.idAttempt, tc.repo.gotAttemptID)
			}
			if tc.assertItems != nil {
				tc.assertItems(t, items)
			}
		})
	}
}
