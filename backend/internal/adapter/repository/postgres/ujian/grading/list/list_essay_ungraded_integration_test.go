package gradingrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListGradingRepo_ListUjianEssayUngraded(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewListGradingRepo(tx, nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &start, testutil.Ptr(start.Add(time.Hour)), testutil.Ptr(start.Add(time.Hour)))
	fixtures.CreateHasilUjian(attempt.ID, nil, testutil.Ptr(70.0), nil, testutil.Ptr(false), nil, exam.Jadwal.ID)

	tingkatKelasID := int(exam.Kelas.ID)
	namaKelasID := int(exam.NamaKelas.ID)
	sesiID := int(exam.Sesi.ID)
	items, err := repo.ListUjianEssayUngraded(ctx, query.ListUjianEssayUngradedFilter{
		Limit:          10,
		TingkatKelasID: &tingkatKelasID,
		NamaKelasID:    &namaKelasID,
		SesiID:         &sesiID,
	})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ujian.ID(exam.Ujian.ID), items[0].IdUjian)
}
