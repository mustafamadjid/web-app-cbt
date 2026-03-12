package ujian_service

import (
	"context"
	"errors"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type GetUjianService struct {
	repo ujian_repo.ListUjianRepository
}

func NewGetujianService(repo ujian_repo.ListUjianRepository) *GetUjianService {
	return &GetUjianService{
		repo: repo,
	}
}

func (r *GetUjianService) GetAllUjianService(ctx context.Context, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	logger := corelog.FromContext(ctx)

	var err error
	filter, err = sanitizeAndValidateListUjianFilter(filter)
	if err != nil {
		logger.Error(ctx, "failed get ujian", "layer", "core.service", "op", "ujian.get.filter", "err", err)
		return nil, coreerror.ErrInvalidInput
	}

	items, err := r.repo.GetAllUjian(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get ujian", "layer", "core.service", "op", "ujian.get", "err", err)
		return nil, err
	}

	items, err = sanitizeAndValidateListUjianItems(items)
	if err != nil {
		logger.Error(ctx, "failed get ujian", "layer", "core.service", "op", "ujian.get.sanitize", "err", err)
		return nil, coreerror.ErrInvalidInput
	}

	return items, nil
}

func (r *GetUjianService) GetUjianByIdService(ctx context.Context, idUjian ujian.ID) (ujian.ListUjian, error) {
	logger := corelog.FromContext(ctx)

	if err := validateUjianID(idUjian); err != nil {
		logger.Error(ctx, "failed get ujian by id", "layer", "core.service", "op", "ujian.get_by_id", "err", err)
		return ujian.ListUjian{}, err
	}

	item, err := r.repo.GetUjianById(ctx, idUjian)
	if err != nil {
		logger.Error(ctx, "failed get ujian by id", "layer", "core.service", "op", "ujian.get_by_id", "err", err)
		return ujian.ListUjian{}, err
	}

	return item, nil
}

var (
	errInvalidListUjian   = errors.New("invalid list ujian")
	errInvalidStatusUjian = errors.New("invalid status ujian")
	errInvalidWaktuUjian  = errors.New("invalid waktu ujian")
	errInvalidNamaUjian   = errors.New("invalid nama ujian")

	errInvalidTanggalUjian  = errors.New("invalid tanggal ujian")
	errInvalidTahun         = errors.New("invalid tahun")
	errInvalidTingkatKelas  = errors.New("invalid tingkat kelas")
	errInvalidRuangUjian    = errors.New("invalid ruang ujian")
	errInvalidKategoriUjian = errors.New("invalid kategori ujian")
	errInvalidJadwalUjian   = errors.New("invalid jadwal ujian")
	errInvalidTokenUjian    = errors.New("invalid token ujian")
)
