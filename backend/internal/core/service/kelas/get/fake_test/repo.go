package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type FakeKelasRepo struct {
	Items     []kelas.FullKelasData
	Err       error
	Called    bool
	GotFilter query.ListKelasFilter
}

func (f *FakeKelasRepo) GetKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	f.Called = true
	f.GotFilter = filter
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Items, nil
}

// Not used
func (f *FakeKelasRepo) CreateTingkatKelas(ctx context.Context, tingkatKelas int)error {
	panic("not used in this test")
}

// Not used
func (f *FakeKelasRepo) CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas) error {
	panic("not used in this test")
}

// Not used
func (f *FakeKelasRepo) ExistTingkatKelas(ctx context.Context, tingkatKelas int) (bool, error) {
	panic("not used in this test")
}

func (f *FakeKelasRepo)ExistNamaKelas(ctx context.Context,namaKelas string)(bool, error){
	panic("not used in this test")
}

func (f *FakeKelasRepo)GetKelasById(ctx context.Context, idTingkatKelas int, idNamaKelas int)(kelas.KelasData, error){
	panic("not used in this test")
}