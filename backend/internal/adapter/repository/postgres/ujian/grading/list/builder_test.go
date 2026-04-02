package gradingrepo

import (
	"testing"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func TestBuildListUjianEssayUngradedQuery_WithFiltersAndPagination(t *testing.T) {
	t.Parallel()

	repo := &ListGradingRepo{}
	queryText, args := repo.buildListUjianEssayUngradedQuery(query.ListUjianEssayUngradedFilter{
		Search:         "esai",
		TanggalUjian:   stringPtr("2026-04-02"),
		Tahun:          stringPtr("2026"),
		Bulan:          stringPtr("04"),
		TingkatKelasID: intPtr(3),
		NamaKelasID:    intPtr(5),
		MapelID:        intPtr(7),
		SesiID:         intPtr(9),
		Limit:          10,
		Offset:         20,
	})

	assert.Contains(t, queryText, "hu.essay_graded = FALSE")
	assert.Contains(t, queryText, "is2.tipe_soal = 'essay'")
	assert.Contains(t, queryText, "u.nama_ujian ILIKE $1")
	assert.Contains(t, queryText, "ju.tanggal_ujian = $2::date")
	assert.Contains(t, queryText, "EXTRACT(YEAR FROM ju.tanggal_ujian)::int = $3::int")
	assert.Contains(t, queryText, "EXTRACT(MONTH FROM ju.tanggal_ujian)::int = $4::int")
	assert.Contains(t, queryText, "u.id_kelas = $5")
	assert.Contains(t, queryText, "u.id_nama_kelas = $6")
	assert.Contains(t, queryText, "bs.id_mapel = $7")
	assert.Contains(t, queryText, "ju.id_sesi = $8")
	assert.Contains(t, queryText, "LIMIT $9 OFFSET $10")

	require.Equal(t, []any{"%esai%", "2026-04-02", "2026", "04", 3, 5, 7, 9, 10, 20}, args)
}

func TestBuildListUjianEssayUngradedQuery_DefaultConditionsAlwaysPresent(t *testing.T) {
	t.Parallel()

	repo := &ListGradingRepo{}
	queryText, args := repo.buildListUjianEssayUngradedQuery(query.ListUjianEssayUngradedFilter{})

	assert.Contains(t, queryText, "JOIN hasil_ujian hu")
	assert.Contains(t, queryText, "hu.essay_graded = FALSE")
	assert.Contains(t, queryText, "FROM isi_soal is2")
	assert.Contains(t, queryText, "is2.tipe_soal = 'essay'")
	assert.NotContains(t, queryText, "LIMIT $")
	require.Empty(t, args)
}
