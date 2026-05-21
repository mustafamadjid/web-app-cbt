export function unique(prefix) {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
}

export const fixtures = {
  password: "Password123!",
  mapel: () => ({
    kode: unique("MAPEL"),
    nama: unique("Mata Pelajaran E2E"),
  }),
  kelas: () => ({
    tingkat: "X",
    nama: unique("Kelas E2E"),
  }),
  ruang: () => ({
    kode: unique("RUANG"),
    nama: unique("Ruang E2E"),
    kapasitas: "30",
  }),
  sesi: () => ({
    nama: unique("Sesi E2E"),
    mulai: "08:00",
    selesai: "09:00",
  }),
};
