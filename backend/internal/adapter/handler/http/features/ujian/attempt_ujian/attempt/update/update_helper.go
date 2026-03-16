package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func parseUpdateAttemptUjianID(ps httprouter.Params) (ujian.ID, error) {
	rawIDAttempt := strings.TrimSpace(ps.ByName("idAttempt"))
	if rawIDAttempt == "" {
		return 0, errors.New("id attempt is required")
	}

	idAttempt, err := strconv.Atoi(rawIDAttempt)
	if err != nil || idAttempt <= 0 {
		return 0, errors.New("id attempt must be a positive number")
	}

	return ujian.ID(idAttempt), nil
}

func toUpdateAttemptUjianPatch(req UpdateAttemptUjianRequest) updatepatch.UpdateAttemptUjianPatch {
	return updatepatch.UpdateAttemptUjianPatch{
		StatusAttempt: toUpdateAttemptStatusPointer(req.StatusAttempt),
		WaktuSubmit:   req.WaktuSubmit,
	}
}

func toUpdateAttemptStatusPointer(value *string) *ujian.StatusAttempt {
	if value == nil {
		return nil
	}

	status := ujian.StatusAttempt(*value)
	return &status
}
