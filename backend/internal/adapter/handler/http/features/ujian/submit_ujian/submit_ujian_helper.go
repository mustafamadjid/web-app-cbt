package httpx

import (
	"errors"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func parseSubmitUjianRequest(ps httprouter.Params) (SubmitUjianRequest, error) {
	rawIDAttempt := strings.TrimSpace(ps.ByName("idAttempt"))
	if rawIDAttempt == "" {
		return SubmitUjianRequest{}, errors.New("id attempt is required")
	}

	idAttempt, err := strconv.Atoi(rawIDAttempt)
	if err != nil {
		return SubmitUjianRequest{}, errors.New("id attempt must be a positive number")
	}

	req := SubmitUjianRequest{
		IDAttempt: idAttempt,
	}

	if err := ValidateSubmitUjianRequest(req); err != nil {
		return SubmitUjianRequest{}, err
	}

	return req, nil
}
