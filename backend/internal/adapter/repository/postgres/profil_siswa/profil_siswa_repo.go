package profilsiswarepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ProfilSiswaRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewProfilSiswaRepo(q pg.Executor, logger corelog.Logger) *ProfilSiswaRepo {
	return &ProfilSiswaRepo{q: q, logger: logger}
}

func (r *ProfilSiswaRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ProfilSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	const query = `
		SELECT 
			ps.id_pengguna,
			ps.id_siswa,
			ps.id_nama_kelas,
			u.username,
			u.email,
			u.nama_lengkap,
			u.jenis_kelamin,
			u.no_hp,
			u.foto,
			r.nama_role,
			u.status_akun,
			ps.nisn,
			ps.no_absen,
			ps.angkatan,
			ps.tempat_lahir,
			ps.tanggal_lahir,
			nk.nama_kelas,
			k.tingkat_kelas
		FROM profil_siswa ps
		JOIN pengguna u ON ps.id_pengguna = u.id_pengguna
		JOIN role r ON u.id_role = r.id_role
		JOIN nama_kelas nk ON ps.id_nama_kelas = nk.id_nama_kelas
		JOIN kelas k ON nk.id_kelas = k.id_kelas
		WHERE ps.id_pengguna = $1
	`

	result, err := scanProfilSiswaRow(r.q.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return user.DataSiswa{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed finding profil siswa", "op", "profil_siswa_repo.find_by_id", "user_id", id, "err", err)
		return user.DataSiswa{}, err
	}

	return result, nil
}

func (r *ProfilSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM profil_siswa WHERE nisn = $1)`

	var exists bool
	if err := r.q.QueryRow(ctx, query, nisn).Scan(&exists); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking nisn existence", "op", "profil_siswa_repo.exists_by_nisn", "nisn", nisn, "err", err)
		return false, err
	}

	return exists, nil
}

func (r *ProfilSiswaRepo) CreateProfilSiswa(ctx context.Context, profilSiswa user.ProfilSiswa) (user.ID, error) {
	const query = `
		INSERT INTO profil_siswa (
			id_pengguna,
			id_nama_kelas,
			nisn,
			no_absen,
			angkatan,
			tempat_lahir,
			tanggal_lahir
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_siswa
	`
	// id_kelas,
	// profilSiswa.IdTingkatKelas,

	var id user.ID
	err := r.q.QueryRow(
		ctx,
		query,
		profilSiswa.IdPengguna,
		profilSiswa.IdNamaKelas,
		string(profilSiswa.Nisn),
		profilSiswa.NoAbsen,
		profilSiswa.Angkatan,
		profilSiswa.TempatLahir,
		profilSiswa.TanggalLahir,
	).Scan(&id)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating profil siswa", "op", "profil_siswa_repo.create", "user_id", profilSiswa.IdPengguna, "err", err)
		return 0, err
	}

	return id, nil
}

func (r *ProfilSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa updatepatch.ProfilSiswa) error {
	set := make([]string, 0, 7)
	args := make([]any, 0, 8)

	add := func(col string, v any) {
		args = append(args, v)
		set = append(set, fmt.Sprintf("%s=$%d", col, len(args)))
	}

	// if profilSiswa.IdTingkatKelas != nil {
	// 	add("id_kelas", *profilSiswa.IdTingkatKelas)
	// }
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
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating profil siswa", "op", "profil_siswa_repo.update", "user_id", idPengguna, "err", err)
	}
	return err
}

func (r *ProfilSiswaRepo) GetListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem, error) {
	queryText, args := r.buildListSiswaQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing siswa", "op", "profil_siswa_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanSiswaListRows(ctx, "profil_siswa_repo.list", rows)
}
