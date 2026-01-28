package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type ProfilSiswaRepo struct {
	q Executor
}

func (r *ProfilSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.ProfilSiswa, error) {
	const query = `
		SELECT id_siswa,
			id_pengguna,
			id_kelas,
			id_nama_kelas,
			nisn,
			no_absen,
			angkatan,
			tempat_lahir,
			tanggal_lahir
		FROM profil_siswa
		WHERE id_pengguna = $1
	`

	var result user.ProfilSiswa
	var nisn string
	err := r.q.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.IdPengguna,
		&result.IdTingkatKelas,
		&result.IdNamaKelas,
		&nisn,
		&result.NoAbsen,
		&result.Angkatan,
		&result.TempatLahir,
		&result.TanggalLahir,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.ProfilSiswa{}, coreerror.ErrNotFound
	}
	if err != nil {
		return user.ProfilSiswa{}, err
	}

	result.Nisn = user.NISN(nisn)
	return result, nil
}

func (r *ProfilSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM profil_siswa WHERE nisn = $1)`

	var exists bool
	if err := r.q.QueryRow(ctx, query, nisn).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ProfilSiswaRepo) CreateProfilSiswa(ctx context.Context, profilSiswa user.ProfilSiswa) (user.ID, error) {
	const query = `
		INSERT INTO profil_siswa (
			id_pengguna,
			id_kelas,
			id_nama_kelas,
			nisn,
			no_absen,
			angkatan,
			tempat_lahir,
			tanggal_lahir
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_siswa
	`

	var id user.ID
	err := r.q.QueryRow(
		ctx,
		query,
		profilSiswa.IdPengguna,
		profilSiswa.IdTingkatKelas,
		profilSiswa.IdNamaKelas,
		string(profilSiswa.Nisn),
		profilSiswa.NoAbsen,
		profilSiswa.Angkatan,
		profilSiswa.TempatLahir,
		profilSiswa.TanggalLahir,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProfilSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa outuser.UpdateProfilSiswaPatch) error {
	set := make([]string, 0, 7)
	args := make([]any, 0, 8)

	add := func(col string, v any) {
		args = append(args, v)
		set = append(set, fmt.Sprintf("%s=$%d", col, len(args)))
	}

	if profilSiswa.IdTingkatKelas != nil {
		add("id_kelas", *profilSiswa.IdTingkatKelas)
	}
	if profilSiswa.IdNamaKelas != nil {
		add("id_nama_kelas", *profilSiswa.IdNamaKelas)
	}
	if profilSiswa.Nisn != nil {
		add("nisn", *profilSiswa.Nisn)
	}
	if profilSiswa.NoAbsen != nil {
		add("no_absen", *profilSiswa.NoAbsen)
	}
	if profilSiswa.Angkatan != nil {
		add("angkatan", *profilSiswa.Angkatan)
	}
	if profilSiswa.TempatLahir != nil {
		add("tempat_lahir", *profilSiswa.TempatLahir)
	}
	if profilSiswa.TanggalLahir != nil {
		add("tanggal_lahir", *profilSiswa.TanggalLahir)
	}

	if len(set) == 0 {
		return nil
	}

	args = append(args, idPengguna)
	q := fmt.Sprintf(`UPDATE profil_siswa SET %s WHERE id_pengguna=$%d`, strings.Join(set, ", "), len(args))

	_, err := r.q.Exec(ctx, q, args...)
	return err
}
