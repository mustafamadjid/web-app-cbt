package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestUpdateGuruBranchCoverage(t *testing.T) {
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
		txm                *fake_test.FakeTxManager
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
			name:  "Branch 1 -> semua patch berhasil",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
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
			name:               "Branch 2 -> aktor bukan admin",
			cmd:                validCmd(),
			actor:              nonAdminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 3 -> id pengguna kosong",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.IdPengguna = 0; return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErrText:        "Id pengguna required",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 4 -> username kosong",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.Username = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            coreerror.ErrUsernameLengthInvalid,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 5 -> nama lengkap kosong",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.NamaLengkap = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErrText:        "nama_lengkap cannot be empty",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 6 -> email tidak valid",
			cmd:                func() user_service.UpdateGuruCmd { c := validCmd(); c.Email = strPtr("invalid"); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name: "Branch 7 -> tidak ada field yang diupdate",
			cmd: func() user_service.UpdateGuruCmd {
				return user_service.UpdateGuruCmd{IdPengguna: 10}
			}(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            coreerror.ErrNoFieldToUpdate,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 8 -> begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{
				BeginErr: testErr,
				Tx: &fake_test.FakeTx{
					UserRepo:       &fake_test.FakeUserRepo{},
					ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
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
			name:  "Branch 9 -> update pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{UpdateErr: testErr},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 10 -> update profil guru gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{UpdateErr: testErr},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name:  "Branch 11 -> commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
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
			name: "Branch 12 -> hanya update profil",
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
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     false,
			wantUpdateProfil:   true,
		},
		{
			name: "Branch 13 -> hanya update pengguna",
			cmd: func() user_service.UpdateGuruCmd {
				c := validCmd()
				c.Nip = nil
				c.Jabatan = nil
				c.BidangStudi = nil
				return c
			}(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
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
			name: "Branch 14 -> hanya update username",
			cmd: user_service.UpdateGuruCmd{
				IdPengguna: 10,
				Username:   strPtr("guruupdate"),
			},
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:       &fake_test.FakeUserRepo{},
				ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantUsernameValue:  "guruupdate",
		},
		{
			name: "Branch 15 -> foto kosong diabaikan",
			cmd: user_service.UpdateGuruCmd{
				IdPengguna: 10,
				Foto:       strPtr(" "),
			},
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
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
			sessionRepo := &fake_test.FakeSessionRepo{}
			deleteFileRepo := &fake_test.FakeDeleteFileRepo{}
			userRepo := &fake_test.FakeUserRepo{}
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

func TestUpdateSiswaBranchCoverage(t *testing.T) {
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
		txm                *fake_test.FakeTxManager
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
			name:  "Branch 1 -> semua patch berhasil",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
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
			name:               "Branch 2 -> aktor bukan admin",
			cmd:                validCmd(),
			actor:              nonAdminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 3 -> id pengguna kosong",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.IdPengguna = 0; return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErrText:        "Id pengguna required",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 4 -> username kosong",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Username = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            coreerror.ErrUsernameLengthInvalid,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 5 -> nama lengkap kosong",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.NamaLengkap = strPtr(" "); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErrText:        "nama_lengkap cannot be empty",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 6 -> email tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Email = strPtr("invalid"); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 7 -> nisn tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Nisn = strPtr("123"); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            user.ErrInvalidNISN,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 8 -> no absen tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.NoAbsen = intPtr(0); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            user.ErrInvalidAbsen,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:               "Branch 9 -> angkatan tidak valid",
			cmd:                func() user_service.UpdateSiswaCmd { c := validCmd(); c.Angkatan = intPtr(2000); return c }(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            user.ErrInvalidAngkatan,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name: "Branch 10 -> tidak ada field yang diupdate",
			cmd: func() user_service.UpdateSiswaCmd {
				return user_service.UpdateSiswaCmd{IdPengguna: 10}
			}(),
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
			wantErr:            coreerror.ErrNoFieldToUpdate,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 11 -> begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{
				BeginErr: testErr,
				Tx: &fake_test.FakeTx{
					UserRepo:        &fake_test.FakeUserRepo{},
					ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
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
			name:  "Branch 12 -> update pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{UpdateErr: testErr},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 13 -> update profil siswa gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{UpdateErr: testErr},
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name:  "Branch 14 -> commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
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
			name: "Branch 15 -> hanya update profil",
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
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     false,
			wantUpdateProfil:   true,
		},
		{
			name: "Branch 16 -> hanya update pengguna",
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
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
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
			name: "Branch 17 -> hanya update username",
			cmd: user_service.UpdateSiswaCmd{
				IdPengguna: 10,
				Username:   strPtr("siswaupdate"),
			},
			actor: adminActor,
			txm: &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
				UserRepo:        &fake_test.FakeUserRepo{},
				ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantUsernameValue:  "siswaupdate",
		},
		{
			name: "Branch 18 -> foto kosong diabaikan",
			cmd: user_service.UpdateSiswaCmd{
				IdPengguna: 10,
				Foto:       strPtr(" "),
			},
			actor:              adminActor,
			txm:                &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{}},
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
			sessionRepo := &fake_test.FakeSessionRepo{}
			deleteFileRepo := &fake_test.FakeDeleteFileRepo{}
			userRepo := &fake_test.FakeUserRepo{}
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

func TestUpdateGuru_DeleteFileErrorBranch(t *testing.T) {
	t.Parallel()

	actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	fotoBaru := "foto-baru.png"
	deleteFileErr := errors.New("delete file failed")

	cmd := user_service.UpdateGuruCmd{
		IdPengguna: 10,
		Foto:       &fotoBaru,
	}

	txm := &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
		UserRepo:       &fake_test.FakeUserRepo{},
		ProfilGuruRepo: &fake_test.FakeProfilGuruRepo{},
	}}
	sessionRepo := &fake_test.FakeSessionRepo{}
	deleteFileRepo := &fake_test.FakeDeleteFileRepo{DeleteErr: deleteFileErr}
	userRepo := &fake_test.FakeUserRepo{
		FindResult: user.Pengguna{ID: cmd.IdPengguna, Foto: "foto-lama.png"},
	}
	service := user_service.NewUpdateUserService(txm, sessionRepo, deleteFileRepo, userRepo)

	err := service.UpdateGuru(context.Background(), cmd, actor)

	assert.ErrorIs(t, err, deleteFileErr)
	assert.True(t, userRepo.FindCalled)
	assert.Equal(t, cmd.IdPengguna, userRepo.LastFindID)
	assert.True(t, deleteFileRepo.DeleteCalled)
	assert.Equal(t, "foto-lama.png", deleteFileRepo.LastPath)
	assert.False(t, txm.BeginCalled)
	assert.False(t, sessionRepo.RevokeAllCalled)
}

func TestUpdateSiswa_DeleteFileErrorBranch(t *testing.T) {
	t.Parallel()

	actor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	fotoBaru := "foto-baru.png"
	deleteFileErr := errors.New("delete file failed")

	cmd := user_service.UpdateSiswaCmd{
		IdPengguna: 10,
		Foto:       &fotoBaru,
	}

	txm := &fake_test.FakeTxManager{Tx: &fake_test.FakeTx{
		UserRepo:        &fake_test.FakeUserRepo{},
		ProfilSiswaRepo: &fake_test.FakeProfilSiswaRepo{},
	}}
	sessionRepo := &fake_test.FakeSessionRepo{}
	deleteFileRepo := &fake_test.FakeDeleteFileRepo{DeleteErr: deleteFileErr}
	userRepo := &fake_test.FakeUserRepo{
		FindResult: user.Pengguna{ID: cmd.IdPengguna, Foto: "foto-lama.png"},
	}
	service := user_service.NewUpdateUserService(txm, sessionRepo, deleteFileRepo, userRepo)

	err := service.UpdateSiswa(context.Background(), cmd, actor)

	assert.ErrorIs(t, err, deleteFileErr)
	assert.True(t, userRepo.FindCalled)
	assert.Equal(t, cmd.IdPengguna, userRepo.LastFindID)
	assert.True(t, deleteFileRepo.DeleteCalled)
	assert.Equal(t, "foto-lama.png", deleteFileRepo.LastPath)
	assert.False(t, txm.BeginCalled)
	assert.False(t, sessionRepo.RevokeAllCalled)
}
