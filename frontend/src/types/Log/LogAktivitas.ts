export type UserRole = "admin" | "guru" | "siswa" | "unknown";

export type LogAktivitasApiItem = {
  id_aktivitas: string;
  id_pengguna: number;
  action: string;
  description: string;
  ip_address: string;
  created_at: string;
};

export type AktivitasLogItem = {
  id: string | number;
  username: string;
  role?: UserRole;
  aksi: string;
  deskripsi: string; 
  waktu?: string; 
};
