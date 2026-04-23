package attemptrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttemptUjianRepo_CRUDAndListSubmitted(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewAttemptUjianRepo(tx, nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)

	err := repo.CreateAttemptUjian(ctx, ujian.AttemptUjian{
		IdPesertaUjian: ujian.ID(exam.Peserta.ID),
		StatusAttempt:  ujian.ATTEMPT_IN_PROGRESS,
		WaktuMulai:     &start,
		DeadlineAt:     testutil.Ptr(start.Add(time.Hour)),
	})
	require.NoError(t, err)

	var idAttempt int
	err = tx.QueryRow(ctx, `SELECT id_attempt FROM attempt_ujian WHERE id_peserta_ujian = $1 AND status_attempt = 'in_progress'`, exam.Peserta.ID).Scan(&idAttempt)
	require.NoError(t, err)

	item, err := repo.GetAttemptUjianById(ctx, ujian.ID(idAttempt))
	require.NoError(t, err)
	assert.Equal(t, ujian.ATTEMPT_IN_PROGRESS, item.StatusAttempt)

	submitTime := start.Add(45 * time.Minute)
	status := ujian.ATTEMPT_SUBMITTED
	err = repo.UpdateAttemptUjian(ctx, ujian.ID(idAttempt), updatepatch.UpdateAttemptUjianPatch{
		StatusAttempt: &status,
		WaktuSubmit:   &submitTime,
	})
	require.NoError(t, err)

	gradedBy := ujian.ID(exam.Guru.ID)
	nilai := 88.0
	passed := true
	fixtures.CreateHasilUjian(int64(idAttempt), testutil.Ptr(userIDFromUjianID(gradedBy)), &nilai, &passed, nil, &submitTime, exam.Jadwal.ID)

	listed, err := repo.ListPesertaUjianAttemptSubmittedByIdJadwalUjian(ctx, ujian.ID(exam.Jadwal.ID))
	require.NoError(t, err)
	require.Len(t, listed, 1)
	assert.Equal(t, ujian.ID(idAttempt), listed[0].IdAttempt)

	err = repo.DeleteAttemptUjian(ctx, ujian.ID(idAttempt))
	require.NoError(t, err)

	_, err = repo.GetAttemptUjianById(ctx, ujian.ID(idAttempt))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func TestAttemptUjianRepo_CreateAttempt_ConflictAndSubmit(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewAttemptUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)

	err := repo.CreateAttemptUjian(scope.Context(), ujian.AttemptUjian{
		IdPesertaUjian: ujian.ID(exam.Peserta.ID),
		StatusAttempt:  ujian.ATTEMPT_IN_PROGRESS,
		WaktuMulai:     &start,
		DeadlineAt:     testutil.Ptr(start.Add(time.Hour)),
	})
	require.NoError(t, err)

	err = repo.CreateAttemptUjian(scope.Context(), ujian.AttemptUjian{
		IdPesertaUjian: ujian.ID(exam.Peserta.ID),
		StatusAttempt:  ujian.ATTEMPT_IN_PROGRESS,
		WaktuMulai:     &start,
		DeadlineAt:     testutil.Ptr(start.Add(2 * time.Hour)),
	})
	assert.ErrorIs(t, err, coreerror.ErrSiswaHasActiveAttempt)

	var idAttempt int64
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT id_attempt
		FROM attempt_ujian
		WHERE id_peserta_ujian = $1 AND status_attempt = 'in_progress'
	`, exam.Peserta.ID).Scan(&idAttempt)
	require.NoError(t, err)

	err = repo.SubmitAttemptUjian(scope.Context(), ujian.ID(idAttempt))
	require.NoError(t, err)

	var status string
	var gradingJobs int
	err = scope.Pool().QueryRow(scope.Context(), `SELECT status_attempt FROM attempt_ujian WHERE id_attempt = $1`, idAttempt).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, string(ujian.ATTEMPT_SUBMITTED), status)

	err = scope.Pool().QueryRow(scope.Context(), `SELECT COUNT(*) FROM grading_jobs WHERE id_attempt = $1`, idAttempt).Scan(&gradingJobs)
	require.NoError(t, err)
	assert.Equal(t, 1, gradingJobs)
}

func userIDFromUjianID(id ujian.ID) user.ID {
	return user.ID(id)
}
