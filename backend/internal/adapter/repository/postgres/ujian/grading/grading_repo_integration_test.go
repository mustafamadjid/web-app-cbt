package gradingrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGradingRepo_ClaimAndUpdateStatusJob(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewGradingRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, nil, nil, nil)
	job := fixtures.CreateGradingJob(attempt.ID, ujian.StatusQueued)

	claimed, err := repo.ClaimQueuedJobs(scope.Context(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, claimed)
	found := false
	for _, item := range claimed {
		if int64(item.IDgradingJob) == job.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	err = repo.UpdateStatusJob(scope.Context(), int(job.ID), ujian.StatusDone, "", "")
	require.NoError(t, err)

	var status string
	err = scope.Pool().QueryRow(scope.Context(), `SELECT status FROM grading_jobs WHERE id_grading_jobs = $1`, job.ID).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, string(ujian.StatusDone), status)
}

func TestGradingRepo_UpsertNilaiAndStatistik(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewGradingRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, nil, nil, nil)
	essayTrue := true
	fixtures.CreateJawaban(attempt.ID, exam.SoalEssay.ID, nil, testutil.Ptr("essay benar"), nil, &essayTrue)

	passed := true
	essayGraded := true
	gradedBy := ujian.ID(exam.Guru.ID)
	err := repo.UpsertNilaiToHasilUjian(scope.Context(), 40, ujian.HasilUjian{
		IdAttempt:   ujian.ID(attempt.ID),
		GradedBy:    &gradedBy,
		Passed:      &passed,
		EssayGraded: &essayGraded,
	})
	require.NoError(t, err)

	var nilai float64
	var total int
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT nilai_akhir
		FROM hasil_ujian
		WHERE id_attempt = $1
	`, attempt.ID).Scan(&nilai)
	require.NoError(t, err)
	assert.Equal(t, 55.0, nilai)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT total_peserta_ujian
		FROM statistik_ujian
		WHERE id_jadwal_ujian = $1
	`, exam.Jadwal.ID).Scan(&total)
	require.NoError(t, err)
	assert.Equal(t, 1, total)

	err = repo.UpsertJawabanBenarToStatistikSoal(scope.Context(), []ujian.StatistikSoal{
		{IDSoal: ujian.ID(exam.SoalPilgan.ID), IDUjian: ujian.ID(exam.Ujian.ID)},
		{IDSoal: ujian.ID(exam.SoalPilgan.ID), IDUjian: ujian.ID(exam.Ujian.ID)},
	})
	require.NoError(t, err)

	err = repo.UpsertJawabanSalahToStatistikSoal(scope.Context(), []ujian.StatistikSoal{
		{IDSoal: ujian.ID(exam.SoalPilgan.ID), IDUjian: ujian.ID(exam.Ujian.ID)},
	})
	require.NoError(t, err)

	var benar int
	var salah int
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT jumlah_jawaban_benar, jumlah_jawaban_salah
		FROM statistik_soal
		WHERE id_soal = $1 AND id_ujian = $2
	`, exam.SoalPilgan.ID, exam.Ujian.ID).Scan(&benar, &salah)
	require.NoError(t, err)
	assert.Equal(t, 2, benar)
	assert.Equal(t, 1, salah)
}

func TestGradingRepo_UpdateAndGradingEssayUjian(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewGradingRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, nil, nil, nil)
	jawabanID := fixtures.CreateJawaban(attempt.ID, exam.SoalEssay.ID, nil, testutil.Ptr("essay"), nil, nil)
	gradedBy := ujian.ID(exam.Guru.ID)
	essayTrue := true

	err := repo.UpdateAndGradingEssayUjian(scope.Context(), []ujian.JawabanUjian{
		{IdJawaban: ujian.ID(jawabanID), EssayIsBenar: &essayTrue},
	}, gradedBy)
	require.NoError(t, err)

	var stored bool
	var essayGraded bool
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT essay_is_benar
		FROM jawaban_ujian_siswa
		WHERE id_jawaban = $1
	`, jawabanID).Scan(&stored)
	require.NoError(t, err)
	assert.True(t, stored)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT essay_graded
		FROM hasil_ujian
		WHERE id_attempt = $1
	`, attempt.ID).Scan(&essayGraded)
	require.NoError(t, err)
	assert.True(t, essayGraded)
}

func TestGradingRepo_UpdateAndGradingEssayUjian_RequiresPool(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewGradingRepo(tx, nil)
	essayTrue := true

	err := repo.UpdateAndGradingEssayUjian(ctx, []ujian.JawabanUjian{{IdJawaban: 1, EssayIsBenar: &essayTrue}}, ujian.ID(user.ID(1)))
	assert.EqualError(t, err, "grading repo requires pgx pool for update transaction")
	assert.NotErrorIs(t, err, coreerror.ErrNotFound)
}
