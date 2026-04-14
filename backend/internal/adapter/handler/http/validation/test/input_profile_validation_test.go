package tests

import (
	"testing"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePersonName(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidatePersonName("  Siti Aisyah, S.Pd.  ", "nama_lengkap")
	require.NoError(t, err)
	assert.Equal(t, "Siti Aisyah, S.Pd.", value)

	_, err = httpx.ValidatePersonName("<script>alert(1)</script>", "nama_lengkap")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)
}

func TestValidateGenderLabel(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidateGenderLabel("Laki-laki", "jenis_kelamin")
	require.NoError(t, err)
	assert.Equal(t, "LAKI_LAKI", value)

	value, err = httpx.ValidateGenderLabel("LAKI_LAKI", "jenis_kelamin")
	require.NoError(t, err)
	assert.Equal(t, "LAKI_LAKI", value)

	_, err = httpx.ValidateGenderLabel("unknown", "jenis_kelamin")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)
}

func TestValidatePhoneNumber(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidatePhoneNumber(" +62 812-3456-7890 ", "no_hp")
	require.NoError(t, err)
	assert.Equal(t, "+62 812-3456-7890", value)

	_, err = httpx.ValidatePhoneNumber("<script>alert(1)</script>", "no_hp")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)
}

func TestValidateSafeLabelText(t *testing.T) {
	t.Parallel()

	value, err := httpx.ValidateSafeLabelText("Matematika & IPA", "bidang_studi")
	require.NoError(t, err)
	assert.Equal(t, "Matematika & IPA", value)

	_, err = httpx.ValidateSafeLabelText("<b>Bandung</b>", "tempat_lahir")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)
}

func TestValidateNIPAndNISNText(t *testing.T) {
	t.Parallel()

	nip, err := httpx.ValidateNIPText("123456789012345678", "nip")
	require.NoError(t, err)
	assert.Equal(t, "123456789012345678", nip)

	_, err = httpx.ValidateNIPText("123<script>", "nip")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)

	nisn, err := httpx.ValidateNISNText("1234567890", "nisn")
	require.NoError(t, err)
	assert.Equal(t, "1234567890", nisn)

	_, err = httpx.ValidateNISNText("123<script>", "nisn")
	assert.ErrorIs(t, err, httpx.ErrInvalidFormat)
}
