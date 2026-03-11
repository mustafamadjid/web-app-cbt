package ujian_repo

import (
	"context"
	"time"
)

type SiswaUjianChecker interface {
	CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int)(bool,int,int,error) // return bool, idPesertaUjian,idJadwalUjian, error
	CheckTokenUjian(ctx context.Context,token string) (bool, error)

	GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error)
}
