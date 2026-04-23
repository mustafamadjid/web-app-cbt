package aktivitasuserrepo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestAktivitasUserRepo_CreateAktivitasUser(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewAktivitasUserRepo(tx, nil)

	idPengguna := fixtures.LookupSeedUserID("myadmin")
	description := fmt.Sprintf("integration-create-%d", time.Now().UnixNano())
	record := aktivitas_user.AktivitasUser{
		IdPengguna:  idPengguna,
		Action:      aktivitas_user.LOGIN,
		Description: description,
		IpAddress:   "127.0.0.1",
	}

	err := repo.CreateAktivitasUser(ctx, record)
	assert.NoError(t, err)

	var storedID aktivitas_user.AktivitasID
	var storedPenggunaID user.ID
	var storedAction aktivitas_user.Action
	var storedDescription string
	var storedIPAddress string

	err = tx.QueryRow(ctx, `
		SELECT id_aktivitas, id_pengguna, action, description, ip_address
		FROM aktivitas_user
		WHERE description = $1
	`, description).Scan(
		&storedID,
		&storedPenggunaID,
		&storedAction,
		&storedDescription,
		&storedIPAddress,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, storedID)
	assert.Equal(t, idPengguna, storedPenggunaID)
	assert.Equal(t, aktivitas_user.LOGIN, storedAction)
	assert.Equal(t, description, storedDescription)
	assert.Equal(t, "127.0.0.1", storedIPAddress)
}

func TestAktivitasUserRepo_CreateAktivitasUser_InvalidForeignKey(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewAktivitasUserRepo(tx, nil)

	description := fmt.Sprintf("integration-invalid-fk-%d", time.Now().UnixNano())
	var err error
	_, err = tx.Exec(ctx, "SAVEPOINT aktivitas_user_invalid_fk")
	assert.NoError(t, err)

	err = repo.CreateAktivitasUser(ctx, aktivitas_user.AktivitasUser{
		IdPengguna:  user.ID(999999999),
		Action:      aktivitas_user.CREATE,
		Description: description,
		IpAddress:   "127.0.0.2",
	})

	assert.Error(t, err)

	_, err = tx.Exec(ctx, "ROLLBACK TO SAVEPOINT aktivitas_user_invalid_fk")
	assert.NoError(t, err)

	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM aktivitas_user
		WHERE description = $1
	`, description).Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestAktivitasUserRepo_GetAktivitasUser(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewAktivitasUserRepo(tx, nil)

	idPengguna := fixtures.LookupSeedUserID("myadmin")
	baseTime := time.Date(2200, time.January, 1, 10, 0, 0, 0, time.UTC)

	insertAktivitasUserFixture(t, ctx, tx, idPengguna, aktivitas_user.LOGIN, "integration-get-older", "10.10.10.10", baseTime)
	insertAktivitasUserFixture(t, ctx, tx, idPengguna, aktivitas_user.UPDATE, "integration-get-newer", "10.10.10.11", baseTime.Add(time.Minute))

	results, err := repo.GetAktivitasUser(ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.Equal(t, "integration-get-newer", results[0].Description)
	assert.Equal(t, "10.10.10.11", results[0].IpAddress)
	assert.Equal(t, aktivitas_user.UPDATE, results[0].Action)
	assert.Equal(t, "myadmin", results[0].Username)
	assert.Equal(t, user.ADMIN, results[0].Role)

	assert.Equal(t, "integration-get-older", results[1].Description)
	assert.Equal(t, "10.10.10.10", results[1].IpAddress)
	assert.Equal(t, aktivitas_user.LOGIN, results[1].Action)
	assert.Equal(t, "myadmin", results[1].Username)
	assert.Equal(t, user.ADMIN, results[1].Role)

	assert.True(t, results[0].CreatedAt.After(results[1].CreatedAt))
}

func TestAktivitasUserRepo_GetAktivitasUser_Limit30(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewAktivitasUserRepo(tx, nil)

	idPengguna := fixtures.LookupSeedUserID("myadmin")
	baseTime := time.Date(2300, time.January, 1, 10, 0, 0, 0, time.UTC)

	for i := 0; i < 35; i++ {
		insertAktivitasUserFixture(
			t,
			ctx,
			tx,
			idPengguna,
			aktivitas_user.Action(aktivitas_user.CREATE),
			fmt.Sprintf("integration-limit-%02d", i),
			fmt.Sprintf("172.16.0.%d", i),
			baseTime.Add(time.Duration(i)*time.Minute),
		)
	}

	results, err := repo.GetAktivitasUser(ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 30)

	for i, item := range results {
		expectedIndex := 34 - i
		assert.Equal(t, fmt.Sprintf("integration-limit-%02d", expectedIndex), item.Description)
		assert.Equal(t, fmt.Sprintf("172.16.0.%d", expectedIndex), item.IpAddress)
		assert.Equal(t, aktivitas_user.CREATE, item.Action)
		assert.Equal(t, "myadmin", item.Username)
		assert.Equal(t, user.ADMIN, item.Role)
	}
}
func insertAktivitasUserFixture(
	t *testing.T,
	ctx context.Context,
	tx pgx.Tx,
	idPengguna user.ID,
	action aktivitas_user.Action,
	description string,
	ipAddress string,
	createdAt time.Time,
) {
	t.Helper()

	_, err := tx.Exec(ctx, `
		INSERT INTO aktivitas_user (id_pengguna, action, description, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, idPengguna, action, description, ipAddress, createdAt)
	if err != nil {
		t.Fatalf("insert aktivitas user fixture %q: %v", description, err)
	}
}
