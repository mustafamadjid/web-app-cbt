package sesi_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type FakeSesiRepo struct {
	GetSesiRet       []sesi.Sesi
	GetSesiErr       error
	GetSesiByIdRet   sesi.Sesi
	GetSesiByIdErr   error
	GetSesiByKodeRet sesi.Sesi
	GetSesiByKodeErr error

	GetSesiCalled       bool
	GetSesiByIdCalled   bool
	GetSesiByKodeCalled bool
	GotFilter           query.ListSesiFilter
	GotId               int
	GotKode             string
}

func (f *FakeSesiRepo) GetSesi(_ context.Context, filter query.ListSesiFilter) ([]sesi.Sesi, error) {
	f.GetSesiCalled = true
	f.GotFilter = filter
	return f.GetSesiRet, f.GetSesiErr
}

func (f *FakeSesiRepo) GetSesiById(_ context.Context, idSesi int) (sesi.Sesi, error) {
	f.GetSesiByIdCalled = true
	f.GotId = idSesi
	return f.GetSesiByIdRet, f.GetSesiByIdErr
}

func (f *FakeSesiRepo) GetSesiByKode(_ context.Context, kodeSesi string) (sesi.Sesi, error) {
	f.GetSesiByKodeCalled = true
	f.GotKode = kodeSesi
	return f.GetSesiByKodeRet, f.GetSesiByKodeErr
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

func (f *FakeSesiRepo) DeleteSesi(_ context.Context, _ int) error {
	panic("not used in this test")
}
