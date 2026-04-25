package worker_test

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	importsoaljobrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/import_soal_job"
	isisoalbatchrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/isi_soal_batch"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/testutil"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	createjob "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/create_job"
	importversion "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/import_version"
	worker "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopLogger struct{}

func (noopLogger) With(...any) corelog.Logger { return noopLogger{} }
func (noopLogger) Info(context.Context, string, ...any)  {}
func (noopLogger) Error(context.Context, string, ...any) {}

func TestImportSoalWorker_ProcessesQueuedJob(t *testing.T) {
	scope := testutil.NewCommittedFixtureScope(t)
	fixtures := scope.Fixtures()
	jobRepo := importsoaljobrepo.NewImportSoalJobRepo(scope.Pool(), nil)
	batchRepo := isisoalbatchrepo.NewIsiSoalBatchRepo(scope.Pool(), nil)

	kelas := fixtures.CreateKelas(94)
	guru := fixtures.CreateUser(user.GURU)
	fixtures.CreateGuruProfile(guru.ID)
	mapel := fixtures.CreateMapel(kelas.ID)
	bank := fixtures.CreateBankSoal(mapel.ID, kelas.ID, guru.ID)
	oldVersion := fixtures.CreateBankSoalVersion(bank.ID, guru.ID, "published", 1)
	fixtures.SetActiveBankSoalVersion(bank.ID, oldVersion.ID)

	docxPath := filepath.Join(t.TempDir(), "import-worker.docx")
	require.NoError(t, os.WriteFile(docxPath, buildWorkerDocx(), 0o644))

	createSvc := createjob.NewCreateJobService(jobRepo)
	res, err := createSvc.Execute(scope.Context(), createjob.CreateJobCmd{
		IDBankSoal: bank.ID,
		IDPengguna: int64(guru.ID),
		FilePath:   docxPath,
	})
	require.NoError(t, err)

	workerSvc := worker.NewWorker(
		jobRepo,
		importversion.NewService(batchRepo),
		t.TempDir(),
		"/uploads/image/bank_soal",
		10*time.Millisecond,
		noopLogger{},
	)

	ctx, cancel := context.WithCancel(scope.Context())
	done := make(chan struct{})
	go func() {
		defer close(done)
		workerSvc.Start(ctx)
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done

	job, err := jobRepo.GetJobByID(scope.Context(), res.IDJob)
	require.NoError(t, err)
	assert.Equal(t, importsoal.StatusCompleted, job.Status)
	assert.Equal(t, 1, job.TotalSoal)

	var activeVersionID int64
	err = scope.Pool().QueryRow(scope.Context(), `SELECT id_bank_soal_version_aktif FROM bank_soal WHERE id_bank_soal = $1`, bank.ID).Scan(&activeVersionID)
	require.NoError(t, err)
	assert.NotEqual(t, oldVersion.ID, activeVersionID)
}

func buildWorkerDocx() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, _ := zw.Create("word/document.xml")
	_, _ = fw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>
<w:p><w:r><w:t>[Q:PG] Apa ibukota Indonesia?</w:t></w:r></w:p>
<w:p><w:r><w:t>[A] Jakarta</w:t></w:r></w:p>
<w:p><w:r><w:t>[B] Bandung</w:t></w:r></w:p>
<w:p><w:r><w:t>[ANS] A</w:t></w:r></w:p>
<w:p><w:r><w:t>[W] 10</w:t></w:r></w:p>
</w:body>
</w:document>`))
	_, _ = zw.Create("word/_rels/document.xml.rels")
	_ = zw.Close()
	return buf.Bytes()
}
