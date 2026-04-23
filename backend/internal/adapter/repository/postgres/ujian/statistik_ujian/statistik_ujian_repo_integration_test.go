package statistikujian_repo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatistikUjianRepo_GetStatistikUjianByIdJadwal(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewStatistikUjianRepo(tx, nil)

	exam := fixtures.CreateExamFixture()
	fixtures.CreateStatistikUjian(exam.Jadwal.ID, 90, 70, 80, 10)

	item, err := repo.GetStatistikUjianByIdJadwal(ctx, ujian.ID(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(exam.Jadwal.ID), item.IDJadwalUjian)
	assert.Equal(t, 90.0, item.NilaiTertinggi)
	assert.Equal(t, 70.0, item.NilaiTerendah)
	assert.Equal(t, 80.0, item.NilaiRataRata)
	assert.Equal(t, 10, item.TotalPesertaUjian)
}

func TestStatistikUjianRepo_GetStatistikUjianByIdJadwal_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewStatistikUjianRepo(tx, nil)

	_, err := repo.GetStatistikUjianByIdJadwal(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
