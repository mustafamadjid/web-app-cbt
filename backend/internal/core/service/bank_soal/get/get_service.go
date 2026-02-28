package bank_soal_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type GetBankSoalService struct {
	repo bank_soal_repo.BankSoalRepository
}

func NewGetBankSoalService(repo bank_soal_repo.BankSoalRepository) *GetBankSoalService {
	return &GetBankSoalService{
		repo: repo,
	}
}

func (r *GetBankSoalService) GetBankSoalService(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	logger := corelog.FromContext(ctx)

	filter = sanitizeGetBankSoalFilter(filter)

	items, err := r.repo.GetBankSoal(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get bank soal", "layer", "core.service", "op", "bank_soal.get", "err", err)
		return nil, err
	}
	return items, nil
}

func (r *GetBankSoalService) GetBankSoalByIdService(ctx context.Context, idBankSoal bank_soal.ID) (bank_soal.BankSoal, error) {
	logger := corelog.FromContext(ctx)

	if err := validateBankSoalID(idBankSoal); err != nil {
		logger.Error(ctx, "failed get bank soal by id", "layer", "core.service", "op", "bank_soal.get_by_id", "err", coreerror.ErrMissingId)
		return bank_soal.BankSoal{}, nil
	}

	item, err := r.repo.GetBankSoalById(ctx, idBankSoal)
	if err != nil {
		logger.Error(ctx, "failed get bank soal by id", "layer", "core.service", "op", "bank_soal.get_by_id", "err", err)
		return bank_soal.BankSoal{}, err
	}
	return item, nil
}

func (r *GetBankSoalService) GetBankSoalByGuruService(ctx context.Context, idGuru bank_soal.ID) ([]bank_soal.BankSoal, error) {
	logger := corelog.FromContext(ctx)

	if err := validateBankSoalID(idGuru); err != nil {
		logger.Error(ctx, "failed get bank soal by guru", "layer", "core.service", "op", "bank_soal.get_by_guru", "err", coreerror.ErrMissingId)
		return nil, err
	}

	items, err := r.repo.GetBankSoalByGuru(ctx, idGuru)
	if err != nil {
		logger.Error(ctx, "failed get bank soal by guru", "layer", "core.service", "op", "bank_soal.get_by_guru", "err", err)
		return nil, err
	}
	return items, nil
}

// -----------------------
// Sanitizer and validator
// -----------------------

func sanitizeGetBankSoalFilter(filter query.BankSoalFilter) query.BankSoalFilter {
	filter.Search = strings.TrimSpace(filter.Search)

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return filter
}

func validateBankSoalID(id bank_soal.ID) error {
	if id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
