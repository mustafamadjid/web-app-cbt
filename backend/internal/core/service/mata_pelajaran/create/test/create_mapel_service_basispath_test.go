package matapelajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
	"github.com/stretchr/testify/assert"
)

func TestCreateMapelService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	existErr := errors.New("exist kode mapel error")
	createErr := errors.New("create mapel error")

	tests := []struct {
		name       string
		repo       *FakeMapelRepo
		input      matapelajaran.MataPelajaran
		wantErr    error
		wantCreate bool
	}{
		{
			name: "Branch 1 -> id kelas kosong",
			repo: &FakeMapelRepo{},
			input: matapelajaran.MataPelajaran{
				KodeMapel: "MTK01",
				NamaMapel: "Matematika",
				Deskripsi: "Deskripsi",
				IdKelas:   0,
			},
			wantErr:    coreerror.ErrInvalidInput,
			wantCreate: false,
		},
		{
			name: "Branch 2 -> cek kode mapel gagal",
			repo: &FakeMapelRepo{ExistKodeMapelErr: existErr},
			input: matapelajaran.MataPelajaran{
				KodeMapel: "MTK01",
				NamaMapel: "Matematika",
				Deskripsi: "Deskripsi",
				IdKelas:   1,
			},
			wantErr:    existErr,
			wantCreate: false,
		},
		{
			name: "Branch 3 -> kode mapel sudah ada",
			repo: &FakeMapelRepo{ExistKodeMapelRet: true},
			input: matapelajaran.MataPelajaran{
				KodeMapel: "MTK01",
				NamaMapel: "Matematika",
				Deskripsi: "Deskripsi",
				IdKelas:   1,
			},
			wantErr:    coreerror.ErrKodeMapelExist,
			wantCreate: false,
		},
		{
			name: "Branch 4 -> create mapel gagal",
			repo: &FakeMapelRepo{CreateMapelErr: createErr},
			input: matapelajaran.MataPelajaran{
				KodeMapel: "MTK01",
				NamaMapel: "Matematika",
				Deskripsi: "Deskripsi",
				IdKelas:   1,
			},
			wantErr:    createErr,
			wantCreate: true,
		},
		{
			name: "Branch 5 -> create mapel berhasil",
			repo: &FakeMapelRepo{},
			input: matapelajaran.MataPelajaran{
				KodeMapel: "  mtk01  ",
				NamaMapel: "  Matematika  ",
				Deskripsi: "  Deskripsi  ",
				IdKelas:   1,
			},
			wantErr:    nil,
			wantCreate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mapel_service.NewMapelService(tt.repo)
			err := svc.CreateMapelService(ctx, tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCreate, tt.repo.CreateMapelCalled)
		})
	}
}
