package httpx

import (
	"errors"
	"testing"
)

func TestValidatePersonName(t *testing.T) {
	value, err := ValidatePersonName("  Siti Aisyah, S.Pd.  ", "nama_lengkap")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "Siti Aisyah, S.Pd." {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidatePersonName_RejectsHTMLLikeInput(t *testing.T) {
	_, err := ValidatePersonName("<script>alert(1)</script>", "nama_lengkap")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestValidateGenderLabel(t *testing.T) {
	value, err := ValidateGenderLabel("Laki-laki", "jenis_kelamin")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "LAKI_LAKI" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidateGenderLabel_EnumValue(t *testing.T) {
	value, err := ValidateGenderLabel("LAKI_LAKI", "jenis_kelamin")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "LAKI_LAKI" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	value, err := ValidatePhoneNumber(" +62 812-3456-7890 ", "no_hp")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "+62 812-3456-7890" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidatePhoneNumber_RejectsHTMLLikeInput(t *testing.T) {
	_, err := ValidatePhoneNumber("<script>alert(1)</script>", "no_hp")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestValidateSafeLabelText(t *testing.T) {
	value, err := ValidateSafeLabelText("Matematika & IPA", "bidang_studi")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "Matematika & IPA" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidateSafeLabelText_RejectsHTMLLikeInput(t *testing.T) {
	_, err := ValidateSafeLabelText("<b>Bandung</b>", "tempat_lahir")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestValidateNIPText(t *testing.T) {
	value, err := ValidateNIPText("123456789012345678", "nip")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "123456789012345678" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidateNIPText_InvalidFormat(t *testing.T) {
	_, err := ValidateNIPText("123<script>", "nip")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestValidateNISNText(t *testing.T) {
	value, err := ValidateNISNText("1234567890", "nisn")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if value != "1234567890" {
		t.Fatalf("unexpected sanitized value: %q", value)
	}
}

func TestValidateNISNText_InvalidFormat(t *testing.T) {
	_, err := ValidateNISNText("123<script>", "nisn")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}
