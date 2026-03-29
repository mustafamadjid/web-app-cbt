package authuserrepo

import (
	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func scanAuthUserRow(row pgx.Row) (user.Pengguna, error) {
	var item user.Pengguna
	if err := row.Scan(&item.ID, &item.Username, &item.PasswordHashed, &item.Role, &item.StatusAkun); err != nil {
		return user.Pengguna{}, err
	}

	item.Role = user.Role(string(item.Role))
	item.StatusAkun = user.StatusAkun(string(item.StatusAkun))

	return item, nil
}
