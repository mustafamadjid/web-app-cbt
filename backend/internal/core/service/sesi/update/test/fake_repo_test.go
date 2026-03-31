package sesi_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type FakeSesiRepo struct {
	UpdateSesiErr error

	UpdateSesiCalled bool
	GotUpdateId      int
	GotUpdatePatch   updatepatch.UpdateSesiPatch
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

func (f *FakeSesiRepo) UpdateSesi(_ context.Context, idSesi int, patch updatepatch.UpdateSesiPatch) error {
	f.UpdateSesiCalled = true
	f.GotUpdateId = idSesi
	f.GotUpdatePatch = patch
	return f.UpdateSesiErr
}

func (f *FakeSesiRepo) DeleteSesi(_ context.Context, _ int) error {
	panic("not used in this test")
}
