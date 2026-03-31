package matapelajaran_service_test

import (
	"context"

	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type FakeMapelRepo struct {
	GetMapelRet     []matapelajaran.MataPelajaran
	GetMapelErr     error
	GetMapelByIdRet matapelajaran.MataPelajaran
	GetMapelByIdErr error

	GetMapelCalled     bool
	GetMapelByIdCalled bool
	GotFilter          query.ListMapelFilter
	GotId              int
}

func (f *FakeMapelRepo) GetMapel(_ context.Context, filter query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	f.GetMapelCalled = true
	f.GotFilter = filter
	return f.GetMapelRet, f.GetMapelErr
}

func (f *FakeMapelRepo) GetMapelById(_ context.Context, idMapel int) (matapelajaran.MataPelajaran, error) {
	f.GetMapelByIdCalled = true
	f.GotId = idMapel
	return f.GetMapelByIdRet, f.GetMapelByIdErr
}

func (f *FakeMapelRepo) CreateMapel(_ context.Context, _ matapelajaran.MataPelajaran) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) UpdateMapel(_ context.Context, _ int, _ updatepatch.UpdateMapelPatch) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) DeleteMapel(_ context.Context, _ int) error {
	panic("not used in this test")
}

func (f *FakeMapelRepo) ExistKodeMapel(_ context.Context, _ string) (bool, error) {
	panic("not used in this test")
}
