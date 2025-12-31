export type JadwalUjianItem = {
    id : number;
    nama_ujian : string;
    pengawas_ujian : string;
    tgl_ujian : string; // Data tangal sudah diolah duluan di server 
    waktu_mulai: string; // Data waktu sudah diolah duluan di server
    sesi_ujian?: number ;
    ruang_ujian?: string;

}

