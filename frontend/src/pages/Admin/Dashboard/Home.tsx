import { LayoutDashboard } from "lucide-react";

// Widgets
import { StatistikWidget } from "@/components/features/widget/Soal/StatistikWidget";
import { SimpleStatWidget } from "@/components/features/widget/Soal/StatistikWidgetSimpler";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import { JadwalUjianWidget } from "@/components/features/widget/JadwalUjian/JadwalUjian";
import { LogAktivitasWidget } from "@/components/features/widget/LogAktivitas/LogAktivitasWidget";
import { TotalSiswaGuruWidget } from "@/components/features/widget/StatistikPengguna/TotalSiswaGuruWidget";
import { UjianBerlangsungWidget } from "@/components/features/widget/UjianBerlangsung/UjianBerlangsungWidget";

// Types
import type { PengumumanItem } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type { AktivitasLogItem } from "@/types/Log/LogAktivitas";
import type { UjianBerlangsungItem } from "@/types/Widget/UjianBerlangsung";

// --- DUMMY DATA ---
const pengumuman: PengumumanItem[] = [
  {
    id: "p1",
    judul: "Informasi Perubahan Ruangan",
    tanggal_rilis_pengumuman: "Senin, 24/11/2025",
    isi_pengumuman:
      "Dengan hormat, kami informasikan bahwa telah terjadi penyesuaian terkait pelaksanaan ujian...\n\nRuang sebelumnya: Ruang 02\nRuangan baru: Kelas ABC",
    dokumen: [
      {
        name: "Surat Perubahan Ruangan.pdf",
        url: "#",
        sizeLabel: "412 KB",
      },
    ],
  },
  {
    id: "p2",
    judul: "Informasi Perubahan Tanggal Ujian",
    tanggal_rilis_pengumuman: "Senin, 24/11/2025",
    isi_pengumuman: "Ujian diundur ke Selasa, 25/11/2025 pukul 09:00.",
    dokumen: null,
  },
];

export const dummyJadwalUjian: JadwalUjianItem[] = [
  {
    id: 1,
    nama_ujian: "UTS Matematika",
    pengawas_ujian: "Budi Santoso",
    tgl_ujian: "Kamis, 20 Nov 2025",
    waktu_mulai: "08:00",
    sesi_ujian: 1,
    ruang_ujian: "R. 101",
  },
  {
    id: 2,
    nama_ujian: "UTS B. Indonesia",
    pengawas_ujian: "Siti Aminah",
    tgl_ujian: "Kamis, 20 Nov 2025",
    waktu_mulai: "10:30",
    sesi_ujian: 2,
    ruang_ujian: "R. 102",
  },
  {
    id: 3,
    nama_ujian: "UAS Fisika",
    pengawas_ujian: "Ahmad Fauzi",
    tgl_ujian: "Jumat, 21 Nov 2025",
    waktu_mulai: "08:00",
    sesi_ujian: 1,
    ruang_ujian: "Lab Fisika",
  },
];

export const dummyAktivitas: AktivitasLogItem[] = [
  {
    id: 1,
    username: "admin01",
    role: "admin",
    aksi: "LOGIN",
    deskripsi: "Masuk ke sistem melalui halaman admin.",
    waktu: "08:12",
  },
  {
    id: 2,
    username: "guru_sri",
    role: "guru",
    aksi: "UPDATE",
    deskripsi: "Mengubah nilai ujian Matematika kelas XI IPA 1.",
    waktu: "09:05",
  },
  {
    id: 3,
    username: "siswa_andi",
    role: "siswa",
    aksi: "CREATE",
    deskripsi: "Mengumpulkan tugas Bahasa Indonesia.",
    waktu: "10:41",
  },
];

export const dummyUjianBerlangsung: UjianBerlangsungItem[] = [
  {
    id: 101,
    nama_ujian: "Penilaian Harian Matematika Wajib",
    mata_pelajaran: "Matematika",
    kelas: ["XI IPA 1", "XI IPA 2"],
    waktu_mulai: "08:00",
    waktu_selesai: "09:30",
    total_siswa: 70,
    siswa_mengerjakan: 20, // Sedang mengerjakan (Kuning)
    siswa_selesai: 45, // Sudah selesai (Hijau)
    // Sisanya (5) belum login (Abu-abu)
  },
  {
    id: 102,
    nama_ujian: "Ujian Susulan Bahasa Inggris",
    mata_pelajaran: "B. Inggris",
    kelas: ["X IPS 3"],
    waktu_mulai: "08:30",
    waktu_selesai: "10:00",
    total_siswa: 35,
    siswa_mengerjakan: 30,
    siswa_selesai: 2,
  },
];

export const Home = () => {
  return (
    <div className="min-h-screen bg-[#ecf1ed]  pb-20">
      <div className="mx-auto max-w-[1920px] space-y-8 p-4 sm:p-6 lg:p-8">
        {/* === HEADER SECTION === */}
        <div className="flex flex-col gap-1">
          <h1 className="flex items-center gap-2 text-2xl font-bold text-slate-800 sm:text-3xl">
            <LayoutDashboard className="h-6 w-6 text-[#397e50]  sm:h-8 sm:w-8" />
            Dashboard
          </h1>
          <p className="text-sm text-slate-500">
            Selamat datang kembali, ringkasan aktivitas akademik hari ini.
          </p>
        </div>

        {/* === TOP STATS (GRID 4 KOLOM) === */}
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
          {/* Card 1: Total User (Card Spesial) */}
          <div className="h-full">
            <TotalSiswaGuruWidget
              totalGuru={24}
              totalSiswa={450}
              className="h-full"
            />
          </div>

          {/* Card 2: Statistik Ujian (Donut) */}
          <div className="h-full">
            <StatistikWidget
              title="Total Ujian Terlaksana"
              value={4}
              footerText="4 Selesai dari 10 Jadwal"
              percent={40}
              centerLabel="Progress"
              className="h-full"
            />
          </div>

          {/* Card 3: Bank Soal (Simple) */}
          <div className="h-full">
            <SimpleStatWidget
              title="Total Bank Soal"
              value="1,240"
              trend="up"
              trendText="+12% bulan ini"
              className="h-full"
            />
          </div>

          {/* Card 4: Mata Pelajaran (Simple) */}
          <div className="h-full">
            <SimpleStatWidget
              title="Mata Pelajaran Aktif"
              value="18"
              trend="neutral"
              trendText="Tidak ada perubahan"
              className="h-full"
            />
          </div>
        </div>

        {/* === MAIN CONTENT === */}
        <div className="grid gap-6 lg:grid-cols-12 items-stretch">
          {/* LEFT COLUMN */}
          <div className="flex flex-col gap-6 lg:col-span-7 xl:col-span-8 h-full min-h-0">
            <PengumumanWidget
              title="Papan Pengumuman"
              items={pengumuman}
              className="flex-1 min-h-0"
            />

            <UjianBerlangsungWidget
              items={dummyUjianBerlangsung}
              className="flex-1 min-h-0"
            />
          </div>

          {/* RIGHT COLUMN */}
          <div className="flex flex-col gap-6 lg:col-span-5 xl:col-span-4 h-full min-h-0">
            <JadwalUjianWidget
              items={dummyJadwalUjian}
              className="flex-1 min-h-0"
              maxHeightClassName="h-full"
            />

            <LogAktivitasWidget
              items={dummyAktivitas}
              lihatSemuaTo="/log-aktivitas"
              className="flex-1 min-h-0"
              maxHeightClassName="h-full"
            />
          </div>
        </div>
      </div>
    </div>
  );
};
