export type AnnouncementDoc = {
  id?: number;
  name: string;
  url: string;
  mimeType?: string;
  sizeLabel?: string;
};

export type PengumumanItem = {
  id: number;
  judul: string;
  isi_pengumuman: string;
  tanggal_rilis_pengumuman: string;
  dokumen?: AnnouncementDoc | AnnouncementDoc[] | null;
};
