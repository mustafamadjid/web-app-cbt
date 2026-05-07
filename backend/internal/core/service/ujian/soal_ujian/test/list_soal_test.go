package ujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
	"github.com/stretchr/testify/assert"
)

type fakeListSoalUjianRepo struct {
	listRet     []ujian.SoalUjianSiswa
	listErr     error
	listCalled  bool
	gotBankSoal ujian.ID
}

func (f *fakeListSoalUjianRepo) GetSoalUjianByBankSoal(_ context.Context, idBankSoal ujian.ID) ([]ujian.SoalUjianSiswa, error) {
	f.listCalled = true
	f.gotBankSoal = idBankSoal
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listRet, nil
}

func (*fakeListSoalUjianRepo) GetSoalUjianByBankSoalForSiswa(context.Context, ujian.ID) ([]ujian.SoalUjianSiswa, bool, error) {
	return nil, false, nil
}

func TestListSoalUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	soal := []ujian.SoalUjianSiswa{
		{IdSoal: 1, Pertanyaan: "Soal 1"},
		{IdSoal: 2, Pertanyaan: "Soal 2"},
	}

	tests := []struct {
		name       string
		idBankSoal ujian.ID
		acakSoal   bool
		repo       *fakeListSoalUjianRepo
		wantErr    error
		wantCalled bool
		assertGot  func(t *testing.T, got []ujian.SoalUjianSiswa)
	}{
		{
			name:       "Path 1 -> id bank soal tidak valid",
			idBankSoal: 0,
			repo:       &fakeListSoalUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Nil(t, got)
			},
		},
		{
			name:       "Path 2 -> repo get soal ujian gagal",
			idBankSoal: 5,
			repo:       &fakeListSoalUjianRepo{listErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Nil(t, got)
			},
		},
		{
			name:       "Path 3 -> berhasil list soal tanpa acak",
			idBankSoal: 5,
			repo:       &fakeListSoalUjianRepo{listRet: soal},
			wantCalled: true,
			assertGot: func(t *testing.T, got []ujian.SoalUjianSiswa) {
				assert.Equal(t, soal, got)
			},
		},
		{
			name:       "Path 4 -> berhasil list soal dengan acak",
			idBankSoal: 5,
			acakSoal:   true,
			repo:       &fakeListSoalUjianRepo{listRet: soal},
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

			svc := ujian_service.NewListSoalUjianService(tc.repo)
			got, err := svc.ListSoalUjian(ctx, tc.idBankSoal, tc.acakSoal)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.listCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.idBankSoal, tc.repo.gotBankSoal)
			}
			tc.assertGot(t, got)
		})
	}
}

func TestListSoalUjianService_LoopCoverage(t *testing.T) {
	t.Parallel()

	baseSoal := []ujian.SoalUjianSiswa{
		{IdSoal: 1, Pertanyaan: "Soal 1"},
		{IdSoal: 2, Pertanyaan: "Soal 2"},
		{IdSoal: 3, Pertanyaan: "Soal 3"},
	}

	tests := []struct {
		name       string
		repoItems  []ujian.SoalUjianSiswa
		acakSoal   bool
		assertGot  func(t *testing.T, expected, got []ujian.SoalUjianSiswa)
	}{
		{
			name:      "Loop 1 -> acak false tidak masuk loop",
			repoItems: append([]ujian.SoalUjianSiswa(nil), baseSoal...),
			acakSoal:  false,
			assertGot: func(t *testing.T, expected, got []ujian.SoalUjianSiswa) {
				assert.Equal(t, expected, got)
			},
		},
		{
			name:      "Loop 2 -> acak true dengan 0 soal tidak masuk loop",
			repoItems: []ujian.SoalUjianSiswa{},
			acakSoal:  true,
			assertGot: func(t *testing.T, expected, got []ujian.SoalUjianSiswa) {
				assert.Empty(t, got)
			},
		},
		{
			name:      "Loop 3 -> acak true dengan 1 soal tidak masuk loop",
			repoItems: []ujian.SoalUjianSiswa{{IdSoal: 1, Pertanyaan: "Soal 1"}},
			acakSoal:  true,
			assertGot: func(t *testing.T, expected, got []ujian.SoalUjianSiswa) {
				assert.Equal(t, expected, got)
			},
		},
		{
			name:      "Loop 4 -> acak true dengan 2 soal masuk loop 1 iterasi",
			repoItems: append([]ujian.SoalUjianSiswa(nil), baseSoal[:2]...),
			acakSoal:  true,
			assertGot: func(t *testing.T, expected, got []ujian.SoalUjianSiswa) {
				assert.Len(t, got, len(expected))
				assert.ElementsMatch(t, expected, got)
			},
		},
		{
			name:      "Loop 5 -> acak true dengan lebih dari 2 soal masuk loop banyak iterasi",
			repoItems: append([]ujian.SoalUjianSiswa(nil), baseSoal...),
			acakSoal:  true,
			assertGot: func(t *testing.T, expected, got []ujian.SoalUjianSiswa) {
				assert.Len(t, got, len(expected))
				assert.ElementsMatch(t, expected, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeListSoalUjianRepo{listRet: append([]ujian.SoalUjianSiswa(nil), tc.repoItems...)}
			svc := ujian_service.NewListSoalUjianService(repo)

			got, err := svc.ListSoalUjian(context.Background(), 9, tc.acakSoal)

			assert.NoError(t, err)
			tc.assertGot(t, tc.repoItems, got)
		})
	}
}
