package profilsiswarepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type siswaScanner interface {
	Scan(dest ...any) error
}

func scanProfilSiswaRow(row siswaScanner) (user.DataSiswa, error) {
	var (
		result       user.DataSiswa
		nisn         sql.NullString
		noAbsen      sql.NullInt32
		angkatan     sql.NullInt32
		tempatLahir  sql.NullString
		tanggalLahir sql.NullTime
		jenisKelamin int16
		email        sql.NullString
		noHp         sql.NullString
		foto         sql.NullString
	)

	err := row.Scan(
		&result.IdPengguna,
		&result.IdSiswa,
		&result.IdNamaKelas,
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
	if err != nil {
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

func (r *ProfilSiswaRepo) buildListSiswaQuery(filter query.ListSiswaFilter) (string, []any) {
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

	if filter.IdNamaKelas != nil {
		args = append(args, *filter.IdNamaKelas)
		where = append(where, fmt.Sprintf("ps.id_nama_kelas = $%d", len(args)))
	}

	if filter.IdTingkatKelas != nil {
		args = append(args, *filter.IdTingkatKelas)
		where = append(where, fmt.Sprintf("k.id_kelas = $%d", len(args)))
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

	return baseQuery, args
}

func (r *ProfilSiswaRepo) scanSiswaListRows(ctx context.Context, op string, rows pgx.Rows) ([]query.SiswaListItem, error) {
	var results []query.SiswaListItem
	for rows.Next() {
		var (
			item         query.SiswaListItem
			jenisKelamin int16
			status       string
			email        sql.NullString
			noHp         sql.NullString
			foto         sql.NullString
			noAbsen      sql.NullInt32
			angkatan     sql.NullInt32
			tempatLahir  sql.NullString
			tanggalLahir sql.NullTime
		)

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
			r.loggerFor(ctx).Error(ctx, "failed scanning siswa list", "op", op, "err", err)
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
		r.loggerFor(ctx).Error(ctx, "failed iterating siswa list", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}

func formatJenisKelamin(value int16) (string, error) {
	switch value {
	case 1:
		return "LAKI_LAKI", nil
	case 2:
		return "PEREMPUAN", nil
	default:
		return "", coreerror.ErrInvalidInput
	}
}
