package postgres

import(
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
)

type ProfilgGuruRepo struct {
	q postgres.Executor
}