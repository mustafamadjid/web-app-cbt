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

type ProfilgGuruRepo struct {
	q Executor
}

func NewProfilgGuruRepo(q Executor) *ProfilgGuruRepo {
	return &ProfilgGuruRepo{q: q}
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

func (r *ProfilgGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru outuser.UpdateProfilGuruPatch) error {
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

	var results []query.GuruListItem
	for rows.Next() {
		var item query.GuruListItem
		var jenisKelamin int16
		var status string
		var email string
		var role string
		var nip string

		if err := rows.Scan(
			&item.IdPengguna,
			&role,
			&item.Username,
			&item.NoHp,
			&email,
			&item.NamaLengkap,
			&status,
			&nip,
			&item.Jabatan,
			&item.BidangStudi,
			&item.Foto,
			&jenisKelamin,
		); err != nil {
			return nil, err
		}

		jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
		if err != nil {
			return nil, err
		}

		item.Email = user.Email(email)
		item.Nip = user.NIP(nip)
		item.JenisKelamin = jenisKelaminValue
		item.StatusAkun = user.StatusAkun(status)
		item.Role = user.Role(role)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
