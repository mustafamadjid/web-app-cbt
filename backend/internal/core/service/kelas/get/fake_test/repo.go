package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type FakeKelasRepo struct {
	Items     []kelas.FullKelasData
	Err       error
	Called    bool
	GotFilter query.ListKelasFilter

	KelasByID         kelas.KelasData
	KelasByIDErr      error
	GetByIDCalled     bool
	GotIDTingkatKelas int
	GotIDNamaKelas    int
}

func (f *FakeKelasRepo) GetKelas(_ context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	f.Called = true
	f.GotFilter = filter
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Items, nil
}

func (f *FakeKelasRepo) GetKelasById(_ context.Context, idTingkatKelas int, idNamaKelas int) (kelas.KelasData, error) {
	f.GetByIDCalled = true
	f.GotIDTingkatKelas = idTingkatKelas
	f.GotIDNamaKelas = idNamaKelas
	if f.KelasByIDErr != nil {
		return kelas.KelasData{}, f.KelasByIDErr
	}
	return f.KelasByID, nil
}

// Not used
func (f *FakeKelasRepo) CreateTingkatKelas(_ context.Context, _ int) error {
	panic("not used in this test")
}

// Not used
func (f *FakeKelasRepo) CreateNamaKelas(_ context.Context, _ kelas.NamaKelas) error {
	panic("not used in this test")
}

// Not used
func (f *FakeKelasRepo) ExistTingkatKelas(_ context.Context, _ int) (bool, error) {
	panic("not used in this test")
}

func (f *FakeKelasRepo) ExistNamaKelas(_ context.Context, _ string) (bool, error) {
	panic("not used in this test")
}

func (f *FakeKelasRepo) UpdateNamaKelas(_ context.Context, _ int, _ updatepatch.NamaKelasPatch) error {
	panic("not used in this test")
}

func (f *FakeKelasRepo) DeleteNamaKelas(_ context.Context, _ int) error {
	panic("not used in this test")
}
