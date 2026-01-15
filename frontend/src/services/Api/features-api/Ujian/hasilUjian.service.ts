import { api, type ApiEnvelope } from "@/services/Api/api";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type {
  HasilUjianSiswa,
  StatistikUjian,
} from "@/types/Ujian/HasilUjian";

export type HasilUjianDetailResponse = {
  statistik: StatistikUjian;
  siswa: HasilUjianSiswa[];
};

const dummyHasilUjianList: JadwalUjianItem[] = [
  {
    id: 21,
    nama_ujian: "Ujian Matematika",
    pengawas_ujian: "Pak Budi Santoso",
    tgl_ujian: "Senin, 12 Februari 2026",
    tanggal_ujian: "2026-02-12",
    waktu_mulai: "07:30",
    waktu_selesai: "09:00",
    sesi_ujian: 1,
    ruang_ujian: "Ruang 101",
    id_ruang: 1,
    status_ujian: "selesai",
    started: 1,
    pembuat_username: "guru_sri",
    pengawas_username: "guru_sri",
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 1",
  },
  {
    id: 22,
    nama_ujian: "Ujian IPA",
    pengawas_ujian: "Pak Andi Pratama",
    tgl_ujian: "Selasa, 13 Februari 2026",
    tanggal_ujian: "2026-02-13",
    waktu_mulai: "08:00",
    waktu_selesai: "09:40",
    sesi_ujian: 1,
    ruang_ujian: "Lab IPA",
    id_ruang: 3,
    status_ujian: "selesai",
    started: 1,
    pembuat_username: "guru_raka",
    pengawas_username: "guru_raka",
    tingkat_kelas: 12,
    tingkat_kelas_id: 3,
    nama_kelas: "12 IPS 1",
  },
  {
    id: 23,
    nama_ujian: "Ujian Bahasa Indonesia",
    pengawas_ujian: "Bu Siti Nuraini",
    tgl_ujian: "Rabu, 14 Februari 2026",
    tanggal_ujian: "2026-02-14",
    waktu_mulai: "10:00",
    waktu_selesai: "11:30",
    sesi_ujian: 2,
    ruang_ujian: "Ruang 102",
    id_ruang: 2,
    status_ujian: "selesai",
    started: 1,
    pembuat_username: "guru_sri",
    pengawas_username: "guru_sri",
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPA 2",
  },
];

