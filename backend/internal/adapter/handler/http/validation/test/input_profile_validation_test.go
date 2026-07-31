package tests

import (
	"testing"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidatePersonName(t *testing.T) {
	tests := []struct {
		name, value, want string
		wantErr           error
	}{
		{name: "Path 1 -> valid person name is trimmed", value: "  Siti Aisyah, S.Pd.  ", want: "Siti Aisyah, S.Pd."},
		{name: "Path 2 -> markup in person name is rejected", value: "<script>alert(1)</script>", wantErr: httpx.ErrInvalidFormat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpx.ValidatePersonName(tt.value, "nama_lengkap")
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateGenderLabel(t *testing.T) {
	tests := []struct {
		name, value, want string
		wantErr           error
	}{
		{name: "Path 1 -> male short label is normalized", value: "L", want: "LAKI_LAKI"},
		{name: "Path 2 -> male hyphenated label is normalized", value: "Laki-laki", want: "LAKI_LAKI"},
		{name: "Path 3 -> female synonym is normalized", value: "wanita", want: "PEREMPUAN"},
		{name: "Path 4 -> invalid pattern fails before normalization", value: "<female>", wantErr: httpx.ErrInvalidFormat},
		{name: "Path 5 -> unknown valid label is rejected", value: "unknown", wantErr: httpx.ErrInvalidFormat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpx.ValidateGenderLabel(tt.value, "jenis_kelamin")
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProfileTextValidators(t *testing.T) {
	tests := []struct {
		name, kind, value, want string
		wantErr                 error
	}{
		{name: "Path 1 -> safe label is accepted", kind: "label", value: "Matematika & IPA", want: "Matematika & IPA"},
		{name: "Path 2 -> markup label is rejected", kind: "label", value: "<b>Bandung</b>", wantErr: httpx.ErrInvalidFormat},
		{name: "Path 3 -> formatted phone is accepted", kind: "phone", value: " +62 812-3456-7890 ", want: "+62 812-3456-7890"},
		{name: "Path 4 -> alphabetic phone is rejected", kind: "phone", value: "telephone", wantErr: httpx.ErrInvalidFormat},
		{name: "Path 5 -> 18 digit NIP is accepted", kind: "nip", value: "123456789012345678", want: "123456789012345678"},
		{name: "Path 6 -> dash sentinel NIP is accepted", kind: "nip", value: "-", want: "-"},
		{name: "Path 7 -> invalid NIP is rejected", kind: "nip", value: "123", wantErr: httpx.ErrInvalidFormat},
		{name: "Path 8 -> 10 digit NISN is accepted", kind: "nisn", value: "1234567890", want: "1234567890"},
		{name: "Path 9 -> dash sentinel NISN is accepted", kind: "nisn", value: "-", want: "-"},
		{name: "Path 10 -> invalid NISN is rejected", kind: "nisn", value: "123", wantErr: httpx.ErrInvalidFormat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var err error
			switch tt.kind {
			case "label":
				got, err = httpx.ValidateSafeLabelText(tt.value, "field")
			case "phone":
				got, err = httpx.ValidatePhoneNumber(tt.value, "field")
			case "nip":
				got, err = httpx.ValidateNIPText(tt.value, "field")
			case "nisn":
				got, err = httpx.ValidateNISNText(tt.value, "field")
			}
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
