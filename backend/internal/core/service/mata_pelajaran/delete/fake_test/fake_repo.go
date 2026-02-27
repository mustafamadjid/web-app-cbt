package fake_test

import (
	"context"

	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type FakeMapelRepo struct {
	DeleteMapelErr error

	DeleteMapelCalled bool
	GotDeleteId       int
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

func (f *FakeMapelRepo) UpdateMapel(_ context.Context, _ int, _ updatepatch.UpdateMapelPatch) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) DeleteMapel(_ context.Context, idMapel int) error {
	f.DeleteMapelCalled = true
	f.GotDeleteId = idMapel
	return f.DeleteMapelErr
}

func (f *FakeMapelRepo) ExistKodeMapel(_ context.Context, _ string) (bool, error) {
	panic("not used in this test")
}
