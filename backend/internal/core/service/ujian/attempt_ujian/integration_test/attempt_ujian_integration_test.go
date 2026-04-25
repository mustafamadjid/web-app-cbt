package attempt_ujian_test

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	attemptrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/attempt"
	ujianrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/ujian_siswa_checker"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptcreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/create"
	attemptdelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/delete"
	attemptget "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/get"
	attemptlist "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/list_peserta_submitted"
	attemptsubmit "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/submit_ujian"
	attemptupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttemptUjianServices_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	checkerRepo := ujianrepo.NewSiswaUjianCheckerRepo(scope.Pool(), nil)
	attemptRepo := attemptrepo.NewAttemptUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	now := time.Date(2099, time.April, 1, 9, 0, 0, 0, time.UTC)

	createSvc := attemptcreate.NewAttemptUjianService(checkerRepo, attemptRepo)
	require.NoError(t, createSvc.AttemptUjian(scope.Context(), int(exam.Siswa.ID), int(exam.Jadwal.ID), exam.Jadwal.Token, now))

	var idAttempt int64
	err := scope.Pool().QueryRow(scope.Context(), `
		SELECT id_attempt
		FROM attempt_ujian
		WHERE id_peserta_ujian = $1
	`, exam.Peserta.ID).Scan(&idAttempt)
	require.NoError(t, err)

	got, err := attemptget.NewGetAttemptUjianService(attemptRepo).GetAttemptUjianById(scope.Context(), ujian.ID(idAttempt))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(exam.Peserta.ID), got.IdPesertaUjian)
	assert.Equal(t, ujian.ATTEMPT_IN_PROGRESS, got.StatusAttempt)

	newStatus := ujian.ATTEMPT_SUBMITTED
	submitTime := now.Add(45 * time.Minute)
	updateSvc := attemptupdate.NewUpdateAttemptUjianService(attemptRepo)
	require.NoError(t, updateSvc.UpdateAttemptUjian(scope.Context(), ujian.ID(idAttempt), updatepatch.UpdateAttemptUjianPatch{
		StatusAttempt: &newStatus,
		WaktuSubmit:   &submitTime,
	}))

	got, err = attemptget.NewGetAttemptUjianService(attemptRepo).GetAttemptUjianById(scope.Context(), ujian.ID(idAttempt))
	require.NoError(t, err)
	assert.Equal(t, ujian.ATTEMPT_SUBMITTED, got.StatusAttempt)
	require.NotNil(t, got.WaktuSubmit)
	assert.True(t, submitTime.Equal(*got.WaktuSubmit))

	listSubmitted, err := attemptlist.NewPesertaUjianSubmittedService(attemptRepo).ListPesertaUjianSubmitted(scope.Context(), int(exam.Jadwal.ID))
	require.NoError(t, err)
	require.Len(t, listSubmitted, 1)
	assert.Equal(t, ujian.ID(idAttempt), listSubmitted[0].IdAttempt)

	submitSvc := attemptsubmit.NewSubmitUjianService(attemptRepo, checkerRepo)
	require.NoError(t, submitSvc.SubmitUjian(scope.Context(), ujian.ID(idAttempt), int(exam.Siswa.ID)))

	var gradingJobs int
	err = scope.Pool().QueryRow(scope.Context(), `SELECT COUNT(*) FROM grading_jobs WHERE id_attempt = $1`, idAttempt).Scan(&gradingJobs)
	require.NoError(t, err)
	assert.Equal(t, 1, gradingJobs)

	deleteSvc := attemptdelete.NewDeleteAttemptUjianService(attemptRepo)
	require.NoError(t, deleteSvc.DeleteAttemptUjian(scope.Context(), ujian.ID(idAttempt)))

	_, err = attemptget.NewGetAttemptUjianService(attemptRepo).GetAttemptUjianById(scope.Context(), ujian.ID(idAttempt))
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}

func TestAttemptUjianCreateService_RejectsInvalidInput(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	repo := ujianrepo.NewSiswaUjianCheckerRepo(scope.Pool(), nil)
	attemptRepo := attemptrepo.NewAttemptUjianRepo(scope.Pool(), nil)
	svc := attemptcreate.NewAttemptUjianService(repo, attemptRepo)

	err := svc.AttemptUjian(scope.Context(), 0, 0, "", time.Time{})
	assert.ErrorIs(t, err, coreerror.ErrMissingId)
}
