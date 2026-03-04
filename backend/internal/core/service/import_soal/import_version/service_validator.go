package import_version

import (
	"fmt"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

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