const dummyHasilUjianDetail: Record<number, HasilUjianDetailResponse> = {
  21: {
    statistik: {
      nilai_terendah: 58,
      nilai_tertinggi: 96,
      rata_rata: 78,
      jumlah_peserta: 28,
    },
    siswa: [
      {
        id: 101,
        role: "SISWA",
        namaLengkap: "Alya Putri",
        username: "alya.putri",
        jenisKelamin: "PEREMPUAN",
        email: "alya.putri@example.com",
        noHp: "081234567890",
        noAbsen: 2,
        angkatan: 2024,
        tempatLahir: "Bandung",
        tanggalLahir: "2008-06-14",
        id_tingkat_kelas: 1,
        id_nama_kelas: "10-IPA-1",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-1.png",
        nilai: 92,
        jumlah_benar: 18,
        jumlah_salah: 2,
        jumlah_kosong: 0,
      },
      {
        id: 102,
        role: "SISWA",
        namaLengkap: "Dimas Pratama",
        username: "dimas.pratama",
        jenisKelamin: "LAKI_LAKI",
        email: "dimas.pratama@example.com",
        noHp: "081298765432",
        noAbsen: 7,
        angkatan: 2024,
        tempatLahir: "Jakarta",
        tanggalLahir: "2008-01-22",
        id_tingkat_kelas: 1,
        id_nama_kelas: "10-IPA-1",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-2.png",
        nilai: 76,
        jumlah_benar: 14,
        jumlah_salah: 6,
        jumlah_kosong: 0,
      },
      {
        id: 103,
        role: "SISWA",
        namaLengkap: "Raka Kusuma",
        username: "raka.kusuma",
        jenisKelamin: "LAKI_LAKI",
        email: "raka.kusuma@example.com",
        noHp: "081277788899",
        noAbsen: 15,
        angkatan: 2024,
        tempatLahir: "Bekasi",
        tanggalLahir: "2008-09-05",
        id_tingkat_kelas: 1,
        id_nama_kelas: "10-IPA-1",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-3.png",
        nilai: 68,
        jumlah_benar: 12,
        jumlah_salah: 7,
        jumlah_kosong: 1,
      },
    ],
  },
  22: {
    statistik: {
      nilai_terendah: 60,
      nilai_tertinggi: 93,
      rata_rata: 81,
      jumlah_peserta: 30,
    },
    siswa: [
      {
        id: 201,
        role: "SISWA",
        namaLengkap: "Nadia Putri",
        username: "nadia.putri",
        jenisKelamin: "PEREMPUAN",
        email: "nadia.putri@example.com",
        noHp: "081355566677",
        noAbsen: 3,
        angkatan: 2023,
        tempatLahir: "Bogor",
        tanggalLahir: "2007-11-18",
        id_tingkat_kelas: 3,
        id_nama_kelas: "12-IPS-1",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-4.png",
        nilai: 88,
        jumlah_benar: 16,
        jumlah_salah: 3,
        jumlah_kosong: 1,
      },
      {
        id: 202,
        role: "SISWA",
        namaLengkap: "Fajar Nugraha",
        username: "fajar.nugraha",
        jenisKelamin: "LAKI_LAKI",
        email: "fajar.nugraha@example.com",
        noHp: "081366677788",
        noAbsen: 11,
        angkatan: 2023,
        tempatLahir: "Depok",
        tanggalLahir: "2007-05-02",
        id_tingkat_kelas: 3,
        id_nama_kelas: "12-IPS-1",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-5.png",
        nilai: 79,
        jumlah_benar: 14,
        jumlah_salah: 5,
        jumlah_kosong: 1,
      },
    ],
  },
  23: {
    statistik: {
      nilai_terendah: 55,
      nilai_tertinggi: 90,
      rata_rata: 74,
      jumlah_peserta: 26,
    },
    siswa: [
      {
        id: 301,
        role: "SISWA",
        namaLengkap: "Salsa Anindya",
        username: "salsa.anindya",
        jenisKelamin: "PEREMPUAN",
        email: "salsa.anindya@example.com",
        noHp: "081399900011",
        noAbsen: 5,
        angkatan: 2024,
        tempatLahir: "Semarang",
        tanggalLahir: "2008-03-21",
        id_tingkat_kelas: 2,
        id_nama_kelas: "11-IPA-2",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-6.png",
        nilai: 84,
        jumlah_benar: 13,
        jumlah_salah: 2,
        jumlah_kosong: 0,
      },
      {
        id: 302,
        role: "SISWA",
        namaLengkap: "Rangga Saputra",
        username: "rangga.saputra",
        jenisKelamin: "LAKI_LAKI",
        email: "rangga.saputra@example.com",
        noHp: "081388811122",
        noAbsen: 9,
        angkatan: 2024,
        tempatLahir: "Solo",
        tanggalLahir: "2008-12-08",
        id_tingkat_kelas: 2,
        id_nama_kelas: "11-IPA-2",
        statusAkun: "aktif",
        urlGambarProfil: "/images/avatar/student-7.png",
        nilai: 69,
        jumlah_benar: 10,
        jumlah_salah: 4,
        jumlah_kosong: 1,
      },
    ],
  },
};

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

const USE_DUMMY = true; // ✅ set false saat BE sudah siap

export type HasilUjianFilterParams = {
  tingkatKelasId?: number;
  tahun?: string;
};

const filterByTingkatKelasId = (
  data: JadwalUjianItem[],
  tingkatKelasId?: number
) => {
  if (tingkatKelasId == null) return data;
  return data.filter((ujian) => ujian.tingkat_kelas_id === tingkatKelasId);
};

const filterByTahun = (data: JadwalUjianItem[], tahun?: string) => {
  if (!tahun) return data;
  return data.filter((ujian) => ujian.tanggal_ujian?.startsWith(`${tahun}-`));
};

export async function getHasilUjianList(
  params: HasilUjianFilterParams = {}
): Promise<JadwalUjianItem[]> {
  if (USE_DUMMY) {
    await sleep(200);
    return filterByTahun(
      filterByTingkatKelasId(dummyHasilUjianList, params.tingkatKelasId),
      params.tahun
    );
  }

  const res = await api<ApiEnvelope<JadwalUjianItem[]>>("/ujian/hasil", {
    method: "GET",
    params: {
      tingkat_kelas_id: params.tingkatKelasId ?? undefined,
      tahun: params.tahun ?? undefined,
    },
  });
  return res.data;
}

export async function getHasilUjianDetail(
  ujianId: number
): Promise<HasilUjianDetailResponse> {
  if (USE_DUMMY) {
    await sleep(200);
    const detail = dummyHasilUjianDetail[ujianId];
    if (!detail) {
      throw new Error("Detail hasil ujian tidak ditemukan.");
    }
    return detail;
  }

  const res = await api<ApiEnvelope<HasilUjianDetailResponse>>(
    `/ujian/hasil/${ujianId}`,
    {
      method: "GET",
    }
  );
  return res.data;
}
