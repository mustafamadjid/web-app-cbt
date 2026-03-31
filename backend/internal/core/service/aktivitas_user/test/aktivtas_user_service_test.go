package aktivitas_user_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	"github.com/stretchr/testify/assert"
)

func TestCreateAktivitasUserService(t *testing.T) {
	repoErr := errors.New("repo error")
	tests := []struct {
		name             string
		cmd              aktivitas_user_service.AktivitasUserCmd
		repoErr          error
		wantErr          error
		wantCreateCalled bool
		wantDesc         string
		wantIP           string
	}{
		{
			name: "Branch 1 ->  semua validasi lolos",
			cmd: aktivitas_user_service.AktivitasUserCmd{
				IdPengguna:  user.ID(1),
				Action:      aktivitas_user.LOGIN,
				Description: "  Login admin ",
				IpAddress:   " 127.0.0.1 ",
			},
			wantCreateCalled: true,
			wantDesc:         "Login admin",
			wantIP:           "127.0.0.1",
		},
		{
			name: "Branch 2 ->  action tidak valid",
			cmd: aktivitas_user_service.AktivitasUserCmd{
				IdPengguna:  user.ID(1),
				Action:      aktivitas_user.Action("INVALID"),
				Description: "Login admin",
				IpAddress:   "127.0.0.1",
			},
			wantErr:          coreerror.ErrInvalidActionActivity,
			wantCreateCalled: false,
		},
		{
			name: "Branch 3 ->  ip address tidak valid",
			cmd: aktivitas_user_service.AktivitasUserCmd{
				IdPengguna:  user.ID(1),
				Action:      aktivitas_user.LOGIN,
				Description: "Login admin",
				IpAddress:   "",
			},
			wantErr:          coreerror.ErrInvalidIpAddress,
			wantCreateCalled: false,
		},
		{
			name: "Branch 4 ->  repository error",
			cmd: aktivitas_user_service.AktivitasUserCmd{
				IdPengguna:  user.ID(1),
				Action:      aktivitas_user.LOGIN,
				Description: "Login admin",
				IpAddress:   "127.0.0.1",
			},
			repoErr:          repoErr,
			wantErr:          repoErr,
			wantCreateCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &FakeAktivitasUserRepo{CreateErr: tt.repoErr}
			svc := aktivitas_user_service.NewAktivitasUserService(repo)

			err := svc.CreateAktivitasUserService(context.Background(), tt.cmd)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCreateCalled, repo.CreateCalled)
			if tt.wantCreateCalled && tt.wantErr == nil {
				assert.Equal(t, tt.wantDesc, repo.CreatedRecord.Description)
				assert.Equal(t, tt.wantIP, repo.CreatedRecord.IpAddress)
			}
		})
	}
}

func TestGetAktivitasUserService(t *testing.T) {
	repoErr := errors.New("repo error")
	tests := []struct {
		name      string
		repoErr   error
		repoData  []aktivitas_user.AktivitasUser
		wantErr   error
		wantCount int
	}{
		{
			name:      "Branch 1 ->  semua validasi lolos",
			repoData:  []aktivitas_user.AktivitasUser{{IdAktivitas: "1"}},
			wantCount: 1,
		},
		{
			name:    "Branch 2 ->  repository error",
			repoErr: repoErr,
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &FakeAktivitasUserRepo{GetErr: tt.repoErr, GetData: tt.repoData}
			svc := aktivitas_user_service.NewAktivitasUserService(repo)

			res, err := svc.GetAktivitasUserService(context.Background())
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCount, len(res))
			assert.True(t, repo.GetCalled)
		})
	}
}
