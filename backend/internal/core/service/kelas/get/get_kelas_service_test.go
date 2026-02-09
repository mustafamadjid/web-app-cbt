package kelas_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
)

type fakeKelasRepo struct {
	items     []kelas.FullKelasData
	err       error
	called    bool
	gotFilter query.ListKelasFilter
}

func (f *fakeKelasRepo) GetKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	f.called = true
	f.gotFilter = filter
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

// Not used
func (f *fakeKelasRepo) CreateTingkatKelas(ctx context.Context, tingkatKelas kelas.TingkatKelas)(kelas.ID,error){
	panic("not used in this test")
}

// Not used
func (f *fakeKelasRepo) CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas)(kelas.ID,error){
	panic("not used in this test")
}

func TestGetKelasService_GetFullKelas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	items := []kelas.FullKelasData{{}}
	repoErr := errors.New("repo error")
	tingkat := 10

	tests := []struct {
		name       string
		filter     query.ListKelasFilter
		repo       *fakeKelasRepo
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
			repo:       &fakeKelasRepo{items: items},
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
			repo:       &fakeKelasRepo{items: items},
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
			repo:       &fakeKelasRepo{err: repoErr},
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

			assert.Equal(t, tt.wantCalled, tt.repo.called)
			if tt.wantCalled {
				assert.Equal(t, tt.wantFilter, tt.repo.gotFilter)
			}
		})
	}
}
