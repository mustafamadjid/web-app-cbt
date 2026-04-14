package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/create"
	"github.com/stretchr/testify/assert"
)

func TestAttemptUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	deadline := now.Add(30 * time.Minute)
	checkErr := errors.New("checker error")
	createErr := errors.New("create error")

	tests := []struct {
		name       string
		idSiswa    int
		idJadwal   int
		token      string
		waktu      time.Time
		checker    *fakeSiswaUjianChecker
		repo       *fakeAttemptUjianRepo
		wantErr    error
		wantCreate bool
	}{
		{
			name:       "path 1 -> id siswa tidak valid",
			idSiswa:    0,
			idJadwal:   20,
			token:      "TOKEN",
			waktu:      now,
			checker:    &fakeSiswaUjianChecker{},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCreate: false,
		},
		{
			name:       "path 2 -> id jadwal ujian tidak valid",
			idSiswa:    10,
			idJadwal:   0,
			token:      "TOKEN",
			waktu:      now,
			checker:    &fakeSiswaUjianChecker{},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCreate: false,
		},
		{
			name:       "path 3 -> token ujian kosong",
			idSiswa:    10,
			idJadwal:   20,
			token:      "   ",
			waktu:      now,
			checker:    &fakeSiswaUjianChecker{},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingTokenUjian,
			wantCreate: false,
		},
		{
			name:       "path 4 -> waktu attempt kosong",
			idSiswa:    10,
			idJadwal:   20,
			token:      "TOKEN",
			checker:    &fakeSiswaUjianChecker{},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrTimeEmpty,
			wantCreate: false,
		},
		{
			name:     "path 5 -> validasi peserta gagal",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidErr: checkErr,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    checkErr,
			wantCreate: false,
		},
		{
			name:     "path 6 -> peserta tidak diizinkan attempt",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet: false,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrPesertaNotAllowedToAttemptJadwal,
			wantCreate: false,
		},
		{
			name:     "path 7 -> deadline ujian gagal didapat",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineErr:  checkErr,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    checkErr,
			wantCreate: false,
		},
		{
			name:     "path 8 -> waktu attempt melewati deadline",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    deadline.Add(1 * time.Minute),
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineRet:  deadline,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrWaktuAttemptPesertaInvalid,
			wantCreate: false,
		},
		{
			name:     "path 9 -> validasi token ujian gagal",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineRet:  deadline,
				checkTokenErr:   checkErr,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    checkErr,
			wantCreate: false,
		},
		{
			name:     "path 10 -> token ujian tidak valid",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineRet:  deadline,
				checkTokenRet:   false,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrTokenUjianInvalid,
			wantCreate: false,
		},
		{
			name:     "path 11 -> create attempt generic error",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineRet:  deadline,
				checkTokenRet:   true,
			},
			repo:       &fakeAttemptUjianRepo{createErr: createErr},
			wantErr:    createErr,
			wantCreate: true,
		},
		{
			name:     "path 12 -> active attempt dianggap resume",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineRet:  deadline,
				checkTokenRet:   true,
			},
			repo:       &fakeAttemptUjianRepo{createErr: coreerror.ErrSiswaHasActiveAttempt},
			wantCreate: true,
		},
		{
			name:     "path 13 -> attempt ujian berhasil",
			idSiswa:  10,
			idJadwal: 20,
			token:    "TOKEN",
			waktu:    now,
			checker: &fakeSiswaUjianChecker{
				checkValidRet:   true,
				checkValidIDRet: 77,
				getDeadlineRet:  deadline,
				checkTokenRet:   true,
			},
			repo:       &fakeAttemptUjianRepo{},
			wantCreate: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewAttemptUjianService(tc.checker, tc.repo)
			err := svc.AttemptUjian(context.Background(), tc.idSiswa, tc.idJadwal, tc.token, tc.waktu)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCreate, tc.repo.createCalled)
		})
	}
}
