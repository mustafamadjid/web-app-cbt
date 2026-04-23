package profilsiswarepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfilSiswaRepo_CRUDAndQueries(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewProfilSiswaRepo(tx, nil)

	kelas := fixtures.CreateKelas(55)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "Kelas Siswa IT")
	createdUser := fixtures.CreateUser(user.SISWA)

	profilID, err := repo.CreateProfilSiswa(ctx, user.ProfilSiswa{
		IdPengguna:   createdUser.ID,
		IdNamaKelas:  user.ID(namaKelas.ID),
		Nisn:         user.NISN("1234567890"),
		NoAbsen:      9,
		Angkatan:     2024,
		TempatLahir:  "Bandung",
		TanggalLahir: time.Date(2007, time.March, 2, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.NotZero(t, profilID)

	exists, err := repo.ExistByNISN(ctx, "1234567890")
	require.NoError(t, err)
	assert.True(t, exists)

	item, err := repo.FindProfilSiswaByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.Username, item.Username)
	assert.Equal(t, "1234567890", string(item.Nisn))
	assert.Equal(t, "Kelas Siswa IT", item.NamaKelas)

	updatedNoAbsen := 11
	err = repo.UpdateProfilSiswa(ctx, createdUser.ID, updatepatch.ProfilSiswa{
		NoAbsen: &updatedNoAbsen,
	})
	require.NoError(t, err)

	updated, err := repo.FindProfilSiswaByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, updatedNoAbsen, updated.NoAbsen)

	idKelasFilter := int(kelas.ID)
	idNamaKelasFilter := int(namaKelas.ID)
	items, err := repo.GetListSiswa(ctx, query.ListSiswaFilter{
		Search:         createdUser.Username,
		Limit:          10,
		IdTingkatKelas: &idKelasFilter,
		IdNamaKelas:    &idNamaKelasFilter,
	})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, createdUser.ID, items[0].IdPengguna)
}

func TestProfilSiswaRepo_FindProfilSiswaByID_NotFound(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := NewProfilSiswaRepo(tx, nil)

	_, err := repo.FindProfilSiswaByID(ctx, 999999)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
