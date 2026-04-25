package statistik_ujian_test

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	statistikujianrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/statistik_ujian"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	statistikservice "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/statistik_ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatistikUjianService_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := statistikujianrepo.NewStatistikUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	fixtures.CreateStatistikUjian(exam.Jadwal.ID, 95, 60, 80, 3)

	item, err := statistikservice.NewStatistikUjianService(repo).GetStatistikUjian(scope.Context(), int(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(exam.Jadwal.ID), item.IDJadwalUjian)
	assert.Equal(t, 95.0, item.NilaiTertinggi)
	assert.Equal(t, 60.0, item.NilaiTerendah)
	assert.Equal(t, 80.0, item.NilaiRataRata)
	assert.Equal(t, 3, item.TotalPesertaUjian)
}

func TestStatistikUjianService_InvalidID(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	repo := statistikujianrepo.NewStatistikUjianRepo(scope.Pool(), nil)

	_, err := statistikservice.NewStatistikUjianService(repo).GetStatistikUjian(scope.Context(), 0)
	assert.ErrorIs(t, err, coreerror.ErrMissingId)
}
