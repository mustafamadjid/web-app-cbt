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

	cases := []struct {
		name               string
		cmd                user_service.UpdateGuruCmd
		actor              user.Actor
		txm                *FakeTxManager
		wantErr            error
		wantErrText        string
		wantBeginCalled    bool
		wantCommitCalled   bool
		wantRollbackCalled bool
		wantUpdateUser     bool
		wantUpdateProfil   bool
		wantEmailValue     string
		wantUsernameValue  string
	}{
		{
			name:  "Path 1 -> semua patch berhasil",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
			wantEmailValue:     "guru@example.com",
			wantUsernameValue:  "guruuser",
		},
		{
			name:               "Path 2 -> aktor bukan admin",
			cmd:                validCmd(),
			actor:              nonAdminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 3 -> id pengguna kosong",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.IdPengguna = 0; return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErrText:        "Id pengguna required",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 4 -> username kosong",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.Username = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrUsernameLengthInvalid,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 5 -> nama lengkap kosong",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.NamaLengkap = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErrText:        "nama_lengkap cannot be empty",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 6 -> email tidak valid",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.Email = strPtr("invalid"); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name: "Path 7 -> tidak ada field yang diupdate",
			cmd: func() user_service.UpdateGuruCmd {
				return user_service.UpdateGuruCmd{IdPengguna: 10}
			}(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrNoFieldToUpdate,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Path 8 -> begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				BeginErr: testErr,
				Tx: &FakeTx{
					UserRepo:       &FakeUserRepo{},
					ProfilGuruRepo: &FakeProfilGuruRepo{},
				},
			},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Path 9 -> update pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{UpdateErr: testErr},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
		},
		{
			name:  "Path 10 -> update profil guru gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{UpdateErr: testErr},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name:  "Path 11 -> commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
				CommitErr:      testErr,
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name: "Path 12 -> hanya update profil",
			cmd: func() user_service.UpdateGuruCmd {
				c := validCmd()
				c.Username = nil
				c.NamaLengkap = nil
				c.Email = nil
				c.JenisKelamin = nil
				c.NoHp = nil
				c.Foto = nil
				c.StatusAkun = nil
				c.Role = nil
				return c
			}(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     false,
			wantUpdateProfil:   true,
		},
		{
			name: "Path 13 -> hanya update pengguna",
			cmd: func() user_service.UpdateGuruCmd {
				c := validCmd()
				c.Nip = nil
				c.Jabatan = nil
				c.BidangStudi = nil
				return c
			}(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantEmailValue:     "guru@example.com",
			wantUsernameValue:  "guruuser",
		},
		{
			name: "Path 14 -> hanya update username",
			cmd: user_service.UpdateGuruCmd{
				IdPengguna: 10,
				Username:   strPtr("guruupdate"),
			},
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantUsernameValue:  "guruupdate",
		},
		{
			name: "Path 15 -> foto kosong diabaikan",
			cmd: user_service.UpdateGuruCmd{
				IdPengguna: 10,
				Foto:       strPtr(" "),
			},
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrNoFieldToUpdate,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sessionRepo := &FakeSessionRepo{}
			deleteFileRepo := &FakeDeleteFileRepo{}
			userRepo := &FakeUserRepo{}
			service := user_service.NewUpdateUserService(tc.txm, sessionRepo, deleteFileRepo, userRepo)

			err := service.UpdateGuru(context.Background(), tc.cmd, tc.actor)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else if tc.wantErrText != "" {
				assert.EqualError(t, err, tc.wantErrText)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantBeginCalled, tc.txm.BeginCalled)
			if tc.txm.Tx != nil {
				assert.Equal(t, tc.wantCommitCalled, tc.txm.Tx.CommitCalled)
				assert.Equal(t, tc.wantRollbackCalled, tc.txm.Tx.RollbackCalled)
				if tc.txm.Tx.UserRepo != nil {
					assert.Equal(t, tc.wantUpdateUser, tc.txm.Tx.UserRepo.UpdateCalled)
					if tc.wantUsernameValue != "" && tc.txm.Tx.UserRepo.LastPatch.Username != nil {
						assert.Equal(t, tc.wantUsernameValue, *tc.txm.Tx.UserRepo.LastPatch.Username)
					}
					if tc.wantEmailValue != "" && tc.txm.Tx.UserRepo.LastPatch.Email != nil {
						assert.Equal(t, tc.wantEmailValue, string(*tc.txm.Tx.UserRepo.LastPatch.Email))
					}
				}
				if tc.txm.Tx.ProfilGuruRepo != nil {
					assert.Equal(t, tc.wantUpdateProfil, tc.txm.Tx.ProfilGuruRepo.UpdateCalled)
				}
			}
		})
	}
}

