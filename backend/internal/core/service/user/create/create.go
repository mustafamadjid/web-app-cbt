package user_service

import (
	"context"
	"strings"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

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
	logger := corelog.FromContext(ctx)
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



	isDashedNip := cmd.Nip == "-"

	nipToStore := user.NIP(cmd.Nip)

	var nipValidated user.NIP
	
	if !isDashedNip {
		v, error := user.CheckNewNip(cmd.Nip)
	if error != nil {
		return CreateGuruRes{}, error
		}
		nipValidated = v
		nipToStore = v
	}


	// Validasi
	emailValidated, error := user.CheckNewEmail(cmd.Email)
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
		logger.Error(ctx, "failed starting transaction", "op", "user.create_guru", "err", err)
		return CreateGuruRes{}, err
	}

	defer func() { _ = tx.Rollback() }()

	existUsername, error := tx.Pengguna().UserExistByUsername(ctx, cmd.Username)
	if error != nil {
		logger.Error(ctx, "failed checking username", "op", "user.create_guru", "err", error)
		return CreateGuruRes{}, error
	}
	if existUsername {
		return CreateGuruRes{}, coreerror.ErrUsernameTaken
	}


	if !isDashedNip {
		existNip, error := tx.ProfilGuru().ExistByNIP(ctx, nipValidated)
	if error != nil {
		logger.Error(ctx, "failed checking nip", "op", "user.create_guru", "err", error)
		return CreateGuruRes{}, error
	}
	if existNip {
		return CreateGuruRes{}, coreerror.ErrNipTaken
	}
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
		logger.Error(ctx, "failed creating user", "op", "user.create_guru", "err", error)
		return CreateGuruRes{}, error
	}

	profilGuruData := user.ProfilGuru{
		IdPengguna:  idPengguna,
		Nip:         nipToStore,
		Jabatan:     cmd.Jabatan,
		BidangStudi: cmd.BidangStudi,
	}

	idProfilGuru, error := tx.ProfilGuru().CreateProfilGuru(ctx, profilGuruData)
	if error != nil {
		logger.Error(ctx, "failed creating profil guru", "op", "user.create_guru", "user_id", idPengguna, "err", error)
		return CreateGuruRes{}, error
	}

	if error := tx.Commit(); error != nil {
		logger.Error(ctx, "failed committing transaction", "op", "user.create_guru", "user_id", idPengguna, "err", error)
		return CreateGuruRes{}, error
	}

	return CreateGuruRes{
		IdPengguna:   idPengguna,
		IdProfilGuru: idProfilGuru,
	}, nil

}

func (uc *CreateTx) CreateSiswa(ctx context.Context, cmd CreateSiswaCmd, actor user.Actor) (CreateSiswaRes, error) {
	logger := corelog.FromContext(ctx)
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

	isDashedNisn := cmd.Nisn == "-"

	nisnToStore := user.NISN(cmd.Nisn)

	var nisnValidated user.NISN

	if !isDashedNisn {
		v,err := user.CheckNewNISN(cmd.Nisn)
		if err != nil {
			return CreateSiswaRes{}, err
		}
		nisnValidated = v
		nisnToStore = v
	}

	emailValidated, err := user.CheckNewEmail(cmd.Email)
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
		logger.Error(ctx, "failed starting transaction", "op", "user.create_siswa", "err", err)
		return CreateSiswaRes{}, err
	}

	defer func() { _ = tx.Rollback() }()

	existUsername, err := tx.Pengguna().UserExistByUsername(ctx, cmd.Username)
	if err != nil {
		logger.Error(ctx, "failed checking username", "op", "user.create_siswa", "err", err)
		return CreateSiswaRes{}, err
	}
	if existUsername {
		return CreateSiswaRes{}, coreerror.ErrUsernameTaken
	}

	if !isDashedNisn {
		existNisn, err := tx.ProfilSiswa().ExistByNISN(ctx, string(nisnValidated))
	if err != nil {
		logger.Error(ctx, "failed checking nisn", "op", "user.create_siswa", "err", err)
		return CreateSiswaRes{}, err
	}
	if existNisn {
		return CreateSiswaRes{}, coreerror.ErrNisnTaken
	}
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
		logger.Error(ctx, "failed creating user", "op", "user.create_siswa", "err", err)
		return CreateSiswaRes{}, err
	}

	profilSiswaData := user.ProfilSiswa{
		IdPengguna:     idPengguna,
		IdTingkatKelas: cmd.IdTingkatKelas,
		IdNamaKelas:    cmd.IdNamaKelas,
		Nisn:           nisnToStore,
		NoAbsen:        cmd.NoAbsen,
		Angkatan:       cmd.Angkatan,
		TempatLahir:    cmd.TempatLahir,
		TanggalLahir:   cmd.TanggalLahir,
	}

	idProfilSiswa, err := tx.ProfilSiswa().CreateProfilSiswa(ctx, profilSiswaData)
	if err != nil {
		logger.Error(ctx, "failed creating profil siswa", "op", "user.create_siswa", "user_id", idPengguna, "err", err)
		return CreateSiswaRes{}, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "failed committing transaction", "op", "user.create_siswa", "user_id", idPengguna, "err", err)
		return CreateSiswaRes{}, err
	}

	logger.Info(ctx, "success creating user", "op", "user.create_siswa", "user_id", idPengguna)
	return CreateSiswaRes{
		IdPengguna:    idPengguna,
		IdProfilSiswa: idProfilSiswa,
	}, nil
}
