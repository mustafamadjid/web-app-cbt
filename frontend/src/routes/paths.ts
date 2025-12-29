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
    tambah_siswa:join("kelola-akun", "/tambah-siswa")
  },
} as const;