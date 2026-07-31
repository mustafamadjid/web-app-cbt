package bcrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cryptobcrypt "golang.org/x/crypto/bcrypt"
)

func TestNewHasher(t *testing.T) {
	tests := []struct {
		name, scenario string
		cost           int
		wantCost       int
	}{
		{name: "Branch 1 -> zero cost selects bcrypt default", cost: 0, wantCost: cryptobcrypt.DefaultCost},
		{name: "Branch 2 -> explicit cost is preserved", cost: cryptobcrypt.MinCost, wantCost: cryptobcrypt.MinCost},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewHasher(tt.cost)
			require.NotNil(t, hasher)
			assert.Equal(t, tt.wantCost, hasher.cost)
		})
	}
}

func TestHasher_GenerateHash(t *testing.T) {
	tests := []struct {
		name, plain string
		cost        int
		wantErr     bool
	}{
		{name: "Branch 1 -> valid cost produces verifiable hash", plain: "Password-123", cost: cryptobcrypt.MinCost},
		{name: "Branch 2 -> cost above bcrypt maximum returns error", plain: "Password-123", cost: cryptobcrypt.MaxCost + 1, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewHasher(tt.cost).GenerateHash(tt.plain)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
				return
			}
			require.NoError(t, err)
			assert.NotEqual(t, tt.plain, hash)
			assert.NoError(t, cryptobcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.plain)))
		})
	}
}

func TestHasher_ComparePaswordAndHashed(t *testing.T) {
	hash, err := NewHasher(cryptobcrypt.MinCost).GenerateHash("Password-123")
	require.NoError(t, err)
	tests := []struct {
		name, hash, plain string
		want              bool
	}{
		{name: "Branch 1 -> matching password returns true", hash: hash, plain: "Password-123", want: true},
		{name: "Branch 2 -> different password returns false", hash: hash, plain: "wrong-password"},
		{name: "Branch 3 -> malformed hash returns false", hash: "not-a-bcrypt-hash", plain: "Password-123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewHasher(cryptobcrypt.MinCost).ComparePaswordAndHashed(tt.hash, tt.plain))
		})
	}
}
