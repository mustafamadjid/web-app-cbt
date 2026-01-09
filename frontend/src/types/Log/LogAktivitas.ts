export type UserRole = "admin" | "guru" | "siswa";
export type AktivitasLogItem = {
  id: number;
  username: string;
  role: UserRole;
  aksi: string;
  deskripsi: string; 
  waktu?: string; 
};
