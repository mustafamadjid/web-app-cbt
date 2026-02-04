package aktivitas_user_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	aktivitas_user_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/aktivitas_user"
)

type AktivitasUserService struct {
	aktivitasUserRepo aktivitas_user_repo.AktivitasUserRepository
}

func (svc *AktivitasUserService)CreateAktivitasUserService(ctx context.Context, aktivitasUser AktivitasUserCmd) error {
	aktivitasUser.Description = strings.TrimSpace(aktivitasUser.Description)
	aktivitasUser.IpAddress = strings.TrimSpace(aktivitasUser.IpAddress)

	if !aktivitasUser.Action.ValidAction() {
		return coreerror.ErrInvalidActionActivity
	}

	if !aktivitas_user.ValidIpAddress(aktivitasUser.IpAddress) {
		return coreerror.ErrInvalidIpAddress
	}

	aktivitasData := aktivitas_user.AktivitasUser{
		IdPengguna:  aktivitasUser.IdPengguna,
		Action:      aktivitasUser.Action,
		Description: aktivitasUser.Description,
		IpAddress:   aktivitasUser.IpAddress,
	}

	if err := svc.aktivitasUserRepo.CreateAktivitasUser(ctx, aktivitasData); err != nil {
		return err
	}

	return nil
}