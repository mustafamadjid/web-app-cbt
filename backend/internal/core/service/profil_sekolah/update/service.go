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
	EmailSekolah  string
	NoTelpSekolah string
	KepalaSekolah string
	WakaSekolah   string
	NamaSekolah   string
	AlamatSekolah string
	LogoSekolah   *string
}

func NewUpdateProfilSekolahService(repo out.ProfilSekolahRepository) *UpdateProfilSekolahService {
	return &UpdateProfilSekolahService{repo: repo}
}

func (s *UpdateProfilSekolahService) UpdateProfilSekolah(ctx context.Context, cmd UpdateProfilSekolahCmd) error {
	if cmd.IDProfil <= 0 {
		return coreerror.ErrInvalidInput
	}

	email, err := normalizeRequired(cmd.EmailSekolah)
	if err != nil {
		return err
	}
	noTelp, err := normalizeRequired(cmd.NoTelpSekolah)
	if err != nil {
		return err
	}
	kepala, err := normalizeRequired(cmd.KepalaSekolah)
	if err != nil {
		return err
	}
	waka, err := normalizeRequired(cmd.WakaSekolah)
	if err != nil {
		return err
	}
	nama, err := normalizeRequired(cmd.NamaSekolah)
	if err != nil {
		return err
	}
	alamat, err := normalizeRequired(cmd.AlamatSekolah)
	if err != nil {
		return err
	}

	var logo *string
	if cmd.LogoSekolah != nil {
		trimmed := strings.TrimSpace(*cmd.LogoSekolah)
		if trimmed == "" {
			return coreerror.ErrInvalidInput
		}
		logo = &trimmed
	}

	profil := profil_sekolah.ProfilSekolah{
		IDProfil:      cmd.IDProfil,
		EmailSekolah:  email,
		NoTelpSekolah: noTelp,
		KepalaSekolah: kepala,
		WakaSekolah:   waka,
		NamaSekolah:   nama,
		AlamatSekolah: alamat,
		LogoSekolah:   logo,
	}

	return s.repo.UpdateProfilSekolah(ctx, cmd.IDProfil, profil)
}

func normalizeRequired(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", coreerror.ErrInvalidInput
	}
	return trimmed, nil
}
