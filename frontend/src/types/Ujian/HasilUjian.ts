import type { DataAkunSiswa } from "../KelolaAkun/AkunSiswa"

export type StatistikUjian = {
    nilai_terendah:number
    nilai_tertinggi:number
    rata_rata:number
    jumlah_peserta:number
}

export type HasilUjianSiswa = DataAkunSiswa & {
    nilai?:number;
    jumlah_benar?:number;
    jumlah_salah?:number;
    jumlah_kosong?:number;
}
