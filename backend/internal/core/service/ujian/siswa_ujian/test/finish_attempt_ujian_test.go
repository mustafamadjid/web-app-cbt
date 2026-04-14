package siswaujian_service_test

import (
	"testing"

	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian"
	"github.com/stretchr/testify/assert"
)

func TestFinishAttemptUjianService_ZeroValue(t *testing.T) {
	t.Parallel()

	var svc siswaujian_service.FinishAttemptUjianService
	assert.Equal(t, siswaujian_service.FinishAttemptUjianService{}, svc)
}
