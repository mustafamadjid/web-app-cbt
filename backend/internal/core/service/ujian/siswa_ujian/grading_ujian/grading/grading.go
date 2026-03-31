package gradingujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	grading_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/grading"
)

type GradingUjianService struct {
	jawabanrepo   ujian_repo.JawabanUjianRepository
	soalUjianRepo ujian_repo.SoalUjianRepository
	bankSoalRepo  bank_soal_repo.BankSoalRepository
	ujianRepo     ujian_repo.UjianRepository
	gradingRepo   grading_repo.GradingUjianRepository
}

func NewGradingUjianService(
	jawabanRepo ujian_repo.JawabanUjianRepository,
	soalUjianRepo ujian_repo.SoalUjianRepository,
	bankSoalRepo bank_soal_repo.BankSoalRepository,
	ujianRepo ujian_repo.UjianRepository,
	gradingRepo grading_repo.GradingUjianRepository) *GradingUjianService {
	return &GradingUjianService{
		jawabanrepo:   jawabanRepo,
		soalUjianRepo: soalUjianRepo,
		bankSoalRepo:  bankSoalRepo,
		ujianRepo:     ujianRepo,
		gradingRepo:   gradingRepo,
	}
}

func (r *GradingUjianService) GradingUjianPilgan(ctx context.Context, idAttempt int) error {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	// Retrieve jawaban ujian by attempt id
	jawabanSiswa, err := r.jawabanrepo.GetJawabanUjianByAttemptId(ctx, ujian.ID(idAttempt))
	if err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	// Retrieve id bank soal by attempt id
	idBankSoalAktif, err := r.bankSoalRepo.GetIdBankSoalByAttemptId(ctx, ujian.ID(idAttempt))
	if err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	// Retrieve soal ujian by bank soal id
	soalUjian, err := r.soalUjianRepo.GetSoalUjianByBankSoal(ctx, ujian.ID(idBankSoalAktif))
	if err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	// Retrieve id ujian by attempt id
	idUjian, err := r.ujianRepo.GetIdUjianByAttempt(ctx, ujian.ID(idAttempt))
	if err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	// Calculate total score
	totalNilai, soalBenar, soalSalah, err := r.TotalScore(jawabanSiswa, soalUjian, idUjian)
	if err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	hasilUjian := ujian.HasilUjian{
		IdAttempt: ujian.ID(idAttempt),
	}

	// Upsert score to hasil ujian table
	if err := r.gradingRepo.UpsertNilaiToHasilUjian(ctx, totalNilai, hasilUjian); err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	// Upserting Soal to Statistik Soal
	if err := r.UpsertingToStatistikSoal(ctx, soalBenar, soalSalah); err != nil {
		logger.Error(ctx, "failed grading ujian", "layer", "core.service", "op", "ujian.grading", "err", err)
		return err
	}

	return nil
}

func (r *GradingUjianService) TotalScore(jawabanSiswa []ujian.JawabanUjian, soalUjian []ujian.SoalUjianSiswa, idUjian ujian.ID) (float64, []ujian.StatistikSoal, []ujian.StatistikSoal, error) {
	if len(soalUjian) == 0 {
		return 0, nil, nil, coreerror.ErrArrayHasNoElement
	}

	type OpsiInfo struct {
		IsBenar   bool
		BobotSoal float64
	}

	opsiMap := make(map[ujian.ID]OpsiInfo)

	for _, soal := range soalUjian {
		for _, opsi := range soal.OpsiJawaban {
			opsiMap[opsi.IdPilihanGanda] = OpsiInfo{
				IsBenar:   opsi.IsBenar,
				BobotSoal: soal.BobotSoal,
			}
		}
	}

	statistikSoalBenar := make([]ujian.StatistikSoal, 0)
	statistikSoalSalah := make([]ujian.StatistikSoal, 0)

	var totalNilai float64

	for _, jawabanSiswa := range jawabanSiswa {
		if jawabanSiswa.IdPilihan == nil {
			continue
		}

		info, ok := opsiMap[*jawabanSiswa.IdPilihan]
		if !ok {
			continue
		}

		if info.IsBenar {
			totalNilai += info.BobotSoal
			statistikSoalBenar = append(statistikSoalBenar, ujian.StatistikSoal{
				IDSoal:  jawabanSiswa.IdSoal,
				IDUjian: idUjian,
			})
		} else {
			statistikSoalSalah = append(statistikSoalSalah, ujian.StatistikSoal{
				IDSoal:  jawabanSiswa.IdSoal,
				IDUjian: idUjian,
			})
		}
	}

	return totalNilai, statistikSoalBenar, statistikSoalSalah, nil
}

func (r *GradingUjianService) UpsertingToStatistikSoal(ctx context.Context, soalBenar []ujian.StatistikSoal, soalSalah []ujian.StatistikSoal) error {
	logger := corelog.FromContext(ctx)

	if len(soalBenar) > 0 {
		err := r.gradingRepo.UpsertJawabanBenarToStatistikSoal(ctx, soalBenar)
		if err != nil {
			logger.Error(ctx, "failed upsert jawaban benar to statistik soal", "layer", "core.service", "op", "ujian.grading", "err", err)
			return err
		}
	}

	if len(soalSalah) > 0 {
		err := r.gradingRepo.UpsertJawabanSalahToStatistikSoal(ctx, soalSalah)
		if err != nil {
			logger.Error(ctx, "failed upsert jawaban salah to statistik soal", "layer", "core.service", "op", "ujian.grading", "err", err)
			return err
		}
	}

	return nil
}
