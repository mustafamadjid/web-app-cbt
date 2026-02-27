package sesi_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestGetSesiService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("get sesi error")
	expectedItems := []sesi.Sesi{
		{IdSesi: 1, KodeSesi: "SESI01", NamaSesi: "Sesi 1"},
	}

	tests := []struct {
		name      string
		repo      *fake_test.FakeSesiRepo
		filter    query.ListSesiFilter
		wantErr   error
		wantItems []sesi.Sesi
	}{
		{
			name:    "Path 1 -> GetSesi repo error",
			repo:    &fake_test.FakeSesiRepo{GetSesiErr: repoErr},
			filter:  query.ListSesiFilter{Limit: 10},
			wantErr: repoErr,
		},
		{
			name:      "Path 2 -> happy path dengan limit/offset clamped",
			repo:      &fake_test.FakeSesiRepo{GetSesiRet: expectedItems},
			filter:    query.ListSesiFilter{Limit: -1, Offset: -5},
			wantErr:   nil,
			wantItems: expectedItems,
		},
		{
			name:      "Path 3 -> limit > 50 di-clamp ke 50",
			repo:      &fake_test.FakeSesiRepo{GetSesiRet: expectedItems},
			filter:    query.ListSesiFilter{Limit: 100},
			wantErr:   nil,
			wantItems: expectedItems,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := sesi_service.NewGetSesiService(tt.repo)
			items, err := svc.GetSesiService(ctx, tt.filter)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantItems, items)
			}
		})
	}
}

func TestGetSesiByIdService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("get sesi by id error")
	expectedItem := sesi.Sesi{IdSesi: 1, KodeSesi: "SESI01", NamaSesi: "Sesi 1"}

	tests := []struct {
		name     string
		idSesi   int
		repo     *fake_test.FakeSesiRepo
		wantErr  error
		wantItem sesi.Sesi
	}{
		{
			name:    "Path 1 -> idSesi <= 0",
			idSesi:  0,
			repo:    &fake_test.FakeSesiRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:    "Path 2 -> GetSesiById error",
			idSesi:  1,
			repo:    &fake_test.FakeSesiRepo{GetSesiByIdErr: repoErr},
			wantErr: repoErr,
		},
		{
			name:     "Path 3 -> happy path",
			idSesi:   1,
			repo:     &fake_test.FakeSesiRepo{GetSesiByIdRet: expectedItem},
			wantErr:  nil,
			wantItem: expectedItem,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := sesi_service.NewGetSesiService(tt.repo)
			item, err := svc.GetSesiByIdService(ctx, tt.idSesi)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantItem, item)
			}
		})
	}
}

func TestGetSesiByKodeService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("get sesi by kode error")
	expectedItem := sesi.Sesi{IdSesi: 1, KodeSesi: "SESI01", NamaSesi: "Sesi 1"}

	tests := []struct {
		name     string
		kodeSesi string
		repo     *fake_test.FakeSesiRepo
		wantErr  error
		wantItem sesi.Sesi
	}{
		{
			name:     "Path 1 -> kodeSesi kosong",
			kodeSesi: "   ",
			repo:     &fake_test.FakeSesiRepo{},
			wantErr:  coreerror.ErrMissingField,
		},
		{
			name:     "Path 2 -> GetSesiByKode error",
			kodeSesi: "SESI01",
			repo:     &fake_test.FakeSesiRepo{GetSesiByKodeErr: repoErr},
			wantErr:  repoErr,
		},
		{
			name:     "Path 3 -> happy path",
			kodeSesi: "  sesi01  ",
			repo:     &fake_test.FakeSesiRepo{GetSesiByKodeRet: expectedItem},
			wantErr:  nil,
			wantItem: expectedItem,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := sesi_service.NewGetSesiService(tt.repo)
			item, err := svc.GetSesiByKodeService(ctx, tt.kodeSesi)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantItem, item)
			}
		})
	}
}
