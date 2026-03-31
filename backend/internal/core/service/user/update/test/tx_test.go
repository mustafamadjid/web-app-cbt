package user_service_test

import (
	"context"

	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	txport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

// ===== fake tx & tx manager =====

type FakeTx struct {
	UserRepo        *FakeUserRepo
	ProfilGuruRepo  *FakeProfilGuruRepo
	ProfilSiswaRepo *FakeProfilSiswaRepo
	kelasRepo       *kelas_repo.KelasRepository

	CommitCalled   bool
	RollbackCalled bool
	CommitErr      error
}

type FakeTxManager struct {
	Tx          *FakeTx
	BeginErr    error
	BeginCalled bool
}

func (m *FakeTxManager) Begin(ctx context.Context) (txport.Tx, error) {
	m.BeginCalled = true
	if m.BeginErr != nil {
		return nil, m.BeginErr
	}
	return m.Tx, nil
}

func (t *FakeTx) Pengguna() outuser.UserRepository {
	return t.UserRepo
}

func (t *FakeTx) ProfilGuru() outuser.ProfilGuruRepository {
	return t.ProfilGuruRepo
}

func (t *FakeTx) ProfilSiswa() outuser.ProfilSiswaRepository {
	return t.ProfilSiswaRepo
}

func (t *FakeTx) Kelas() kelas_repo.KelasRepository {
	return *t.kelasRepo
}

func (t *FakeTx) Commit() error {
	t.CommitCalled = true
	if t.CommitErr != nil {
		return t.CommitErr
	}
	return nil
}

func (t *FakeTx) Rollback() error {
	t.RollbackCalled = true
	return nil
}
