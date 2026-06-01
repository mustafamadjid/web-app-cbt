package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
	"github.com/stretchr/testify/assert"
)

func TestUpdateGuruBasisPath(t *testing.T) {
	t.Parallel()

	adminActor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	nonAdminActor := user.Actor{IdPengguna: 2, Role: user.GURU}

	strPtr := func(value string) *string {
		return &value
	}
	statusAktif := user.AKTIF
	roleGuru := user.GURU

	validCmd := func() user_service.UpdateGuruCmd {
		return user_service.UpdateGuruCmd{
			IdPengguna:   10,
			Username:     strPtr("guruuser"),
			Email:        strPtr("guru@example.com"),
			NamaLengkap:  strPtr("Guru Test"),
			JenisKelamin: strPtr("L"),
			NoHp:         strPtr("08123456789"),
			Foto:         strPtr("foto.png"),
			StatusAkun:   &statusAktif,
			Role:         &roleGuru,
			Nip:          strPtr("123456789012345678"),
			Jabatan:      strPtr("Kepala"),
			BidangStudi:  strPtr("Matematika"),
		}
	}

	testErr := errors.New("test error")

	t.Run("P1 -> aktor tidak memiliki hak untuk memperbarui data guru", func(t *testing.T) {
		t.Parallel()
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateGuru(context.Background(), validCmd(), nonAdminActor)

		assert.ErrorIs(t, err, coreerror.ErrForbidden)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P2 -> ID pengguna tidak valid", func(t *testing.T) {
		t.Parallel()
		cmd := validCmd()
		cmd.IdPengguna = 0
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateGuru(context.Background(), cmd, adminActor)

		assert.EqualError(t, err, "Id pengguna required")
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P3 -> data update guru tidak valid setelah sanitasi dan validasi", func(t *testing.T) {
		t.Parallel()
		cmd := validCmd()
		cmd.Email = strPtr("invalid")
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateGuru(context.Background(), cmd, adminActor)

		assert.ErrorIs(t, err, user.ErrInvalidEmail)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P4 -> foto baru dikirim tetapi foto lama gagal dihapus", func(t *testing.T) {
		t.Parallel()
		txm := &FakeTxManager{Tx: &FakeTx{}}
		deleteFileRepo := &FakeDeleteFileRepo{DeleteErr: testErr}
		userRepo := &FakeUserRepo{FindResult: user.Pengguna{ID: 10, Foto: "foto-lama.png"}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, deleteFileRepo, userRepo)

		err := service.UpdateGuru(context.Background(), validCmd(), adminActor)

		assert.ErrorIs(t, err, testErr)
		assert.True(t, userRepo.FindCalled)
		assert.True(t, deleteFileRepo.DeleteCalled)
		assert.Equal(t, "foto-lama.png", deleteFileRepo.LastPath)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P5 -> foto kosong setelah sanitasi dan tidak ada field yang diperbarui", func(t *testing.T) {
		t.Parallel()
		cmd := user_service.UpdateGuruCmd{
			IdPengguna: 10,
			Foto:       strPtr(" "),
		}
		txm := &FakeTxManager{Tx: &FakeTx{}}
		deleteFileRepo := &FakeDeleteFileRepo{}
		userRepo := &FakeUserRepo{}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, deleteFileRepo, userRepo)

		err := service.UpdateGuru(context.Background(), cmd, adminActor)

		assert.ErrorIs(t, err, coreerror.ErrNoFieldToUpdate)
		assert.False(t, userRepo.FindCalled)
		assert.False(t, deleteFileRepo.DeleteCalled)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P6 -> tidak ada foto baru dan tidak ada field yang diperbarui", func(t *testing.T) {
		t.Parallel()
		cmd := user_service.UpdateGuruCmd{IdPengguna: 10}
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateGuru(context.Background(), cmd, adminActor)

		assert.ErrorIs(t, err, coreerror.ErrNoFieldToUpdate)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P7 -> data guru berhasil diperbarui", func(t *testing.T) {
		t.Parallel()
		txm := &FakeTxManager{Tx: &FakeTx{
			UserRepo:       &FakeUserRepo{},
			ProfilGuruRepo: &FakeProfilGuruRepo{},
		}}
		deleteFileRepo := &FakeDeleteFileRepo{}
		userRepo := &FakeUserRepo{FindResult: user.Pengguna{ID: 10, Foto: "foto-lama.png"}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, deleteFileRepo, userRepo)

		err := service.UpdateGuru(context.Background(), validCmd(), adminActor)

		assert.NoError(t, err)
		assert.True(t, userRepo.FindCalled)
		assert.True(t, deleteFileRepo.DeleteCalled)
		assert.True(t, txm.BeginCalled)
		assert.True(t, txm.Tx.UserRepo.UpdateCalled)
		assert.True(t, txm.Tx.ProfilGuruRepo.UpdateCalled)
		assert.True(t, txm.Tx.CommitCalled)
		assert.True(t, txm.Tx.RollbackCalled)
		assert.Equal(t, "guruuser", *txm.Tx.UserRepo.LastPatch.Username)
		assert.Equal(t, "guru@example.com", string(*txm.Tx.UserRepo.LastPatch.Email))
	})
}

