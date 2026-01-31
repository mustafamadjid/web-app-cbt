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
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ProfilSiswaRepo struct {
	q Executor
}

func NewProfilSiswaRepo(q Executor) *ProfilSiswaRepo {
	return &ProfilSiswaRepo{q: q}
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

func (r *ProfilSiswaRepo) GetListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem, error) {
	sortColumns := map[string]string{
		"nama_lengkap": "p.nama_lengkap",
		"created_at":   "p.created_at",
		"username":     "p.username",
		"nisn":         "ps.nisn",
	}

	baseQuery := `
		SELECT p.id_pengguna,
			p.username,
			p.email,
			p.nama_lengkap,
			p.jenis_kelamin,
			p.no_hp,
			p.foto,
			p.status_akun,
			nk.nama_kelas,
			k.tingkat_kelas,
			ps.angkatan
		FROM pengguna p
		JOIN profil_siswa ps ON ps.id_pengguna = p.id_pengguna
		JOIN kelas k ON ps.id_kelas = k.id_kelas
		JOIN nama_kelas nk ON ps.id_nama_kelas = nk.id_nama_kelas
	`

	where := make([]string, 0, 6)
	args := make([]any, 0, 8)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`(p.nama_lengkap ILIKE $%d OR p.username ILIKE $%d OR p.email ILIKE $%d OR ps.nisn ILIKE $%d)`, idx, idx, idx, idx))
	}

	if filter.Status != nil {
		args = append(args, string(*filter.Status))
		where = append(where, fmt.Sprintf("p.status_akun = $%d", len(args)))
	}

	if filter.Angkatan != nil {
		args = append(args, *filter.Angkatan)
		where = append(where, fmt.Sprintf("ps.angkatan = $%d", len(args)))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("k.tingkat_kelas = $%d", len(args)))
	}

	if filter.JenisKelamin != nil {
		args = append(args, *filter.JenisKelamin)
		where = append(where, fmt.Sprintf("p.jenis_kelamin = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	sortColumn, ok := sortColumns[filter.SortBy]
	if !ok {
		sortColumn = "p.created_at"
	}

	direction := "ASC"
	if filter.SortDesc {
		direction = "DESC"
	}

	baseQuery = fmt.Sprintf("%s ORDER BY %s %s", baseQuery, sortColumn, direction)

	args = append(args, filter.Limit)
	limitIndex := len(args)
	args = append(args, filter.Offset)
	offsetIndex := len(args)
	baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)

	rows, err := r.q.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.SiswaListItem
	for rows.Next() {
		var item query.SiswaListItem
		var jenisKelamin int16
		var status string
		var email string

		if err := rows.Scan(
			&item.IdPengguna,
			&item.Username,
			&email,
			&item.NamaLengkap,
			&jenisKelamin,
			&item.NoHp,
			&item.Foto,
			&status,
			&item.NamaKelas,
			&item.TingkatKelas,
			&item.Angkatan,
		); err != nil {
			return nil, err
		}

		jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
		if err != nil {
			return nil, err
		}

		item.Email = user.Email(email)
		item.JenisKelamin = jenisKelaminValue
		item.StatusAkun = user.StatusAkun(status)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
