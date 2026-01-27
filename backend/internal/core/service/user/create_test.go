package user_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	txport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	userport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

// ===== minimal fake user repo =====

type fakeUserRepo struct {
	existByUsername bool

	createID  user.ID
	createErr error

	// optional tracking
	lastUsername string
	lastCreated  user.Pengguna
}

func (r *fakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	r.lastUsername = username
	return r.existByUsername, nil
}

func (r *fakeUserRepo) CreateUser(ctx context.Context, p user.Pengguna) (user.ID, error) {
	r.lastCreated = p
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

func (r *fakeProfilGuruRepo)FindProfilGuruByID(ctx context.Context, id user.ID) (user.ProfilGuru, error) {
	panic("not used in this test")
}

func (r *fakeProfilGuruRepo)UpdateProfilGuru(ctx context.Context, profilGuru user.ProfilGuru) (user.ID, error) {
	panic("not used in this test")
}



type fakeProfilGuruRepo struct {
	existsNip bool
	createID  user.ID
	createErr error
}

func (r *fakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip string) (bool, error) {
	return r.existsNip, nil
}

func (r *fakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	if r.createErr != nil {
		return 0, r.createErr
	}
	return r.createID, nil
}



type fakeTx struct {
	userRepo      *fakeUserRepo
	profilGuruRepo *fakeProfilGuruRepo

	commitCalled   bool
	rollbackCalled bool
}

type fakeTxManager struct {
	tx *fakeTx
}

func (m *fakeTxManager) Begin(ctx context.Context) (txport.Tx, error) {
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
	return nil
}

func (t *fakeTx) Rollback() error {
	t.rollbackCalled = true
	return nil
}
