package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list_soal_ujian_selesai"
	"github.com/stretchr/testify/assert"
)

type fakeListUjianSelesaiSiswaRepo struct {
	submittedRet    []ujian.ListUjian
	submittedErr    error
	submittedCalled bool
	gotSiswaID      int
}

func (*fakeListUjianSelesaiSiswaRepo) GetAllUjian(context.Context, query.ListUjianFilter) ([]ujian.ListUjian, error) {
	return nil, nil
}

func (*fakeListUjianSelesaiSiswaRepo) GetUjianById(context.Context, ujian.ID) (ujian.ListUjian, error) {
	return ujian.ListUjian{}, nil
}

func (f *fakeListUjianSelesaiSiswaRepo) GetAllUjianSubmittedByIdSiswa(_ context.Context, idSiswa int) ([]ujian.ListUjian, error) {
	f.submittedCalled = true
	f.gotSiswaID = idSiswa
	if f.submittedErr != nil {
		return nil, f.submittedErr
	}
	return f.submittedRet, nil
}

func TestListUjianSelesaiSiswaService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := []ujian.ListUjian{
		{IdUjian: 1, IdJadwalUjian: 11, IdAttempt: 21, NamaUjian: "UTS"},
		{IdUjian: 2, IdJadwalUjian: 12, IdAttempt: 22, NamaUjian: "UAS"},
	}

	tests := []struct {
		name       string
		idSiswa    int
		repo       *fakeListUjianSelesaiSiswaRepo
		wantErr    error
		wantCalled bool
		wantItems  []ujian.ListUjian
	}{
		{
			name:       "Branch 1 -> id siswa tidak valid",
			idSiswa:    0,
			repo:       &fakeListUjianSelesaiSiswaRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Branch 2 -> repo list ujian selesai siswa gagal",
			idSiswa:    7,
			repo:       &fakeListUjianSelesaiSiswaRepo{submittedErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "Branch 3 -> berhasil list ujian selesai siswa",
			idSiswa:    7,
			repo:       &fakeListUjianSelesaiSiswaRepo{submittedRet: expected},
			wantCalled: true,
			wantItems:  expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := siswaujian_service.NewListUjianSelesaiSiswaService(tc.repo)
			got, err := svc.ListUjianSelesaiSiswa(ctx, tc.idSiswa)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.submittedCalled)
			assert.Equal(t, tc.wantItems, got)
			if tc.wantCalled {
				assert.Equal(t, tc.idSiswa, tc.repo.gotSiswaID)
			}
		})
	}
}
