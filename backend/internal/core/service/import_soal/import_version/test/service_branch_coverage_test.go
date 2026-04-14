package import_version_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	importversion "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/import_version"
	"github.com/stretchr/testify/assert"
)

func TestImportVersionService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	validPayload := []importsoal.ParsedSoal{
		{
			TipeSoal:     "pilihan_ganda",
			Pertanyaan:   "Soal 1",
			BobotSoal:    10,
			KunciJawaban: "A",
			Opsi: []importsoal.ParsedOpsi{
				{Label: "A", Isi: "Jawaban A", IsBenar: true},
				{Label: "B", Isi: "Jawaban B", IsBenar: false},
			},
		},
	}

	tests := []struct {
		name       string
		cmd        importversion.Cmd
		repo       *FakeIsiSoalBatchRepo
		wantErr    error
		wantCalled bool
		wantID     int64
	}{
		{
			name:    "branch 1 -> bank id tidak valid",
			cmd:     importversion.Cmd{BankID: 0, UserID: 1, Payload: validPayload},
			repo:    &FakeIsiSoalBatchRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:    "branch 2 -> user id tidak valid",
			cmd:     importversion.Cmd{BankID: 1, UserID: 0, Payload: validPayload},
			repo:    &FakeIsiSoalBatchRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:    "branch 3 -> payload kosong",
			cmd:     importversion.Cmd{BankID: 1, UserID: 2},
			repo:    &FakeIsiSoalBatchRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name: "branch 4 -> lebih dari satu opsi benar",
			cmd: importversion.Cmd{
				BankID: 1,
				UserID: 2,
				Payload: []importsoal.ParsedSoal{
					{
						TipeSoal:     "pilihan_ganda",
						Pertanyaan:   "Soal 1",
						BobotSoal:    10,
						KunciJawaban: "A",
						Opsi: []importsoal.ParsedOpsi{
							{Label: "A", Isi: "Jawaban A", IsBenar: true},
							{Label: "B", Isi: "Jawaban B", IsBenar: true},
						},
					},
				},
			},
			repo:    &FakeIsiSoalBatchRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name: "branch 5 -> repo import version gagal",
			cmd:  importversion.Cmd{BankID: 1, UserID: 2, Payload: validPayload},
			repo: &FakeIsiSoalBatchRepo{
				ImportBankSoalVersionFn: func(context.Context, int64, int64, importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
					return 0, repoErr
				},
			},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name: "branch 6 -> berhasil import version",
			cmd:  importversion.Cmd{BankID: 1, UserID: 2, Payload: validPayload},
			repo: &FakeIsiSoalBatchRepo{
				ImportBankSoalVersionFn: func(context.Context, int64, int64, importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
					return 99, nil
				},
			},
			wantCalled: true,
			wantID:     99,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := importversion.NewService(tc.repo)
			result, err := svc.Execute(ctx, tc.cmd)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.ImportBankSoalVersionCalled)
			assert.Equal(t, tc.wantID, result.VersionID)
		})
	}
}
