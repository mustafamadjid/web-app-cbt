package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	outauth "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/auth_port_out"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
)

type InfraModule struct {
	Pool *pgxpool.Pool
	Txm  txout.TxManager

	Sessions out.SessionRepository
	AuthUsers outauth.AuthUserrepository
	users 	outuser.UserRepository
}

func BuildInfraModule(pool *pgxpool.Pool) *InfraModule {
	txm := pg.NewTxManager(pool)

	return &InfraModule{
		Pool:     pool,
		Txm:      txm,
		Sessions: pg.NewSessionRepo(pool),
		AuthUsers: pg.NewAuthUserRepo(pool),
		users : pg.NewUserRepo(pool),
	}
}
