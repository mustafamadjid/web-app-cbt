package banksoalrepo

import (
	"testing"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int {
	return &v
}

func TestBuildListBankSoalQuery_WithFiltersAndPagination(t *testing.T) {
	t.Parallel()

	repo := &BankSoalRepo{}
	queryText, args := repo.buildListBankSoalQuery(query.BankSoalFilter{
		Search:       "trigonometri",
		TingkatKelas: intPtr(11),
		Mapel:        intPtr(4),
		Limit:        20,
		Offset:       40,
	}, true)

	assert.Contains(t, queryText, "b.id_bank_soal_version_aktif IS NOT NULL")
	assert.Contains(t, queryText, "b.nama_bank_soal ILIKE $1")
	assert.Contains(t, queryText, "b.id_kelas = $2")
	assert.Contains(t, queryText, "b.id_mapel = $3")
	assert.Contains(t, queryText, "LIMIT $4 OFFSET $5")

	require.Equal(t, []any{"%trigonometri%", 11, 4, 20, 40}, args)
}
