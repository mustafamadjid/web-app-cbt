package profil_sekolah_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
)

type UpdateProfilSekolahService struct {
	repo out.ProfilSekolahRepository
}

type UpdateProfilSekolahCmd struct {
	IDProfil      profil_sekolah.IDProfil
	EmailSekolah  *string
	NoTelpSekolah *string
	KepalaSekolah *string
	WakaSekolah   *string
	NamaSekolah   *string
	AlamatSekolah *string
	LogoSekolah   *string
}

func NewUpdateProfilSekolahService(repo out.ProfilSekolahRepository) *UpdateProfilSekolahService {
	return &UpdateProfilSekolahService{repo: repo}
}

func (s *UpdateProfilSekolahService) UpdateProfilSekolah(ctx context.Context, cmd UpdateProfilSekolahCmd) error {
	if cmd.IDProfil <= 0 {
		return coreerror.ErrInvalidInput
	}

	if cmd.EmailSekolah == nil &&
		cmd.NoTelpSekolah == nil &&
		cmd.KepalaSekolah == nil &&
		cmd.WakaSekolah == nil &&
		cmd.NamaSekolah == nil &&
		cmd.AlamatSekolah == nil &&
		cmd.LogoSekolah == nil {
		return coreerror.ErrNoFieldToUpdate
	}

	profil, err := s.repo.GetProfilSekolah(ctx)
	if err != nil {
		return err
	}

	if cmd.EmailSekolah != nil {
		email, err := normalizeOptional(cmd.EmailSekolah)
		if err != nil {
			return err
		}
		profil.EmailSekolah = *email
	}
	if cmd.NoTelpSekolah != nil {
		noTelp, err := normalizeOptional(cmd.NoTelpSekolah)
		if err != nil {
			return err
		}
		profil.NoTelpSekolah = *noTelp
	}
	if cmd.KepalaSekolah != nil {
		kepala, err := normalizeOptional(cmd.KepalaSekolah)
		if err != nil {
			return err
		}
		profil.KepalaSekolah = *kepala
	}
	if cmd.WakaSekolah != nil {
		waka, err := normalizeOptional(cmd.WakaSekolah)
		if err != nil {
			return err
		}
		profil.WakaSekolah = *waka
	}
	if cmd.NamaSekolah != nil {
		nama, err := normalizeOptional(cmd.NamaSekolah)
		if err != nil {
			return err
		}
		profil.NamaSekolah = *nama
	}
	if cmd.AlamatSekolah != nil {
		alamat, err := normalizeOptional(cmd.AlamatSekolah)
		if err != nil {
			return err
		}
		profil.AlamatSekolah = *alamat
	}
	if cmd.LogoSekolah != nil {
		logo, err := normalizeOptional(cmd.LogoSekolah)
		if err != nil {
			return err
		}
		profil.LogoSekolah = logo
	}

	profil.IDProfil = cmd.IDProfil

	return s.repo.UpdateProfilSekolah(ctx, cmd.IDProfil, profil)
}

func normalizeOptional(value *string) (*string, error) {
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, coreerror.ErrInvalidInput
	}
	return &trimmed, nil
}
