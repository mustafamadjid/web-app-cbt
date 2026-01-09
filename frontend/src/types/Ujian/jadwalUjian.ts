export type JadwalUjianItem = {
    id: number;
    nama_ujian : string;
    pengawas_ujian : string;
    tgl_ujian : string; // Data tangal sudah diolah duluan di server 
    waktu_mulai: string; // Data waktu sudah diolah duluan di server
    sesi_ujian?: number ;
    ruang_ujian?: string;
    id_ruang?: number;
    status_ujian?: string;
    tingkat_kelas?: number;
    tingkat_kelas_id?: number;
    nama_kelas?: string;
}

export type JadwalUjianFilterParams = {
  q?: string;
  tanggal?: string; // "YYYY-MM-DD"
  tingkatKelasId?: number;
  ruangUjianId?: number;
};
