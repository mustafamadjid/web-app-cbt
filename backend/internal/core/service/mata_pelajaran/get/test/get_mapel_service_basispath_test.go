package matapelajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
	"github.com/stretchr/testify/assert"
)

func TestGetMapelService_BranchCoverage(t *testing.T) {
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
		repo      *FakeMapelRepo
		filter    query.ListMapelFilter
		wantErr   error
		wantItems []matapelajaran.MataPelajaran
	}{
		{
			name:    "Branch 1 -> nama mapel kosong",
			repo:    &FakeMapelRepo{},
			filter:  query.ListMapelFilter{NamaMapel: &emptyNama, Limit: 10},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:    "Branch 2 -> tingkat kelas tidak valid",
			repo:    &FakeMapelRepo{},
			filter:  query.ListMapelFilter{TingkatKelas: &invalidTingkat, Limit: 10},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:    "Branch 3 -> repo list mapel gagal",
			repo:    &FakeMapelRepo{GetMapelErr: repoErr},
			filter:  query.ListMapelFilter{Limit: 10},
			wantErr: repoErr,
		},
		{
			name:      "Branch 4 -> filter valid dengan limit offset dibersihkan",
			repo:      &FakeMapelRepo{GetMapelRet: expectedItems},
			filter:    query.ListMapelFilter{NamaMapel: &validNama, Limit: -1, Offset: -5},
			wantErr:   nil,
			wantItems: expectedItems,
		},
		{
			name:      "Branch 5 -> limit lebih dari 50 di-clamp",
			repo:      &FakeMapelRepo{GetMapelRet: expectedItems},
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

func TestGetMapelByIdService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	genericErr := errors.New("generic get by id error")
	expectedItem := matapelajaran.MataPelajaran{IdMapel: 1, NamaMapel: "Matematika"}

	tests := []struct {
		name     string
		idMapel  int
		repo     *FakeMapelRepo
		wantErr  error
		wantItem matapelajaran.MataPelajaran
	}{
		{
			name:    "Branch 1 -> id mapel tidak valid",
			idMapel: 0,
			repo:    &FakeMapelRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:    "Branch 2 -> mapel tidak ditemukan",
			idMapel: 1,
			repo:    &FakeMapelRepo{GetMapelByIdErr: coreerror.ErrNotFound},
			wantErr: coreerror.ErrNotFound,
		},
		{
			name:    "Branch 3 -> repo get by id gagal",
			idMapel: 1,
			repo:    &FakeMapelRepo{GetMapelByIdErr: genericErr},
			wantErr: genericErr,
		},
		{
			name:     "Branch 4 -> get by id berhasil",
			idMapel:  1,
			repo:     &FakeMapelRepo{GetMapelByIdRet: expectedItem},
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
