package aktivitas_user_repo

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
)

type AktivitasUserRepository interface {
	CreateAktivitasUser(ctx context.Context, aktivitasUser aktivitas_user.AktivitasUser)error
}