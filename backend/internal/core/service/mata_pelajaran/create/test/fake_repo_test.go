package matapelajaran_service_test

import (
	"context"

	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type FakeMapelRepo struct {
	ExistKodeMapelRet bool
	ExistKodeMapelErr error
	CreateMapelErr    error

	ExistKodeMapelCalled bool
	CreateMapelCalled    bool
	GotKodeMapel         string
	GotCreateMapel       matapelajaran.MataPelajaran
}

func (f *FakeMapelRepo) GetMapel(_ context.Context, _ query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	panic("not used in this test")
}

func (f *FakeMapelRepo) GetMapelById(_ context.Context, _ int) (matapelajaran.MataPelajaran, error) {
	panic("not used in this test")
}

func (f *FakeMapelRepo) CreateMapel(_ context.Context, mapel matapelajaran.MataPelajaran) error {
	f.CreateMapelCalled = true
	f.GotCreateMapel = mapel
	return f.CreateMapelErr
}

func (f *FakeMapelRepo) UpdateMapel(_ context.Context, _ int, _ updatepatch.UpdateMapelPatch) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) DeleteMapel(_ context.Context, _ int) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) ExistKodeMapel(_ context.Context, kodeMapel string) (bool, error) {
	f.ExistKodeMapelCalled = true
	f.GotKodeMapel = kodeMapel
	return f.ExistKodeMapelRet, f.ExistKodeMapelErr
}
