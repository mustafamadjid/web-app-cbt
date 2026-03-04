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

func (r *GetUjianService) GetUjianByIdService(ctx context.Context, id ujian.ID) (ujian.Ujian, error) {
	logger := corelog.FromContext(ctx)

	if err := validateUjianID(id); err != nil {
		logger.Error(ctx, "failed get ujian by id", "layer", "core.service", "op", "ujian.get_by_id", "err", coreerror.ErrMissingId)
		return ujian.Ujian{}, err
	}

	item, err := r.repo.GetUjianById(ctx, id)
	if err != nil {
		logger.Error(ctx, "failed get ujian by id", "layer", "core.service", "op", "ujian.get_by_id", "err", err)
		return ujian.Ujian{}, err
	}
	return item, nil
}

func (r *GetUjianService) GetJadwalUjianByIdService(ctx context.Context, id ujian.ID) (ujian.JadwalUjian, error) {
	logger := corelog.FromContext(ctx)

	if err := validateUjianID(id); err != nil {
		logger.Error(ctx, "failed get jadwal ujian by id", "layer", "core.service", "op", "ujian.get_jadwal_by_id", "err", coreerror.ErrMissingId)
		return ujian.JadwalUjian{}, err
	}

	item, err := r.repo.GetJadwalUjianById(ctx, id)
	if err != nil {
		logger.Error(ctx, "failed get jadwal ujian by id", "layer", "core.service", "op", "ujian.get_jadwal_by_id", "err", err)
		return ujian.JadwalUjian{}, err
	}
	item = sanitizeJadwalUjian(item)
	if err := validateJadwalUjian(item); err != nil {
		logger.Error(ctx, "failed get jadwal ujian by id", "layer", "core.service", "op", "ujian.get_jadwal_by_id.validate", "err", err)
		return ujian.JadwalUjian{}, coreerror.ErrInvalidInput
	}
	return item, nil
}

func (r *GetUjianService) GetAllPesertaUjianService(ctx context.Context, peserta ujian.PesertaUjian) ([]ujian.PesertaUjian, error) {
	logger := corelog.FromContext(ctx)

	if err := validatePesertaFilter(peserta); err != nil {
		logger.Error(ctx, "failed get peserta ujian", "layer", "core.service", "op", "ujian.get_all_peserta.filter", "err", err)
		return nil, coreerror.ErrInvalidInput
	}

	items, err := r.repo.GetAllPesertaUjian(ctx, peserta)
	if err != nil {
		logger.Error(ctx, "failed get peserta ujian", "layer", "core.service", "op", "ujian.get_all_peserta", "err", err)
		return nil, err
	}
	return items, nil
}

func (r *GetUjianService) GetPesertaUjianBySiswaService(ctx context.Context, idSiswa ujian.ID, peserta ujian.PesertaUjian) (ujian.PesertaUjian, error) {
	logger := corelog.FromContext(ctx)

	if err := validateUjianID(idSiswa); err != nil {
		logger.Error(ctx, "failed get peserta ujian by siswa", "layer", "core.service", "op", "ujian.get_peserta_by_siswa", "err", coreerror.ErrMissingId)
		return ujian.PesertaUjian{}, err
	}

	if err := validatePesertaFilter(peserta); err != nil {
		logger.Error(ctx, "failed get peserta ujian by siswa", "layer", "core.service", "op", "ujian.get_peserta_by_siswa.filter", "err", err)
		return ujian.PesertaUjian{}, coreerror.ErrInvalidInput
	}

	item, err := r.repo.GetPesertaUjianBySiswa(ctx, idSiswa, peserta)
	if err != nil {
		logger.Error(ctx, "failed get peserta ujian by siswa", "layer", "core.service", "op", "ujian.get_peserta_by_siswa", "err", err)
		return ujian.PesertaUjian{}, err
	}
	return item, nil
}

func (r *GetUjianService) GetAllJawabanUjianSiswaService(ctx context.Context, jawaban ujian.JawabanUjianSiswa) ([]ujian.JawabanUjianSiswa, error) {
	logger := corelog.FromContext(ctx)

	jawaban, err := sanitizeAndValidateJawabanFilter(jawaban)
	if err != nil {
		logger.Error(ctx, "failed get jawaban ujian siswa", "layer", "core.service", "op", "ujian.get_all_jawaban.filter", "err", err)
		return nil, coreerror.ErrInvalidInput
	}

	items, err := r.repo.GetAllJawabanUjianSiswa(ctx, jawaban)
	if err != nil {
		logger.Error(ctx, "failed get jawaban ujian siswa", "layer", "core.service", "op", "ujian.get_all_jawaban", "err", err)
		return nil, err
	}
	return items, nil
}

func (r *GetUjianService) GetJawabanBySiswaService(ctx context.Context, idSiswa ujian.ID, jawaban ujian.JawabanUjianSiswa) (ujian.JawabanUjianSiswa, error) {
	logger := corelog.FromContext(ctx)

	if err := validateUjianID(idSiswa); err != nil {
		logger.Error(ctx, "failed get jawaban by siswa", "layer", "core.service", "op", "ujian.get_jawaban_by_siswa", "err", coreerror.ErrMissingId)
		return ujian.JawabanUjianSiswa{}, err
	}

	jawaban, err := sanitizeAndValidateJawabanFilter(jawaban)
	if err != nil {
		logger.Error(ctx, "failed get jawaban by siswa", "layer", "core.service", "op", "ujian.get_jawaban_by_siswa.filter", "err", err)
		return ujian.JawabanUjianSiswa{}, coreerror.ErrInvalidInput
	}

	item, err := r.repo.GetJawabanBySiswa(ctx, idSiswa, jawaban)
	if err != nil {
		logger.Error(ctx, "failed get jawaban by siswa", "layer", "core.service", "op", "ujian.get_jawaban_by_siswa", "err", err)
		return ujian.JawabanUjianSiswa{}, err
	}
	return item, nil
}

var (
	errInvalidListUjian   = errors.New("invalid list ujian")
	errInvalidStatusUjian = errors.New("invalid status ujian")
	errInvalidWaktuUjian  = errors.New("invalid waktu ujian")
	errInvalidNamaUjian   = errors.New("invalid nama ujian")

	errInvalidTanggalUjian = errors.New("invalid tanggal ujian")
	errInvalidTahun        = errors.New("invalid tahun")
	errInvalidTingkatKelas = errors.New("invalid tingkat kelas")
	errInvalidRuangUjian   = errors.New("invalid ruang ujian")
	errInvalidPeserta      = errors.New("invalid peserta ujian filter")
	errInvalidJawaban      = errors.New("invalid jawaban ujian filter")
	errInvalidJadwalUjian  = errors.New("invalid jadwal ujian")
	errInvalidTokenUjian   = errors.New("invalid token ujian")
)
