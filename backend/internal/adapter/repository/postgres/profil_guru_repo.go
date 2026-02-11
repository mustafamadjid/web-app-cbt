package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ProfilgGuruRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewProfilgGuruRepo(q Executor, logger corelog.Logger) *ProfilgGuruRepo {
	return &ProfilgGuruRepo{q: q, logger: logger}
}

func (r *ProfilgGuruRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ProfilgGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	const query = `
		SELECT 
			p.id_pengguna,
			p.id_guru,
			u.username,
			u.email,
			u.nama_lengkap,
			u.jenis_kelamin,
			u.no_hp,
			u.foto,
			r.nama_role,
			u.status_akun,
			p.nip,
			p.jabatan,
			p.bidang_studi
		FROM profil_guru p
		JOIN pengguna u ON p.id_pengguna = u.id_pengguna
		JOIN role r ON u.id_role = r.id_role
		WHERE p.id_pengguna = $1
	`

	var result user.DataGuru

	var nip sql.NullString
	var jabatan sql.NullString
	var bidangStudi sql.NullString
	var jenisKelamin int16

	err := r.q.QueryRow(ctx, query, id).Scan(
		&result.IdPengguna,
		&result.IdGuru,
		&result.Username,
		&result.Email,
		&result.NamaLengkap,
		&jenisKelamin,
		&result.NoHp,
		&result.Foto,
		&result.Role,
		&result.StatusAkun,
		&nip,
		&jabatan,
		&bidangStudi,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return user.DataGuru{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed finding profil guru", "op", "profil_guru_repo.find_by_id", "id_guru", id, "err", err)
		return user.DataGuru{}, err
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

func (r *ProfilgGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM profil_guru WHERE nip = $1)`

	var exists bool
	if err := r.q.QueryRow(ctx, query, string(nip)).Scan(&exists); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking nip existence", "op", "profil_guru_repo.exists_by_nip", "nip", nip, "err", err)
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
		r.loggerFor(ctx).Error(ctx, "failed creating profil guru", "op", "profil_guru_repo.create", "user_id", profilGuru.IdPengguna, "err", err)
		return 0, err
	}

	return id, nil
}

func (r *ProfilgGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru updatepatch.ProfilGuru) error {
	set := make([]string, 0, 4)
	args := make([]any, 0, 5)

	add := func(col string, v any) {
		args = append(args, v)
		set = append(set, fmt.Sprintf("%s=$%d", col, len(args)))
	}

	if profilGuru.Nip != nil {
		add("nip", *profilGuru.Nip)
	}
	if profilGuru.Jabatan != nil {
		add("jabatan", *profilGuru.Jabatan)
	}
	if profilGuru.BidangStudi != nil {
		add("bidang_studi", *profilGuru.BidangStudi)
	}

	if len(set) == 0 {
		return nil
	}

	args = append(args, idPengguna)
	q := fmt.Sprintf(`UPDATE profil_guru SET %s WHERE id_pengguna=$%d`, strings.Join(set, ", "), len(args))

	_, err := r.q.Exec(ctx, q, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating profil guru", "op", "profil_guru_repo.update", "user_id", idPengguna, "err", err)
	}
	return err
}

func (r *ProfilgGuruRepo) GetListGuru(ctx context.Context, filter query.ListGuruFilter) ([]query.GuruListItem, error) {
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

	rows, err := r.q.Query(ctx, baseQuery, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing guru", "op", "profil_guru_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	var results []query.GuruListItem
	for rows.Next() {
		var item query.GuruListItem
		var jenisKelamin int16
		var status string
		var email sql.NullString
		var role string
		var nip sql.NullString
		var jabatan sql.NullString
		var bidangStudi sql.NullString
		var noHp sql.NullString
		var foto sql.NullString

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
			r.loggerFor(ctx).Error(ctx, "failed scanning guru list", "op", "profil_guru_repo.list_scan", "err", err)
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
		r.loggerFor(ctx).Error(ctx, "failed iterating guru list", "op", "profil_guru_repo.list_iter", "err", err)
		return nil, err
	}

	return results, nil
}
