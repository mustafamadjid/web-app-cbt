package fake_test

import (
	"context"

	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type FakeMapelRepo struct {
	ExistKodeMapelRet bool
	ExistKodeMapelErr error
	UpdateMapelErr    error

	ExistKodeMapelCalled bool
	UpdateMapelCalled    bool
	GotKodeMapel         string
	GotUpdateId          int
	GotUpdatePatch       updatepatch.UpdateMapelPatch
}

func (f *FakeMapelRepo) GetMapel(_ context.Context, _ query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	panic("not used in this test")
}

func (f *FakeMapelRepo) GetMapelById(_ context.Context, _ int) (matapelajaran.MataPelajaran, error) {
	panic("not used in this test")
}

func (f *FakeMapelRepo) CreateMapel(_ context.Context, _ matapelajaran.MataPelajaran) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) UpdateMapel(_ context.Context, idMapel int, mapel updatepatch.UpdateMapelPatch) error {
	f.UpdateMapelCalled = true
	f.GotUpdateId = idMapel
	f.GotUpdatePatch = mapel
	return f.UpdateMapelErr
}

func (f *FakeMapelRepo) DeleteMapel(_ context.Context, _ int) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) ExistKodeMapel(_ context.Context, kodeMapel string) (bool, error) {
	f.ExistKodeMapelCalled = true
	f.GotKodeMapel = kodeMapel
	return f.ExistKodeMapelRet, f.ExistKodeMapelErr
}
