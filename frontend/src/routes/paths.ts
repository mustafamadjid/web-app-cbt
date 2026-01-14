const join = (base: string, path: string) =>
  `${base}${path.startsWith("/") ? path : `/${path}`}`;

const DASHBOARD_ADMIN = "/dashboard/administrator";

const djoin = (path: string) => join(DASHBOARD_ADMIN, path);

export const paths = {
  public: {
    home: "/",
    login: "/login",
  },

  dashboard: {
    home_admin: DASHBOARD_ADMIN,

    kelola_akun_siswa: djoin(join("kelola-akun", "/siswa")),
    kelola_akun_guru: djoin(join("kelola-akun", "/guru")),
    tambah_guru: djoin(join("kelola-akun", "/tambah-guru")),
    tambah_siswa: djoin(join("kelola-akun", "/tambah-siswa")),

    data_master_mapel: djoin(join("data-master", "/mapel")),
    data_master_kelas: djoin(join("data-master", "/kelas")),
    data_master_ruang: djoin(join("data-master", "/ruang")),
    data_master_sesi: djoin(join("data-master", "/sesi")),

    tambah_data_master_mapel: djoin(join("data-master", "/tambah-mapel")),
    tambah_data_master_kelas: djoin(join("data-master", "/tambah-kelas")),
    tambah_data_master_ruang: djoin(join("data-master", "/tambah-ruang")),
    tambah_data_master_sesi: djoin(join("data-master", "/tambah-sesi")),

    bank_soal: djoin("bank-soal"),
    tambah_bank_soal: djoin(join("bank-soal", "/tambah")),

    jadwal_ujian: djoin(join("ujian", "/jadwal")),
    detail_ujian: djoin(join("ujian", "/jadwal/:id")),
    buat_ujian: djoin(join("ujian", "/buat-ujian")),
    hasil_ujian: djoin(join("ujian", "/hasil")),
    hasil_ujian_detail: djoin(join("ujian", "/hasil/:id")),
    hasil_ujian_siswa_detail: djoin(join("ujian", "/hasil/:id/siswa/:siswaId")),


    cetak: djoin("cetak"),

    pengaturan: djoin("pengaturan"),
  },
} as const;
