package httpx

type UpdatePengumumanRequest struct {
	JudulPengumuman          *string
	IsiPengumuman            *string
	TanggalRilisPengumuman   *string
	TanggalSelesaiPengumuman *string
}

func (r UpdatePengumumanRequest) IsEmpty() bool {
	return r.JudulPengumuman == nil &&
		r.IsiPengumuman == nil &&
		r.TanggalRilisPengumuman == nil &&
		r.TanggalSelesaiPengumuman == nil
}
