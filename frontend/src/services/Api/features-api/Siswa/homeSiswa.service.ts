import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type {
  PengumumanItem,
  AnnouncementDoc,
} from "@/types/Widget/Pengumuman";
import type {
  SiswaDashboardSummary,
  SiswaProfile,
  SiswaSemesterAverage,
} from "@/types/Widget/SiswaDashboard";
import { getJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const dummyProfile: SiswaProfile = {
  id: 14,
  nama: "Nabila Putri",
  kelas: "XI IPS 2",
  tingkat_kelas_id: 2,
};

const dummyPengumumanDocs: AnnouncementDoc[] = [
  {
    id: 11,
    name: "Panduan Ujian Matematika.pdf",
    url: "#",
    sizeLabel: "320 KB",
  },
];

const dummyPengumuman: PengumumanItem[] = [
  {
    id: 1,
    judul: "Simulasi Ujian Tengah Semester",
    isi_pengumuman:
      "Simulasi UTS akan dilaksanakan minggu depan. Pastikan kamu sudah mengunduh kisi-kisi yang tersedia di LMS.",
    tanggal_rilis_pengumuman: "Senin, 24 Februari 2025",
    dokumen: dummyPengumumanDocs,
  },
  {
    id: 2,
    judul: "Pengumpulan Tugas Praktikum",
    isi_pengumuman:
      "Tugas praktikum IPA dikumpulkan paling lambat Jumat pukul 16.00 WIB.",
    tanggal_rilis_pengumuman: "Rabu, 26 Februari 2025",
    dokumen: null,
  },
  {
    id: 3,
    judul: "Kelas Tambahan Bahasa Inggris",
    isi_pengumuman:
      "Kelas tambahan akan dibuka setiap Selasa pukul 15.00 di Ruang Multimedia.",
    tanggal_rilis_pengumuman: "Jumat, 28 Februari 2025",
    dokumen: null,
  },
];

const dummyRataRata: SiswaSemesterAverage[] = [
  { semester: "Ganjil 2024/2025", rata_rata: 84.5, target: 85 },
  { semester: "Genap 2024/2025", rata_rata: 88.2, target: 87 },
];

export async function getSiswaDashboardSummary(): Promise<SiswaDashboardSummary> {
  await sleep(200);
  const jadwal = await getJadwalUjian({
    tingkatKelasId: dummyProfile.tingkat_kelas_id,
  });
  const ujianSelesai = jadwal.filter(
    (item) => item.status_ujian === "selesai"
  ).length;

  return {
    ujian_selesai: ujianSelesai,
    total_ujian: jadwal.length,
    rata_rata_semester: dummyRataRata,
  };
}

export async function getSiswaJadwalUjian(): Promise<JadwalUjianItem[]> {
  await sleep(150);
  const jadwal = await getJadwalUjian({
    tingkatKelasId: dummyProfile.tingkat_kelas_id,
  });
  return jadwal.filter((item) => item.status_ujian !== "selesai");
}

export async function getSiswaPengumuman(): Promise<PengumumanItem[]> {
  await sleep(120);
  return dummyPengumuman;
}

export async function getSiswaRataRataSemester(): Promise<
  SiswaSemesterAverage[]
> {
  await sleep(80);
  return dummyRataRata;
}
