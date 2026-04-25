package import_soal_test

import (
	"archive/zip"
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	importsoaljobrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/import_soal_job"
	isisoalbatchrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/isi_soal_batch"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	importversion "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/import_version"
	createjob "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	getjob "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/get_job"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportSoalJobAndVersionServices_Integration(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	jobRepo := importsoaljobrepo.NewImportSoalJobRepo(scope.Pool(), nil)
	batchRepo := isisoalbatchrepo.NewIsiSoalBatchRepo(scope.Pool(), nil)

	kelas := fixtures.CreateKelas(92)
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	oldVersion := fixtures.CreateBankSoalVersion(bank.ID, guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(bank.ID, oldVersion.ID)

	createSvc := createjob.NewCreateJobService(jobRepo)
	res, err := createSvc.Execute(scope.Context(), createjob.CreateJobCmd{
		IDBankSoal: bank.ID,
		IDPengguna: int64(guru.ID),
		FilePath:   filepath.Join(t.TempDir(), "bank-soal.docx"),
	})
	require.NoError(t, err)
	assert.NotZero(t, res.IDJob)

	gotJob, err := getjob.NewGetJobService(jobRepo).GetByID(scope.Context(), res.IDJob)
	require.NoError(t, err)
	assert.Equal(t, bank.ID, gotJob.IDBankSoal)
	assert.Equal(t, importsoal.StatusPending, gotJob.Status)

	listedByBank, err := getjob.NewGetJobService(jobRepo).GetByBankSoal(scope.Context(), bank.ID)
	require.NoError(t, err)
	require.Len(t, listedByBank, 1)

	payload := []importsoal.ParsedSoal{
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
	}
	versionRes, err := importversion.NewService(batchRepo).Execute(scope.Context(), importversion.Cmd{
		BankID:  bank.ID,
		UserID:  int64(guru.ID),
		Payload: payload,
	})
	require.NoError(t, err)
	assert.NotZero(t, versionRes.VersionID)
	assert.NotEqual(t, oldVersion.ID, versionRes.VersionID)

	var activeVersionID int64
	err = scope.Pool().QueryRow(scope.Context(), `SELECT id_bank_soal_version_aktif FROM bank_soal WHERE id_bank_soal = $1`, bank.ID).Scan(&activeVersionID)
	require.NoError(t, err)
	assert.Equal(t, versionRes.VersionID, activeVersionID)
}

func TestImportSoalVersionService_RejectsInvalidPayload(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	kelas := scope.Fixtures().CreateKelas(93)
	guru := scope.Fixtures().CreateUser(user.GURU)
	scope.Fixtures().CreateGuruProfile(guru.ID)
	mapel := scope.Fixtures().CreateMapel(kelas.ID)
	bank := scope.Fixtures().CreateBankSoal(mapel.ID, kelas.ID, guru.ID)

	_, err := importversion.NewService(isisoalbatchrepo.NewIsiSoalBatchRepo(scope.Pool(), nil)).Execute(scope.Context(), importversion.Cmd{
		BankID:  bank.ID,
		UserID:  int64(guru.ID),
		Payload: nil,
	})
	assert.Error(t, err)
}

func TestDocxParserIntegrationHelpers(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, _ := zw.Create("word/document.xml")
	_, _ = fw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>[Q:PG] Halo?</w:t></w:r></w:p><w:p><w:r><w:t>[A] A</w:t></w:r></w:p><w:p><w:r><w:t>[B] B</w:t></w:r></w:p><w:p><w:r><w:t>[ANS] A</w:t></w:r></w:p><w:p><w:r><w:t>[W] 10</w:t></w:r></w:p></w:body></w:document>`))
	require.NoError(t, zw.Close())

	paragraphs, err := parser.ExtractParagraphs(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, paragraphs, 5)
}

func TestImportSoalJobRepoIntegrationUsesContext(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	repo := importsoaljobrepo.NewImportSoalJobRepo(scope.Pool(), nil)

	_, err := repo.GetJobByID(context.Background(), 0)
	assert.Error(t, err)
}
