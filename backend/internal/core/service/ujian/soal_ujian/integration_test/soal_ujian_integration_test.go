package soal_ujian_test

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujianlistrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/list"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	soalservice "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSoalUjianService_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := ujianlistrepo.NewListSoalUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	items, err := soalservice.NewListSoalUjianService(repo).ListSoalUjian(scope.Context(), ujian.ID(exam.BankSoal.ID), false)
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, ujian.ID(exam.SoalPilgan.ID), items[0].IdSoal)
}

func TestListSoalUjianService_InvalidID(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	repo := ujianlistrepo.NewListSoalUjianRepo(scope.Pool(), nil)

	_, err := soalservice.NewListSoalUjianService(repo).ListSoalUjian(scope.Context(), 0, false)
	assert.ErrorIs(t, err, coreerror.ErrMissingId)
}
