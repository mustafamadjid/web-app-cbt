package sesi_service_test

import (
	"context"
	"time"

	authsession "github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
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

type FakeSessionRepo struct {
	GetAllActiveSessionRet []authsession.SessionWithUser
	GetAllActiveSessionErr error

	GetAllActiveSessionCalled bool
}

func (f *FakeSessionRepo) GetSession(_ context.Context, _ string) (authsession.Session, error) {
	return authsession.Session{}, nil
}

func (f *FakeSessionRepo) GetSessionByUserId(_ context.Context, _ user.ID) (authsession.Session, error) {
	return authsession.Session{}, nil
}

func (f *FakeSessionRepo) GetAllActiveSession(_ context.Context) ([]authsession.SessionWithUser, error) {
	f.GetAllActiveSessionCalled = true
	return f.GetAllActiveSessionRet, f.GetAllActiveSessionErr
}

func (f *FakeSessionRepo) CreateSession(_ context.Context, _ user.ID, _ user.Role, _ time.Time) (string, error) {
	return "", nil
}

func (f *FakeSessionRepo) RevokeSession(_ context.Context, _ string) error {
	return nil
}

func (f *FakeSessionRepo) RevokeSessionAllbyUser(_ context.Context, _ user.ID) error {
	return nil
}

func (f *FakeSessionRepo) RevokeExpiredSessions(_ context.Context, _ user.ID) (bool, error) {
	return false, nil
}

func (f *FakeSessionRepo) HasActiveSession(_ context.Context, _ user.ID) (bool, error) {
	return false, nil
}
