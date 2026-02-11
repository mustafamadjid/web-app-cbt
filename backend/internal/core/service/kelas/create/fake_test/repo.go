package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type FakeKelasRepo struct {
	ExistTingkatKelasRet  bool
	ExistTingkatKelasErr  error
	CreateTingkatKelasErr error

	ExistNamaKelasRet  bool
	ExistNamaKelasErr  error
	CreateNamaKelasErr error

	ExistTingkatKelasCalled  bool
	CreateTingkatKelasCalled bool
	ExistNamaKelasCalled     bool
	CreateNamaKelasCalled    bool

	GotExistTingkatKelas  int
	GotCreateTingkatKelas int
	GotExistNamaKelas     string
	GotCreateNamaKelas    kelas.NamaKelas
}

func (f *FakeKelasRepo) GetKelas(_ context.Context, _ query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	return nil, nil
}

func (f *FakeKelasRepo) CreateTingkatKelas(_ context.Context, tingkatKelas int) error {
	f.CreateTingkatKelasCalled = true
	f.GotCreateTingkatKelas = tingkatKelas
	return f.CreateTingkatKelasErr
}

func (f *FakeKelasRepo) CreateNamaKelas(_ context.Context, namaKelas kelas.NamaKelas) error {
	f.CreateNamaKelasCalled = true
	f.GotCreateNamaKelas = namaKelas
	return f.CreateNamaKelasErr
}

func (f *FakeKelasRepo) ExistTingkatKelas(_ context.Context, tingkatKelas int) (bool, error) {
	f.ExistTingkatKelasCalled = true
	f.GotExistTingkatKelas = tingkatKelas
	return f.ExistTingkatKelasRet, f.ExistTingkatKelasErr
}

func (f *FakeKelasRepo) ExistNamaKelas(_ context.Context, namaKelas string) (bool, error) {
	f.ExistNamaKelasCalled = true
	f.GotExistNamaKelas = namaKelas
	return f.ExistNamaKelasRet, f.ExistNamaKelasErr
}


func (f *FakeKelasRepo)GetKelasById(ctx context.Context, idTingkatKelas int, idNamaKelas int)(kelas.KelasData, error){
	panic("not used in this test")
}

func (f *FakeKelasRepo)UpdateNamaKelas(ctx context.Context, idNamaKelas int, dataUpdate updatepatch.NamaKelasPatch)error{
	panic("not used in this test")
}

func(f *FakeKelasRepo)DeleteNamaKelas(ctx context.Context, idNamaKelas int) error{
	panic("not used in this test")
}