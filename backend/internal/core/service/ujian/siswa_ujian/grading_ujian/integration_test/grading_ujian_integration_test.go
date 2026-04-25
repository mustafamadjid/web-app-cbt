package grading_ujian_test

import (
	"context"
	"testing"
	"time"

	banksoalrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/bank_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	gradinglistrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/grading/list"
	gradingrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/grading"
	jawabanujianrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/siswa/jawaban_ujian"
	ujianlistrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/list"
	ujianrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian"
	attemptrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/attempt"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	gradinglistsvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/list/list_ujian_essay_ungraded"
	gradingstatsvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/statistik_ujian"
	gradingworker "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/worker"
	essaygradingsvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading"
	gradingservice "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGradingExecutor struct {
	pilganCalls int
	statsCalls  int
}

func (f *fakeGradingExecutor) GradingUjianPilgan(context.Context, int) error {
	f.pilganCalls++
	return nil
}

func (f *fakeGradingExecutor) StatistikUjian(context.Context, int) error {
	f.statsCalls++
	return nil
}

func TestGradingUjianService_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	jawabanRepo := jawabanujianrepo.NewJawabanUjianRepo(scope.Pool(), nil)
	soalRepo := ujianlistrepo.NewListSoalUjianRepo(scope.Pool(), nil)
	bankRepo := banksoalrepo.NewBankSoalRepo(scope.Pool(), nil)
	ujianRepo := ujianrepo.NewUjianRepo(scope.Pool(), nil, scope.Pool())
	gradingRepo := gradingrepo.NewGradingRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	now := time.Date(2099, time.April, 3, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &now, testutil.Ptr(now.Add(45*time.Minute)), testutil.Ptr(now.Add(time.Hour)))
	require.NoError(t, jawabanRepo.SaveJawabanUjian(scope.Context(), ujian.ID(attempt.ID), []ujian.JawabanUjian{
		{IdSoal: ujian.ID(exam.SoalPilgan.ID), IdPilihan: testutil.Ptr(ujian.ID(exam.OpsiBenar.ID)), WaktuJawab: testutil.Ptr(now.Add(10 * time.Minute))},
	}))

	svc := gradingservice.NewGradingUjianService(jawabanRepo, soalRepo, bankRepo, ujianRepo, gradingRepo)
	require.NoError(t, svc.GradingUjianPilgan(scope.Context(), int(attempt.ID)))
	require.NoError(t, gradingstatsvc.NewStatistikUjianService(gradingRepo).StatistikUjian(scope.Context(), int(attempt.ID)))

	var hasilCount int
	err := scope.Pool().QueryRow(scope.Context(), `SELECT COUNT(*) FROM hasil_ujian WHERE id_attempt = $1`, attempt.ID).Scan(&hasilCount)
	require.NoError(t, err)
	assert.Equal(t, 1, hasilCount)

	var benar int
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT COALESCE(jumlah_jawaban_benar, 0)
		FROM statistik_soal
		WHERE id_soal = $1 AND id_ujian = $2
	`, exam.SoalPilgan.ID, exam.Ujian.ID).Scan(&benar)
	require.NoError(t, err)
	assert.Equal(t, 1, benar)

	var totalPeserta int
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT total_peserta_ujian
		FROM statistik_ujian
		WHERE id_jadwal_ujian = $1
	`, exam.Jadwal.ID).Scan(&totalPeserta)
	require.NoError(t, err)
	assert.Equal(t, 1, totalPeserta)
}

func TestEssayGradingAndUngradedList_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	gradingRepo := gradingrepo.NewGradingRepo(scope.Pool(), nil)
	listRepo := gradinglistrepo.NewListGradingRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	now := time.Date(2099, time.April, 3, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &now, testutil.Ptr(now.Add(45*time.Minute)), testutil.Ptr(now.Add(time.Hour)))
	jawabanID := fixtures.CreateJawaban(attempt.ID, exam.SoalEssay.ID, nil, testutil.Ptr("essay integration"), nil, nil)
	essayTrue := true

	svc := essaygradingsvc.NewEssayGradingUjianService(gradingRepo)
	require.NoError(t, svc.EssayGrading(scope.Context(), []ujian.JawabanUjian{{IdJawaban: ujian.ID(jawabanID), EssayIsBenar: &essayTrue}}, ujian.ID(exam.Guru.ID)))

	var stored bool
	err := scope.Pool().QueryRow(scope.Context(), `SELECT essay_is_benar FROM jawaban_ujian_siswa WHERE id_jawaban = $1`, jawabanID).Scan(&stored)
	require.NoError(t, err)
	assert.True(t, stored)

	var essayGraded bool
	err = scope.Pool().QueryRow(scope.Context(), `SELECT essay_graded FROM hasil_ujian WHERE id_attempt = $1`, attempt.ID).Scan(&essayGraded)
	require.NoError(t, err)
	assert.True(t, essayGraded)

	examUngraded := fixtures.CreateExamFixture()
	nowUngraded := now.Add(time.Hour)
	attemptUngraded := fixtures.CreateAttempt(examUngraded.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &nowUngraded, testutil.Ptr(nowUngraded.Add(45*time.Minute)), testutil.Ptr(nowUngraded.Add(time.Hour)))
	fixtures.CreateHasilUjian(attemptUngraded.ID, nil, testutil.Ptr(70.0), nil, testutil.Ptr(false), nil, examUngraded.Jadwal.ID)
	tingkatKelasID := int(exam.Kelas.ID)
	namaKelasID := int(exam.NamaKelas.ID)
	sesiID := int(exam.Sesi.ID)
	items, err := gradinglistsvc.NewListUjianEssayUngradedService(listRepo).ListUjianEssayUngraded(scope.Context(), query.ListUjianEssayUngradedFilter{
		Limit:          10,
		TingkatKelasID: &tingkatKelasID,
		NamaKelasID:    &namaKelasID,
		SesiID:         &sesiID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
	assert.Equal(t, ujian.ID(examUngraded.Ujian.ID), items[0].IdUjian)
}

func TestStatistikUjianAndWorker_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	gradingRepo := gradingrepo.NewGradingRepo(scope.Pool(), nil)
	attemptRepo := attemptrepo.NewAttemptUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	now := time.Date(2099, time.April, 3, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &now, testutil.Ptr(now.Add(45*time.Minute)), testutil.Ptr(now.Add(time.Hour)))
	fixtures.CreateGradingJob(attempt.ID, ujian.StatusQueued)

	statsSvc := gradingstatsvc.NewStatistikUjianService(gradingRepo)
	require.NoError(t, statsSvc.StatistikUjian(scope.Context(), int(attempt.ID)))

	workerSvc := gradingworker.NewGradingUjianWorkerService(gradingRepo, &fakeGradingExecutor{}, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(scope.Context())
	done := make(chan struct{})
	go func() {
		defer close(done)
		workerSvc.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	var status string
	err := scope.Pool().QueryRow(scope.Context(), `SELECT status FROM grading_jobs WHERE id_attempt = $1`, attempt.ID).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, string(ujian.StatusDone), status)

	gotAttempt, err := attemptRepo.GetAttemptUjianById(scope.Context(), ujian.ID(attempt.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ATTEMPT_SUBMITTED, gotAttempt.StatusAttempt)
}
