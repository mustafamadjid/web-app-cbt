export const routes = {
  login: "/login",
  adminDashboard: "/dashboard/administrator",
  guruDashboard: "/dashboard/guru",
  siswaDashboard: "/dashboard/siswa",

  adminGuru: "/dashboard/administrator/kelola-akun/guru",
  adminSiswa: "/dashboard/administrator/kelola-akun/siswa",
  kelolaSesi: "/dashboard/administrator/kelola-sesi",
  tambahGuru: "/dashboard/administrator/kelola-akun/tambah-guru",
  tambahSiswa: "/dashboard/administrator/kelola-akun/tambah-siswa",

  mapel: "/dashboard/administrator/data-master/mapel",
  kelas: "/dashboard/administrator/data-master/kelas",
  ruang: "/dashboard/administrator/data-master/ruang",
  sesi: "/dashboard/administrator/data-master/sesi",
  tambahMapel: "/dashboard/administrator/data-master/tambah-mapel",
  tambahKelas: "/dashboard/administrator/data-master/tambah-kelas",
  tambahRuang: "/dashboard/administrator/data-master/tambah-ruang",
  tambahSesi: "/dashboard/administrator/data-master/tambah-sesi",

  bankSoalAdmin: "/dashboard/administrator/bank-soal",
  bankSoalGuru: "/dashboard/guru/bank-soal",
  buatBankSoalAdmin: "/dashboard/administrator/bank-soal/buat",
  buatBankSoalGuru: "/dashboard/guru/bank-soal/buat",

  jadwalUjianAdmin: "/dashboard/administrator/ujian/jadwal",
  jadwalUjianGuru: "/dashboard/guru/ujian/jadwal",
  buatUjianAdmin: "/dashboard/administrator/ujian/buat-ujian",
  buatUjianGuru: "/dashboard/guru/ujian/buat-ujian",
  hasilUjianAdmin: "/dashboard/administrator/ujian/hasil",
  hasilUjianGuru: "/dashboard/guru/ujian/hasil",

  pengumumanAdmin: "/dashboard/administrator/pengumuman",
  pengumumanGuru: "/dashboard/guru/pengumuman",
  pengumumanSiswa: "/dashboard/siswa/pengumuman",

  cetakAdmin: "/dashboard/administrator/cetak",
  cetakGuru: "/dashboard/guru/cetak",
  pengaturan: "/dashboard/administrator/pengaturan",

  ujianSiswa: "/dashboard/siswa/ujian",
  hasilUjianSiswa: "/dashboard/siswa/ujian/hasil",
};

export function url(path) {
  return path;
}
