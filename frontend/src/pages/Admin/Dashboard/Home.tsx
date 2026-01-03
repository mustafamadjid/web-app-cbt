import { StatistikWidget } from "@/components/features/widget/Soal/StatistikWidget";
import { SimpleStatWidget } from "@/components/features/widget/Soal/StatistikWidgetSimpler";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import { JadwalUjianWidget } from "@/components/features/widget/JadwalUjian/JadwalUjian";
import { LogAktivitasWidget } from "@/components/features/widget/LogAktivitas/LogAktivitasWidget";

import type { PengumumanItem } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type { AktivitasLogItem } from "@/types/Log/LogAktivitas";

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
        url: "https://example.com/docs/surat-perubahan-ruangan.pdf",
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
    nama_ujian: "Ujian Tengah Semester Matematika",
    pengawas_ujian: "Budi Santoso",
    tgl_ujian: "Kamis, 20 November 2025",
    waktu_mulai: "08:00",
    sesi_ujian: 1,
    ruang_ujian: "Ruang 101",
  },
  {
    id: 2,
    nama_ujian: "Ujian Tengah Semester Bahasa Indonesia",
    pengawas_ujian: "Siti Aminah",
    tgl_ujian: "Kamis, 20 November 2025",
    waktu_mulai: "10:30",
    sesi_ujian: 2,
    ruang_ujian: "Ruang 102",
  },
  {
    id: 3,
    nama_ujian: "Ujian Akhir Semester Fisika",
    pengawas_ujian: "Ahmad Fauzi",
    tgl_ujian: "Jumat, 21 November 2025",
    waktu_mulai: "08:00",
    sesi_ujian: 1,
    ruang_ujian: "Lab Fisika",
  },
  {
    id: 4,
    nama_ujian: "Ujian Akhir Semester Kimia",
    pengawas_ujian: "Dewi Lestari",
    tgl_ujian: "Jumat, 21 November 2025",
    waktu_mulai: "13:00",
    sesi_ujian: 2,
    ruang_ujian: "Lab Kimia",
  },
  {
    id: 4,
    nama_ujian: "Ujian Akhir Semester Kimia",
    pengawas_ujian: "Dewi Lestari",
    tgl_ujian: "Jumat, 21 November 2025",
    waktu_mulai: "13:00",
    sesi_ujian: 2,
    ruang_ujian: "Lab Kimia",
  },
  {
    id: 4,
    nama_ujian: "Ujian Akhir Semester Kimia",
    pengawas_ujian: "Dewi Lestari",
    tgl_ujian: "Jumat, 21 November 2025",
    waktu_mulai: "13:00",
    sesi_ujian: 2,
    ruang_ujian: "Lab Kimia",
  },
];

export const dummyAktivitas: AktivitasLogItem[] = [
  {
    id: 1,
    username: "admin01",
    role: "admin",
    aksi: "LOGIN",
    deskripsi: "Masuk ke sistem melalui halaman admin.",
    waktu: "Kamis, 20 November 2025 • 08:12",
  },
  {
    id: 2,
    username: "guru_sri",
    role: "guru",
    aksi: "UPDATE",
    deskripsi: "Mengubah nilai ujian Matematika kelas XI IPA 1.",
    waktu: "Kamis, 20 November 2025 • 09:05",
  },
  {
    id: 3,
    username: "siswa_andi",
    role: "siswa",
    aksi: "CREATE",
    deskripsi: "Mengumpulkan tugas Bahasa Indonesia (Bab 3).",
    waktu: "Kamis, 20 November 2025 • 10:41",
  },
  {
    id: 3,
    username: "siswa_andi",
    role: "siswa",
    aksi: "CREATE",
    deskripsi: "Mengumpulkan tugas Bahasa Indonesia (Bab 3).",
    waktu: "Kamis, 20 November 2025 • 10:41",
  },
  {
    id: 3,
    username: "siswa_andi",
    role: "siswa",
    aksi: "CREATE",
    deskripsi: "Mengumpulkan tugas Bahasa Indonesia (Bab 3).",
    waktu: "Kamis, 20 November 2025 • 10:41",
  },
];



export const Home = ()=> {
    return (
      <>
        {/* Widget dan konten */}
        <div className="flex flex-col gap-5">
          <div className="flex items-center gap-3 px-8 ">
            {/* Statistik Ujian Terlaksana */}
            <div className="w-1/2 sm:w-1/2">
              <StatistikWidget
                title="Total Ujian Terlaksana"
                value={4}
                footerText="Dari total 10 ujian"
                percent={75}
                centerLabel="Selesai"
              />
            </div>

            {/* Statistik Bank soal */}
            <div className="w-1/3 sm:w-1/4">
              <SimpleStatWidget
                title="Total Bank Soal"
                value="100K"
                trend="up"
                trendText="Naik 25% dari Nov 2025"
              />
            </div>

            {/* Statistik Mata Pelajaran */}
            <div className="w-1/3 sm:w-1/4">
              <SimpleStatWidget
                title="Total Mata Pelajaran"
                value="10"
                trend="up"
                trendText="Naik 25% dari Nov 2025"
              />
            </div>
          </div>

          <div className="sm:flex sm:flex-row flex flex-col gap-4 sm:gap-4 px-8">
            {/* Pengumuman */}
            <div className="w-full">
              <PengumumanWidget title="Pengumuman" items={pengumuman} />
            </div>

            <div className="w-full flex flex-col gap-4">
              {/* Jadwal Ujian */}
              <div className="w-full ">
                <JadwalUjianWidget items={dummyJadwalUjian} />
              </div>

              {/* Log Aktivitas */}
              <div className="w-full">
                <LogAktivitasWidget
                  items={dummyAktivitas}
                  lihatSemuaTo="/log-aktivitas"
                />
              </div>
            </div>
          </div>
        </div>
      </>
    );
}