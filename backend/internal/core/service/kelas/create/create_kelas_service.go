package kelas_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type CreateKelasService struct {
	kelasRepo kelas_repo.KelasRepository
}

func NewCreateKelasService(kelasRepo kelas_repo.KelasRepository) *CreateKelasService {
	return &CreateKelasService{
		kelasRepo: kelasRepo,
	}
}

func (r *CreateKelasService) CreateTingkatKelas(ctx context.Context, cmd CreateTingkatKelasCmd) error {
	logger := corelog.FromContext(ctx)

	exist, err := r.kelasRepo.ExistTingkatKelas(ctx, cmd.TingkatKelas)
	if err != nil {
		logger.Error(ctx, "failed creating tingkat kelas", "layer", "core.service", "op", "kelas.create_tingkat_kelas.existTingkatKelas", "err", err)
		return err
	}

	if exist {
		return coreerror.ErrTingkatKelasExist
	}

	if err := r.kelasRepo.CreateTingkatKelas(ctx, cmd.TingkatKelas); err != nil {
		logger.Error(ctx, "failed creating tingkat kelas", "layer", "core.service", "op", "kelas.create_tingkat_kela.CreateTingkatKelas", "err", err)
		return err
	}

	return nil
}

func (r *CreateKelasService) CreateNamaKelas(ctx context.Context, cmd CreateNamaKelasCmd) error {
	logger := corelog.FromContext(ctx)

	cmd = sanitizeCreateNamaKelasCmd(cmd)

	exist, err := r.kelasRepo.ExistNamaKelas(ctx, cmd.NamaKelas)
	if err != nil {
		logger.Error(ctx, "failed creating nama kelas", "layer", "core.service", "op", "kelas.create_nama_kelas.existNamaKelas", "err", err)
		return err
	}

	if exist {
		return coreerror.ErrNamaKelasExist
	}

	cmdData := kelas.NamaKelas{
		IdTingkatKelas: cmd.IdTingkatKelas,
		NamaKelas:      cmd.NamaKelas,
	}

	if err := r.kelasRepo.CreateNamaKelas(ctx, cmdData); err != nil {
		logger.Error(ctx, "failed creating nama kelas", "layer", "core.service", "op", "kelas.create_nama_kela.CreateNamaKelas", "err", err)
		return err
	}

	return nil
}
