package siswaujiancheckerrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSiswaUjianCheckerRepo_Checkers(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewSiswaUjianCheckerRepo(tx, nil)

	exam := fixtures.CreateExamFixture()
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_IN_PROGRESS, nil, nil, nil)

	valid, pesertaID, err := repo.CheckValidSiswaInPesertaUjianById(ctx, int(exam.Siswa.ID), int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.True(t, valid)
	assert.Equal(t, int(exam.Peserta.ID), pesertaID)

	tokenValid, err := repo.CheckTokenUjian(ctx, exam.Jadwal.Token, int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.True(t, tokenValid)

	ownership, err := repo.CheckAttemptOwnershipBySiswa(ctx, int(exam.Siswa.ID), ujian.ID(attempt.ID))
	require.NoError(t, err)
	assert.True(t, ownership)

	deadline, err := repo.GetDeadlineUjian(ctx, int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.True(t, deadline.Equal(exam.Jadwal.WaktuSelesai))
}

func TestSiswaUjianCheckerRepo_GetDeadlineUjian_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewSiswaUjianCheckerRepo(tx, nil)

	_, err := repo.GetDeadlineUjian(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
