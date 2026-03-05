import { api } from "@/services/Api/api";
import { useFetch } from "@/hooks/fetch";
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
        id_pengguna: 101,
        role: "SISWA",
        nama_lengkap: "Alya Putri",
        username: "alya.putri",
        nisn: "1234567890",
        jenis_kelamin: "PEREMPUAN",
        email: "alya.putri@example.com",
        no_hp: "081234567890",
        no_absen: 2,
        angkatan: 2024,
        tempat_lahir: "Bandung",
        tanggal_lahir: "2008-06-14",
        tingkat_kelas: 1,
        nama_kelas: "10-IPA-1",
        kelas: "10-IPA-1",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-1.png",
        nilai: 92,
        jumlah_benar: 18,
        jumlah_salah: 2,
        jumlah_kosong: 0,
      },
      {
        id_pengguna: 102,
        role: "SISWA",
        nama_lengkap: "Dimas Pratama",
        username: "dimas.pratama",
        nisn: "2234567890",
        jenis_kelamin: "LAKI_LAKI",
        email: "dimas.pratama@example.com",
        no_hp: "081298765432",
        no_absen: 7,
        angkatan: 2024,
        tempat_lahir: "Jakarta",
        tanggal_lahir: "2008-01-22",
        tingkat_kelas: 1,
        nama_kelas: "10-IPA-1",
        kelas: "10-IPA-1",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-2.png",
        nilai: 76,
        jumlah_benar: 14,
        jumlah_salah: 6,
        jumlah_kosong: 0,
      },
      {
        id_pengguna: 103,
        role: "SISWA",
        nama_lengkap: "Raka Kusuma",
        username: "raka.kusuma",
        nisn: "3234567890",
        jenis_kelamin: "LAKI_LAKI",
        email: "raka.kusuma@example.com",
        no_hp: "081277788899",
        no_absen: 15,
        angkatan: 2024,
        tempat_lahir: "Bekasi",
        tanggal_lahir: "2008-09-05",
        tingkat_kelas: 1,
        nama_kelas: "10-IPA-1",
        kelas: "10-IPA-1",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-3.png",
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
        id_pengguna: 201,
        role: "SISWA",
        nama_lengkap: "Nadia Putri",
        username: "nadia.putri",
        nisn: "4234567890",
        jenis_kelamin: "PEREMPUAN",
        email: "nadia.putri@example.com",
        no_hp: "081355566677",
        no_absen: 3,
        angkatan: 2023,
        tempat_lahir: "Bogor",
        tanggal_lahir: "2007-11-18",
        tingkat_kelas: 3,
        nama_kelas: "12-IPS-1",
        kelas: "12-IPS-1",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-4.png",
        nilai: 88,
        jumlah_benar: 16,
        jumlah_salah: 3,
        jumlah_kosong: 1,
      },
      {
        id_pengguna: 202,
        role: "SISWA",
        nama_lengkap: "Fajar Nugraha",
        username: "fajar.nugraha",
        nisn: "5234567890",
        jenis_kelamin: "LAKI_LAKI",
        email: "fajar.nugraha@example.com",
        no_hp: "081366677788",
        no_absen: 11,
        angkatan: 2023,
        tempat_lahir: "Depok",
        tanggal_lahir: "2007-05-02",
        tingkat_kelas: 3,
        nama_kelas: "12-IPS-1",
        kelas: "12-IPS-1",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-5.png",
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
        id_pengguna: 301,
        role: "SISWA",
        nama_lengkap: "Salsa Anindya",
        username: "salsa.anindya",
        nisn: "6234567890",
        jenis_kelamin: "PEREMPUAN",
        email: "salsa.anindya@example.com",
        no_hp: "081399900011",
        no_absen: 5,
        angkatan: 2024,
        tempat_lahir: "Semarang",
        tanggal_lahir: "2008-03-21",
        tingkat_kelas: 2,
        nama_kelas: "11-IPA-2",
        kelas: "11-IPA-2",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-6.png",
        nilai: 84,
        jumlah_benar: 13,
        jumlah_salah: 2,
        jumlah_kosong: 0,
      },
      {
        id_pengguna: 302,
        role: "SISWA",
        nama_lengkap: "Rangga Saputra",
        username: "rangga.saputra",
        nisn: "7234567890",
        jenis_kelamin: "LAKI_LAKI",
        email: "rangga.saputra@example.com",
        no_hp: "081388811122",
        no_absen: 9,
        angkatan: 2024,
        tempat_lahir: "Solo",
        tanggal_lahir: "2008-12-08",
        tingkat_kelas: 2,
        nama_kelas: "11-IPA-2",
        kelas: "11-IPA-2",
        status_akun: "AKTIF",
        foto_profil: "/images/avatar/student-7.png",
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

  const queryParams: Record<string, string | number | undefined> = {
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    tahun: params.tahun ?? undefined,
  };

  const res = await api<JadwalUjianItem[]>("/ujian/hasil", {
    method: "GET",
    params: queryParams,
  });
  return res;
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

  const res = await api<HasilUjianDetailResponse>(
    `/ujian/hasil/${ujianId}`,
    {
      method: "GET",
    }
  );
  return res;
}

// =====================
// Hook Wrappers
// =====================

export function useGetHasilUjianList(params: HasilUjianFilterParams = {}) {
  return useFetch(
    () => getHasilUjianList(params),
    [params.tingkatKelasId, params.tahun],
  );
}

export function useGetHasilUjianDetail(ujianId: number) {
  return useFetch(() => getHasilUjianDetail(ujianId), [ujianId]);
}
