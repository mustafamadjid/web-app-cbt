package attemptujian_service_test

import (
	"context"
	"testing"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	"github.com/stretchr/testify/assert"
)

func TestExpireAttemptUjianService(t *testing.T) {
	t.Parallel()

	repo := &FakeAttemptUjianRepo{}
	updater := attemptujian_service.NewUpdateAttemptUjianService(repo)
	svc := attemptujian_service.NewExpireAttemptUjianService(updater)

	err := svc.ExpireAttemptUjian(context.Background(), 14)

	assert.NoError(t, err)
	assert.True(t, repo.UpdateAttemptUjianCalled)
	assert.Equal(t, ujian.ID(14), repo.GotUpdateAttemptID)
	assert.Nil(t, repo.GotUpdatePatch.WaktuSubmit)
	if assert.NotNil(t, repo.GotUpdatePatch.StatusAttempt) {
		assert.Equal(t, ujian.ATTEMPT_EXPIRED, *repo.GotUpdatePatch.StatusAttempt)
	}
}
