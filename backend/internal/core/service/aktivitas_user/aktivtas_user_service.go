package aktivitas_user_service

import (
	"context"
	"strings"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/aktivitas_user"
)

type AktivitasUserService struct {
	aktivitasUserRepo aktivitas_user_repo.AktivitasUserRepository
}

func NewAktivitasUserService(aktivitasUserRepo aktivitas_user_repo.AktivitasUserRepository) *AktivitasUserService {
	return &AktivitasUserService{aktivitasUserRepo: aktivitasUserRepo}
}

func (svc *AktivitasUserService)CreateAktivitasUserService(ctx context.Context, aktivitasUser AktivitasUserCmd) error {
	logger := corelog.FromContext(ctx)

	aktivitasUser.Description = strings.TrimSpace(aktivitasUser.Description)
	aktivitasUser.IpAddress = strings.TrimSpace(aktivitasUser.IpAddress)

	if !aktivitasUser.Action.ValidAction() {
		logger.Error(ctx,"Invalid action activity","op","aktivitas_user_service.CreateAktivitasUserService","err",coreerror.ErrInvalidActionActivity)
		return coreerror.ErrInvalidActionActivity
	}

	if !aktivitas_user.ValidIpAddress(aktivitasUser.IpAddress) {
		logger.Error(ctx,"Invalid ip address","op","aktivitas_user_service.CreateAktivitasUserService","err",coreerror.ErrInvalidIpAddress)
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

	logger.Info(ctx,"Success create aktivitas user","op","aktivitas_user_service.CreateAktivitasUserService")
	return nil
}

func (svc *AktivitasUserService)GetAktivitasUserService(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	getAktivitas,err := svc.aktivitasUserRepo.GetAktivitasUser(ctx)
	if err != nil {
		return nil, err
	}

	return getAktivitas, nil
}