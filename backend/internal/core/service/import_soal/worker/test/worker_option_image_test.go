package worker_test

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkerProcessJobs_SavesOptionImages(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	docxPath := filepath.Join(tmpDir, "option-image.docx")
	require.NoError(t, os.WriteFile(docxPath, buildWorkerDocxWithOptionImage(), 0o644))

	var (
		mu       sync.Mutex
		done     = make(chan struct{})
		returned bool
		payload  importsoalrepo.ImportBankSoalVersionPayload
	)

	jobRepo := &FakeImportSoalJobRepo{
		GetPendingJobsFn: func(_ context.Context, _ int) ([]importsoal.ImportSoalJob, error) {
			mu.Lock()
			defer mu.Unlock()
			if returned {
				return nil, nil
			}
			returned = true
			return []importsoal.ImportSoalJob{
				{IDJob: 1, IDBankSoal: 10, IDPengguna: 20, FilePath: docxPath},
			}, nil
		},
		UpdateJobStatusFn: func(_ context.Context, _ int64, status importsoal.JobStatus, _, _ string, _ int) error {
			if status == importsoal.StatusCompleted {
				close(done)
			}
			return nil
		},
	}
	batchRepo := &FakeIsiSoalBatchRepo{
		ImportBankSoalVersionFn: func(_ context.Context, _, _ int64, got importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
			mu.Lock()
			payload = got
			mu.Unlock()
			return 1, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go newWorkerForTest(jobRepo, batchRepo, tmpDir, 10*time.Millisecond).Start(ctx)

	select {
	case <-done:
		cancel()
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for import job completion")
	}

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, payload.SoalList, 1)
	require.Len(t, payload.SoalList[0].Opsi, 2)

	image := payload.SoalList[0].Opsi[0].IsiContent.Blocks[0].Children[0]
	assert.Equal(t, "image", image.Type)
	assert.Equal(t, "/uploads/image/image1.png", image.Src)
	_, err := os.Stat(filepath.Join(tmpDir, "image1.png"))
	require.NoError(t, err)
}

func buildWorkerDocxWithOptionImage() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	document := `<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<w:body>
	<w:p><w:r><w:drawing><a:blip r:embed="rId1"/></w:drawing></w:r></w:p>
	<w:p><w:r><w:t>[Q:PG] Pilih gambar yang benar</w:t></w:r></w:p>
	<w:p><w:r><w:t>[A] [IMG] Gambar pilihan</w:t></w:r></w:p>
	<w:p><w:r><w:t>[B] Jawaban biasa</w:t></w:r></w:p>
	<w:p><w:r><w:t>[ANS] A</w:t></w:r></w:p>
</w:body></w:document>`
	fw, _ := zw.Create("word/document.xml")
	_, _ = fw.Write([]byte(document))

	rels, _ := zw.Create("word/_rels/document.xml.rels")
	_, _ = rels.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image1.png"/>
</Relationships>`))

	media, _ := zw.Create("word/media/image1.png")
	_, _ = media.Write([]byte("fake image"))

	_ = zw.Close()
	return buf.Bytes()
}
