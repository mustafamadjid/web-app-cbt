package matapelajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestCreateMapelService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	existErr := errors.New("exist kode mapel error")
	createErr := errors.New("create mapel error")

	tests := []struct {
		name       string
		repo       *fake_test.FakeMapelRepo
		input      matapelajaran.MataPelajaran
		wantErr    error
		wantCreate bool
	}{
		{
			name: "Path 1 -> IdKelas kosong (0)",
			repo: &fake_test.FakeMapelRepo{},
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
			name: "Path 2 -> gagal cek ExistKodeMapel",
			repo: &fake_test.FakeMapelRepo{ExistKodeMapelErr: existErr},
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
			name: "Path 3 -> kode mapel sudah exist",
			repo: &fake_test.FakeMapelRepo{ExistKodeMapelRet: true},
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
			name: "Path 4 -> gagal CreateMapel",
			repo: &fake_test.FakeMapelRepo{CreateMapelErr: createErr},
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
			name: "Path 5 -> berhasil create mapel",
			repo: &fake_test.FakeMapelRepo{},
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
