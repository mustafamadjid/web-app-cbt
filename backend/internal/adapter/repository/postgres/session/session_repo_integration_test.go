package sessionrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionRepo_CreateAndQuerySession(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewSessionRepo(tx, nil)

	createdUser := fixtures.CreateUser(user.GURU)
	expiresAt := time.Now().Add(2 * time.Hour)

	sessionID, err := repo.CreateSession(ctx, createdUser.ID, user.GURU, expiresAt)
	require.NoError(t, err)
	assert.NotEmpty(t, sessionID)

	got, err := repo.GetSession(ctx, sessionID)
	require.NoError(t, err)
	assert.Equal(t, sessionID, got.SessionID)
	assert.Equal(t, createdUser.ID, got.UserID)
	assert.Equal(t, user.GURU, got.Role)
	assert.False(t, got.Revoked)

	gotByUser, err := repo.GetSessionByUserId(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, sessionID, gotByUser.SessionID)

	hasActive, err := repo.HasActiveSession(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.True(t, hasActive)

	active, err := repo.GetAllActiveSession(ctx)
	require.NoError(t, err)
	assert.Len(t, active, 1)
	assert.Equal(t, createdUser.Username, active[0].Pengguna.Username)
}

func TestSessionRepo_CreateSession_ConflictAndRevoke(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewSessionRepo(scope.Pool(), nil)

	createdUser := fixtures.CreateUser(user.SISWA)
	_, err := repo.CreateSession(scope.Context(), createdUser.ID, user.SISWA, time.Now().Add(time.Hour))
	require.NoError(t, err)

	_, err = repo.CreateSession(scope.Context(), createdUser.ID, user.SISWA, time.Now().Add(2*time.Hour))
	assert.ErrorIs(t, err, coreerror.ErrHasSession)

	active, err := repo.GetSessionByUserId(scope.Context(), createdUser.ID)
	require.NoError(t, err)

	err = repo.RevokeSession(scope.Context(), active.SessionID)
	require.NoError(t, err)

	_, err = repo.GetSessionByUserId(scope.Context(), createdUser.ID)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func TestSessionRepo_RevokeExpiredAndRevokeAllByUser(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewSessionRepo(tx, nil)

	expiredUser := fixtures.CreateUser(user.ADMIN)
	fixtures.CreateSession(expiredUser.ID, user.ADMIN, time.Now().Add(-time.Hour), nil)

	revoked, err := repo.RevokeExpiredSessions(ctx, expiredUser.ID)
	require.NoError(t, err)
	assert.True(t, revoked)

	activeUser := fixtures.CreateUser(user.ADMIN)
	sessionOne := fixtures.CreateSession(activeUser.ID, user.ADMIN, time.Now().Add(time.Hour), nil)
	sessionTwo := fixtures.CreateSession(activeUser.ID, user.ADMIN, time.Now().Add(2*time.Hour), nil)
	assert.NotEqual(t, sessionOne.ID, sessionTwo.ID)

	err = repo.RevokeSessionAllbyUser(ctx, activeUser.ID)
	require.NoError(t, err)

	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM sessions
		WHERE id_pengguna = $1 AND revoked_at IS NULL
	`, activeUser.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
