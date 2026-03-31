package aktivitas_user_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
)

type FakeAktivitasUserRepo struct {
	CreateErr     error
	CreateCalled  bool
	CreatedRecord aktivitas_user.AktivitasUser

	GetErr    error
	GetCalled bool
	GetData   []aktivitas_user.AktivitasUser
}

func (r *FakeAktivitasUserRepo) CreateAktivitasUser(ctx context.Context, aktivitasUser aktivitas_user.AktivitasUser) error {
	r.CreateCalled = true
	r.CreatedRecord = aktivitasUser
	return r.CreateErr
}

func (r *FakeAktivitasUserRepo) GetAktivitasUser(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	r.GetCalled = true
	if r.GetErr != nil {
		return nil, r.GetErr
	}
	return r.GetData, nil
}
