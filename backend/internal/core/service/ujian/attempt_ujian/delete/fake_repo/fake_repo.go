package fake_repo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

type FakeAttemptUjianRepo struct {
	DeleteAttemptUjianErr error

	DeleteAttemptUjianCalled bool
	GotDeleteAttemptID       ujian.ID
}

func (f *FakeAttemptUjianRepo) DeleteAttemptUjian(_ context.Context, idAttempt ujian.ID) error {
	f.DeleteAttemptUjianCalled = true
	f.GotDeleteAttemptID = idAttempt
	return f.DeleteAttemptUjianErr
}