func TestUpdateSiswaBasisPath(t *testing.T) {
	t.Parallel()

	adminActor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	nonAdminActor := user.Actor{IdPengguna: 2, Role: user.GURU}

	strPtr := func(value string) *string {
		return &value
	}
	intPtr := func(value int) *int {
		return &value
	}
	idPtr := func(value user.ID) *user.ID {
		return &value
	}
	timePtr := func(value time.Time) *time.Time {
		return &value
	}

	statusAktif := user.AKTIF
	roleSiswa := user.SISWA
	validTanggal := time.Date(2005, time.March, 10, 0, 0, 0, 0, time.UTC)

	validCmd := func() user_service.UpdateSiswaCmd {
		return user_service.UpdateSiswaCmd{
			IdPengguna:     10,
			Username:       strPtr("siswauser"),
			Email:          strPtr("siswa@example.com"),
			NamaLengkap:    strPtr("Siswa Test"),
			JenisKelamin:   strPtr("L"),
			NoHp:           strPtr("08123456789"),
			Foto:           strPtr("foto.png"),
			StatusAkun:     &statusAktif,
			Role:           &roleSiswa,
			IdTingkatKelas: idPtr(3),
			IdNamaKelas:    idPtr(4),
			Nisn:           strPtr("1234567890"),
			NoAbsen:        intPtr(5),
			Angkatan:       intPtr(2020),
			TempatLahir:    strPtr("Bandung"),
			TanggalLahir:   timePtr(validTanggal),
		}
	}

	testErr := errors.New("test error")

	t.Run("P1 -> aktor tidak memiliki hak untuk memperbarui data siswa", func(t *testing.T) {
		t.Parallel()
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateSiswa(context.Background(), validCmd(), nonAdminActor)

		assert.ErrorIs(t, err, coreerror.ErrForbidden)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P2 -> ID pengguna tidak valid", func(t *testing.T) {
		t.Parallel()
		cmd := validCmd()
		cmd.IdPengguna = 0
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateSiswa(context.Background(), cmd, adminActor)

		assert.EqualError(t, err, "Id pengguna required")
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P3 -> data update siswa tidak valid setelah sanitasi dan validasi", func(t *testing.T) {
		t.Parallel()
		cmd := validCmd()
		cmd.Nisn = strPtr("123")
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateSiswa(context.Background(), cmd, adminActor)

		assert.ErrorIs(t, err, user.ErrInvalidNISN)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P4 -> foto baru dikirim tetapi foto lama gagal dihapus", func(t *testing.T) {
		t.Parallel()
		txm := &FakeTxManager{Tx: &FakeTx{}}
		deleteFileRepo := &FakeDeleteFileRepo{DeleteErr: testErr}
		userRepo := &FakeUserRepo{FindResult: user.Pengguna{ID: 10, Foto: "foto-lama.png"}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, deleteFileRepo, userRepo)

		err := service.UpdateSiswa(context.Background(), validCmd(), adminActor)

		assert.ErrorIs(t, err, testErr)
		assert.True(t, userRepo.FindCalled)
		assert.True(t, deleteFileRepo.DeleteCalled)
		assert.Equal(t, "foto-lama.png", deleteFileRepo.LastPath)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P5 -> foto kosong setelah sanitasi dan tidak ada field yang diperbarui", func(t *testing.T) {
		t.Parallel()
		cmd := user_service.UpdateSiswaCmd{
			IdPengguna: 10,
			Foto:       strPtr(" "),
		}
		txm := &FakeTxManager{Tx: &FakeTx{}}
		deleteFileRepo := &FakeDeleteFileRepo{}
		userRepo := &FakeUserRepo{}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, deleteFileRepo, userRepo)

		err := service.UpdateSiswa(context.Background(), cmd, adminActor)

		assert.ErrorIs(t, err, coreerror.ErrNoFieldToUpdate)
		assert.False(t, userRepo.FindCalled)
		assert.False(t, deleteFileRepo.DeleteCalled)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P6 -> tidak ada foto baru dan tidak ada field yang diperbarui", func(t *testing.T) {
		t.Parallel()
		cmd := user_service.UpdateSiswaCmd{IdPengguna: 10}
		txm := &FakeTxManager{Tx: &FakeTx{}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, &FakeDeleteFileRepo{}, &FakeUserRepo{})

		err := service.UpdateSiswa(context.Background(), cmd, adminActor)

		assert.ErrorIs(t, err, coreerror.ErrNoFieldToUpdate)
		assert.False(t, txm.BeginCalled)
	})

	t.Run("P7 -> data siswa berhasil diperbarui", func(t *testing.T) {
		t.Parallel()
		txm := &FakeTxManager{Tx: &FakeTx{
			UserRepo:        &FakeUserRepo{},
			ProfilSiswaRepo: &FakeProfilSiswaRepo{},
		}}
		deleteFileRepo := &FakeDeleteFileRepo{}
		userRepo := &FakeUserRepo{FindResult: user.Pengguna{ID: 10, Foto: "foto-lama.png"}}
		service := user_service.NewUpdateUserService(txm, &FakeSessionRepo{}, deleteFileRepo, userRepo)

		err := service.UpdateSiswa(context.Background(), validCmd(), adminActor)

		assert.NoError(t, err)
		assert.True(t, userRepo.FindCalled)
		assert.True(t, deleteFileRepo.DeleteCalled)
		assert.True(t, txm.BeginCalled)
		assert.True(t, txm.Tx.UserRepo.UpdateCalled)
		assert.True(t, txm.Tx.ProfilSiswaRepo.UpdateCalled)
		assert.True(t, txm.Tx.CommitCalled)
		assert.True(t, txm.Tx.RollbackCalled)
		assert.Equal(t, "siswauser", *txm.Tx.UserRepo.LastPatch.Username)
		assert.Equal(t, "siswa@example.com", string(*txm.Tx.UserRepo.LastPatch.Email))
		assert.Equal(t, "1234567890", *txm.Tx.ProfilSiswaRepo.LastPatch.Nisn)
	})
}

