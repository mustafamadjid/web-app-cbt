
// Widgets
import { useEffect, useState } from "react";

import StatistikWidget from "@/components/features/widget/Soal/StatistikWidget";
import SimpleStatWidget from "@/components/features/widget/Soal/StatistikWidgetSimpler";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import JadwalUjianWidget from "@/components/features/widget/JadwalUjian/JadwalUjian";
import LogAktivitasWidget from "@/components/features/widget/LogAktivitas/LogAktivitasWidget";
import TotalSiswaGuruWidget from "@/components/features/widget/StatistikPengguna/TotalSiswaGuruWidget";
import UjianBerlangsungWidget from "@/components/features/widget/UjianBerlangsung/UjianBerlangsungWidget";

// Types
import type { PengumumanItem } from "@/types/Widget/Pengumuman";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type {
  AktivitasLogItem,
  LogAktivitasApiItem,
} from "@/types/Log/LogAktivitas";
import type { UjianBerlangsungItem } from "@/types/Widget/UjianBerlangsung";
import type { Role } from "@/types/Sidebar/SidebarMenu";
import { useAuth } from "@/contexts/AuthContext";
import { getLogAktivitas } from "@/services/Api/features-api/LogAktivitas/log-aktivitas.service";

// --- DUMMY DATA ---
const pengumuman: PengumumanItem[] = [
  {
    id: 1,
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
    id: 2,
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
    status_ujian: "belum_dimulai",
    started: 0,
    sesi_ujian: 1,
    ruang_ujian: "R. 101",
  },
  {
    id: 2,
    nama_ujian: "UTS B. Indonesia",
    pengawas_ujian: "Siti Aminah",
    tgl_ujian: "Kamis, 20 Nov 2025",
    waktu_mulai: "10:30",
    status_ujian: "belum_dimulai",
    started: 0,
    sesi_ujian: 2,
    ruang_ujian: "R. 102",
  },
  {
    id: 3,
    nama_ujian: "UAS Fisika",
    pengawas_ujian: "Ahmad Fauzi",
    tgl_ujian: "Jumat, 21 Nov 2025",
    waktu_mulai: "08:00",
    status_ujian: "belum_dimulai",
    started: 0,
    sesi_ujian: 1,
    ruang_ujian: "Lab Fisika",
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

const mapLogAktivitasToWidget = (
  item: LogAktivitasApiItem
): AktivitasLogItem => ({
  id: item.id_aktivitas,
  username: `Pengguna ${item.id_pengguna}`,
  role: "unknown",
  aksi: item.action,
  deskripsi: item.description || "-",
  waktu: item.created_at,
});

export const Home = () => {
  const {user} = useAuth(); 
  const role = (user?.role ?? "SISWA") as Role;
  const [logAktivitas, setLogAktivitas] = useState<AktivitasLogItem[]>([]);
  const [logAktivitasError, setLogAktivitasError] = useState<string | null>(
    null
  );

  const logAktivitasTitle = logAktivitasError
    ? "Log Aktivitas (Gagal Memuat)"
    : "Log Aktivitas";

  useEffect(() => {
    if (role !== "ADMIN") return;

    let isActive = true;
    setLogAktivitasError(null);

    void getLogAktivitas()
      .then((data) => {
        if (!isActive) return;
        const mapped = data.map(mapLogAktivitasToWidget);
        setLogAktivitas(mapped);
      })
      .catch((error: Error) => {
        if (!isActive) return;
        setLogAktivitas([]);
        setLogAktivitasError(
          error.message || "Gagal memuat log aktivitas."
        );
      });

    return () => {
      isActive = false;
    };
  }, [role]);


  return (
    <div className="min-h-screen bg-[#ecf1ed]  pb-20">
      <div className="mx-auto max-w-[1920px] space-y-8 p-4 sm:p-6 lg:p-8">
        {/* === TOP STATS (GRID 4 KOLOM) === */}
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-5">
          {/* Card 1: Total User (Dibuat Lebih Lebar) */}
          {role === "ADMIN" && (
            <div className="h-full sm:col-span-2 lg:col-span-2">
              <TotalSiswaGuruWidget
                totalGuru={24}
                totalSiswa={450}
                className="h-full"
              />
            </div>
          )}

          {/* Card 2: Statistik Ujian */}
          <div
            className={`h-full  ${
              role === "ADMIN" ? "lg:col-span-1" : "lg:col-span-3"
            }`}
          >
            <StatistikWidget
              title="Total Ujian Terlaksana"
              value={4}
              footerText="4 Selesai dari 10 Jadwal"
              percent={40}
              centerLabel="Progress"
              className="h-full"
            />
          </div>

          {/* Card 3: Bank Soal */}
          <div className="h-full lg:col-span-1">
            <SimpleStatWidget
              title="Total Bank Soal"
              value="1,240"
              trend="up"
              trendText="+12% bulan ini"
              className="h-full"
            />
          </div>

          {/* Card 4: Mata Pelajaran */}
          <div className="h-full lg:col-span-1">
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
          <div
            className={`flex flex-col gap-6 xl:col-span-8 h-full min-h-0 lg:col-span-6 ${
              role === "ADMIN" ? "" : "lg:h-300"
            }`}
          >
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

            {role === "ADMIN" && (
              <LogAktivitasWidget
                title={logAktivitasTitle}
                items={logAktivitas}
                lihatSemuaTo="/log-aktivitas"
                className="flex-1 min-h-0"
                maxHeightClassName="h-full"
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
