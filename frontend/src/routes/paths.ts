const join = (base: string, path: string) =>
  `${base}${path.startsWith("/") ? path : `/${path}`}`;

const DASHBOARD_ADMIN = "/dashboard/administrator";
const DASHBOARD_GURU = "/dashboard/guru";
const DASHBOARD_SISWA = "/dashboard/siswa";


const djoinAdmin = (path: string) => join(DASHBOARD_ADMIN, path);
const djoinGuru = (path: string) => join(DASHBOARD_GURU, path);
const djoinSiswa = (path: string) => join(DASHBOARD_SISWA, path);


export const paths = {
  public: {
    home: "/",
    login: "/login",
  },

  dashboard: {
    home_admin: DASHBOARD_ADMIN,
    home_guru: DASHBOARD_GURU,
    home_siswa: DASHBOARD_SISWA,

    kelola_akun_siswa: djoinAdmin(join("kelola-akun", "/siswa")),
    kelola_akun_guru: djoinAdmin(join("kelola-akun", "/guru")),
    tambah_guru: djoinAdmin(join("kelola-akun", "/tambah-guru")),
    tambah_siswa: djoinAdmin(join("kelola-akun", "/tambah-siswa")),

    data_master_mapel: djoinAdmin(join("data-master", "/mapel")),
    data_master_kelas: djoinAdmin(join("data-master", "/kelas")),
    data_master_ruang: djoinAdmin(join("data-master", "/ruang")),
    data_master_sesi: djoinAdmin(join("data-master", "/sesi")),

    tambah_data_master_mapel: djoinAdmin(join("data-master", "/tambah-mapel")),
    tambah_data_master_kelas: djoinAdmin(join("data-master", "/tambah-kelas")),
    tambah_data_master_ruang: djoinAdmin(join("data-master", "/tambah-ruang")),
    tambah_data_master_sesi: djoinAdmin(join("data-master", "/tambah-sesi")),

    bank_soal: djoinAdmin("bank-soal"),
    tambah_bank_soal: djoinAdmin(join("bank-soal", "/tambah")),

    jadwal_ujian: djoinAdmin(join("ujian", "/jadwal")),
    detail_ujian: djoinAdmin(join("ujian", "/jadwal/:id")),
    buat_ujian: djoinAdmin(join("ujian", "/buat-ujian")),
    hasil_ujian: djoinAdmin(join("ujian", "/hasil")),
    hasil_ujian_detail: djoinAdmin(join("ujian", "/hasil/:id")),
    hasil_ujian_siswa_detail: djoinAdmin(
      join("ujian", "/hasil/:id/siswa/:siswaId")
    ),

    bank_soal_guru: djoinGuru("bank-soal"),
    tambah_bank_soal_guru: djoinGuru(join("bank-soal", "/tambah")),

    jadwal_ujian_guru: djoinGuru(join("ujian", "/jadwal")),
    detail_ujian_guru: djoinGuru(join("ujian", "/jadwal/:id")),
    buat_ujian_guru: djoinGuru(join("ujian", "/buat-ujian")),
    hasil_ujian_guru: djoinGuru(join("ujian", "/hasil")),
    hasil_ujian_detail_guru: djoinGuru(join("ujian", "/hasil/:id")),
    hasil_ujian_siswa_detail_guru: djoinGuru(
      join("ujian", "/hasil/:id/siswa/:siswaId")
    ),

    cetak: djoinAdmin("cetak"),
    cetak_guru: djoinGuru("cetak"),


    pengaturan: djoinAdmin("pengaturan"),

    ujian_siswa : djoinSiswa("ujian"),
    hasil_ujian_siswa : djoinSiswa(join("ujian", "/hasil")),


  },
} as const;
