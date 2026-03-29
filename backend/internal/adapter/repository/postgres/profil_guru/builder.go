package profilgururepo

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

type guruScanner interface {
	Scan(dest ...any) error
}

func scanProfilGuruRow(row guruScanner) (user.DataGuru, error) {
	var (
		result       user.DataGuru
		email        sql.NullString
		noHp         sql.NullString
		foto         sql.NullString
		nip          sql.NullString
		jabatan      sql.NullString
		bidangStudi  sql.NullString
		jenisKelamin int16
	)

	err := row.Scan(
		&result.IdPengguna,
		&result.IdGuru,
		&result.Username,
		&email,
		&result.NamaLengkap,
		&jenisKelamin,
		&noHp,
		&foto,
		&result.Role,
		&result.StatusAkun,
		&nip,
		&jabatan,
		&bidangStudi,
	)
	if err != nil {
		return user.DataGuru{}, err
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
	if nip.Valid {
		result.Nip = nip.String
	}
	if jabatan.Valid {
		result.Jabatan = jabatan.String
	}
	if bidangStudi.Valid {
		result.BidangStudi = bidangStudi.String
	}

	jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
	if err != nil {
		return user.DataGuru{}, err
	}
	result.JenisKelamin = jenisKelaminValue

	return result, nil
}

func (r *ProfilgGuruRepo) buildListGuruQuery(filter query.ListGuruFilter) (string, []any) {
	sortColumns := map[string]string{
		"nama_lengkap": "p.nama_lengkap",
		"created_at":   "p.created_at",
		"username":     "p.username",
		"nip":          "pg.nip",
	}

	baseQuery := `
		SELECT p.id_pengguna,
			r.nama_role,
			p.username,
			p.no_hp,
			p.email,
			p.nama_lengkap,
			p.status_akun,
			pg.nip,
			pg.jabatan,
			pg.bidang_studi,
			p.foto,
			p.jenis_kelamin
		FROM pengguna p
		JOIN profil_guru pg ON pg.id_pengguna = p.id_pengguna
		JOIN role r ON p.id_role = r.id_role
	`

	where := make([]string, 0, 4)
	args := make([]any, 0, 6)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf(`(p.nama_lengkap ILIKE $%d OR p.username ILIKE $%d OR p.email ILIKE $%d OR pg.nip ILIKE $%d)`, idx, idx, idx, idx))
	}

	if filter.Status != nil {
		args = append(args, string(*filter.Status))
		where = append(where, fmt.Sprintf("p.status_akun = $%d", len(args)))
	}

	if filter.Bidang != nil {
		args = append(args, "%"+*filter.Bidang+"%")
		where = append(where, fmt.Sprintf("pg.bidang_studi ILIKE $%d", len(args)))
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

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)
	} else if filter.Offset > 0 {
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s OFFSET $%d", baseQuery, offsetIndex)
	}

	return baseQuery, args
}

func (r *ProfilgGuruRepo) scanGuruListRows(ctx context.Context, op string, rows pgx.Rows) ([]query.GuruListItem, error) {
	var results []query.GuruListItem
	for rows.Next() {
		var (
			item         query.GuruListItem
			jenisKelamin int16
			status       string
			email        sql.NullString
			role         string
			nip          sql.NullString
			jabatan      sql.NullString
			bidangStudi  sql.NullString
			noHp         sql.NullString
			foto         sql.NullString
		)

		if err := rows.Scan(
			&item.IdPengguna,
			&role,
			&item.Username,
			&noHp,
			&email,
			&item.NamaLengkap,
			&status,
			&nip,
			&jabatan,
			&bidangStudi,
			&foto,
			&jenisKelamin,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning guru list", "op", op, "err", err)
			return nil, err
		}

		jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
		if err != nil {
			return nil, err
		}

		if email.Valid {
			item.Email = user.Email(email.String)
		}
		if nip.Valid {
			item.Nip = user.NIP(nip.String)
		}
		item.JenisKelamin = jenisKelaminValue
		item.StatusAkun = user.StatusAkun(status)
		item.Role = user.Role(role)
		if jabatan.Valid {
			item.Jabatan = jabatan.String
		}
		if bidangStudi.Valid {
			item.BidangStudi = bidangStudi.String
		}
		if noHp.Valid {
			item.NoHp = noHp.String
		}
		if foto.Valid {
			item.Foto = foto.String
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating guru list", "op", op, "err", err)
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
