package kelas_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
)

func TestUpdateKelasService_UpdateNamaKelas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	idTingkatKelas := kelas.ID(3)
	repoErr := errors.New("repo error")

	tests := []struct {
		name            string
		idNamaKelas     int
		patch           updatepatch.NamaKelasPatch
		repo            *FakeKelasRepo
		wantErr         error
		wantRepoCalled  bool
		wantPatchedName *string
	}{
		{
			name:           "Branch 1 -> id nama kelas kosong",
			idNamaKelas:    0,
			patch:          updatepatch.NamaKelasPatch{},
			repo:           &FakeKelasRepo{},
			wantErr:        coreerror.ErrMissingId,
			wantRepoCalled: false,
		},
		{
			name:        "Branch 2 -> nama kelas setelah trim kosong",
			idNamaKelas: 10,
			patch: updatepatch.NamaKelasPatch{
				NamaKelas: ptrString("   "),
			},
			repo:           &FakeKelasRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:        "Branch 3 -> id tingkat kelas kosong",
			idNamaKelas: 10,
			patch: updatepatch.NamaKelasPatch{
				IdTingkatKelas: ptrID(0),
			},
			repo:           &FakeKelasRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:        "Branch 4 -> repo not found",
			idNamaKelas: 10,
			patch: updatepatch.NamaKelasPatch{
				NamaKelas:      ptrString("  IPA 2  "),
				IdTingkatKelas: &idTingkatKelas,
			},
			repo:            &FakeKelasRepo{UpdateErr: coreerror.ErrNotFound},
			wantErr:         coreerror.ErrNotFound,
			wantRepoCalled:  true,
			wantPatchedName: ptrString("IPA 2"),
		},
		{
			name:        "Branch 5 -> repo error lain",
			idNamaKelas: 10,
			patch: updatepatch.NamaKelasPatch{
				NamaKelas: ptrString("XI IPA 1"),
			},
			repo:            &FakeKelasRepo{UpdateErr: repoErr},
			wantErr:         repoErr,
			wantRepoCalled:  true,
			wantPatchedName: ptrString("XI IPA 1"),
		},
		{
			name:        "Branch 6 -> sukses update",
			idNamaKelas: 11,
			patch: updatepatch.NamaKelasPatch{
				NamaKelas:      ptrString("  XII IPA 3"),
				IdTingkatKelas: &idTingkatKelas,
			},
			repo:            &FakeKelasRepo{},
			wantErr:         nil,
			wantRepoCalled:  true,
			wantPatchedName: ptrString("XII IPA 3"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := kelas_service.NewUpdateKelasService(tt.repo)
			err := svc.UpdateNamaKelas(ctx, tt.idNamaKelas, tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantRepoCalled, tt.repo.UpdateCalled)
			if tt.wantRepoCalled {
				assert.Equal(t, tt.idNamaKelas, tt.repo.GotIDNamaKelas)
				if tt.wantPatchedName != nil {
					if assert.NotNil(t, tt.repo.GotPatch.NamaKelas) {
						assert.Equal(t, *tt.wantPatchedName, *tt.repo.GotPatch.NamaKelas)
					}
				}
			}
		})
	}
}

func ptrString(s string) *string { return &s }

func ptrID(id kelas.ID) *kelas.ID { return &id }
