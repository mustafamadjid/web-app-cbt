package isisoalbatchrepo

import (
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsiSoalBatchRepo_ImportBankSoalVersion(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	repo := NewIsiSoalBatchRepo(scope.Pool(), nil)

	kelas := fixtures.CreateKelas(75)
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	oldVersion := fixtures.CreateBankSoalVersion(bank.ID, guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(bank.ID, oldVersion.ID)

	payload := importsoalrepo.ImportBankSoalVersionPayload{
		SoalList: []importsoal.ParsedSoal{
			{
				Pertanyaan: "Soal PG",
				TipeSoal:   "pilihan_ganda",
				BobotSoal:  10,
				NoUrutSoal: 1,
				Opsi: []importsoal.ParsedOpsi{
					{Label: "A", Isi: "A", IsBenar: true},
					{Label: "B", Isi: "B", IsBenar: false},
				},
			},
			{
				Pertanyaan: "Soal Essay",
				TipeSoal:   "essay",
				BobotSoal:  15,
				NoUrutSoal: 2,
			},
		},
	}

	newVersionID, err := repo.ImportBankSoalVersion(scope.Context(), bank.ID, int64(guru.ID), payload)
	require.NoError(t, err)
	assert.NotZero(t, newVersionID)
	assert.NotEqual(t, oldVersion.ID, newVersionID)

	var activeVersionID int64
	var oldStatus string
	var newStatus string
	var soalCount int
	var opsiCount int
	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT id_bank_soal_version_aktif
		FROM bank_soal
		WHERE id_bank_soal = $1
	`, bank.ID).Scan(&activeVersionID)
	require.NoError(t, err)
	assert.Equal(t, newVersionID, activeVersionID)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT status
		FROM bank_soal_version
		WHERE id_bank_soal_version = $1
	`, oldVersion.ID).Scan(&oldStatus)
	require.NoError(t, err)
	assert.Equal(t, "archived", oldStatus)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT status
		FROM bank_soal_version
		WHERE id_bank_soal_version = $1
	`, newVersionID).Scan(&newStatus)
	require.NoError(t, err)
	assert.Equal(t, "published", newStatus)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT COUNT(*)
		FROM isi_soal
		WHERE id_bank_soal_version = $1
	`, newVersionID).Scan(&soalCount)
	require.NoError(t, err)
	assert.Equal(t, 2, soalCount)

	err = scope.Pool().QueryRow(scope.Context(), `
		SELECT COUNT(*)
		FROM opsi_pilihan_ganda op
		JOIN isi_soal s ON s.id_soal = op.id_soal
		WHERE s.id_bank_soal_version = $1
	`, newVersionID).Scan(&opsiCount)
	require.NoError(t, err)
	assert.Equal(t, 2, opsiCount)
}

func TestIsiSoalBatchRepo_ImportBankSoalVersion_InvalidInput(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	repo := NewIsiSoalBatchRepo(scope.Pool(), nil)

	_, err := repo.ImportBankSoalVersion(scope.Context(), 0, 0, importsoalrepo.ImportBankSoalVersionPayload{})
	assert.ErrorIs(t, err, coreerror.ErrInvalidInput)
}
