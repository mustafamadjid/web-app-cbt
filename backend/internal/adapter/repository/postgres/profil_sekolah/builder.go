package profilsekolahrepo

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
)

func scanProfilSekolahRow(row pgx.Row) (profil_sekolah.ProfilSekolah, error) {
	var (
		item          profil_sekolah.ProfilSekolah
		emailSekolah  sql.NullString
		noTelpSekolah sql.NullString
		kepalaSekolah sql.NullString
		wakaSekolah   sql.NullString
		namaSekolah   sql.NullString
		alamatSekolah sql.NullString
	)

	if err := row.Scan(
		&item.IDProfil,
		&emailSekolah,
		&noTelpSekolah,
		&kepalaSekolah,
		&wakaSekolah,
		&namaSekolah,
		&alamatSekolah,
		&item.LogoSekolah,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return profil_sekolah.ProfilSekolah{}, err
	}

	if emailSekolah.Valid {
		item.EmailSekolah = emailSekolah.String
	}
	if noTelpSekolah.Valid {
		item.NoTelpSekolah = noTelpSekolah.String
	}
	if kepalaSekolah.Valid {
		item.KepalaSekolah = kepalaSekolah.String
	}
	if wakaSekolah.Valid {
		item.WakaSekolah = wakaSekolah.String
	}
	if namaSekolah.Valid {
		item.NamaSekolah = namaSekolah.String
	}
	if alamatSekolah.Valid {
		item.AlamatSekolah = alamatSekolah.String
	}

	return item, nil
}
