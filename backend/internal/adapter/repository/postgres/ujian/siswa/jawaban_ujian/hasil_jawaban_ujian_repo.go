package jawabanujian_repo

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

func(r *JawabanUjianRepo)ListHasilJawabanUjianByIdAttempt(ctx context.Context, idAttempt ujian.ID) ([]ujian.HasilJawabanUjian, error){
	
}