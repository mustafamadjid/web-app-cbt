package import_version

import (
	"context"
	"fmt"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type Cmd struct {
	BankID  int64
	UserID  int64
	Payload []importsoal.ParsedSoal
}

type Result struct {
	VersionID int64
}

type Service struct {
	repo importsoalrepo.IsiSoalBatchRepo
}

func NewService(repo importsoalrepo.IsiSoalBatchRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Execute(ctx context.Context, cmd Cmd) (Result, error) {
	logger := corelog.FromContext(ctx)

	if cmd.BankID <= 0 {
		return Result{}, fmt.Errorf("%w: bank id must be greater than 0", coreerror.ErrInvalidInput)
	}
	if cmd.UserID <= 0 {
		return Result{}, fmt.Errorf("%w: user id must be greater than 0", coreerror.ErrInvalidInput)
	}
	if len(cmd.Payload) == 0 {
		return Result{}, fmt.Errorf("%w: payload must not be empty", coreerror.ErrInvalidInput)
	}
	if err := validateExactlyOneCorrectOption(cmd.Payload); err != nil {
		return Result{}, err
	}

	versionID, err := s.repo.ImportBankSoalVersion(ctx, cmd.BankID, cmd.UserID, importsoalrepo.ImportBankSoalVersionPayload{
		SoalList: cmd.Payload,
	})
	if err != nil {
		logger.Error(ctx, "failed importing bank soal version", "layer", "core.service", "op", "import_soal.import_version", "bank_id", cmd.BankID, "err", err)
		return Result{}, err
	}

	return Result{
		VersionID: versionID,
	}, nil
}

func validateExactlyOneCorrectOption(soalList []importsoal.ParsedSoal) error {
	for i, soal := range soalList {
		if soal.TipeSoal != "pilihan_ganda" {
			continue
		}

		correctCount := 0
		for _, opsi := range soal.Opsi {
			if opsi.IsBenar {
				correctCount++
			}
		}

		if correctCount != 1 {
			return fmt.Errorf("%w: soal ke-%d harus memiliki tepat satu opsi benar", coreerror.ErrInvalidInput, i+1)
		}
	}
	return nil
}
