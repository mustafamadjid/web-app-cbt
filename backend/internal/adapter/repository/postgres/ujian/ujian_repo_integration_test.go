package ujianrepo

import (
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUjianRepo_GetIdUjianByAttempt(t *testing.T) {
	ctx, tx := testutil.BeginRollbackTx(t)
	fixtures := testutil.NewFixtures(t, ctx, tx)
	repo := NewUjianRepo(tx, nil, nil)

	exam := fixtures.CreateExamFixture()
	start := time.Date(2099, time.January, 1, 9, 0, 0, 0, time.UTC)
	attempt := fixtures.CreateAttempt(exam.Peserta.ID, ujian.ATTEMPT_IN_PROGRESS, &start, nil, testutil.Ptr(start.Add(time.Hour)))

	idUjian, err := repo.GetIdUjianByAttempt(ctx, ujian.ID(attempt.ID))
	require.NoError(t, err)
	assert.Equal(t, ujian.ID(exam.Ujian.ID), idUjian)
}

func TestUjianRepo_CreateUpdateDeleteUjian(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewUjianRepo(scope.Pool(), nil, scope.Pool())

	kelas := fixtures.CreateKelas(81)
	namaKelas := fixtures.CreateNamaKelas(kelas.ID, "Kelas Ujian Repo")
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	pengawas := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(pengawas.ID)
	siswa := fixtures.CreateUser(user.SISWA)
	fixtures.CreateSiswaProfile(siswa.ID, namaKelas.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	version := fixtures.CreateBankSoalVersion(bank.ID, guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(bank.ID, version.ID)
	base := time.Date(2099, time.February, 10, 9, 0, 0, 0, time.UTC)

	err := repo.CreateUjian(scope.Context(), ujian.PenjadwalanUjian{
		Ujian: ujian.Ujian{
			IdBankSoal:     ujian.ID(bank.ID),
			IdKelas:        ujian.ID(kelas.ID),
			IdNamaKelas:    testutil.Ptr(ujian.ID(namaKelas.ID)),
			IdGuru:         ujian.ID(guru.ID),
			NamaUjian:      "Repo Create Ujian",
			DeskripsiUjian: testutil.Ptr("deskripsi repo"),
			AcakSoal:       false,
		},
		JadwalUjian: ujian.JadwalUjian{
			IdSesi:       ujian.ID(fixtures.CreateSesi().ID),
			IdRuangan:    ujian.ID(fixtures.CreateRuangUjian().ID),
			IdPengawas:   ujian.ID(pengawas.ID),
			TanggalUjian: base,
			WaktuMulai:   base,
			WaktuSelesai: base.Add(time.Hour),
			Token:        "UJIANREPO",
			StatusUjian:  ujian.BELUM_MULAI,
		},
	})
	require.NoError(t, err)

	var idUjian int64
	err = scope.Pool().QueryRow(scope.Context(), `SELECT id_ujian FROM ujian WHERE nama_ujian = 'Repo Create Ujian'`).Scan(&idUjian)
	require.NoError(t, err)

	scope.AddCleanupQuery(`DELETE FROM jadwal_ujian WHERE id_ujian = $1`, idUjian)
	scope.AddCleanupQuery(`DELETE FROM ujian WHERE id_ujian = $1`, idUjian)

	var pesertaCount int
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT COUNT(*)
		FROM peserta_ujian pu
		JOIN jadwal_ujian ju ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
		WHERE ju.id_ujian = $1
	`, idUjian).Scan(&pesertaCount)
	require.NoError(t, err)
	assert.Equal(t, 1, pesertaCount)

	updatedName := "Repo Updated Ujian"
	updatedStatus := ujian.SELESAI
	err = repo.UpdateUjian(scope.Context(), ujian.ID(idUjian), updatepatch.UpdatePenjadwalanUjian{
		Ujian: updatepatch.UpdateUjianPatch{
			NamaUjian: &updatedName,
		},
		JadwalUjian: updatepatch.UpdateJadwalUjianPatch{
			StatusUjian: &updatedStatus,
		},
	})
	require.NoError(t, err)

	var storedName string
	var storedStatus string
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT u.nama_ujian, ju.status_ujian
		FROM ujian u
		JOIN jadwal_ujian ju ON ju.id_ujian = u.id_ujian
		WHERE u.id_ujian = $1
	`, idUjian).Scan(&storedName, &storedStatus)
	require.NoError(t, err)
	assert.Equal(t, updatedName, storedName)
	assert.Equal(t, string(updatedStatus), storedStatus)

	err = repo.DeleteUjian(scope.Context(), ujian.ID(idUjian))
	require.NoError(t, err)

	var count int
	err = scope.Pool().QueryRow(scope.Context(), `SELECT COUNT(*) FROM ujian WHERE id_ujian = $1`, idUjian).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestUjianRepo_CreateUjian_ConflictSchedule(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewUjianRepo(scope.Pool(), nil, scope.Pool())

	exam := fixtures.CreateExamFixture()
	otherGuru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(otherGuru.ID)
	otherMapel := fixtures.CreateMapel(exam.Kelas.ID)
	otherBank := fixtures.CreateBankSoal(otherMapel.ID, exam.Kelas.ID, otherGuru.ID)
	otherVersion := fixtures.CreateBankSoalVersion(otherBank.ID, otherGuru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(otherBank.ID, otherVersion.ID)

	err := repo.CreateUjian(scope.Context(), ujian.PenjadwalanUjian{
		Ujian: ujian.Ujian{
			IdBankSoal:     ujian.ID(otherBank.ID),
			IdKelas:        ujian.ID(exam.Kelas.ID),
			IdNamaKelas:    testutil.Ptr(ujian.ID(exam.NamaKelas.ID)),
			IdGuru:         ujian.ID(otherGuru.ID),
			NamaUjian:      "Conflict Ujian",
			DeskripsiUjian: testutil.Ptr("conflict"),
		},
		JadwalUjian: ujian.JadwalUjian{
			IdSesi:       ujian.ID(exam.Sesi.ID),
			IdRuangan:    ujian.ID(exam.Ruang.ID),
			IdPengawas:   ujian.ID(exam.Pengawas.ID),
			TanggalUjian: exam.Jadwal.WaktuMulai,
			WaktuMulai:   exam.Jadwal.WaktuMulai.Add(10 * time.Minute),
			WaktuSelesai: exam.Jadwal.WaktuSelesai.Add(10 * time.Minute),
			Token:        "CONFLICT",
			StatusUjian:  ujian.BELUM_MULAI,
		},
	})
	assert.ErrorIs(t, err, coreerror.ErrConflict)
}
