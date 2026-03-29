package gradingujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	bank_soal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type GradingUjianService struct {
	jawabanrepo ujian_repo.JawabanUjianRepository
	soalUjianRepo  ujian_repo.SoalUjianRepository
	bankSoalRepo bank_soal_repo.BankSoalRepository
}

func NewGradingUjianService(
	jawabanRepo ujian_repo.JawabanUjianRepository, 
	soalUjianRepo ujian_repo.SoalUjianRepository, 
	bankSoalRepo bank_soal_repo.BankSoalRepository) *GradingUjianService {
	return &GradingUjianService{
		jawabanrepo: jawabanRepo,
		soalUjianRepo: soalUjianRepo,
		bankSoalRepo: bankSoalRepo,
	}
}

func(r *GradingUjianService) GradingUjian(ctx context.Context,idAttempt int) {
	logger := corelog.FromContext(ctx)

	if idAttempt <= 0 {
		logger.Error(ctx,"failed grading ujian","layer","core.service","op","ujian.grading","err",coreerror.ErrMissingId)
		return
	}

	// Retrieve jawaban ujian by attempt id
	jawabanSiswa, err := r.jawabanrepo.GetJawabanUjianByAttemptId(ctx,ujian.ID(idAttempt))
	if err != nil {
		logger.Error(ctx,"failed grading ujian","layer","core.service","op","ujian.grading","err",err)
		return
	}

	// Retrieve id bank soal by attempt id
	idBankSoalAktif,err := r.bankSoalRepo.GetIdBankSoalByAttemptId(ctx,idAttempt)
	if err != nil {
		logger.Error(ctx,"failed grading ujian","layer","core.service","op","ujian.grading","err",err)
		return
	}

	// Retrieve soal ujian by bank soal id
	soalUjian, err := r.soalUjianRepo.GetSoalUjianByBankSoal(ctx,ujian.ID(idBankSoalAktif))
	if err != nil {
		logger.Error(ctx,"failed grading ujian","layer","core.service","op","ujian.grading","err",err)
		return
	}
	

	// Grading process

	
}

func(r *GradingUjianService) TotalScore(jawabanSiswa []ujian.JawabanUjian, soalUjian []ujian.SoalUjianSiswa)float64 {
	if len(jawabanSiswa) != len(soalUjian) || len(jawabanSiswa) == 0 || len(soalUjian) == 0 {
		return 0
	}

	type OpsiInfo struct {
		IsBenar   bool
		BobotSoal float64
	}

	opsiMap := make(map[ujian.ID]OpsiInfo)
	
	for _,soal := range soalUjian {
		for _, opsi := range soal.OpsiJawaban {
			opsiMap[opsi.IdPilihanGanda] = OpsiInfo{
				IsBenar:   opsi.IsBenar,
				BobotSoal: soal.BobotSoal,
			}
		}
	}

	var totalNilai float64
	
	for _,jawabanSiswa := range jawabanSiswa {
		if jawabanSiswa.IdPilihan == nil {
			continue
		}

		info, ok := opsiMap[*jawabanSiswa.IdPilihan]
		if !ok {
			continue
		}

		if info.IsBenar {
			totalNilai += info.BobotSoal
		}
	}
	return totalNilai
}