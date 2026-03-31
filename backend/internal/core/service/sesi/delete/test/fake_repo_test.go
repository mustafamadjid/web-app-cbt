package sesi_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type FakeSesiRepo struct {
	DeleteSesiErr error

	DeleteSesiCalled bool
	GotDeleteId      int
}

func (f *FakeSesiRepo) GetSesi(_ context.Context, _ query.ListSesiFilter) ([]sesi.Sesi, error) {
	panic("not used in this test")
}

func (f *FakeSesiRepo) GetSesiById(_ context.Context, _ int) (sesi.Sesi, error) {
	panic("not used in this test")
}

func (f *FakeSesiRepo) GetSesiByKode(_ context.Context, _ string) (sesi.Sesi, error) {
	panic("not used in this test")
}

func (f *FakeSesiRepo) ExistByKodeSesi(_ context.Context, _ string) (bool, error) {
	panic("not used in this test")
}

func (f *FakeSesiRepo) CreateSesi(_ context.Context, _ sesi.Sesi) error {
	panic("not used in this test")
}

func (f *FakeSesiRepo) UpdateSesi(_ context.Context, _ int, _ updatepatch.UpdateSesiPatch) error {
	panic("not used in this test")
}

func (f *FakeSesiRepo) DeleteSesi(_ context.Context, idSesi int) error {
	f.DeleteSesiCalled = true
	f.GotDeleteId = idSesi
	return f.DeleteSesiErr
}
