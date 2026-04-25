package ujian_penjadwalan_test

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	ujianrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian"
	ujianlistrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/ujian/list"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujiandelete "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/delete"
	ujianget "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
	ujiancreate "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	ujianupdate "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGetUpdateDeleteUjianService_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := ujianrepo.NewUjianRepo(scope.Pool(), nil, scope.Pool())

	kelas := fixtures.CreateKelas(91)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "X IPA 1")
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	pengawas := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(pengawas.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bankSoal := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	version := fixtures.CreateBankSoalVersion(bankSoal.ID, guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(bankSoal.ID, version.ID)
	sesi := fixtures.CreateSesi()
	ruang := fixtures.CreateRuangUjian()
	now := time.Date(2099, time.March, 1, 9, 0, 0, 0, time.UTC)

	createSvc := ujiancreate.NewCreateUjianService(repo)
	err := createSvc.CreateUjianService(scope.Context(), ujian.PenjadwalanUjian{
		Ujian: ujian.Ujian{
			IdBankSoal:     ujian.ID(bankSoal.ID),
			IdKelas:        ujian.ID(kelas.ID),
			IdNamaKelas:    testutil.Ptr(ujian.ID(namaKelas.ID)),
			IdGuru:         ujian.ID(guru.ID),
			NamaUjian:      "  Ujian Integrasi  ",
			DeskripsiUjian: testutil.Ptr("  Deskripsi Integrasi  "),
			AcakSoal:       true,
		},
		JadwalUjian: ujian.JadwalUjian{
			IdSesi:       ujian.ID(sesi.ID),
			IdRuangan:    ujian.ID(ruang.ID),
			IdPengawas:   ujian.ID(pengawas.ID),
			TanggalUjian: now,
			WaktuMulai:   now,
			WaktuSelesai: now.Add(time.Hour),
			Token:        "  uj-int-01  ",
			StatusUjian:  ujian.BELUM_MULAI,
		},
	})
	require.NoError(t, err)

	var idUjian int64
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT id_ujian
		FROM ujian
		WHERE nama_ujian = 'Ujian Integrasi'
	`).Scan(&idUjian)
	require.NoError(t, err)

	getSvc := ujianget.NewGetujianService(ujianlistrepo.NewListUjianRepo(scope.Pool(), nil))
	items, err := getSvc.GetAllUjianService(scope.Context(), query.ListUjianFilter{
		Limit:          10,
		TingkatKelasID: &[]int{int(kelas.ID)}[0],
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
	assert.Equal(t, ujian.ID(idUjian), items[0].IdUjian)
	assert.Equal(t, "Ujian Integrasi", items[0].NamaUjian)

	item, err := getSvc.GetUjianByIdService(scope.Context(), ujian.ID(idUjian))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(idUjian), item.IdUjian)
	assert.Equal(t, ujian.ID(sesi.ID), item.IdSesi)
	assert.Equal(t, ujian.ID(ruang.ID), item.IdRuangan)

	updatedName := "Ujian Integrasi Updated"
	updatedStatus := ujian.MULAI
	updatedToken := "token-updated"
	waktuMulai := now.Add(10 * time.Minute)
	waktuSelesai := now.Add(70 * time.Minute)
	updateSvc := ujianupdate.NewUpdateUjianService(repo)
	err = updateSvc.UpdateUjianService(scope.Context(), ujian.ID(idUjian), updatepatch.UpdatePenjadwalanUjian{
		Ujian: updatepatch.UpdateUjianPatch{
			NamaUjian: &updatedName,
		},
		JadwalUjian: updatepatch.UpdateJadwalUjianPatch{
			StatusUjian:  &updatedStatus,
			Token:        &updatedToken,
			WaktuMulai:   &waktuMulai,
			WaktuSelesai: &waktuSelesai,
		},
	})
	require.NoError(t, err)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT u.nama_ujian, ju.status_ujian, ju.token, ju.waktu_mulai, ju.waktu_selesai
		FROM ujian u
		JOIN jadwal_ujian ju ON ju.id_ujian = u.id_ujian
		WHERE u.id_ujian = $1
	`, idUjian).Scan(&item.NamaUjian, &item.StatusUjian, &item.Token, &item.WaktuMulai, &item.WaktuSelesai)
	require.NoError(t, err)
	assert.Equal(t, updatedName, item.NamaUjian)
	require.NotNil(t, item.StatusUjian)
	assert.Equal(t, updatedStatus, *item.StatusUjian)
	assert.Equal(t, "TOKEN-UPDATED", item.Token)
	assert.True(t, waktuMulai.Equal(item.WaktuMulai))
	assert.True(t, waktuSelesai.Equal(item.WaktuSelesai))

	deleteSvc := ujiandelete.NewDeleteUjianService(repo)
	err = deleteSvc.DeleteUjianService(scope.Context(), ujian.ID(idUjian))
	require.NoError(t, err)

	var count int
	err = scope.Pool().QueryRow(scope.Context(), `SELECT COUNT(*) FROM ujian WHERE id_ujian = $1`, idUjian).Scan(&count)
	require.NoError(t, err)
	assert.Zero(t, count)
}

func TestDeleteUjianService_InvalidID(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	repo := ujianrepo.NewUjianRepo(scope.Pool(), nil, scope.Pool())
	svc := ujiandelete.NewDeleteUjianService(repo)

	err := svc.DeleteUjianService(scope.Context(), 0)
	assert.ErrorIs(t, err, coreerror.ErrMissingId)
}