func TestUpdateUser_DeleteFileErrorBasisPath(t *testing.T) {
	t.Parallel()

	actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	fotoBaru := "foto-baru.png"
	deleteFileErr := errors.New("delete file failed")

	tests := []struct {
		name           string
		runGuru        bool
		wantDeletePath string
		wantFindID     user.ID
	}{
		{
			name:           "Path 1 -> update guru gagal saat hapus foto lama",
			runGuru:        true,
			wantDeletePath: "foto-lama.png",
			wantFindID:     10,
		},
		{
			name:           "Path 2 -> update siswa gagal saat hapus foto lama",
			runGuru:        false,
			wantDeletePath: "foto-lama.png",
			wantFindID:     10,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			txm := &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilGuruRepo:  &FakeProfilGuruRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}}
			sessionRepo := &FakeSessionRepo{}
			deleteFileRepo := &FakeDeleteFileRepo{DeleteErr: deleteFileErr}
			userRepo := &FakeUserRepo{
				FindResult: user.Pengguna{ID: tc.wantFindID, Foto: "foto-lama.png"},
			}
			service := user_service.NewUpdateUserService(txm, sessionRepo, deleteFileRepo, userRepo)

			var err error
			if tc.runGuru {
				err = service.UpdateGuru(context.Background(), user_service.UpdateGuruCmd{
					IdPengguna: tc.wantFindID,
					Foto:       &fotoBaru,
				}, actor)
			} else {
				err = service.UpdateSiswa(context.Background(), user_service.UpdateSiswaCmd{
					IdPengguna: tc.wantFindID,
					Foto:       &fotoBaru,
				}, actor)
			}

			assert.ErrorIs(t, err, deleteFileErr)
			assert.True(t, userRepo.FindCalled)
			assert.Equal(t, tc.wantFindID, userRepo.LastFindID)
			assert.True(t, deleteFileRepo.DeleteCalled)
			assert.Equal(t, tc.wantDeletePath, deleteFileRepo.LastPath)
			assert.False(t, txm.BeginCalled)
			assert.False(t, sessionRepo.RevokeAllCalled)
		})
	}
}
