package ujianlistrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUjianRepo_ListAndDetailQueries(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewListUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_SUBMITTED, &start, testutil.Ptr(start.Add(time.Hour)), testutil.Ptr(start.Add(time.Hour)))
	assert.NotZero(t, attempt.ID)

	tingkatKelasID := int(exam.Kelas.ID)
	items, err := repo.GetAllUjian(scope.Context(), query.ListUjianFilter{TingkatKelasID: &tingkatKelasID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, ujian.ID(exam.Ujian.ID), items[0].IdUjian)

	item, err := repo.GetUjianById(scope.Context(), ujian.ID(exam.Ujian.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(exam.Jadwal.ID), item.IdJadwalUjian)
	assert.Equal(t, ujian.ID(exam.Ruang.ID), item.IdRuangan)

	submitted, err := repo.GetAllUjianSubmittedByIdSiswa(scope.Context(), int(exam.Siswa.ID))
	require.NoError(t, err)
	require.Len(t, submitted, 1)
	assert.Equal(t, ujian.ID(attempt.ID), submitted[0].IdAttempt)
}

func TestListSoalUjianRepo_ListSoalAndOpsi(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewListSoalUjianRepo(scope.Pool(), nil)

	exam := fixtures.CreateExamFixture()

	soal, err := repo.GetSoalUjianByBankSoal(scope.Context(), ujian.ID(exam.BankSoal.ID))
	require.NoError(t, err)
	require.Len(t, soal, 2)
	assert.Equal(t, 2, len(soal[0].OpsiJawaban))

	soalSiswa, acakSoal, err := repo.GetSoalUjianByBankSoalForSiswa(scope.Context(), ujian.ID(exam.Jadwal.ID))
	require.NoError(t, err)
	assert.False(t, acakSoal)
	require.Len(t, soalSiswa, 2)

	opsi, err := repo.GetOpsiPilihanGandaByBankSoal(scope.Context(), int(exam.Version.ID))
	require.NoError(t, err)
	require.Len(t, opsi, 2)
}
