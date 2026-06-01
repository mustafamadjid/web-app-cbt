package user_service

import (
	"context"
	"errors"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

func (uc *CreateTx) createGuruTx(ctx context.Context, data ValidatedGuruCreate, hashedPassword string) (CreateGuruRes, error) {
	logger := corelog.FromContext(ctx)
	tx, err := uc.txm.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", opCreateGuru, "err", err)
		return CreateGuruRes{}, err
	}

	defer func() { _ = tx.Rollback() }()

	if err := ensureUsernameAvailable(ctx, tx, data.cmd.Username, opCreateGuru); err != nil {
		return CreateGuruRes{}, err
	}

	if !data.isDashedNip {
		existNip, err := tx.ProfilGuru().ExistByNIP(ctx, data.nipValidated)
		if err != nil {
			logger.Error(ctx, "failed checking nip", "layer", "core.service", "op", opCreateGuru, "err", err)
			return CreateGuruRes{}, err
		}
		if existNip {
			logger.Error(ctx, "failed checking nip", "layer", "core.service", "op", opCreateGuru, "err", coreerror.ErrNipTaken)
			return CreateGuruRes{}, coreerror.ErrNipTaken
		}
	}

	idPengguna, err := tx.Pengguna().CreateUser(ctx, buildCreateGuruUser(data.cmd, data.email, hashedPassword))
	if err != nil {
		if isCreateUserConstraintErr(err) {
			return CreateGuruRes{}, err
		}

		logger.Error(ctx, "failed creating user", "layer", "core.service", "op", opCreateGuru, "err", err)
		return CreateGuruRes{}, err
	}

	idProfilGuru, err := tx.ProfilGuru().CreateProfilGuru(ctx, buildProfilGuru(data.cmd, idPengguna, data.nipToStore))
	if err != nil {
		logger.Error(ctx, "failed creating profil guru", "layer", "core.service", "op", opCreateGuru, "user_id", idPengguna, "err", err)
		return CreateGuruRes{}, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "failed committing transaction", "layer", "core.service", "op", opCreateGuru, "user_id", idPengguna, "err", err)
		return CreateGuruRes{}, err
	}

	return CreateGuruRes{
		IdPengguna:   idPengguna,
		IdProfilGuru: idProfilGuru,
	}, nil
}

func (uc *CreateTx) createSiswaTx(ctx context.Context, data ValidatedSiswaCreate, hashedPassword string) (CreateSiswaRes, error) {
	logger := corelog.FromContext(ctx)
	tx, err := uc.txm.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", opCreateSiswa, "err", err)
		return CreateSiswaRes{}, err
	}

	defer func() { _ = tx.Rollback() }()

	if err := ensureUsernameAvailable(ctx, tx, data.cmd.Username, opCreateSiswa); err != nil {
		return CreateSiswaRes{}, err
	}

	if !data.isDashedNisn {
		existNisn, err := tx.ProfilSiswa().ExistByNISN(ctx, string(data.nisnValidated))
		if err != nil {
			logger.Error(ctx, "failed checking nisn", "layer", "core.service", "op", opCreateSiswa, "err", err)
			return CreateSiswaRes{}, err
		}
		if existNisn {
			return CreateSiswaRes{}, coreerror.ErrNisnTaken
		}
	}

	idPengguna, err := tx.Pengguna().CreateUser(ctx, buildCreateSiswaUser(data.cmd, data.email, hashedPassword))
	if err != nil {
		if isCreateUserConstraintErr(err) {
			return CreateSiswaRes{}, err
		}

		logger.Error(ctx, "failed creating user", "layer", "core.service", "op", opCreateSiswa, "err", err)
		return CreateSiswaRes{}, err
	}

	idProfilSiswa, err := tx.ProfilSiswa().CreateProfilSiswa(ctx, buildProfilSiswa(data.cmd, idPengguna, data.nisnToStore))
	if err != nil {
		logger.Error(ctx, "failed creating profil siswa", "layer", "core.service", "op", opCreateSiswa, "user_id", idPengguna, "err", err)
		return CreateSiswaRes{}, err
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "failed committing transaction", "layer", "core.service", "op", opCreateSiswa, "user_id", idPengguna, "err", err)
		return CreateSiswaRes{}, err
	}

	return CreateSiswaRes{
		IdPengguna:    idPengguna,
		IdProfilSiswa: idProfilSiswa,
	}, nil
}

func ensureUsernameAvailable(ctx context.Context, tx txout.Tx, username string, op string) error {
	logger := corelog.FromContext(ctx)
	existUsername, err := tx.Pengguna().UserExistByUsername(ctx, username)
	if err != nil {
		logger.Error(ctx, "failed checking username", "layer", "core.service", "op", op, "err", err)
		return err
	}
	if existUsername {
		return coreerror.ErrUsernameTaken
	}
	return nil
}

func isCreateUserConstraintErr(err error) bool {
	return errors.Is(err, coreerror.ErrUsernameTaken) ||
		errors.Is(err, coreerror.ErrEmailTaken) ||
		errors.Is(err, coreerror.ErrNoHpTaken)
}
