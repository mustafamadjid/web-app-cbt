package dashboard_repo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardStatistikRepo_GetDashboardStatistik(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewDashboardStatistikRepo(tx, nil)

	before, err := repo.GetDashboardStatistik(ctx)
	require.NoError(t, err)

	kelas := fixtures.CreateKelas(77)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "Dashboard Integration")
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	siswa := fixtures.CreateUser(user.SISWA)
	fixtures.CreateSiswaProfile(siswa.ID, namaKelas.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	version := fixtures.CreateBankSoalVersion(bank.ID, guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(bank.ID, version.ID)
	fixtures.CreateSoal(version.ID, "essay", "Dashboard soal", 10, 1)
	sesi := fixtures.CreateSesi()
	ruang := fixtures.CreateRuangUjian()
	pengawas := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(pengawas.ID)
	uj := fixtures.CreateUjian(bank.ID, kelas.ID, testutil.Ptr(namaKelas.ID), guru.ID, "Dashboard Ujian")
	base := time.Date(2000, time.January, 1, 9, 0, 0, 0, time.UTC)
	fixtures.CreateJadwalUjian(uj.ID, sesi.ID, ruang.ID, pengawas.ID, base, base, base.Add(time.Hour), ujian.SELESAI)

	after, err := repo.GetDashboardStatistik(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, after.TotalGuru, before.TotalGuru+2)
	assert.GreaterOrEqual(t, after.TotalSiswa, before.TotalSiswa+1)
	assert.GreaterOrEqual(t, after.TotalMapelAktif, before.TotalMapelAktif+1)
	assert.GreaterOrEqual(t, after.TotalBankSoal, before.TotalBankSoal+1)
	assert.GreaterOrEqual(t, after.TotalUjianTerlaksana, before.TotalUjianTerlaksana+1)
}
