package postgres

import(
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
)

type UserRepo struct {
	q postgres.Executor
}