package user_service_test

import "context"

type FakeDeleteFileRepo struct {
	DeleteErr error

	DeleteCalled bool
	LastPath     string
}

func (r *FakeDeleteFileRepo) DeleteFile(ctx context.Context, filePath string) error {
	r.DeleteCalled = true
	r.LastPath = filePath
	if r.DeleteErr != nil {
		return r.DeleteErr
	}
	return nil
}