func TestUpdateSiswaBasisPath(t *testing.T) {
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

	cases := []struct {
		name               string
		cmd                user_service.UpdateSiswaCmd
		actor              user.Actor
		txm                *FakeTxManager
		wantErr            error
		wantErrText        string
		wantBeginCalled    bool
		wantCommitCalled   bool
		wantRollbackCalled bool
		wantUpdateUser     bool
		wantUpdateProfil   bool
		wantEmailValue     string
		wantNisnValue      string
		wantUsernameValue  string
	}{
		{
			name:  "Path 1 -> semua patch berhasil",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
			wantEmailValue:     "siswa@example.com",
			wantNisnValue:      "1234567890",
			wantUsernameValue:  "siswauser",
		},
		{
			name:               "Path 2 -> aktor bukan admin",
			cmd:                validCmd(),
			actor:              nonAdminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 3 -> id pengguna kosong",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.IdPengguna = 0; return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErrText:        "Id pengguna required",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 4 -> username kosong",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Username = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrUsernameLengthInvalid,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 5 -> nama lengkap kosong",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.NamaLengkap = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErrText:        "nama_lengkap cannot be empty",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 6 -> email tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Email = strPtr("invalid"); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 7 -> nisn tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Nisn = strPtr("123"); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            user.ErrInvalidNISN,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 8 -> no absen tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.NoAbsen = intPtr(0); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            user.ErrInvalidAbsen,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Path 9 -> angkatan tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Angkatan = intPtr(2000); return c }(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            user.ErrInvalidAngkatan,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name: "Path 10 -> tidak ada field yang diupdate",
			cmd: func() user_service.UpdateSiswaCmd {
				return user_service.UpdateSiswaCmd{IdPengguna: 10}
			}(),
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrNoFieldToUpdate,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Path 11 -> begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				BeginErr: testErr,
				Tx: &FakeTx{
					UserRepo:        &FakeUserRepo{},
					ProfilSiswaRepo: &FakeProfilSiswaRepo{},
				},
			},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Path 12 -> update pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{UpdateErr: testErr},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
		},
		{
			name:  "Path 13 -> update profil siswa gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{UpdateErr: testErr},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name:  "Path 14 -> commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
				CommitErr:       testErr,
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name: "Path 15 -> hanya update profil",
			cmd: func() user_service.UpdateSiswaCmd {
				c := validCmd()
				c.Username = nil
				c.NamaLengkap = nil
				c.Email = nil
				c.JenisKelamin = nil
				c.NoHp = nil
				c.Foto = nil
				c.StatusAkun = nil
				c.Role = nil
				return c
			}(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     false,
			wantUpdateProfil:   true,
		},
		{
			name: "Path 16 -> hanya update pengguna",
			cmd: func() user_service.UpdateSiswaCmd {
				c := validCmd()
				c.IdTingkatKelas = nil
				c.IdNamaKelas = nil
				c.Nisn = nil
				c.NoAbsen = nil
				c.Angkatan = nil
				c.TempatLahir = nil
				c.TanggalLahir = nil
				return c
			}(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantEmailValue:     "siswa@example.com",
			wantUsernameValue:  "siswauser",
		},
		{
			name: "Path 17 -> hanya update username",
			cmd: user_service.UpdateSiswaCmd{
				IdPengguna: 10,
				Username:   strPtr("siswaupdate"),
			},
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantUsernameValue:  "siswaupdate",
		},
		{
			name: "Path 18 -> foto kosong diabaikan",
			cmd: user_service.UpdateSiswaCmd{
				IdPengguna: 10,
				Foto:       strPtr(" "),
			},
			actor:              adminActor,
			txm:                &FakeTxManager{Tx: &FakeTx{}},
			wantErr:            coreerror.ErrNoFieldToUpdate,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sessionRepo := &FakeSessionRepo{}
			deleteFileRepo := &FakeDeleteFileRepo{}
			userRepo := &FakeUserRepo{}
			service := user_service.NewUpdateUserService(tc.txm, sessionRepo, deleteFileRepo, userRepo)

			err := service.UpdateSiswa(context.Background(), tc.cmd, tc.actor)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else if tc.wantErrText != "" {
				assert.EqualError(t, err, tc.wantErrText)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantBeginCalled, tc.txm.BeginCalled)
			if tc.txm.Tx != nil {
				assert.Equal(t, tc.wantCommitCalled, tc.txm.Tx.CommitCalled)
				assert.Equal(t, tc.wantRollbackCalled, tc.txm.Tx.RollbackCalled)
				if tc.txm.Tx.UserRepo != nil {
					assert.Equal(t, tc.wantUpdateUser, tc.txm.Tx.UserRepo.UpdateCalled)
					if tc.wantUsernameValue != "" && tc.txm.Tx.UserRepo.LastPatch.Username != nil {
						assert.Equal(t, tc.wantUsernameValue, *tc.txm.Tx.UserRepo.LastPatch.Username)
					}
					if tc.wantEmailValue != "" && tc.txm.Tx.UserRepo.LastPatch.Email != nil {
						assert.Equal(t, tc.wantEmailValue, string(*tc.txm.Tx.UserRepo.LastPatch.Email))
					}
				}
				if tc.txm.Tx.ProfilSiswaRepo != nil {
					assert.Equal(t, tc.wantUpdateProfil, tc.txm.Tx.ProfilSiswaRepo.UpdateCalled)
					if tc.wantNisnValue != "" && tc.txm.Tx.ProfilSiswaRepo.LastPatch.Nisn != nil {
						assert.Equal(t, tc.wantNisnValue, *tc.txm.Tx.ProfilSiswaRepo.LastPatch.Nisn)
					}
				}
			}
		})
	}
}

func TestUpdateUser_DeleteFileErrorBasisPath(t *testing.T) {
	t.Parallel()

	actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	fotoBaru := "foto-baru.png"
	deleteFileErr := errors.New("delete file failed")

	tests := []struct {
		name            string
		runGuru         bool
		wantDeletePath  string
		wantFindID      user.ID
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
