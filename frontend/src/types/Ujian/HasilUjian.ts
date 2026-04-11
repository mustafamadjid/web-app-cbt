import type { DataAkunSiswa } from "../KelolaAkun/AkunSiswa";

export type { StatistikUjian } from "./StatistikUjian";

export type HasilUjianSiswa = DataAkunSiswa & {
  nilai?: number;
  jumlah_benar?: number;
  jumlah_salah?: number;
  jumlah_kosong?: number;
};
