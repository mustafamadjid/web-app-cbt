package user_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

func (uc *CreateTx) CreateGuru(ctx context.Context, cmd CreateGuruCmd, actor user.Actor) (CreateGuruRes, error) {
	logger := corelog.FromContext(ctx)

	if err := validateCreateGuruActor(actor); err != nil {
		return CreateGuruRes{}, err
	}

	cmd = sanitizeCreateGuruCmd(cmd)
	if err := user.CheckUsernameLength(cmd.Username); err != nil {
		return CreateGuruRes{}, err
	}

	isDashedNip := cmd.Nip == "-"
	nipToStore := user.NIP(cmd.Nip)
	var nipValidated user.NIP
	if !isDashedNip {
		v, err := user.CheckNewNip(cmd.Nip)
		if err != nil {
			return CreateGuruRes{}, err
		}
		nipValidated = v
		nipToStore = v
	}

	email, err := validateCreateEmail(cmd.Email)
	if err != nil {
		return CreateGuruRes{}, err
	}

	data := ValidatedGuruCreate{
		cmd:          cmd,
		email:        email,
		nipToStore:   nipToStore,
		nipValidated: nipValidated,
		isDashedNip:  isDashedNip,
	}

	hashedPassword, err := uc.hasher.GenerateHash(data.cmd.Password)
	if err != nil {
		logger.Error(ctx, "failed hashing password", "layer", "core.service", "op", opCreateGuru, "err", err)
		return CreateGuruRes{}, err
	}

	return uc.createGuruTx(ctx, data, hashedPassword)
}

func (uc *CreateTx) CreateSiswa(ctx context.Context, cmd CreateSiswaCmd, actor user.Actor) (CreateSiswaRes, error) {
	logger := corelog.FromContext(ctx)

	if err := validateCreateSiswaActor(actor); err != nil {
		return CreateSiswaRes{}, err
	}

	cmd = sanitizeCreateSiswaCmd(cmd)
	if err := user.CheckUsernameLength(cmd.Username); err != nil {
		return CreateSiswaRes{}, err
	}

	isDashedNisn := cmd.Nisn == "-"
	nisnToStore := user.NISN(cmd.Nisn)
	var nisnValidated user.NISN
	if !isDashedNisn {
		v, err := user.CheckNewNISN(cmd.Nisn)
		if err != nil {
			return CreateSiswaRes{}, err
		}
		nisnValidated = v
		nisnToStore = v
	}

	email, err := validateCreateEmail(cmd.Email)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	if err := user.CheckAbsen(cmd.NoAbsen); err != nil {
		return CreateSiswaRes{}, err
	}

	if err := user.CheckAngkatan(cmd.Angkatan); err != nil {
		return CreateSiswaRes{}, err
	}

	data := ValidatedSiswaCreate{
		cmd:           cmd,
		email:         email,
		nisnToStore:   nisnToStore,
		nisnValidated: nisnValidated,
		isDashedNisn:  isDashedNisn,
	}

	hashedPassword, err := uc.hasher.GenerateHash(data.cmd.Password)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	res, err := uc.createSiswaTx(ctx, data, hashedPassword)
	if err != nil {
		return CreateSiswaRes{}, err
	}

	logger.Info(ctx, "success creating user", "layer", "core.service", "op", opCreateSiswa, "user_id", res.IdPengguna)
	return res, nil
}
