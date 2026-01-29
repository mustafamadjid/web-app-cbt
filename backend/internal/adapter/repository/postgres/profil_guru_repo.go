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

func (r *ProfilgGuruRepo) UpdateProfilGuru(ctx context.Context,idPengguna user.ID, profilGuru outuser.UpdateProfilGuruPatch) error{
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
