package aktivitas_user_service

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	aktivitas_user_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/aktivitas_user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type AktivitasUserService struct {
	aktivitasUserRepo aktivitas_user_repo.AktivitasUserRepository
}

func NewAktivitasUserService(aktivitasUserRepo aktivitas_user_repo.AktivitasUserRepository) *AktivitasUserService {
	return &AktivitasUserService{aktivitasUserRepo: aktivitasUserRepo}
}

func (svc *AktivitasUserService) CreateAktivitasUserService(ctx context.Context, aktivitasUser AktivitasUserCmd) error {
	logger := corelog.FromContext(ctx)

	aktivitasUser = sanitizeAktivitasUserCmd(aktivitasUser)

	if err := validateAktivitasAction(aktivitasUser); err != nil {
		logger.Error(ctx, "Invalid action activity", "layer", "core.service", "op", "aktivitas_user_service.CreateAktivitasUserService", "err", coreerror.ErrInvalidActionActivity)
		return err
	}

	if err := validateAktivitasIPAddress(aktivitasUser); err != nil {
		logger.Error(ctx, "Invalid ip address", "layer", "core.service", "op", "aktivitas_user_service.CreateAktivitasUserService", "err", coreerror.ErrInvalidIpAddress)
		return err
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

func (svc *AktivitasUserService) GetAktivitasUserService(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	getAktivitas, err := svc.aktivitasUserRepo.GetAktivitasUser(ctx)
	if err != nil {
		return nil, err
	}

	return getAktivitas, nil
}

// -----------------------
// Sanitizer and validator
// -----------------------

func sanitizeAktivitasUserCmd(cmd AktivitasUserCmd) AktivitasUserCmd {
	cmd.Description = strings.TrimSpace(cmd.Description)
	cmd.IpAddress = strings.TrimSpace(cmd.IpAddress)
	return cmd
}

func validateAktivitasAction(cmd AktivitasUserCmd) error {
	if !cmd.Action.ValidAction() {
		return coreerror.ErrInvalidActionActivity
	}
	return nil
}

func validateAktivitasIPAddress(cmd AktivitasUserCmd) error {
	if !aktivitas_user.ValidIpAddress(cmd.IpAddress) {
		return coreerror.ErrInvalidIpAddress
	}
	return nil
}
