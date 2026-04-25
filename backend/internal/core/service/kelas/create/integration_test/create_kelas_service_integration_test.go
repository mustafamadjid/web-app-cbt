package integration_test

import (
	"fmt"
	"testing"
	"time"

	kelasrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/kelas"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateKelasService_CreateTingkatAndNamaKelas(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	svc := kelas_service.NewCreateKelasService(repo)

	tingkat := 9000 + int(time.Now().UnixNano()%1000)
	err := svc.CreateTingkatKelas(ctx, kelas_service.CreateTingkatKelasCmd{TingkatKelas: tingkat})
	require.NoError(t, err)

	exists, err := repo.ExistTingkatKelas(ctx, tingkat)
	require.NoError(t, err)
	assert.True(t, exists)

	var idKelas int
	err = tx.QueryRow(ctx, `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idKelas)
	require.NoError(t, err)

	namaKelas := fmt.Sprintf("  Integration Kelas %d  ", tingkat)
	err = svc.CreateNamaKelas(ctx, kelas_service.CreateNamaKelasCmd{
		IdTingkatKelas: kelas.ID(idKelas),
		NamaKelas:      namaKelas,
	})
	require.NoError(t, err)

	var storedNama string
	err = tx.QueryRow(ctx, `
		SELECT nama_kelas
		FROM nama_kelas
		WHERE id_kelas = $1
	`, idKelas).Scan(&storedNama)
	require.NoError(t, err)
	assert.Equal(t, "Integration Kelas "+fmt.Sprint(tingkat), storedNama)
}

func TestCreateKelasService_DuplicateValuesReturnDomainError(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	repo := kelasrepo.NewKelasRepo(tx, nil)
	svc := kelas_service.NewCreateKelasService(repo)

	tingkat := 9500 + int(time.Now().UnixNano()%1000)
	require.NoError(t, svc.CreateTingkatKelas(ctx, kelas_service.CreateTingkatKelasCmd{TingkatKelas: tingkat}))

	err := svc.CreateTingkatKelas(ctx, kelas_service.CreateTingkatKelasCmd{TingkatKelas: tingkat})
	assert.ErrorIs(t, err, coreerror.ErrTingkatKelasExist)

	var idKelas int
	err = tx.QueryRow(ctx, `SELECT id_kelas FROM kelas WHERE tingkat_kelas = $1`, tingkat).Scan(&idKelas)
	require.NoError(t, err)

	nama := fmt.Sprintf("Duplicate Kelas %d", tingkat)
	require.NoError(t, svc.CreateNamaKelas(ctx, kelas_service.CreateNamaKelasCmd{
		IdTingkatKelas: kelas.ID(idKelas),
		NamaKelas:      nama,
	}))

	err = svc.CreateNamaKelas(ctx, kelas_service.CreateNamaKelasCmd{
		IdTingkatKelas: kelas.ID(idKelas),
		NamaKelas:      nama,
	})
	assert.ErrorIs(t, err, coreerror.ErrNamaKelasExist)
}
