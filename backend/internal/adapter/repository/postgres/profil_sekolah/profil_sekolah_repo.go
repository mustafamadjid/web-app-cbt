package profilsekolahrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
)

type ProfilSekolahRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewProfilSekolahRepo(q pg.Executor, logger corelog.Logger) *ProfilSekolahRepo {
	return &ProfilSekolahRepo{q: q, logger: logger}
}

func (r *ProfilSekolahRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, idProfil profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error {
	const query = `
		INSERT INTO profil_sekolah (
			id_profil,
			email_sekolah,
			no_telp_sekolah,
			kepala_sekolah,
			waka_sekolah,
			nama_sekolah,
			alamat_sekolah,
			logo_sekolah,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, now()
		)
		ON CONFLICT (id_profil)
		DO UPDATE SET
			email_sekolah = EXCLUDED.email_sekolah,
			no_telp_sekolah = EXCLUDED.no_telp_sekolah,
			kepala_sekolah = EXCLUDED.kepala_sekolah,
			waka_sekolah = EXCLUDED.waka_sekolah,
			nama_sekolah = EXCLUDED.nama_sekolah,
			alamat_sekolah = EXCLUDED.alamat_sekolah,
			logo_sekolah = EXCLUDED.logo_sekolah,
			updated_at = now()
	`

	_, err := r.q.Exec(
		ctx,
		query,
		idProfil,
		profil.EmailSekolah,
		profil.NoTelpSekolah,
		profil.KepalaSekolah,
		profil.WakaSekolah,
		profil.NamaSekolah,
		profil.AlamatSekolah,
		profil.LogoSekolah,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating profil sekolah", "op", "profil_sekolah_repo.upsert", "profil_id", idProfil, "err", err)
	}
	return err
}

func (r *ProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	const query = `
		SELECT id_profil,
			email_sekolah,
			no_telp_sekolah,
			kepala_sekolah,
			waka_sekolah,
			nama_sekolah,
			alamat_sekolah,
			logo_sekolah,
			created_at,
			updated_at
		FROM profil_sekolah
		WHERE id_profil = 1
	`

	profil, err := scanProfilSekolahRow(r.q.QueryRow(ctx, query))
	if errors.Is(err, pgx.ErrNoRows) {
		return profil_sekolah.ProfilSekolah{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting profil sekolah", "op", "profil_sekolah_repo.get", "err", err)
		return profil_sekolah.ProfilSekolah{}, err
	}

	return profil, nil
}

var _ out.ProfilSekolahRepository = (*ProfilSekolahRepo)(nil)
