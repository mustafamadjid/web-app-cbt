package user_service

import (
	"context"
	"strings"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

type CreateGuruCmd struct {
	Username     string
	Email        string
	Password     string
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string

	Nip         string
	Jabatan     string
	BidangStudi string
}

type CreateGuruRes struct {
	IdPengguna   user.ID
	IdProfilGuru user.ID
}

type CreateSiswaCmd struct {
	Username     string
	Email        string
	Password     string
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string

	IdTingkatKelas user.ID
	IdNamaKelas    user.ID
	Nisn           string
	NoAbsen        int
	Angkatan       int
	TempatLahir    string
	TanggalLahir   time.Time
}

type CreateSiswaRes struct {
	IdPengguna    user.ID
	IdProfilSiswa user.ID
}

type CreateTx struct {
	txm    txout.TxManager
	hasher out.PasswordHasher
}

func NewCreateGuruService(txm txout.TxManager, hasher out.PasswordHasher) *CreateTx {
	return &CreateTx{txm: txm, hasher: hasher}
}

func NewCreateSiswaService(txm txout.TxManager, hasher out.PasswordHasher) *CreateTx {
	return &CreateTx{txm: txm, hasher: hasher}
}

func (uc *CreateTx) CreateGuru(ctx context.Context, cmd CreateGuruCmd, actor user.Actor) (CreateGuruRes, error) {
	// Cek role
	if actor.Role != user.ADMIN {
		return CreateGuruRes{}, coreerror.ErrForbidden
	}

	// Normalisasi
	cmd.Username = strings.TrimSpace(cmd.Username)
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.Password = strings.TrimSpace(cmd.Password)
	cmd.NamaLengkap = strings.TrimSpace(cmd.NamaLengkap)
	cmd.JenisKelamin = strings.TrimSpace(cmd.JenisKelamin)
	cmd.NoHp = strings.TrimSpace(cmd.NoHp)
	cmd.Nip = strings.TrimSpace(cmd.Nip)
	cmd.Foto = strings.TrimSpace(cmd.Foto)
	cmd.Jabatan = strings.TrimSpace(cmd.Jabatan)
	cmd.BidangStudi = strings.TrimSpace(cmd.BidangStudi)

	// Validasi
	emailValidated, error := user.CheckNewEmail(cmd.Email)
	if error != nil {
		return CreateGuruRes{}, error
	}

	nipValidated, error := user.CheckNewNip(cmd.Nip)
	if error != nil {
		return CreateGuruRes{}, error
	}

	// Hash password
	hashedPassword, error := uc.hasher.GenerateHash(cmd.Password)
	if error != nil {
		return CreateGuruRes{}, error
	}

	// --- Transaksi ----
	tx, err := uc.txm.Begin(ctx)
	if err != nil {
		return CreateGuruRes{}, err
	}

	defer func() { _ = tx.Rollback() }()

	existUsername, error := tx.Pengguna().UserExistByUsername(ctx, cmd.Username)
	if error != nil {
		return CreateGuruRes{}, error
	}
	if existUsername {
		return CreateGuruRes{}, coreerror.ErrUsernameTaken
	}

	existNip, error := tx.ProfilGuru().ExistByNIP(ctx, nipValidated)
	if error != nil {
		return CreateGuruRes{}, error
	}
	if existNip {
		return CreateGuruRes{}, coreerror.ErrNipTaken
	}

	userData := user.Pengguna{
		Username:       cmd.Username,
		Email:          emailValidated,
		PasswordHashed: hashedPassword,
		NamaLengkap:    cmd.NamaLengkap,
		JenisKelamin:   cmd.JenisKelamin,
		NoHp:           cmd.NoHp,
		Foto:           cmd.Foto,
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	}

	idPengguna, error := tx.Pengguna().CreateUser(ctx, userData)
	if error != nil {
		return CreateGuruRes{}, error
	}

	profilGuruData := user.ProfilGuru{
		IdPengguna:  idPengguna,
		Nip:         nipValidated,
		Jabatan:     cmd.Jabatan,
		BidangStudi: cmd.BidangStudi,
	}

	idProfilGuru, error := tx.ProfilGuru().CreateProfilGuru(ctx, profilGuruData)
	if error != nil {
		return CreateGuruRes{}, error
	}

	if error := tx.Commit(); error != nil {
		return CreateGuruRes{}, error
	}

	return CreateGuruRes{
		IdPengguna:   idPengguna,
		IdProfilGuru: idProfilGuru,
	}, nil

}

func (uc *CreateTx) CreateSiswa(ctx context.Context, cmd CreateSiswaCmd, actor user.Actor) (CreateSiswaRes, error) {
	if actor.Role != user.ADMIN {
		return CreateSiswaRes{}, coreerror.ErrForbidden
	}

	cmd.Username = strings.TrimSpace(cmd.Username)
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.Password = strings.TrimSpace(cmd.Password)
	cmd.NamaLengkap = strings.TrimSpace(cmd.NamaLengkap)
	cmd.JenisKelamin = strings.TrimSpace(cmd.JenisKelamin)
	cmd.NoHp = strings.TrimSpace(cmd.NoHp)
	cmd.Foto = strings.TrimSpace(cmd.Foto)
	cmd.Nisn = strings.TrimSpace(cmd.Nisn)
	cmd.TempatLahir = strings.TrimSpace(cmd.TempatLahir)

	emailValidated, err := user.CheckNewEmail(cmd.Email)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	nisnValidated, err := user.CheckNewNISN(cmd.Nisn)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	if err := user.CheckAbsen(cmd.NoAbsen); err != nil {
		return CreateSiswaRes{}, err
	}

	if err := user.CheckAngkatan(cmd.Angkatan); err != nil {
		return CreateSiswaRes{}, err
	}

	hashedPassword, err := uc.hasher.GenerateHash(cmd.Password)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	tx, err := uc.txm.Begin(ctx)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	defer func() { _ = tx.Rollback() }()

	existUsername, err := tx.Pengguna().UserExistByUsername(ctx, cmd.Username)
	if err != nil {
		return CreateSiswaRes{}, err
	}
	if existUsername {
		return CreateSiswaRes{}, coreerror.ErrUsernameTaken
	}

	existNisn, err := tx.ProfilSiswa().ExistByNISN(ctx, string(nisnValidated))
	if err != nil {
		return CreateSiswaRes{}, err
	}
	if existNisn {
		return CreateSiswaRes{}, coreerror.ErrNisnTaken
	}

	userData := user.Pengguna{
		Username:       cmd.Username,
		Email:          emailValidated,
		PasswordHashed: hashedPassword,
		NamaLengkap:    cmd.NamaLengkap,
		JenisKelamin:   cmd.JenisKelamin,
		NoHp:           cmd.NoHp,
		Foto:           cmd.Foto,
		Role:           user.SISWA,
		StatusAkun:     user.AKTIF,
	}

	idPengguna, err := tx.Pengguna().CreateUser(ctx, userData)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	profilSiswaData := user.ProfilSiswa{
		IdPengguna:     idPengguna,
		IdTingkatKelas: cmd.IdTingkatKelas,
		IdNamaKelas:    cmd.IdNamaKelas,
		Nisn:           nisnValidated,
		NoAbsen:        cmd.NoAbsen,
		Angkatan:       cmd.Angkatan,
		TempatLahir:    cmd.TempatLahir,
		TanggalLahir:   cmd.TanggalLahir,
	}

	idProfilSiswa, err := tx.ProfilSiswa().CreateProfilSiswa(ctx, profilSiswaData)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	if err := tx.Commit(); err != nil {
		return CreateSiswaRes{}, err
	}

	return CreateSiswaRes{
		IdPengguna:    idPengguna,
		IdProfilSiswa: idProfilSiswa,
	}, nil
}
