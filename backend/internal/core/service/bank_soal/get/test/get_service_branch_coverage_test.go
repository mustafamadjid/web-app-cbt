package bank_soal_service_test

import (
	"context"
	"testing"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBankSoalService_FilterBranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		input      query.BankSoalFilter
		wantSearch string
		wantLimit  int
		wantOffset int
	}{
		{
			name:       "branch 1 -> limit <= 0 dan offset negatif",
			input:      query.BankSoalFilter{Search: "  matematika  ", Limit: 0, Offset: -1},
			wantSearch: "matematika",
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "branch 2 -> limit di atas 50",
			input:      query.BankSoalFilter{Search: "  fisika  ", Limit: 99, Offset: 10},
			wantSearch: "fisika",
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "branch 3 -> limit dan offset valid dipertahankan",
			input:      query.BankSoalFilter{Search: "  kimia  ", Limit: 25, Offset: 5},
			wantSearch: "kimia",
			wantLimit:  25,
			wantOffset: 5,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &FakeBankSoalRepo{}
			svc := bank_soal_service.NewGetBankSoalService(repo)

			_, err := svc.GetBankSoalService(ctx, tt.input)
			require.NoError(t, err)
			assert.True(t, repo.GetCalled)
			assert.Equal(t, tt.wantSearch, repo.GotFilter.Search)
			assert.Equal(t, tt.wantLimit, repo.GotFilter.Limit)
			assert.Equal(t, tt.wantOffset, repo.GotFilter.Offset)
		})
	}
}

func TestGetBankSoalUploadedService_FilterBranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		input      query.BankSoalFilter
		wantSearch string
		wantLimit  int
		wantOffset int
	}{
		{
			name:       "branch 1 -> limit <= 0 dan offset negatif",
			input:      query.BankSoalFilter{Search: "  matematika  ", Limit: 0, Offset: -1},
			wantSearch: "matematika",
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "branch 2 -> limit di atas 50",
			input:      query.BankSoalFilter{Search: "  fisika  ", Limit: 99, Offset: 10},
			wantSearch: "fisika",
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "branch 3 -> limit dan offset valid dipertahankan",
			input:      query.BankSoalFilter{Search: "  kimia  ", Limit: 25, Offset: 5},
			wantSearch: "kimia",
			wantLimit:  25,
			wantOffset: 5,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &FakeBankSoalRepo{}
			svc := bank_soal_service.NewGetBankSoalService(repo)

			_, err := svc.GetBankSoalUploadedService(ctx, tt.input)
			require.NoError(t, err)
			assert.True(t, repo.GetUploadedCalled)
			assert.Equal(t, tt.wantSearch, repo.GotUploadedFilter.Search)
			assert.Equal(t, tt.wantLimit, repo.GotUploadedFilter.Limit)
			assert.Equal(t, tt.wantOffset, repo.GotUploadedFilter.Offset)
		})
	}
}
