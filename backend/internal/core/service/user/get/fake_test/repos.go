package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type FakeGetGuruRepo struct {
	Items     []query.GuruListItem
	Err       error
	Called    bool
	GotFilter query.ListGuruFilter
}

func (f *FakeGetGuruRepo) GetListGuru(ctx context.Context, q query.ListGuruFilter) ([]query.GuruListItem, error) {
	f.Called = true
	f.GotFilter = q
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Items, nil
}

func (f *FakeGetGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	return user.DataGuru{}, nil
}

func (f *FakeGetGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	return false, nil
}

func (f *FakeGetGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	return 0, nil
}

func (f *FakeGetGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru updatepatch.ProfilGuru) error {
	return nil
}

type FakeGetSiswaRepo struct {
	Items     []query.SiswaListItem
	Err       error
	Called    bool
	GotFilter query.ListSiswaFilter
}

func (f *FakeGetSiswaRepo) GetListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem, error) {
	f.Called = true
	f.GotFilter = filter
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Items, nil
}

func (f *FakeGetSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	return user.DataSiswa{}, nil
}

func (f *FakeGetSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	return false, nil
}

func (f *FakeGetSiswaRepo) CreateProfilSiswa(ctx context.Context, g user.ProfilSiswa) (user.ID, error) {
	return 0, nil
}

func (f *FakeGetSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa updatepatch.ProfilSiswa) error {
	return nil
}

type FakeProfilGuruRepo struct {
	Result user.DataGuru
	Err    error
	Called bool
	GotID  user.ID
}

func (f *FakeProfilGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	f.Called = true
	f.GotID = id
	if f.Err != nil {
		return user.DataGuru{}, f.Err
	}
	return f.Result, nil
}

func (f *FakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	return false, nil
}

func (f *FakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	return 0, nil
}

func (f *FakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru updatepatch.ProfilGuru) error {
	return nil
}

type FakeProfilSiswaRepo struct {
	Result user.DataSiswa
	Err    error
	Called bool
	GotID  user.ID
}

func (f *FakeProfilSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	f.Called = true
	f.GotID = id
	if f.Err != nil {
		return user.DataSiswa{}, f.Err
	}
	return f.Result, nil
}

func (f *FakeProfilSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	return false, nil
}

func (f *FakeProfilSiswaRepo) CreateProfilSiswa(ctx context.Context, g user.ProfilSiswa) (user.ID, error) {
	return 0, nil
}

func (f *FakeProfilSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa updatepatch.ProfilSiswa) error {
	return nil
}
