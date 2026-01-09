import { api, type ApiEnvelope } from "../../api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

import type { BankSoalItem } from "@/types/DataMaster/BankSoal";
import type {
  BankSoalOption,
  BuatUjianFormValues,
  BuatUjianSubmitResponse,
  GuruPengawasOption,
  RuangUjianOption,
  SesiUjianOption,
  SiswaPreviewItem,
} from "@/types/Ujian/BuatUjian";
import { getBankSoalByKelas } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { DUMMY_SISWA } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";

const DUMMY_GURU_PENGAWAS: GuruPengawasOption[] = [
  { id: 1, nama: "Budi Santoso", nip: "198305232010011001", mapel: "Matematika" },
  { id: 2, nama: "Siti Rahmawati", nip: "198912112012012002", mapel: "Bahasa Indonesia" },
  { id: 3, nama: "Dedi Pratama", nip: "198706142011021003", mapel: "Fisika" },
];

const DUMMY_SESI: SesiUjianOption[] = [
  { id: 1, kode: "S1", nama: "Sesi 1 (Pagi)" },
  { id: 2, kode: "S2", nama: "Sesi 2 (Siang)" },
  { id: 3, kode: "S3", nama: "Sesi 3 (Sore)" },
];

const DUMMY_RUANG: RuangUjianOption[] = [
  { id: 1, nama: "Lab Komputer 1" },
  { id: 2, nama: "Lab Komputer 2" },
  { id: 3, nama: "Ruang Multimedia" },
  { id: 4, nama: "Aula Utama" },
];

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const getTotalSoal = (item: BankSoalItem) => {
  if (typeof item.total_soal === "number") return item.total_soal;
  const pg = item.jumlah_soal_pg ?? 0;
  const essay = item.jumlah_soal_essay ?? 0;
  return pg + essay;
};

export async function getUjianGuruPengawasOptions(): Promise<GuruPengawasOption[]> {
  await sleep(180);
  return [...DUMMY_GURU_PENGAWAS];
}

export async function getUjianSesiOptions(): Promise<SesiUjianOption[]> {
  await sleep(180);
  return [...DUMMY_SESI];
}

export async function getUjianRuangOptions(): Promise<RuangUjianOption[]> {
  await sleep(180);
  return [...DUMMY_RUANG];
}

export async function getUjianBankSoalOptions(params: {
  tingkatKelasId?: number;
}): Promise<BankSoalOption[]> {
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  const kelasId = tingkatKelas ?? undefined;
  const data = await getBankSoalByKelas({
    tingkatKelasId: params.tingkatKelasId,
    kelasId,
  });

  return data.map((item) => ({
    id: item.id,
    nama: item.nama_banksoal ?? "-",
    mata_pelajaran: item.mata_pelajaran,
    materi: item.materi,
    kelas: item.kelas,
    jumlah_soal_pg: item.jumlah_soal_pg,
    jumlah_soal_essay: item.jumlah_soal_essay,
    total_soal: getTotalSoal(item),
  }));
}
export async function getUjianSiswaPreview(params: {
  // Data dari server harusnya di sorting duluan berdasarkan nama nomor absen dan kelas
  tingkatKelasId?: number;
}): Promise<SiswaPreviewItem[]> {
  await sleep(220);
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  if (!tingkatKelas) return [];

  return DUMMY_SISWA.filter((siswa) => siswa.__kelasId === tingkatKelas).map(
    (siswa) => ({
    id: siswa.id,
    nama: siswa.namaLengkap,
    username: siswa.username,
    no_absen: siswa.noAbsen,
    kelas: siswa.kelas,
    status_akun: siswa.statusAkun,
  })
  );
}

export async function submitBuatUjian(values: BuatUjianFormValues) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<BuatUjianSubmitResponse>>("/ujian", {
    method: "POST",
    data,
  });

  return res.data;
}
