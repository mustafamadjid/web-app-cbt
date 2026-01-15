import { createBrowserRouter } from "react-router";
import { paths } from "./paths";

// Components

// Login Page
import LoginPage from "../pages/Auth/LoginPage";

// Header
import HeaderLayout from "@/layouts/MainLayout/HeaderLayout/HeaderLayout";

// Dashboard Admin
import { Home } from "@/pages/Admin/Dashboard/Home";
import KelolaAkunGuru from "@/pages/Admin/Dashboard/KelolaAkun/AkunGuru";
import { KelolaAkunSiswa } from "@/pages/Admin/Dashboard/KelolaAkun/AkunSiswa";
import TambahGuru from "@/pages/Admin/Dashboard/KelolaAkun/TambahAkun/TambahGuru";
import TambahSiswa from "@/pages/Admin/Dashboard/KelolaAkun/TambahAkun/TambahSiswa";

import MataPelajaran from "@/pages/Admin/Dashboard/DataMaster/MataPelajaran";
import DataKelas from "@/pages/Admin/Dashboard/DataMaster/DataKelas";
import RuangUjian from "@/pages/Admin/Dashboard/DataMaster/RuangUjian";
import DataSesi from "@/pages/Admin/Dashboard/DataMaster/DataSesi";
import TambahMataPelajaran from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahMapel";
import TambahKelas from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahKelas";
import TambahRuang from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahRuangUjian";
import TambahSesi from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahSesi";

import BankSoal from "@/pages/Admin/Dashboard/BankSoal/BankSoal";
import TambahBankSoal from "@/pages/Admin/Dashboard/BankSoal/TambahBankSoal";

import BuatUjian from "@/pages/Admin/Dashboard/Ujian/BuatUjian";
import JadwalUjian from "@/pages/Admin/Dashboard/Ujian/JadwalUjian";
import DetailUjian from "@/pages/Admin/Dashboard/Ujian/DetailUjian/DetailUjian";

import Cetak from "@/pages/Admin/Cetak/Cetak";

import AdminLayout from "@/layouts/MainLayout/AdminLayout/AdminLayout";

import PengaturanProfil from "@/pages/Admin/Dashboard/Pengaturan/Pengaturan";
import HasilUjian from "@/pages/Admin/Dashboard/Ujian/HasilUjian";
import HasilUjianDetail from "@/pages/Admin/Dashboard/Ujian/DetailHasilUjian/HasilUjianDetail";
import HasilUjianSiswaDetail from "@/pages/Admin/Dashboard/Ujian/DetailHasilUjian/HasilUjianSiswaDetail";

export const router = createBrowserRouter([
  // Login Page
  {
    path: paths.public.login,
    element: <LoginPage />,
  },
  {
    path: "/dashboard/administrator",
    element: <AdminLayout />,
    children: [
      {
        element: <HeaderLayout />,
        children: [
          { index: true, element: <Home /> },
          {
            path: paths.dashboard.kelola_akun_guru,
            element: <KelolaAkunGuru />,
          },
          {
            path: paths.dashboard.kelola_akun_siswa,
            element: <KelolaAkunSiswa />,
          },
          {
            path: paths.dashboard.data_master_mapel,
            element: <MataPelajaran />,
          },
          { path: paths.dashboard.data_master_kelas, element: <DataKelas /> },
          { path: paths.dashboard.data_master_ruang, element: <RuangUjian /> },
          { path: paths.dashboard.pengaturan, element: <PengaturanProfil /> },
          { path: paths.dashboard.data_master_sesi, element: <DataSesi /> },
          { path: paths.dashboard.bank_soal, element: <BankSoal /> },
          { path: paths.dashboard.buat_ujian, element: <BuatUjian /> },
          { path: paths.dashboard.jadwal_ujian, element: <JadwalUjian /> },
          { path: paths.dashboard.detail_ujian, element: <DetailUjian /> },
          { path: paths.dashboard.hasil_ujian, element: <HasilUjian /> },
          {
            path: paths.dashboard.hasil_ujian_detail,
            element: <HasilUjianDetail />,
          },
          {
            path: paths.dashboard.hasil_ujian_siswa_detail,
            element: <HasilUjianSiswaDetail />,
          },
          { path: paths.dashboard.cetak, element: <Cetak /> },
        ],
      },

      { path: paths.dashboard.tambah_guru, element: <TambahGuru /> },
      { path: paths.dashboard.tambah_siswa, element: <TambahSiswa /> },
      {
        path: paths.dashboard.tambah_data_master_mapel,
        element: <TambahMataPelajaran />,
      },
      {
        path: paths.dashboard.tambah_data_master_kelas,
        element: <TambahKelas />,
      },
      {
        path: paths.dashboard.tambah_data_master_ruang,
        element: <TambahRuang />,
      },
      {
        path: paths.dashboard.tambah_data_master_sesi,
        element: <TambahSesi />,
      },
      {
        path: paths.dashboard.tambah_bank_soal,
        element: <TambahBankSoal />,
      },
    ],
  },
  {
    path: "/dashboard/guru",
    element: <AdminLayout />,
    children: [
      {
        element: <HeaderLayout />,
        children: [
          { index: true, element: <Home /> },

        ]
      }
    ]
  }
]);
