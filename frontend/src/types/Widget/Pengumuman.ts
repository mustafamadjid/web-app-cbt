export type PengumumanGetResponse = {
  id_pengumuman: number;
  id_pengguna: number;
  judul_pengumuman: string;
  isi_pengumuman: string;
  tanggal_rilis_pengumuman: string;
  tanggal_selesai_pengumuman: string;
  dokumen_pengumuman?: string;
};

export type PengumumanFormValues = {
  judul_pengumuman: string;
  isi_pengumuman: string;
  tanggal_rilis_pengumuman: string;
  tanggal_selesai_pengumuman: string;
  dokumen_pengumuman: File | null;
};

export type PengumumanCreatePayload = PengumumanFormValues;

export type PengumumanUpdatePayload = Partial<{
  judul_pengumuman: string;
  isi_pengumuman: string;
  tanggal_rilis_pengumuman: string;
  tanggal_selesai_pengumuman: string;
  dokumen_pengumuman: File | null;
}>;

export type PengumumanStatusKey = "incoming" | "active" | "non-active";
