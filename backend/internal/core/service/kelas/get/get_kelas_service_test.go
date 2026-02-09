package kelas_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get/fake_test"
)

func TestGetKelasService_GetFullKelas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	items := []kelas.FullKelasData{{}}
	repoErr := errors.New("repo error")
	tingkat := 10

	tests := []struct {
		name       string
		filter     query.ListKelasFilter
		repo       *fake_test.FakeKelasRepo
		wantItems  []kelas.FullKelasData
		wantErr    error
		wantCalled bool
		wantFilter query.ListKelasFilter
	}{
		{
			name: "Branch 1 -> trim search, limit default, offset default",
			filter: query.ListKelasFilter{
				Search:       "  kelas ",
				Limit:        0,
				Offset:       -1,
				TingkatKelas: &tingkat,
			},
			repo:       &fake_test.FakeKelasRepo{Items: items},
			wantItems:  items,
			wantCalled: true,
			wantFilter: query.ListKelasFilter{
				Search:       "kelas",
				Limit:        20,
				Offset:       0,
				TingkatKelas: &tingkat,
			},
		},
		{
			name: "Branch 2 -> limit capped",
			filter: query.ListKelasFilter{
				Limit:  100,
				Offset: 2,
			},
			repo:       &fake_test.FakeKelasRepo{Items: items},
			wantItems:  items,
			wantCalled: true,
			wantFilter: query.ListKelasFilter{
				Limit:  50,
				Offset: 2,
			},
		},
		{
			name: "Branch 3 -> repo error",
			filter: query.ListKelasFilter{
				Limit: 10,
			},
			repo:       &fake_test.FakeKelasRepo{Err: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
			wantFilter: query.ListKelasFilter{
				Limit:  10,
				Offset: 0,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := kelas_service.NewGetKelasService(tt.repo)
			result, err := svc.GetFullKelas(ctx, tt.filter)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantItems, result)
			}

			assert.Equal(t, tt.wantCalled, tt.repo.Called)
			if tt.wantCalled {
				assert.Equal(t, tt.wantFilter, tt.repo.GotFilter)
			}
		})
	}
}
