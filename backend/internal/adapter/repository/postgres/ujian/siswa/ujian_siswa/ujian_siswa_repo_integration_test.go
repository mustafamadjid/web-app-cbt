package ujiansiswarepo

import (
	"database/sql"
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUjianSiswaRepo_ListAndGetAttempt(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewUjianSiswaRepo(tx, nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_IN_PROGRESS, &start, nil, testutil.Ptr(start.Add(time.Hour)))

	tingkatKelasID := int(exam.Kelas.ID)
	items, err := repo.ListUjianSiswa(ctx, int(exam.Siswa.ID), query.ListUjianFilter{TingkatKelasID: &tingkatKelasID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ujian.ID(exam.Ujian.ID), items[0].IdUjian)

	waktuSelesai, err := repo.GetWaktuSelesaiUjian(ctx, int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.True(t, waktuSelesai.Equal(exam.Jadwal.WaktuSelesai))

	active, err := repo.GetActiveUjianAttemptBySiswa(ctx, int(exam.Siswa.ID), int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(attempt.ID), active.IdAttempt)
}

func TestUjianSiswaRepo_GetActiveUjianAttemptBySiswa_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewUjianSiswaRepo(tx, nil)

	exam := fixtures.CreateExamFixture()

	_, err := repo.GetActiveUjianAttemptBySiswa(ctx, int(exam.Siswa.ID), int(exam.Jadwal.ID))
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
