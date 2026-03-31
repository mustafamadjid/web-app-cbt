package kelas_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type FakeKelasRepo struct {
	DeleteErr error

	DeleteCalled   bool
	GotIDNamaKelas int
}

func (f *FakeKelasRepo) DeleteNamaKelas(_ context.Context, idNamaKelas int) error {
	f.DeleteCalled = true
	f.GotIDNamaKelas = idNamaKelas
	return f.DeleteErr
}

func (f *FakeKelasRepo) GetKelas(_ context.Context, _ query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	panic("not used in this test")
}

func (f *FakeKelasRepo) GetKelasById(_ context.Context, _ int, _ int) (kelas.KelasData, error) {
	panic("not used in this test")
}

func (f *FakeKelasRepo) CreateTingkatKelas(_ context.Context, _ int) error {
	panic("not used in this test")
}

func (f *FakeKelasRepo) CreateNamaKelas(_ context.Context, _ kelas.NamaKelas) error {
	panic("not used in this test")
}

func (f *FakeKelasRepo) UpdateNamaKelas(_ context.Context, _ int, _ updatepatch.NamaKelasPatch) error {
	panic("not used in this test")
}

func (f *FakeKelasRepo) ExistTingkatKelas(_ context.Context, _ int) (bool, error) {
	panic("not used in this test")
}

func (f *FakeKelasRepo) ExistNamaKelas(_ context.Context, _ string) (bool, error) {
	panic("not used in this test")
}
