package jawabanujian_repo

import (
	"context"
	"testing"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func toIDPointer(v ujian.ID) *ujian.ID {
	return &v
}

func toStringPointer(v string) *string {
	return &v
}

func TestSplitSaveJawabanItems(t *testing.T) {
	t.Parallel()

	items := []ujian.JawabanUjian{
		{
			IdSoal:       10,
			JawabanEssay: toStringPointer("  jawaban essay  "),
		},
		{
			IdSoal:       11,
			JawabanEssay: toStringPointer("   "),
		},
		{
			IdSoal:    12,
			IdPilihan: toIDPointer(33),
		},
		{
			IdSoal: 13,
		},
	}

	upsertItems, clearSoalIDs := splitSaveJawabanItems(items)

	require.Len(t, upsertItems, 2)
	assert.Equal(t, []int64{11, 13}, clearSoalIDs)

	require.NotNil(t, upsertItems[0].JawabanEssay)
	assert.Equal(t, "jawaban essay", *upsertItems[0].JawabanEssay)
	assert.Equal(t, ujian.ID(10), upsertItems[0].IdSoal)

	assert.Equal(t, ujian.ID(12), upsertItems[1].IdSoal)
	require.NotNil(t, upsertItems[1].IdPilihan)
	assert.Equal(t, ujian.ID(33), *upsertItems[1].IdPilihan)
	assert.Nil(t, upsertItems[1].JawabanEssay)
}

func TestSaveJawabanUjianRequiresPoolForTransaction(t *testing.T) {
	t.Parallel()

	repo := NewJawabanUjianRepo(nil, nil)

	err := repo.SaveJawabanUjian(context.Background(), 1, []ujian.JawabanUjian{
		{
			IdSoal:    10,
			IdPilihan: toIDPointer(20),
		},
	})

	require.Error(t, err)
	assert.Equal(t, "jawaban ujian repo requires pgx pool for save transaction", err.Error())
}
