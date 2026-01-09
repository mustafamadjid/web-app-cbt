export type UjianBerlangsungItem = {
  id: number;
  nama_ujian: string;
  mata_pelajaran: string;
  kelas: string[]; // Bisa multi kelas, misal ["X-IPA-1", "X-IPA-2"]
  waktu_mulai: string; // Format jam, misal "08:00"
  waktu_selesai: string; // Format jam, misal "10:00"
  total_siswa: number;
  siswa_mengerjakan: number; // Sedang login/mengerjakan
  siswa_selesai: number; // Sudah submit
};
