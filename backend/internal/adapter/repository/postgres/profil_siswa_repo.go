package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ProfilSiswaRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewProfilSiswaRepo(q Executor, logger corelog.Logger) *ProfilSiswaRepo {
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
		JOIN kelas k ON ps.id_kelas = k.id_kelas
		JOIN nama_kelas nk ON ps.id_nama_kelas = nk.id_nama_kelas
		WHERE ps.id_pengguna = $1
	`

	var result user.DataSiswa

	var nisn sql.NullString
	var noAbsen sql.NullInt32
	var angkatan sql.NullInt32
	var tempatLahir sql.NullString
	var tanggalLahir sql.NullTime
	var jenisKelamin int16
	var email sql.NullString
	var noHp sql.NullString
	var foto sql.NullString

	err := r.q.QueryRow(ctx, query, id).Scan(
		&result.IdPengguna,
		&result.IdSiswa,
		&result.Username,
		&email,
		&result.NamaLengkap,
		&jenisKelamin,
		&noHp,
		&foto,
		&result.Role,
		&result.StatusAkun,
		&nisn,
		&noAbsen,
		&angkatan,
		&tempatLahir,
		&tanggalLahir,
		&result.NamaKelas,
		&result.TingkatKelas,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.DataSiswa{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed finding profil siswa", "op", "profil_siswa_repo.find_by_id", "user_id", id, "err", err)
		return user.DataSiswa{}, err
	}

	if email.Valid {
		result.Email = email.String
	}
	if noHp.Valid {
		result.NoHp = noHp.String
	}
	if foto.Valid {
		result.Foto = foto.String
	}
	if nisn.Valid {
		result.Nisn = user.NISN(nisn.String)
	}
	if noAbsen.Valid {
		result.NoAbsen = int(noAbsen.Int32)
	}
	if angkatan.Valid {
		result.Angkatan = int(angkatan.Int32)
	}
	if tempatLahir.Valid {
		result.TempatLahir = tempatLahir.String
	}
	if tanggalLahir.Valid {
		result.TanggalLahir = tanggalLahir.Time
	}

	jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
	if err != nil {
		return user.DataSiswa{}, err
	}
	result.JenisKelamin = jenisKelaminValue

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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
			ps.no_absen,
			ps.angkatan,
			ps.tempat_lahir,
			ps.tanggal_lahir,
			ps.nisn
		FROM pengguna p
		JOIN profil_siswa ps ON ps.id_pengguna = p.id_pengguna
		JOIN nama_kelas nk ON ps.id_nama_kelas = nk.id_nama_kelas
		JOIN kelas k ON nk.id_kelas = k.id_kelas
	`
	// JOIN kelas k ON ps.id_kelas = k.id_kelas

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
		r.loggerFor(ctx).Error(ctx, "failed listing siswa", "op", "profil_siswa_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	var results []query.SiswaListItem
	for rows.Next() {
		var item query.SiswaListItem
		var jenisKelamin int16
		var status string
		var email sql.NullString
		var noHp sql.NullString
		var foto sql.NullString
		var noAbsen sql.NullInt32
		var angkatan sql.NullInt32
		var tempatLahir sql.NullString
		var tanggalLahir sql.NullTime

		if err := rows.Scan(
			&item.IdPengguna,
			&item.Username,
			&email,
			&item.NamaLengkap,
			&jenisKelamin,
			&noHp,
			&foto,
			&status,
			&item.NamaKelas,
			&item.TingkatKelas,
			&noAbsen,
			&angkatan,
			&tempatLahir,
			&tanggalLahir,
			&item.Nisn,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning siswa list", "op", "profil_siswa_repo.list_scan", "err", err)
			return nil, err
		}

		jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
		if err != nil {
			return nil, err
		}

		if email.Valid {
			item.Email = user.Email(email.String)
		}
		item.JenisKelamin = jenisKelaminValue
		item.StatusAkun = user.StatusAkun(status)
		if noHp.Valid {
			item.NoHp = noHp.String
		}
		if foto.Valid {
			item.Foto = foto.String
		}
		if noAbsen.Valid {
			item.NoAbsen = int(noAbsen.Int32)
		}
		if angkatan.Valid {
			item.Angkatan = int(angkatan.Int32)
		}
		if tempatLahir.Valid {
			item.TempatLahir = tempatLahir.String
		}
		if tanggalLahir.Valid {
			item.TanggalLahir = tanggalLahir.Time
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating siswa list", "op", "profil_siswa_repo.list_iter", "err", err)
		return nil, err
	}

	return results, nil
}
