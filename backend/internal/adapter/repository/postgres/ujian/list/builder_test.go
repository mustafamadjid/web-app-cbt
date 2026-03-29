package ujianlistrepo

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

func TestBuildListUjianQuery_WithFiltersAndPagination(t *testing.T) {
	t.Parallel()

	repo := &ListUjianRepo{}
	queryText, args := repo.buildListUjianQuery(query.ListUjianFilter{
		Search:         "matematika",
		Tahun:          stringPtr("2026"),
		TingkatKelasID: intPtr(3),
		RuangUjian:     intPtr(7),
		KategoriUjian:  query.BERLANGSUNG,
		Limit:          15,
		Offset:         30,
	})

	assert.Contains(t, queryText, "u.nama_ujian ILIKE $1")
	assert.Contains(t, queryText, "EXTRACT(YEAR FROM ju.tanggal_ujian)::text = $2")
	assert.Contains(t, queryText, "u.id_kelas = $3")
	assert.Contains(t, queryText, "ju.id_ruangan = $4")
	assert.Contains(t, queryText, "ju.waktu_mulai <= NOW() AND ju.waktu_selesai >= NOW()")
	assert.Contains(t, queryText, "LIMIT $5 OFFSET $6")

	require.Equal(t, []any{"%matematika%", "2026", 3, 7, 15, 30}, args)
}
