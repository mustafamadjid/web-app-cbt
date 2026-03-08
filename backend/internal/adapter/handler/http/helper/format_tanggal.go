package httpx

import "time"

func FormatTanggalIndonesia(t time.Time) string {
	hariNama := namaHariIndonesia(t.Weekday())
	bulanNama := namaBulanIndonesia(t.Month())
	if hariNama == "" || bulanNama == "" {
		return t.Format("2006-01-02")
	}

	return hariNama + ", " + t.Format("02") + " " + bulanNama + " " + t.Format("2006")
}

func FormatTanggalWaktuIndonesia(t time.Time) string {
	hariNama := namaHariIndonesia(t.Weekday())
	bulanNama := namaBulanIndonesia(t.Month())
	if hariNama == "" || bulanNama == "" {
		return t.Format("2006-01-02 15:04")
	}

	return hariNama + " " + t.Format("02") + " " + bulanNama + " " + t.Format("15.04")
}

func FormatDateOnly(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func FormatTimeOnly(t time.Time) string {
	return t.Format("15:04")
}

func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

func FormatRFC3339Ptr(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.Format(time.RFC3339)
	return &formatted
}

func namaHariIndonesia(day time.Weekday) string {
	switch day {
	case time.Sunday:
		return "Minggu"
	case time.Monday:
		return "Senin"
	case time.Tuesday:
		return "Selasa"
	case time.Wednesday:
		return "Rabu"
	case time.Thursday:
		return "Kamis"
	case time.Friday:
		return "Jumat"
	case time.Saturday:
		return "Sabtu"
	default:
		return ""
	}
}

func namaBulanIndonesia(month time.Month) string {
	switch month {
	case time.January:
		return "Januari"
	case time.February:
		return "Februari"
	case time.March:
		return "Maret"
	case time.April:
		return "April"
	case time.May:
		return "Mei"
	case time.June:
		return "Juni"
	case time.July:
		return "Juli"
	case time.August:
		return "Agustus"
	case time.September:
		return "September"
	case time.October:
		return "Oktober"
	case time.November:
		return "November"
	case time.December:
		return "Desember"
	default:
		return ""
	}
}
