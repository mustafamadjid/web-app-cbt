package user_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	txport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
	"github.com/stretchr/testify/assert"
)

// ===== minimal fake user repo =====

type fakeUserRepo struct {
	updateErr error

	updateCalled bool
	lastID       user.ID
	lastPatch    outuser.UpdatePenggunaPatch
}

func (r *fakeUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna outuser.UpdatePenggunaPatch) error {
	r.updateCalled = true
	r.lastID = idPengguna
	r.lastPatch = pengguna
	if r.updateErr != nil {
		return r.updateErr
	}
	return nil
}

// ---- methods below are NOT used by this test; keep minimal ----

func (r *fakeUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	panic("not used in this test")
}

func (r *fakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	panic("not used in this test")
}

func (r *fakeUserRepo) CreateUser(ctx context.Context, p user.Pengguna) (user.ID, error) {
	panic("not used in this test")
}

func (r *fakeUserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	panic("not used in this test")
}

func (r *fakeUserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	panic("not used in this test")
}

// ===== minimal fake profil guru repo =====

type fakeProfilGuruRepo struct {
	updateErr error

	updateCalled bool
	lastID       user.ID
	lastPatch    outuser.UpdateProfilGuruPatch
}

func (r *fakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru outuser.UpdateProfilGuruPatch) error {
	r.updateCalled = true
	r.lastID = idPengguna
	r.lastPatch = profilGuru
	if r.updateErr != nil {
		return r.updateErr
	}
	return nil
}

func (r *fakeProfilGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.ProfilGuru, error) {
	panic("not used in this test")
}

func (r *fakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	panic("not used in this test")
}

func (r *fakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	panic("not used in this test")
}

// ===== fake tx & tx manager =====

type fakeTx struct {
	userRepo       *fakeUserRepo
	profilGuruRepo *fakeProfilGuruRepo

	commitCalled   bool
	rollbackCalled bool
	commitErr      error
}

type fakeTxManager struct {
	tx          *fakeTx
	beginErr    error
	beginCalled bool
}

func (m *fakeTxManager) Begin(ctx context.Context) (txport.Tx, error) {
	m.beginCalled = true
	if m.beginErr != nil {
		return nil, m.beginErr
	}
	return m.tx, nil
}

func (t *fakeTx) Pengguna() outuser.UserRepository {
	return t.userRepo
}

func (t *fakeTx) ProfilGuru() outuser.ProfilGuruRepository {
	return t.profilGuruRepo
}

func (t *fakeTx) ProfilSiswa() outuser.ProfilSiswaRepository {
	panic("not used in this test")
}

func (t *fakeTx) Commit() error {
	t.commitCalled = true
	if t.commitErr != nil {
		return t.commitErr
	}
	return nil
}

func (t *fakeTx) Rollback() error {
	t.rollbackCalled = true
	return nil
}

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
		txm                *fakeTxManager
		wantErr            error
		wantErrText        string
		wantBeginCalled    bool
		wantCommitCalled   bool
		wantRollbackCalled bool
		wantUpdateUser     bool
		wantUpdateProfil   bool
		wantEmailValue     string
	}{
		{
			name:  "Branch 1 -> semua patch berhasil",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
			wantEmailValue:     "guru@example.com",
		},
		{
			name:  "Branch 2 -> aktor bukan admin",
			cmd:   validCmd(),
			actor: nonAdminActor,
			txm:   &fakeTxManager{tx: &fakeTx{}},
			wantErr:            coreerror.ErrForbidden,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 3 -> id pengguna kosong",
			cmd:   func() user_service.UpdateGuruCmd { c := validCmd(); c.IdPengguna = 0; return c }(),
			actor: adminActor,
			txm:   &fakeTxManager{tx: &fakeTx{}},
			wantErrText:        "Id pengguna required",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 4 -> username kosong",
			cmd:   func() user_service.UpdateGuruCmd { c := validCmd(); c.Username = strPtr(" "); return c }(),
			actor: adminActor,
			txm:   &fakeTxManager{tx: &fakeTx{}},
			wantErrText:        "username cannot be empty",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 5 -> nama lengkap kosong",
			cmd:   func() user_service.UpdateGuruCmd { c := validCmd(); c.NamaLengkap = strPtr(" "); return c }(),
			actor: adminActor,
			txm:   &fakeTxManager{tx: &fakeTx{}},
			wantErrText:        "nama_lengkap cannot be empty",
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 6 -> email tidak valid",
			cmd:   func() user_service.UpdateGuruCmd { c := validCmd(); c.Email = strPtr("invalid"); return c }(),
			actor: adminActor,
			txm:   &fakeTxManager{tx: &fakeTx{}},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantUpdateUser:     false,
			wantUpdateProfil:   false,
		},
		{
			name:  "Branch 7 -> tidak ada field yang diupdate",
			cmd:   func() user_service.UpdateGuruCmd { return user_service.UpdateGuruCmd{IdPengguna: 10, Username: strPtr("guru")} }(),
			actor: adminActor,
			txm:   &fakeTxManager{tx: &fakeTx{}},
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
			txm: &fakeTxManager{
				beginErr: testErr,
				tx: &fakeTx{
					userRepo:       &fakeUserRepo{},
					profilGuruRepo: &fakeProfilGuruRepo{},
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
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{updateErr: testErr},
				profilGuruRepo: &fakeProfilGuruRepo{},
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
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{updateErr: testErr},
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
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
				commitErr:      testErr,
			}},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   true,
		},
		{
			name:  "Branch 12 -> hanya update profil",
			cmd:   func() user_service.UpdateGuruCmd { c := validCmd(); c.NamaLengkap = nil; c.Email = nil; c.NoHp = nil; c.Foto = nil; c.StatusAkun = nil; c.Role = nil; return c }(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     false,
			wantUpdateProfil:   true,
		},
		{
			name:  "Branch 13 -> hanya update pengguna",
			cmd:   func() user_service.UpdateGuruCmd { c := validCmd(); c.Nip = nil; c.Jabatan = nil; c.BidangStudi = nil; return c }(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			wantBeginCalled:    true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantUpdateUser:     true,
			wantUpdateProfil:   false,
			wantEmailValue:     "guru@example.com",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			service := user_service.NewUpdateGuruService(tc.txm)

			err := service.UpdateGuru(context.Background(), tc.cmd, tc.actor)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else if tc.wantErrText != "" {
				assert.EqualError(t, err, tc.wantErrText)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantBeginCalled, tc.txm.beginCalled)
			if tc.txm.tx != nil {
				assert.Equal(t, tc.wantCommitCalled, tc.txm.tx.commitCalled)
				assert.Equal(t, tc.wantRollbackCalled, tc.txm.tx.rollbackCalled)
				if tc.txm.tx.userRepo != nil {
					assert.Equal(t, tc.wantUpdateUser, tc.txm.tx.userRepo.updateCalled)
					if tc.wantEmailValue != "" && tc.txm.tx.userRepo.lastPatch.Email != nil {
						assert.Equal(t, tc.wantEmailValue, string(*tc.txm.tx.userRepo.lastPatch.Email))
					}
				}
				if tc.txm.tx.profilGuruRepo != nil {
					assert.Equal(t, tc.wantUpdateProfil, tc.txm.tx.profilGuruRepo.updateCalled)
				}
			}
		})
	}
}
