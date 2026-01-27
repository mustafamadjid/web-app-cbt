package user_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	txport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	userport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user"
	"github.com/stretchr/testify/assert"
)

// ===== minimal fake user repo =====

type fakeUserRepo struct {
	existByUsername bool
	userExistErr    error

	createID  user.ID
	createErr error

	userExistCalled bool
	createCalled    bool
}

func (r *fakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	r.userExistCalled = true
	if r.userExistErr != nil {
		return false, r.userExistErr
	}
	return r.existByUsername, nil
}

func (r *fakeUserRepo) CreateUser(ctx context.Context, p user.Pengguna) (user.ID, error) {
	r.createCalled = true
	if r.createErr != nil {
		return 0, r.createErr
	}
	return r.createID, nil
}

// ---- methods below are NOT used by this test; keep minimal ----

func (r *fakeUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	panic("not used in this test")
}

func (r *fakeUserRepo) UpdateUser(ctx context.Context, p user.Pengguna) (user.ID, error) {
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
	existsNip bool
	existErr  error

	createID  user.ID
	createErr error

	existCalled  bool
	createCalled bool
}

func (r *fakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	r.existCalled = true
	if r.existErr != nil {
		return false, r.existErr
	}
	return r.existsNip, nil
}

func (r *fakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	r.createCalled = true
	if r.createErr != nil {
		return 0, r.createErr
	}
	return r.createID, nil
}

func (r *fakeProfilGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.ProfilGuru, error) {
	panic("not used in this test")
}

func (r *fakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, profilGuru user.ProfilGuru) (user.ID, error) {
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

func (t *fakeTx) Pengguna() userport.UserRepository {
	return t.userRepo
}

func (t *fakeTx) ProfilGuru() userport.ProfilGuruRepository {
	return t.profilGuruRepo
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

// ===== fake hasher =====

type fakeHasher struct {
	hash      string
	err       error
	called    bool
	lastPlain string
}

func (fakeHash *fakeHasher) ComparePaswordAndHashed(hash string, plain string) bool {
	return hash == plain
}

func (fakeHash *fakeHasher) GenerateHash(plain string) (string, error) {
	fakeHash.called = true
	fakeHash.lastPlain = plain
	if fakeHash.err != nil {
		return "", fakeHash.err
	}
	return fakeHash.hash, nil
}

func TestCreateGuru(t *testing.T) {
	validCmd := func() user_service.CreateGuruCmd {
		return user_service.CreateGuruCmd{
			Username:     "guruuser",
			Email:        "guru@example.com",
			Password:     "password",
			NamaLengkap:  "Guru Test",
			JenisKelamin: "L",
			NoHp:         "08123456789",
			Foto:         "foto.png",
			Nip:          "123456789012345678",
			Jabatan:      "Kepala",
			BidangStudi:  "Matematika",
		}
	}

	adminActor := user.Actor{IdPengguna: 1, Role: user.ADMIN}

	testErr := errors.New("test error")

	cases := []struct {
		name               string
		cmd                user_service.CreateGuruCmd
		actor              user.Actor
		txm                *fakeTxManager
		hasher             *fakeHasher
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
			name:  "semua validasi lolos",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{createID: 10},
				profilGuruRepo: &fakeProfilGuruRepo{createID: 20},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
			wantResult:         user_service.CreateGuruRes{IdPengguna: 10, IdProfilGuru: 20},
		},
		{
			name:  "username taken",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{existByUsername: true},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            coreerror.ErrUsernameTaken,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "NIP taken",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{existByUsername: false},
				profilGuruRepo: &fakeProfilGuruRepo{existsNip: true},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            coreerror.ErrNipTaken,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "Create pengguna gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{createErr: testErr},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   false,
		},
		{
			name:  "Create profil guru gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{createID: 10},
				profilGuruRepo: &fakeProfilGuruRepo{createErr: testErr},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
		},
		{
			name:  "Commit gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{createID: 10},
				profilGuruRepo: &fakeProfilGuruRepo{createID: 20},
				commitErr:      testErr,
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   true,
			wantRollbackCalled: true,
			wantCreateUser:     true,
			wantCreateProfil:   true,
		},
		{
			name:  "invalid email",
			cmd:   func() user_service.CreateGuruCmd { c := validCmd(); c.Email = "not-an-email"; return c }(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            user.ErrInvalidEmail,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "invalid nip",
			cmd:   func() user_service.CreateGuruCmd { c := validCmd(); c.Nip = "123"; return c }(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            user.ErrInvalidNIP,
			wantBeginCalled:    false,
			wantHasherCalled:   false,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "hash password gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{tx: &fakeTx{
				userRepo:       &fakeUserRepo{},
				profilGuruRepo: &fakeProfilGuruRepo{},
			}},
			hasher:             &fakeHasher{err: testErr},
			wantErr:            testErr,
			wantBeginCalled:    false,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
			wantCreateUser:     false,
			wantCreateProfil:   false,
		},
		{
			name:  "begin tx gagal",
			cmd:   validCmd(),
			actor: adminActor,
			txm: &fakeTxManager{
				beginErr: testErr,
				tx: &fakeTx{
					userRepo:       &fakeUserRepo{},
					profilGuruRepo: &fakeProfilGuruRepo{},
				},
			},
			hasher:             &fakeHasher{hash: "hashed"},
			wantErr:            testErr,
			wantBeginCalled:    true,
			wantHasherCalled:   true,
			wantCommitCalled:   false,
			wantRollbackCalled: false,
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
			assert.Equal(t, tc.wantBeginCalled, tc.txm.beginCalled)
			assert.Equal(t, tc.wantHasherCalled, tc.hasher.called)

			if tc.txm.tx != nil {
				assert.Equal(t, tc.wantCommitCalled, tc.txm.tx.commitCalled)
				assert.Equal(t, tc.wantRollbackCalled, tc.txm.tx.rollbackCalled)
				if tc.txm.tx.userRepo != nil {
					assert.Equal(t, tc.wantCreateUser, tc.txm.tx.userRepo.createCalled)
				}
				if tc.txm.tx.profilGuruRepo != nil {
					assert.Equal(t, tc.wantCreateProfil, tc.txm.tx.profilGuruRepo.createCalled)
				}
			}
		})
	}
}

var _ out.PasswordHasher = (*fakeHasher)(nil)
