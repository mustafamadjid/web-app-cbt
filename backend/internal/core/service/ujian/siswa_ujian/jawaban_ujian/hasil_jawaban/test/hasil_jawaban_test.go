package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/hasil_jawaban"
	"github.com/stretchr/testify/assert"
)

type fakeHasilJawabanRepo struct {
	listRet      []ujian.HasilJawabanUjian
	listErr      error
	listCalled   bool
	gotAttemptID ujian.ID
}

func (*fakeHasilJawabanRepo) GetJawabanUjianByAttemptId(context.Context, ujian.ID) ([]ujian.JawabanUjian, error) {
	return nil, nil
}

func (*fakeHasilJawabanRepo) SaveJawabanUjian(context.Context, ujian.ID, []ujian.JawabanUjian) error {
	return nil
}

func (f *fakeHasilJawabanRepo) ListHasilJawabanUjianByIdAttempt(_ context.Context, idAttempt ujian.ID) ([]ujian.HasilJawabanUjian, error) {
	f.listCalled = true
	f.gotAttemptID = idAttempt
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listRet, nil
}

func TestHasilJawabanUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	score := 95.0
	expected := []ujian.HasilJawabanUjian{
		{NilaiAkhir: &score},
	}

	tests := []struct {
		name       string
		idAttempt  int
		repo       *fakeHasilJawabanRepo
		wantErr    error
		wantCalled bool
		wantItems  []ujian.HasilJawabanUjian
	}{
		{
			name:       "Path 1 -> id attempt tidak valid",
			idAttempt:  0,
			repo:       &fakeHasilJawabanRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantItems:  []ujian.HasilJawabanUjian{},
			wantCalled: false,
		},
		{
			name:       "Path 2 -> repo list hasil jawaban gagal",
			idAttempt:  9,
			repo:       &fakeHasilJawabanRepo{listErr: repoErr},
			wantErr:    repoErr,
			wantItems:  []ujian.HasilJawabanUjian{},
			wantCalled: true,
		},
		{
			name:       "Path 3 -> berhasil list hasil jawaban",
			idAttempt:  9,
			repo:       &fakeHasilJawabanRepo{listRet: expected},
			wantItems:  expected,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := siswaujian_service.NewHasilJawabanUjianService(tc.repo)
			got, err := svc.ListHasilJawabanUjianByAttempt(ctx, tc.idAttempt)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.listCalled)
			if tc.wantCalled {
				assert.Equal(t, ujian.ID(tc.idAttempt), tc.repo.gotAttemptID)
			}
			assert.Equal(t, tc.wantItems, got)
		})
	}
}
