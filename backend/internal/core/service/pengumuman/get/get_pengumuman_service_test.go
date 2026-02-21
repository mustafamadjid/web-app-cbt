package pengumuman_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/get"
)

type fakeGetPengumumanRepo struct {
	getActiveFn     func(context.Context) ([]pengumuman.Pengumuman, error)
	getNonActiveFn  func(context.Context) ([]pengumuman.Pengumuman, error)
	getIncomingFn   func(context.Context) ([]pengumuman.Pengumuman, error)
	getByIDFn       func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error)
	getByIDCalled   bool
	gotIDPengumuman pengumuman.ID
}

func (f *fakeGetPengumumanRepo) GetPengumumanActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	if f.getActiveFn != nil {
		return f.getActiveFn(ctx)
	}
	return nil, nil
}

func (f *fakeGetPengumumanRepo) GetPengumumanNonActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	if f.getNonActiveFn != nil {
		return f.getNonActiveFn(ctx)
	}
	return nil, nil
}

func (f *fakeGetPengumumanRepo) GetPengumumanIncoming(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	if f.getIncomingFn != nil {
		return f.getIncomingFn(ctx)
	}
	return nil, nil
}

func (f *fakeGetPengumumanRepo) GetPengumumanById(ctx context.Context, idPengumuman pengumuman.ID) (pengumuman.Pengumuman, error) {
	f.getByIDCalled = true
	f.gotIDPengumuman = idPengumuman
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, idPengumuman)
	}
	return pengumuman.Pengumuman{}, nil
}

func (f *fakeGetPengumumanRepo) CreatePengumuman(context.Context, pengumuman.Pengumuman) error {
	return nil
}

func (f *fakeGetPengumumanRepo) UpdatePengumuman(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error {
	return nil
}

func (f *fakeGetPengumumanRepo) DeletePengumuman(context.Context, pengumuman.ID) error {
	return nil
}

func TestGetPengumumanActiveService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := []pengumuman.Pengumuman{{IdPengumuman: 1, JudulPengumuman: "A"}}

	tests := []struct {
		name      string
		repo      *fakeGetPengumumanRepo
		want      []pengumuman.Pengumuman
		expectErr error
	}{
		{
			name:      "branch 1 -> gagal get pengumuman aktif",
			repo:      &fakeGetPengumumanRepo{getActiveFn: func(context.Context) ([]pengumuman.Pengumuman, error) { return nil, repoErr }},
			expectErr: repoErr,
		},
		{
			name:      "happy path -> berhasil get pengumuman aktif",
			repo:      &fakeGetPengumumanRepo{getActiveFn: func(context.Context) ([]pengumuman.Pengumuman, error) { return expected, nil }},
			want:      expected,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewGetPengumumanService(tt.repo)
			got, err := svc.GetPengumumanActiveService(ctx)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPengumumanNonActiveService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := []pengumuman.Pengumuman{{IdPengumuman: 2, JudulPengumuman: "B"}}

	tests := []struct {
		name      string
		repo      *fakeGetPengumumanRepo
		want      []pengumuman.Pengumuman
		expectErr error
	}{
		{
			name:      "branch 1 -> gagal get pengumuman non aktif",
			repo:      &fakeGetPengumumanRepo{getNonActiveFn: func(context.Context) ([]pengumuman.Pengumuman, error) { return nil, repoErr }},
			expectErr: repoErr,
		},
		{
			name:      "happy path -> berhasil get pengumuman non aktif",
			repo:      &fakeGetPengumumanRepo{getNonActiveFn: func(context.Context) ([]pengumuman.Pengumuman, error) { return expected, nil }},
			want:      expected,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewGetPengumumanService(tt.repo)
			got, err := svc.GetPengumumanNonActiveService(ctx)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPengumumanIncomingService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := []pengumuman.Pengumuman{{IdPengumuman: 3, JudulPengumuman: "C"}}

	tests := []struct {
		name      string
		repo      *fakeGetPengumumanRepo
		want      []pengumuman.Pengumuman
		expectErr error
	}{
		{
			name:      "branch 1 -> gagal get pengumuman incoming",
			repo:      &fakeGetPengumumanRepo{getIncomingFn: func(context.Context) ([]pengumuman.Pengumuman, error) { return nil, repoErr }},
			expectErr: repoErr,
		},
		{
			name:      "happy path -> berhasil get pengumuman incoming",
			repo:      &fakeGetPengumumanRepo{getIncomingFn: func(context.Context) ([]pengumuman.Pengumuman, error) { return expected, nil }},
			want:      expected,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewGetPengumumanService(tt.repo)
			got, err := svc.GetPengumumanIncomingService(ctx)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPengumumanByIdService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := pengumuman.Pengumuman{IdPengumuman: 5, JudulPengumuman: "Item"}

	tests := []struct {
		name           string
		idPengumuman   pengumuman.ID
		repo           *fakeGetPengumumanRepo
		want           pengumuman.Pengumuman
		expectErr      error
		wantRepoCalled bool
	}{
		{
			name:           "branch 1 -> id pengumuman tidak valid",
			idPengumuman:   0,
			repo:           &fakeGetPengumumanRepo{},
			expectErr:      coreerror.ErrMissingId,
			wantRepoCalled: false,
		},
		{
			name:         "branch 2 -> repo error",
			idPengumuman: 5,
			repo: &fakeGetPengumumanRepo{
				getByIDFn: func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
					return pengumuman.Pengumuman{}, repoErr
				},
			},
			expectErr:      repoErr,
			wantRepoCalled: true,
		},
		{
			name:         "happy path -> berhasil get by id",
			idPengumuman: 5,
			repo: &fakeGetPengumumanRepo{
				getByIDFn: func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
					return expected, nil
				},
			},
			want:           expected,
			expectErr:      nil,
			wantRepoCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewGetPengumumanService(tt.repo)
			got, err := svc.GetPengumumanByIdService(ctx, tt.idPengumuman)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantRepoCalled, tt.repo.getByIDCalled)
			if tt.wantRepoCalled {
				assert.Equal(t, tt.idPengumuman, tt.repo.gotIDPengumuman)
			}
		})
	}
}
