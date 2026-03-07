package ujianrepo

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/stretchr/testify/require"
)

type fakeExecutor struct {
	querySQL  string
	queryArgs []any
	rows      pgx.Rows
	err       error
}

func (f *fakeExecutor) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (f *fakeExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	f.querySQL = sql
	f.queryArgs = args
	if f.err != nil {
		return nil, f.err
	}
	return f.rows, nil
}

func (f *fakeExecutor) QueryRow(context.Context, string, ...any) pgx.Row {
	return nil
}

type fakeRows struct {
	rows   [][]any
	index  int
	closed bool
	err    error
}

func (f *fakeRows) Close() {
	f.closed = true
}

func (f *fakeRows) Err() error {
	return f.err
}

func (f *fakeRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (f *fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (f *fakeRows) Next() bool {
	if f.index >= len(f.rows) {
		f.closed = true
		return false
	}

	f.index++
	return true
}

func (f *fakeRows) Scan(dest ...any) error {
	row := f.rows[f.index-1]
	for i, value := range row {
		target := reflect.ValueOf(dest[i])
		target.Elem().Set(reflect.ValueOf(value))
	}
	return nil
}

func (f *fakeRows) Values() ([]any, error) {
	return f.rows[f.index-1], nil
}

func (f *fakeRows) RawValues() [][]byte {
	return nil
}

func (f *fakeRows) Conn() *pgx.Conn {
	return nil
}

func TestListSoalUjianRepo_GetSoalUjianByBankSoal(t *testing.T) {
	rows := &fakeRows{
		rows: [][]any{
			{
				ujian.ID(20),
				ujian.ID(1001),
				"PILIHAN_GANDA",
				"Soal pertama",
				sql.NullString{String: "gambar-1.png", Valid: true},
				5,
				1,
				sql.NullInt64{Int64: 200, Valid: true},
				sql.NullString{String: "Opsi A", Valid: true},
				sql.NullBool{Bool: true, Valid: true},
			},
			{
				ujian.ID(20),
				ujian.ID(1001),
				"PILIHAN_GANDA",
				"Soal pertama",
				sql.NullString{String: "gambar-1.png", Valid: true},
				5,
				1,
				sql.NullInt64{Int64: 201, Valid: true},
				sql.NullString{String: "Opsi B", Valid: true},
				sql.NullBool{Bool: false, Valid: true},
			},
			{
				ujian.ID(10),
				ujian.ID(1002),
				"ESSAY",
				"Soal kedua",
				sql.NullString{},
				10,
				2,
				sql.NullInt64{},
				sql.NullString{},
				sql.NullBool{},
			},
		},
	}
	executor := &fakeExecutor{rows: rows}
	repo := NewListSoalUjianRepo(executor, corelog.FromContext(context.Background()))

	items, err := repo.GetSoalUjianByBankSoal(context.Background(), ujian.ID(7))

	require.NoError(t, err)
	require.Equal(t, []any{ujian.ID(7)}, executor.queryArgs)
	require.Len(t, items, 2)
	require.Equal(t, ujian.ID(20), items[0].IdSoal)
	require.Equal(t, ujian.ID(10), items[1].IdSoal)
	require.Len(t, items[0].OpsiJawaban, 2)
	require.Len(t, items[1].OpsiJawaban, 0)
	require.Equal(t, "gambar-1.png", items[0].Gambar)
}

var _ pg.Executor = (*fakeExecutor)(nil)
