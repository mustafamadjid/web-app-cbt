export type UserRole = "ADMIN" | "GURU" | "SISWA";

export type AktivitasLogItem = {
  id_aktivitas: string;
  id_pengguna: number;
  username: string;
  role: UserRole;
  action: string;
  description: string;
  ip_address: string;
  created_at: string;
};
