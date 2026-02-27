package matapelajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestGetMapelService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("get mapel error")
	emptyNama := ""
	validNama := "Matematika"
	invalidTingkat := 0
	expectedItems := []matapelajaran.MataPelajaran{
		{IdMapel: 1, NamaMapel: "Matematika"},
	}

	tests := []struct {
		name      string
		repo      *fake_test.FakeMapelRepo
		filter    query.ListMapelFilter
		wantErr   error
		wantItems []matapelajaran.MataPelajaran
	}{
		{
			name:    "Path 1 -> NamaMapel not nil tapi kosong",
			repo:    &fake_test.FakeMapelRepo{},
			filter:  query.ListMapelFilter{NamaMapel: &emptyNama, Limit: 10},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:    "Path 2 -> TingkatKelas not nil tapi <= 0",
			repo:    &fake_test.FakeMapelRepo{},
			filter:  query.ListMapelFilter{TingkatKelas: &invalidTingkat, Limit: 10},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:    "Path 3 -> GetMapel repo error",
			repo:    &fake_test.FakeMapelRepo{GetMapelErr: repoErr},
			filter:  query.ListMapelFilter{Limit: 10},
			wantErr: repoErr,
		},
		{
			name:      "Path 4 -> happy path dengan limit/offset clamped",
			repo:      &fake_test.FakeMapelRepo{GetMapelRet: expectedItems},
			filter:    query.ListMapelFilter{NamaMapel: &validNama, Limit: -1, Offset: -5},
			wantErr:   nil,
			wantItems: expectedItems,
		},
		{
			name:      "Path 5 -> limit > 50 di-clamp ke 50",
			repo:      &fake_test.FakeMapelRepo{GetMapelRet: expectedItems},
			filter:    query.ListMapelFilter{Limit: 100},
			wantErr:   nil,
			wantItems: expectedItems,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mapel_service.NewGetMapelService(tt.repo)
			items, err := svc.GetMapelService(ctx, tt.filter)

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

func TestGetMapelByIdService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	genericErr := errors.New("generic get by id error")
	expectedItem := matapelajaran.MataPelajaran{IdMapel: 1, NamaMapel: "Matematika"}

	tests := []struct {
		name     string
		idMapel  int
		repo     *fake_test.FakeMapelRepo
		wantErr  error
		wantItem matapelajaran.MataPelajaran
	}{
		{
			name:    "Path 1 -> idMapel <= 0",
			idMapel: 0,
			repo:    &fake_test.FakeMapelRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:    "Path 2 -> GetMapelById err ErrNotFound",
			idMapel: 1,
			repo:    &fake_test.FakeMapelRepo{GetMapelByIdErr: coreerror.ErrNotFound},
			wantErr: coreerror.ErrNotFound,
		},
		{
			name:    "Path 3 -> GetMapelById generic error",
			idMapel: 1,
			repo:    &fake_test.FakeMapelRepo{GetMapelByIdErr: genericErr},
			wantErr: genericErr,
		},
		{
			name:     "Path 4 -> happy path",
			idMapel:  1,
			repo:     &fake_test.FakeMapelRepo{GetMapelByIdRet: expectedItem},
			wantErr:  nil,
			wantItem: expectedItem,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mapel_service.NewGetMapelService(tt.repo)
			item, err := svc.GetMapelById(ctx, tt.idMapel)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantItem, item)
			}
		})
	}
}
