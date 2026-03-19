package httpx

import "errors"

func ValidateSubmitUjianRequest(req SubmitUjianRequest) error {
	if req.IDAttempt <= 0 {
		return errors.New("id attempt must be a positive number")
	}

	return nil
}
