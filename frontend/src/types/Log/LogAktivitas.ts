export type UserRole = "admin" | "guru" | "siswa";
export type AktivitasLogItem = {
  id: string | number;
  username: string;
  role: UserRole;
  aksi: string;
  deskripsi: string; 
  waktu?: string; 
};
