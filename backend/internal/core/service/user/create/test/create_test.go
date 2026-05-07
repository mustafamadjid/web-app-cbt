package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	"github.com/stretchr/testify/assert"
)

func strPtr(v string) *string { return &v }

func TestCreateGuruBasisPath(t *testing.T) {
	validCmd := func() user_service.CreateGuruCmd {
		return user_service.CreateGuruCmd{
			Username:     "guruuser",
			Email:        strPtr("guru@example.com"),
			Password:     "password",
			NamaLengkap:  "Guru Test",
			JenisKelamin: "L",
			NoHp:         strPtr("08123456789"),
			Foto:         "foto.png",
			Nip:          "123456789012345678",
			Jabatan:      "Kepala",
			BidangStudi:  "Matematika",
		}
	}

	adminActor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	nonAdminActor := user.Actor{IdPengguna: 2, Role: user.GURU}

	testErr := errors.New("test error")

	cases := []struct {
		name               string
		cmd                user_service.CreateGuruCmd
		actor              user.Actor
		txm                *FakeTxManager
		hasher             *FakeHasher
		wantErr            error
		wantBeginCalled    bool
		wantHasherCalled   bool
		wantCommitCalled   bool
		wantRollbackCalled bool
		wantCreateUser     bool
		wantCreateProfil   bool
		wantResult         user_service.CreateGuruRes
	}{
		{
			name:  "Path 1 -> semua validasi lolos",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{CreateID: 10},
				ProfilGuruRepo: &FakeProfilGuruRepo{CreateID: 20},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
			wantResult:         user_service.CreateGuruRes{IdPengguna: 10, IdProfilGuru: 20},
		},
		{
			name: "Path 2 -> foto guru opsional",
			cmd: func() user_service.CreateGuruCmd {
				c := validCmd()
				c.Foto = ""
				return c
			}(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{CreateID: 11},
				ProfilGuruRepo: &FakeProfilGuruRepo{CreateID: 21},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
			wantResult:         user_service.CreateGuruRes{IdPengguna: 11, IdProfilGuru: 21},
		},
		{
			name:  "Path 3 -> username taken",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{ExistByUsername: true},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrUsernameTaken,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 4 -> panjang username tidak valid",
			cmd:   func() user_service.CreateGuruCmd { c := validCmd(); c.Username = "guru"; return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrUsernameLengthInvalid,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 5 -> NIP taken",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{ExistByUsername: false},
				ProfilGuruRepo: &FakeProfilGuruRepo{ExistsNip: true},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrNipTaken,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 6 -> create pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{CreateErr: testErr},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 7 -> create profil guru gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{CreateID: 10},
				ProfilGuruRepo: &FakeProfilGuruRepo{CreateErr: testErr},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
		},
		{
			name:  "Path 8 -> commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{CreateID: 10},
				ProfilGuruRepo: &FakeProfilGuruRepo{CreateID: 20},
				CommitErr:      testErr,
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
		},
		{
			name:  "Path 9 -> invalid email",
			cmd:   func() user_service.CreateGuruCmd { c := validCmd(); c.Email = strPtr("not-an-email"); return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 10 -> invalid nip",
			cmd:   func() user_service.CreateGuruCmd { c := validCmd(); c.Nip = "123"; return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            user.ErrInvalidNIP,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 11 -> hash password gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:       &FakeUserRepo{},
				ProfilGuruRepo: &FakeProfilGuruRepo{},
			}},
			hasher:             &FakeHasher{Err: testErr},
			wantErr:            testErr,
			wantBeginCalled:    false,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 12 -> begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				BeginErr: testErr,
				Tx: &FakeTx{
					UserRepo:       &FakeUserRepo{},
					ProfilGuruRepo: &FakeProfilGuruRepo{},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 13 -> aktor bukan admin",
			cmd:   validCmd(),
			actor: nonAdminActor,
			txm: &FakeTxManager{
				Tx: &FakeTx{
					UserRepo:       &FakeUserRepo{},
					ProfilGuruRepo: &FakeProfilGuruRepo{},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 14 -> cek username gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				Tx: &FakeTx{
					UserRepo:       &FakeUserRepo{UserExistErr: coreerror.ErrDbError},
					ProfilGuruRepo: &FakeProfilGuruRepo{},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrDbError,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 15 -> cek nip gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				Tx: &FakeTx{
					UserRepo:       &FakeUserRepo{},
					ProfilGuruRepo: &FakeProfilGuruRepo{ExistErr: coreerror.ErrDbError},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrDbError,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := user_service.NewCreateGuruService(tc.txm, tc.hasher)

			res, err := service.CreateGuru(context.Background(), tc.cmd, tc.actor)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantResult, res)
			assert.Equal(t, tc.wantBeginCalled, tc.txm.BeginCalled)
			assert.Equal(t, tc.wantHasherCalled, tc.hasher.Called)

			if tc.txm.Tx != nil {
				assert.Equal(t, tc.wantCommitCalled, tc.txm.Tx.CommitCalled)
				assert.Equal(t, tc.wantRollbackCalled, tc.txm.Tx.RollbackCalled)
				if tc.txm.Tx.UserRepo != nil {
					assert.Equal(t, tc.wantCreateUser, tc.txm.Tx.UserRepo.CreateCalled)
				}
				if tc.txm.Tx.ProfilGuruRepo != nil {
					assert.Equal(t, tc.wantCreateProfil, tc.txm.Tx.ProfilGuruRepo.CreateCalled)
				}
			}
		})
	}
}

func TestCreateSiswaBasisPath(t *testing.T) {
	validCmd := func() user_service.CreateSiswaCmd {
		return user_service.CreateSiswaCmd{
			Username:     "siswauser",
			Email:        strPtr("siswa@example.com"),
			Password:     "password",
			NamaLengkap:  "Siswa Test",
			JenisKelamin: "L",
			NoHp:         strPtr("08123456789"),
			Foto:         "foto.png",
			IdNamaKelas:  2,
			Nisn:         "1234567890",
			NoAbsen:      1,
			Angkatan:     2021,
			TempatLahir:  "Bandung",
			TanggalLahir: time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC),
		}
	}

	adminActor := user.Actor{IdPengguna: 1, Role: user.ADMIN}
	nonAdminActor := user.Actor{IdPengguna: 1, Role: user.SISWA}

	testErr := errors.New("test error")

	cases := []struct {
		name               string
		cmd                user_service.CreateSiswaCmd
		actor              user.Actor
		txm                *FakeTxManager
		hasher             *FakeHasher
		wantErr            error
		wantBeginCalled    bool
		wantHasherCalled   bool
		wantCommitCalled   bool
		wantRollbackCalled bool
		wantCreateUser     bool
		wantCreateProfil   bool
		wantResult         user_service.CreateSiswaRes
	}{
		{
			name:  "Path 1 -> semua validasi lolos",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{CreateID: 10},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{CreateID: 20},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
			wantResult:         user_service.CreateSiswaRes{IdPengguna: 10, IdProfilSiswa: 20},
		},
		{
			name: "Path 2 -> foto siswa opsional",
			cmd: func() user_service.CreateSiswaCmd {
				c := validCmd()
				c.Foto = ""
				return c
			}(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{CreateID: 11},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{CreateID: 21},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
			wantResult:         user_service.CreateSiswaRes{IdPengguna: 11, IdProfilSiswa: 21},
		},
		{
			name:  "Path 3 -> username taken",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{ExistByUsername: true},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrUsernameTaken,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 4 -> panjang username tidak valid",
			cmd:   func() user_service.CreateSiswaCmd { c := validCmd(); c.Username = "sisw"; return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrUsernameLengthInvalid,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 5 -> NISN taken",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{ExistByUsername: false},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{ExistsNisn: true},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrNisnTaken,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 6 -> create pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{CreateErr: testErr},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 7 -> create profil siswa gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{CreateID: 10},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{CreateErr: testErr},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
		},
		{
			name:  "Path 8 -> commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{CreateID: 10},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{CreateID: 20},
				CommitErr:       testErr,
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
		},
		{
			name:  "Path 9 -> invalid email",
			cmd:   func() user_service.CreateSiswaCmd { c := validCmd(); c.Email = strPtr("not-an-email"); return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 10 -> invalid nisn",
			cmd:   func() user_service.CreateSiswaCmd { c := validCmd(); c.Nisn = "123"; return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            user.ErrInvalidNISN,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 11 -> invalid absen",
			cmd:   func() user_service.CreateSiswaCmd { c := validCmd(); c.NoAbsen = 0; return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            user.ErrInvalidAbsen,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 12 -> invalid angkatan",
			cmd:   func() user_service.CreateSiswaCmd { c := validCmd(); c.Angkatan = 2000; return c }(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            user.ErrInvalidAngkatan,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 13 -> hash password gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{Tx: &FakeTx{
				UserRepo:        &FakeUserRepo{},
				ProfilSiswaRepo: &FakeProfilSiswaRepo{},
			}},
			hasher:             &FakeHasher{Err: testErr},
			wantErr:            testErr,
			wantBeginCalled:    false,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 14 -> begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				BeginErr: testErr,
				Tx: &FakeTx{
					UserRepo:        &FakeUserRepo{},
					ProfilSiswaRepo: &FakeProfilSiswaRepo{},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 15 -> aktor non admin",
			cmd:   validCmd(),
			actor: nonAdminActor,
			txm: &FakeTxManager{
				Tx: &FakeTx{
					UserRepo:        &FakeUserRepo{},
					ProfilSiswaRepo: &FakeProfilSiswaRepo{},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 16 -> cek username gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				Tx: &FakeTx{
					UserRepo:        &FakeUserRepo{UserExistErr: coreerror.ErrDbError},
					ProfilSiswaRepo: &FakeProfilSiswaRepo{},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrDbError,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Path 17 -> cek nisn gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &FakeTxManager{
				Tx: &FakeTx{
					UserRepo:        &FakeUserRepo{},
					ProfilSiswaRepo: &FakeProfilSiswaRepo{ExistNisnErr: coreerror.ErrDbError},
				},
			},
			hasher:             &FakeHasher{Hash: "hashed"},
			wantErr:            coreerror.ErrDbError,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := user_service.NewCreateSiswaService(tc.txm, tc.hasher)

			res, err := service.CreateSiswa(context.Background(), tc.cmd, tc.actor)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantResult, res)
			assert.Equal(t, tc.wantBeginCalled, tc.txm.BeginCalled)
			assert.Equal(t, tc.wantHasherCalled, tc.hasher.Called)

			if tc.txm.Tx != nil {
				assert.Equal(t, tc.wantCommitCalled, tc.txm.Tx.CommitCalled)
				assert.Equal(t, tc.wantRollbackCalled, tc.txm.Tx.RollbackCalled)
				if tc.txm.Tx.UserRepo != nil {
					assert.Equal(t, tc.wantCreateUser, tc.txm.Tx.UserRepo.CreateCalled)
				}
				if tc.txm.Tx.ProfilSiswaRepo != nil {
					assert.Equal(t, tc.wantCreateProfil, tc.txm.Tx.ProfilSiswaRepo.CreateCalled)
				}
			}
		})
	}
}

var _ out.PasswordHasher = (*FakeHasher)(nil)
