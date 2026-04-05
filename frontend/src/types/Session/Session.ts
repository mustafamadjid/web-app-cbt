import type { Role } from "@/types/Sidebar/SidebarMenu";

export type AccountStatus = "AKTIF" | "NONAKTIF";

export type ActiveSessionInfo = {
  session_id: string;
  id_pengguna: number;
  role: Role;
  revoked: boolean;
  expires_at: string;
};

export type ActiveSessionUser = {
  id_pengguna: number;
  username: string;
  email: string | null;
  nama_lengkap: string;
  jenis_kelamin: string;
  no_hp: string | null;
  role: Role;
  status_akun: AccountStatus;
  foto_profil: string;
};

export type ActiveSessionRow = {
  session: ActiveSessionInfo;
  pengguna: ActiveSessionUser;
};

export type ListActiveSessionsResponse = {
  items: ActiveSessionRow[];
};

export type AdminRevokeSessionPayload = {
  session_id: string;
};
