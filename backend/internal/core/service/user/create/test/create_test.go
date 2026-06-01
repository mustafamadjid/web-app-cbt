package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	"github.com/stretchr/testify/assert"
)

func strPtr(v string) *string { return &v }

func validCreateGuruCmd() user_service.CreateGuruCmd {
	return user_service.CreateGuruCmd{
		Username:     "guruuser",
		Email:        strPtr("guru@example.com"),
		Password:     "password",
		NamaLengkap:  "Guru Test",
		JenisKelamin: "L",
		NoHp:         strPtr("08123456789"),
		Foto:         "foto.png",
		Nip:          "123456789012345678",
		Jabatan:      "Kepala",
		BidangStudi:  "Matematika",
	}
}

func validCreateSiswaCmd() user_service.CreateSiswaCmd {
	return user_service.CreateSiswaCmd{
		Username:     "siswauser",
		Email:        strPtr("siswa@example.com"),
		Password:     "password",
		NamaLengkap:  "Siswa Test",
		JenisKelamin: "L",
		NoHp:         strPtr("08123456789"),
		Foto:         "foto.png",
		IdNamaKelas:  2,
		Nisn:         "1234567890",
		NoAbsen:      1,
		Angkatan:     2021,
		TempatLahir:  "Bandung",
		TanggalLahir: time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

func adminActor() user.Actor {
	return user.Actor{IdPengguna: 1, Role: user.ADMIN}
}

func guruActor() user.Actor {
	return user.Actor{IdPengguna: 2, Role: user.GURU}
}

func siswaActor() user.Actor {
	return user.Actor{IdPengguna: 3, Role: user.SISWA}
}

func newGuruTx(userRepo *FakeUserRepo, guruRepo *FakeProfilGuruRepo) *FakeTxManager {
	return &FakeTxManager{Tx: &FakeTx{
		UserRepo:       userRepo,
		ProfilGuruRepo: guruRepo,
	}}
}

func newSiswaTx(userRepo *FakeUserRepo, siswaRepo *FakeProfilSiswaRepo) *FakeTxManager {
	return &FakeTxManager{Tx: &FakeTx{
		UserRepo:        userRepo,
		ProfilSiswaRepo: siswaRepo,
	}}
}

func assertBeforeHash(t *testing.T, txm *FakeTxManager, hasher *FakeHasher) {
	t.Helper()

	assert.False(t, hasher.Called)
	assert.False(t, txm.BeginCalled)
}

func assertHashBeforeTx(t *testing.T, txm *FakeTxManager, hasher *FakeHasher) {
	t.Helper()

	assert.True(t, hasher.Called)
	assert.False(t, txm.BeginCalled)
}

func TestCreateGuruBasisPath(t *testing.T) {
	t.Parallel()

	t.Run("P1 -> aktor tidak memiliki hak untuk membuat data guru", func(t *testing.T) {
		t.Parallel()
		txm := newGuruTx(&FakeUserRepo{}, &FakeProfilGuruRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), validCreateGuruCmd(), guruActor())

		assert.ErrorIs(t, err, coreerror.ErrForbidden)
		assert.Equal(t, user_service.CreateGuruRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P2 -> username tidak sesuai aturan panjang karakter", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateGuruCmd()
		cmd.Username = "guru"
		txm := newGuruTx(&FakeUserRepo{}, &FakeProfilGuruRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, coreerror.ErrUsernameLengthInvalid)
		assert.Equal(t, user_service.CreateGuruRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P3 -> NIP diisi tetapi format NIP tidak valid", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateGuruCmd()
		cmd.Nip = "123"
		txm := newGuruTx(&FakeUserRepo{}, &FakeProfilGuruRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, user.ErrInvalidNIP)
		assert.Equal(t, user_service.CreateGuruRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P4 -> email tidak valid setelah NIP berhasil divalidasi", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateGuruCmd()
		cmd.Email = strPtr("not-an-email")
		txm := newGuruTx(&FakeUserRepo{}, &FakeProfilGuruRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, user.ErrInvalidEmail)
		assert.Equal(t, user_service.CreateGuruRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P5 -> NIP dash tetapi hashing password gagal", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		cmd := validCreateGuruCmd()
		cmd.Nip = "-"
		txm := newGuruTx(&FakeUserRepo{}, &FakeProfilGuruRepo{})
		hasher := &FakeHasher{Err: testErr}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, testErr)
		assert.Equal(t, user_service.CreateGuruRes{}, res)
		assertHashBeforeTx(t, txm, hasher)
	})

	t.Run("P6 -> pembuatan guru berhasil dengan NIP dash", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateGuruCmd()
		cmd.Nip = "-"
		userRepo := &FakeUserRepo{CreateID: 10}
		guruRepo := &FakeProfilGuruRepo{CreateID: 20}
		txm := newGuruTx(userRepo, guruRepo)
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), cmd, adminActor())

		assert.NoError(t, err)
		assert.Equal(t, user_service.CreateGuruRes{IdPengguna: 10, IdProfilGuru: 20}, res)
		assert.True(t, hasher.Called)
		assert.True(t, txm.BeginCalled)
		assert.False(t, guruRepo.ExistCalled)
		assert.True(t, userRepo.CreateCalled)
		assert.True(t, guruRepo.CreateCalled)
		assert.True(t, txm.Tx.CommitCalled)
	})

	t.Run("P7 -> pembuatan guru berhasil dengan NIP valid", func(t *testing.T) {
		t.Parallel()
		userRepo := &FakeUserRepo{CreateID: 11}
		guruRepo := &FakeProfilGuruRepo{CreateID: 21}
		txm := newGuruTx(userRepo, guruRepo)
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateGuruService(txm, hasher)

		res, err := service.CreateGuru(context.Background(), validCreateGuruCmd(), adminActor())

		assert.NoError(t, err)
		assert.Equal(t, user_service.CreateGuruRes{IdPengguna: 11, IdProfilGuru: 21}, res)
		assert.True(t, hasher.Called)
		assert.True(t, txm.BeginCalled)
		assert.True(t, guruRepo.ExistCalled)
		assert.True(t, userRepo.CreateCalled)
		assert.True(t, guruRepo.CreateCalled)
		assert.True(t, txm.Tx.CommitCalled)
	})
}

func TestCreateSiswaBasisPath(t *testing.T) {
	t.Parallel()

	t.Run("P1 -> aktor tidak memiliki hak untuk membuat data siswa", func(t *testing.T) {
		t.Parallel()
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), validCreateSiswaCmd(), siswaActor())

		assert.ErrorIs(t, err, coreerror.ErrForbidden)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P2 -> username tidak sesuai aturan panjang karakter", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateSiswaCmd()
		cmd.Username = "sisw"
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, coreerror.ErrUsernameLengthInvalid)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P3 -> NISN diisi tetapi format NISN tidak valid", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateSiswaCmd()
		cmd.Nisn = "123"
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, user.ErrInvalidNISN)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P4 -> email tidak valid setelah NISN berhasil divalidasi", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateSiswaCmd()
		cmd.Email = strPtr("not-an-email")
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, user.ErrInvalidEmail)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P5 -> nomor absen tidak valid", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateSiswaCmd()
		cmd.NoAbsen = 0
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, user.ErrInvalidAbsen)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P6 -> angkatan tidak valid", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateSiswaCmd()
		cmd.Angkatan = 2000
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, user.ErrInvalidAngkatan)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertBeforeHash(t, txm, hasher)
	})

	t.Run("P7 -> proses hashing password gagal", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		cmd := validCreateSiswaCmd()
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		hasher := &FakeHasher{Err: testErr}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.ErrorIs(t, err, testErr)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assertHashBeforeTx(t, txm, hasher)
	})

	t.Run("P8 -> transaksi pembuatan siswa gagal", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		txm := newSiswaTx(&FakeUserRepo{}, &FakeProfilSiswaRepo{})
		txm.BeginErr = testErr
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), validCreateSiswaCmd(), adminActor())

		assert.ErrorIs(t, err, testErr)
		assert.Equal(t, user_service.CreateSiswaRes{}, res)
		assert.True(t, hasher.Called)
		assert.True(t, txm.BeginCalled)
	})

	t.Run("P9 -> pembuatan siswa berhasil dengan NISN dash", func(t *testing.T) {
		t.Parallel()
		cmd := validCreateSiswaCmd()
		cmd.Nisn = "-"
		userRepo := &FakeUserRepo{CreateID: 10}
		siswaRepo := &FakeProfilSiswaRepo{CreateID: 20}
		txm := newSiswaTx(userRepo, siswaRepo)
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), cmd, adminActor())

		assert.NoError(t, err)
		assert.Equal(t, user_service.CreateSiswaRes{IdPengguna: 10, IdProfilSiswa: 20}, res)
		assert.True(t, hasher.Called)
		assert.True(t, txm.BeginCalled)
		assert.False(t, siswaRepo.ExistCalled)
		assert.True(t, userRepo.CreateCalled)
		assert.True(t, siswaRepo.CreateCalled)
		assert.True(t, txm.Tx.CommitCalled)
	})

	t.Run("P10 -> pembuatan siswa berhasil dengan NISN valid", func(t *testing.T) {
		t.Parallel()
		userRepo := &FakeUserRepo{CreateID: 11}
		siswaRepo := &FakeProfilSiswaRepo{CreateID: 21}
		txm := newSiswaTx(userRepo, siswaRepo)
		hasher := &FakeHasher{Hash: "hashed"}
		service := user_service.NewCreateSiswaService(txm, hasher)

		res, err := service.CreateSiswa(context.Background(), validCreateSiswaCmd(), adminActor())

		assert.NoError(t, err)
		assert.Equal(t, user_service.CreateSiswaRes{IdPengguna: 11, IdProfilSiswa: 21}, res)
		assert.True(t, hasher.Called)
		assert.True(t, txm.BeginCalled)
		assert.True(t, siswaRepo.ExistCalled)
		assert.True(t, userRepo.CreateCalled)
		assert.True(t, siswaRepo.CreateCalled)
		assert.True(t, txm.Tx.CommitCalled)
	})
}

var _ out.PasswordHasher = (*FakeHasher)(nil)
