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

func(r *GetBankSoalService)GetBankSoalService(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	logger := corelog.FromContext(ctx)

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

	items, err := r.repo.GetBankSoal(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get bank soal", "layer", "core.service", "op", "bank_soal.get", "err", err)
		return nil, err
	}
	return items, nil
}

func(r *GetBankSoalService)GetBankSoalByIdService(ctx context.Context, idBankSoal bank_soal.ID) (bank_soal.BankSoal, error) {
	logger := corelog.FromContext(ctx)

	if idBankSoal <= 0 {
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

func(r *GetBankSoalService)GetBankSoalByGuruService(ctx context.Context, idGuru bank_soal.ID) ([]bank_soal.BankSoal, error) {
	logger := corelog.FromContext(ctx)

	if idGuru <= 0 {
		logger.Error(ctx, "failed get bank soal by guru", "layer", "core.service", "op", "bank_soal.get_by_guru", "err", coreerror.ErrMissingId)
		return nil, coreerror.ErrMissingId
	}

	items, err := r.repo.GetBankSoalByGuru(ctx, idGuru)
	if err != nil {
		logger.Error(ctx, "failed get bank soal by guru", "layer", "core.service", "op", "bank_soal.get_by_guru", "err", err)
		return nil, err
	}
	return items, nil
}