package attemptujian_service_test

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type FakeAttemptUjianRepo struct {
	UpdateAttemptUjianErr error

	UpdateAttemptUjianCalled bool
	GotUpdateAttemptID       ujian.ID
	GotUpdatePatch           updatepatch.UpdateAttemptUjianPatch
}

func (f *FakeAttemptUjianRepo) UpdateAttemptUjian(_ context.Context, idAttempt ujian.ID, patch updatepatch.UpdateAttemptUjianPatch) error {
	f.UpdateAttemptUjianCalled = true
	f.GotUpdateAttemptID = idAttempt
	f.GotUpdatePatch = patch
	return f.UpdateAttemptUjianErr
}
