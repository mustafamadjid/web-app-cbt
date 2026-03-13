package fake_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type FakeAttemptUjianRepo struct {
	GetAttemptUjianByIdErr error
	GetAttemptUjianByIdRet ujian.AttemptUjian

	GetAttemptUjianByIdCalled bool
	GotGetAttemptID           ujian.ID
}

func (f *FakeAttemptUjianRepo) GetAttemptUjianById(_ context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error) {
	f.GetAttemptUjianByIdCalled = true
	f.GotGetAttemptID = idAttempt
	return f.GetAttemptUjianByIdRet, f.GetAttemptUjianByIdErr
}
