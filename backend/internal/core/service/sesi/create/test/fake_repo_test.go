package sesi_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type FakeSesiRepo struct {
	ExistByKodeSesiRet bool
	ExistByKodeSesiErr error
	CreateSesiErr      error

	ExistByKodeSesiCalled bool
	CreateSesiCalled      bool
	GotKodeSesi           string
	GotCreateSesi         sesi.Sesi
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

func (f *FakeSesiRepo) ExistByKodeSesi(_ context.Context, kodeSesi string) (bool, error) {
	f.ExistByKodeSesiCalled = true
	f.GotKodeSesi = kodeSesi
	return f.ExistByKodeSesiRet, f.ExistByKodeSesiErr
}

func (f *FakeSesiRepo) CreateSesi(_ context.Context, s sesi.Sesi) error {
	f.CreateSesiCalled = true
	f.GotCreateSesi = s
	return f.CreateSesiErr
}

func (f *FakeSesiRepo) UpdateSesi(_ context.Context, _ int, _ updatepatch.UpdateSesiPatch) error {
	panic("not used in this test")
}

func (f *FakeSesiRepo) DeleteSesi(_ context.Context, _ int) error {
	panic("not used in this test")
}
