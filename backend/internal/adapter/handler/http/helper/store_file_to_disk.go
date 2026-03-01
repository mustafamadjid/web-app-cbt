package httpx

import (
	"errors"
	"mime/multipart"
	"net/http"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

type SaveUploadedFileFn func(file multipart.File, fh *multipart.FileHeader) (string, error)

func StoreFileToDisk(r *http.Request, fieldName string, required bool, saver SaveUploadedFileFn) (*string, error) {
	file, fh, err := r.FormFile(fieldName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			if required {
				return nil, coreerror.ErrMissingField
			}
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	relativePath, err := saver(file, fh)
	if err != nil {
		return nil, err
	}

	return &relativePath, nil
}
