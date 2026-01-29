package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type UserRepo struct {
	q Executor
}
func NewUserRepo(q Executor) *UserRepo {
	return &UserRepo{q: q}
}


func (r *UserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	const query = `
		SELECT p.id_pengguna,
			p.username,
			p.email,
			p.password,
			p.nama_lengkap,
			p.jenis_kelamin,
			p.no_hp,
			r.nama_role,
			p.status_akun,
			p.foto
		FROM pengguna p
		JOIN role r ON p.id_role = r.id_role
		WHERE p.id_pengguna = $1
	`

	var result user.Pengguna
	var jenisKelamin int16
	var roleName string
	var status string
	var email string

	err := r.q.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Username,
		&email,
		&result.PasswordHashed,
		&result.NamaLengkap,
		&jenisKelamin,
		&result.NoHp,
		&roleName,
		&status,
		&result.Foto,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	if err != nil {
		return user.Pengguna{}, err
	}

	jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
	if err != nil {
		return user.Pengguna{}, err
	}

	result.Email = user.Email(email)
	result.JenisKelamin = jenisKelaminValue
	result.Role = user.Role(roleName)
	result.StatusAkun = user.StatusAkun(status)

	return result, nil
}

func (r *UserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM pengguna WHERE username = $1)`

	var exists bool
	if err := r.q.QueryRow(ctx, query, username).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error) {
	const query = `
		INSERT INTO pengguna (
			foto,
			nama_lengkap,
			jenis_kelamin,
			username,
			password,
			email,
			no_hp,
			id_role,
			status_akun
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			(SELECT id_role FROM role WHERE nama_role = $8),
			$9
		)
		RETURNING id_pengguna
	`

	jenisKelamin, err := parseJenisKelamin(pengguna.JenisKelamin)
	if err != nil {
		return 0, err
	}

	var id user.ID
	err = r.q.QueryRow(
		ctx,
		query,
		pengguna.Foto,
		pengguna.NamaLengkap,
		jenisKelamin,
		pengguna.Username,
		pengguna.PasswordHashed,
		string(pengguna.Email),
		pengguna.NoHp,
		string(pengguna.Role),
		string(pengguna.StatusAkun),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context,idPengguna user.ID,pengguna outuser.UpdatePenggunaPatch) error{
	set := make([]string, 0, 6)
	args := make([]any, 0, 7)

	add := func(col string, v any) {
		args = append(args, v)
		set = append(set, fmt.Sprintf("%s=$%d", col, len(args)))
	}

	if pengguna.NamaLengkap != nil {
		add("nama_lengkap", *pengguna.NamaLengkap)
	}
	if pengguna.Email != nil {
		add("email", string(*pengguna.Email))
	}
	if pengguna.NoHp != nil {
		add("no_hp", *pengguna.NoHp)
	}
	if pengguna.Foto != nil {
		add("foto", *pengguna.Foto)
	}
	if pengguna.StatusAkun != nil {
		add("status_akun", string(*pengguna.StatusAkun))
	}
	if pengguna.Role != nil {
		add("id_role", *pengguna.Role)
	}

	if len(set) == 0 {
		return nil
	}

	args = append(args, idPengguna)
	q := fmt.Sprintf(`UPDATE pengguna SET %s WHERE id_pengguna=$%d`, strings.Join(set, ", "), len(args))

	_, err := r.q.Exec(ctx, q, args...)
	return err
}

func (r *UserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	const query = `DELETE FROM pengguna WHERE id_pengguna = $1`

	result, err := r.q.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *UserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	const query = `
		SELECT p.id_pengguna,
			p.username,
			p.email,
			p.password,
			p.nama_lengkap,
			p.jenis_kelamin,
			p.no_hp,
			r.nama_role,
			p.status_akun,
			p.foto
		FROM pengguna p
		JOIN role r ON p.id_role = r.id_role
		ORDER BY p.id_pengguna
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []user.Pengguna
	for rows.Next() {
		var item user.Pengguna
		var jenisKelamin int16
		var roleName string
		var status string
		var email string

		if err := rows.Scan(
			&item.ID,
			&item.Username,
			&email,
			&item.PasswordHashed,
			&item.NamaLengkap,
			&jenisKelamin,
			&item.NoHp,
			&roleName,
			&status,
			&item.Foto,
		); err != nil {
			return nil, err
		}

		jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
		if err != nil {
			return nil, err
		}

		item.Email = user.Email(email)
		item.JenisKelamin = jenisKelaminValue
		item.Role = user.Role(roleName)
		item.StatusAkun = user.StatusAkun(status)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func parseJenisKelamin(value string) (int16, error) {
	normalized := strings.TrimSpace(strings.ToUpper(value))
	if normalized == "" {
		return 0, coreerror.ErrInvalidInput
	}

	if numeric, err := strconv.Atoi(normalized); err == nil {
		if numeric == 1 || numeric == 2 {
			return int16(numeric), nil
		}
	}

	switch normalized {
	case "L", "LK", "LAKI", "LAKI_LAKI", "LAKI-LAKI", "PRIA", "MALE":
		return 1, nil
	case "P", "PR", "PEREMPUAN", "WANITA", "FEMALE":
		return 2, nil
	default:
		return 0, coreerror.ErrInvalidInput
	}
}

func formatJenisKelamin(value int16) (string, error) {
	switch value {
	case 1:
		return "L", nil
	case 2:
		return "P", nil
	default:
		return "", coreerror.ErrInvalidInput
	}
}
