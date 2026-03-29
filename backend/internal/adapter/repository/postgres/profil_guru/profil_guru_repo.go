package profilgururepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

type ProfilgGuruRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewProfilgGuruRepo(q pg.Executor, logger corelog.Logger) *ProfilgGuruRepo {
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

	result, err := scanProfilGuruRow(r.q.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return user.DataGuru{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed finding profil guru", "op", "profil_guru_repo.find_by_id", "id_guru", id, "err", err)
		return user.DataGuru{}, err
	}

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
	queryText, args := r.buildListGuruQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing guru", "op", "profil_guru_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanGuruListRows(ctx, "profil_guru_repo.list", rows)
}
