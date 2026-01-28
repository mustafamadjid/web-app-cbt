package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type ProfilgGuruRepo struct {
	q Executor
}

func (r *ProfilgGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.ProfilGuru, error) {
	const query = `
		SELECT id_guru,
			id_pengguna,
			nip,
			jabatan,
			bidang_studi
		FROM profil_guru
		WHERE id_pengguna = $1
	`

	var result user.ProfilGuru
	var nip string
	err := r.q.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.IdPengguna,
		&nip,
		&result.Jabatan,
		&result.BidangStudi,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.ProfilGuru{}, coreerror.ErrNotFound
	}
	if err != nil {
		return user.ProfilGuru{}, err
	}

	result.Nip = user.NIP(nip)
	return result, nil
}

func (r *ProfilgGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM profil_guru WHERE nip = $1)`

	var exists bool
	if err := r.q.QueryRow(ctx, query, string(nip)).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ProfilgGuruRepo) CreateProfilGuru(ctx context.Context, profilGuru user.ProfilGuru) (user.ID, error) {
	const query = `
		INSERT INTO profil_guru (
			id_pengguna,
			nip,
			jabatan,
			bidang_studi
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id_guru
	`

	var id user.ID
	err := r.q.QueryRow(
		ctx,
		query,
		profilGuru.IdPengguna,
		string(profilGuru.Nip),
		profilGuru.Jabatan,
		profilGuru.BidangStudi,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProfilgGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru user.ProfilGuru) (user.ID, error) {
	const query = `
		UPDATE profil_guru
		SET nip = $1,
			jabatan = $2,
			bidang_studi = $3,
			updated_at = now()
		WHERE id_pengguna = $4
		RETURNING id_guru
	`

	var id user.ID
	err := r.q.QueryRow(
		ctx,
		query,
		string(profilGuru.Nip),
		profilGuru.Jabatan,
		profilGuru.BidangStudi,
		idPengguna,
	).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, coreerror.ErrNotFound
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}
