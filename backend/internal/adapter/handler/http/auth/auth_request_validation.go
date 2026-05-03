package httpx

import validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"

func sanitizeAndValidateLoginRequest(req LoginRequest) (LoginRequest, error) {
	username, err := validator.ValidateUsername(req.Username)
	if err != nil {
		return LoginRequest{}, err
	}

	password, err := validator.ValidatePassword(req.Password)
	if err != nil {
		return LoginRequest{}, err
	}

	req.Username = username
	req.Password = password
	return req, nil
}

func sanitizeAndValidateAdminRevokeRequest(req AdminRevokeRequest) (AdminRevokeRequest, error) {
	sessionID, err := validator.ValidateRequiredPrintableText(req.SessionId, "session_id")
	if err != nil {
		return AdminRevokeRequest{}, err
	}

	req.SessionId = sessionID
	return req, nil
}
