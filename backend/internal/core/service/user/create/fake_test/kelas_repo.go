package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

// ===== Minimal fake kelas repo =====

type FakeKelasRepo struct {
	ExistsKodeKelas bool
	ExistErr        error

	CreateID  kelas.ID
	CreateErr error

	ExistCalled  bool
	CreateCalled bool
}

func (r *FakeKelasRepo) GetKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	panic("not used in this")
}

func (r *FakeKelasRepo) CreateTingkatKelas(ctx context.Context, tingkatKelas kelas.TingkatKelas) (kelas.ID, error) {
	panic("not used in this test")
}

func (r *FakeKelasRepo) CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas) (kelas.ID, error) {
	panic("not used in this test")
}
