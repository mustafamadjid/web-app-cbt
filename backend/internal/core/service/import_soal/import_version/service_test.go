package import_version_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	faketest "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/fake_test"
	importversion "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/import_version"
	"github.com/stretchr/testify/assert"
)

func TestServiceExecute(t *testing.T) {
	t.Parallel()

	validPayload := []importsoal.ParsedSoal{
		{
			TipeSoal:   "pilihan_ganda",
			Pertanyaan: "2 + 2 = ?",
			BobotSoal:  1,
			Opsi: []importsoal.ParsedOpsi{
				{Label: "A", Isi: "3", IsBenar: false},
				{Label: "B", Isi: "4", IsBenar: true},
			},
		},
	}

	t.Run("invalid payload: PG has more than one correct option", func(t *testing.T) {
		t.Parallel()

		repo := &faketest.FakeIsiSoalBatchRepo{}
		svc := importversion.NewService(repo)

		_, err := svc.Execute(context.Background(), importversion.Cmd{
			BankID: 1,
			UserID: 1,
			Payload: []importsoal.ParsedSoal{
				{
					TipeSoal:   "pilihan_ganda",
					Pertanyaan: "x",
					BobotSoal:  1,
					Opsi: []importsoal.ParsedOpsi{
						{Label: "A", Isi: "1", IsBenar: true},
						{Label: "B", Isi: "2", IsBenar: true},
					},
				},
			},
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, coreerror.ErrInvalidInput)
	})

	t.Run("repo error propagates", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("repo failed")
		repo := &faketest.FakeIsiSoalBatchRepo{
			ImportBankSoalVersionFn: func(_ context.Context, _ int64, _ int64, _ importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
				return 0, expectedErr
			},
		}

		svc := importversion.NewService(repo)
		_, err := svc.Execute(context.Background(), importversion.Cmd{
			BankID:  10,
			UserID:  20,
			Payload: validPayload,
		})

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("success returns version id", func(t *testing.T) {
		t.Parallel()

		repo := &faketest.FakeIsiSoalBatchRepo{
			ImportBankSoalVersionFn: func(_ context.Context, bankID, userID int64, payload importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
				assert.Equal(t, int64(5), bankID)
				assert.Equal(t, int64(8), userID)
				assert.Len(t, payload.SoalList, 1)
				return 99, nil
			},
		}

		svc := importversion.NewService(repo)
		result, err := svc.Execute(context.Background(), importversion.Cmd{
			BankID:  5,
			UserID:  8,
			Payload: validPayload,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(99), result.VersionID)
	})
}
