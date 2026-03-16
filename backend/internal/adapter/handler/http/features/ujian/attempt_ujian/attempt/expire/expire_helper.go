package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func parseExpireAttemptUjianID(ps httprouter.Params) (ujian.ID, error) {
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
