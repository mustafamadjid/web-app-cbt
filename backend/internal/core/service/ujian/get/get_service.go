package ujian_service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

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

func (r *GetUjianService) GetAllUjianService(ctx context.Context, filter query.ListUjianFilter) ([]ujian.Ujian, error) {
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

// -----------------------
// Sanitizer and validator
// -----------------------

var (
	errInvalidTanggalUjian = errors.New("invalid tanggal ujian")
	errInvalidTahun        = errors.New("invalid tahun")
	errInvalidTingkatKelas = errors.New("invalid tingkat kelas")
	errInvalidRuangUjian   = errors.New("invalid ruang ujian")
	errInvalidPeserta      = errors.New("invalid peserta ujian filter")
	errInvalidJawaban      = errors.New("invalid jawaban ujian filter")
)


func sanitizeAndValidateListUjianFilter(filter query.ListUjianFilter) (query.ListUjianFilter, error) {
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

	if filter.TanggalUjian != nil {
		tanggalUjian := strings.TrimSpace(*filter.TanggalUjian)
		if tanggalUjian == "" {
			return filter, errInvalidTanggalUjian
		}

		if _, err := time.Parse("2006-01-02", tanggalUjian); err != nil {
			return filter, errInvalidTanggalUjian
		}

		filter.TanggalUjian = &tanggalUjian
	}

	if filter.Tahun != nil {
		tahun := strings.TrimSpace(*filter.Tahun)
		if tahun == "" {
			return filter, errInvalidTahun
		}

		if len(tahun) != 4 {
			return filter, errInvalidTahun
		}

		tahunInt, err := strconv.Atoi(tahun)
		if err != nil || tahunInt <= 0 {
			return filter, errInvalidTahun
		}

		filter.Tahun = &tahun
	}

	if filter.TingkatKelas != nil && *filter.TingkatKelas <= 0 {
		return filter, errInvalidTingkatKelas
	}

	if filter.RuangUjian != nil && *filter.RuangUjian <= 0 {
		return filter, errInvalidRuangUjian
	}

	return filter, nil
}

func sanitizeAndValidateJawabanFilter(filter ujian.JawabanUjianSiswa) (ujian.JawabanUjianSiswa, error) {
	if filter.IdJawaban < 0 || filter.IdPesertaUjian < 0 || filter.IdSoal < 0 {
		return filter, errInvalidJawaban
	}

	if filter.IdPilihan != nil && *filter.IdPilihan <= 0 {
		return filter, errInvalidJawaban
	}

	if filter.JawabanEssay != nil {
		jawabanEssay := strings.TrimSpace(*filter.JawabanEssay)
		if jawabanEssay == "" {
			return filter, errInvalidJawaban
		}
		filter.JawabanEssay = &jawabanEssay
	}

	return filter, nil
}

func validateUjianID(id ujian.ID) error {
	if id <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validatePesertaFilter(filter ujian.PesertaUjian) error {
	if filter.IdPesertaUjian < 0 || filter.IdJadwalUjian < 0 || filter.IdSiswa < 0 {
		return errInvalidPeserta
	}
	return nil
}
