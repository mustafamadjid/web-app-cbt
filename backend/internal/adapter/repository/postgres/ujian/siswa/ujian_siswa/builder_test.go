package ujiansiswarepo

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

func TestBuildListUjianSiswaQuery_WithFiltersAndPagination(t *testing.T) {
	t.Parallel()

	repo := &UjianSiswaRepo{}
	queryText, args := repo.buildListUjianSiswaQuery(99, query.ListUjianFilter{
		Search:         "fisika",
		TanggalUjian:   stringPtr("2026-03-29"),
		Bulan:          stringPtr("03"),
		TingkatKelasID: intPtr(2),
		RuangUjian:     intPtr(5),
		IDMapel:        intPtr(8),
		KategoriUjian:  query.SELESAI,
		Limit:          10,
		Offset:         20,
	})

	assert.Contains(t, queryText, "pu.id_siswa = $1")
	assert.Contains(t, queryText, "NOT EXISTS")
	assert.Contains(t, queryText, "u.nama_ujian ILIKE $2")
	assert.Contains(t, queryText, "ju.tanggal_ujian = $3::date")
	assert.Contains(t, queryText, "EXTRACT(MONTH FROM ju.tanggal_ujian)::int = $4::int")
	assert.Contains(t, queryText, "u.id_kelas = $5")
	assert.Contains(t, queryText, "ju.id_ruangan = $6")
	assert.Contains(t, queryText, "bs.id_mapel = $7")
	assert.Contains(t, queryText, "ju.waktu_selesai < NOW()")
	assert.Contains(t, queryText, "LIMIT $8 OFFSET $9")

	require.Equal(t, []any{99, "%fisika%", "2026-03-29", "03", 2, 5, 8, 10, 20}, args)
}
