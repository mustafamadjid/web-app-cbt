package jawabanujian_repo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJawabanUjianRepo_SaveGetAndListHasil(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewJawabanUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_IN_PROGRESS, &start, nil, testutil.Ptr(start.Add(time.Hour)))

	essayText := "jawaban essay integration"
	err := repo.SaveJawabanUjian(scope.Context(), ujian.ID(attempt.ID), []ujian.JawabanUjian{
		{
			IdSoal:     ujian.ID(exam.SoalPilgan.ID),
			IdPilihan:  testutil.Ptr(ujian.ID(exam.OpsiBenar.ID)),
			WaktuJawab: testutil.Ptr(start.Add(10 * time.Minute)),
		},
		{
			IdSoal:       ujian.ID(exam.SoalEssay.ID),
			JawabanEssay: &essayText,
			WaktuJawab:   testutil.Ptr(start.Add(20 * time.Minute)),
		},
	})
	require.NoError(t, err)

	got, err := repo.GetJawabanUjianByAttemptId(scope.Context(), ujian.ID(attempt.ID))
	require.NoError(t, err)
	require.Len(t, got, 2)

	gradedBy := user.ID(exam.Guru.ID)
	nilai := 95.0
	passed := true
	essayGraded := false
	fixtures.CreateHasilUjian(attempt.ID, &gradedBy, &nilai, &passed, &essayGraded, testutil.Ptr(start.Add(time.Hour)), exam.Jadwal.ID)

	items, err := repo.ListHasilJawabanUjianByIdAttempt(scope.Context(), ujian.ID(attempt.ID))
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, essayText, *items[1].JawabanSiswa.JawabanEssay)
	require.NotNil(t, items[0].NilaiAkhir)
	assert.Equal(t, nilai, *items[0].NilaiAkhir)
}
