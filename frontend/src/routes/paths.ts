const join = (base: string, path: string) => `${base}${path.startsWith('/')? path : `/${path}`}`

export const paths = {
  public: {
    home: "/",
    login: "/login",
  },

  dashboard: {
    // home_admin: join("/dashboard", "/administrator"),
    kelola_akun_siswa: join("kelola-akun", "/siswa"),
    kelola_akun_guru: join("kelola-akun", "/guru"),
    tambah_guru:join("kelola-akun", "/tambah-guru"),
    tambah_siswa:join("kelola-akun", "/tambah-siswa"),

    data_master_mapel : join("data-master", "/mapel"),
    data_master_kelas : join("data-master", "/kelas"),
    data_master_ruang : join("data-master", "/ruang"),
    data_master_sesi : join("data-master", "/sesi"),
    tambah_data_master_mapel : join("data-master", "/tambah-mapel"),
    tambah_data_master_kelas : join("data-master", "/tambah-kelas"),
    tambah_data_master_ruang : join("data-master", "/tambah-ruang"),
    tambah_data_master_sesi : join("data-master", "/tambah-sesi"),

    pengaturan:"pengaturan",
  },
} as const;