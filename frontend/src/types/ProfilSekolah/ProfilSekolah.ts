export type ProfilSekolahFormValues = {
  nama_sekolah: string;
  alamat_sekolah: string;
  no_telp_sekolah: string;
  email_sekolah: string;
  kepala_sekolah: string;
  waka_sekolah: string;
  logo_sekolah: File | null;
};

export type ProfilSekolahUpdatePayload = {
  email_sekolah: string;
  no_telp_sekolah: string;
  kepala_sekolah: string;
  waka_sekolah: string;
  nama_sekolah: string;
  alamat_sekolah: string;
};

export type ProfilSekolahResponse = ProfilSekolahUpdatePayload & {
  id_profil: number;
  logo_sekolah: string | null;
  created_at: string;
  updated_at: string;
};
