import { createBrowserRouter, Navigate } from "react-router";
import { paths } from "./paths";

// Components

// Login Page
import LoginPage from "../pages/Auth/LoginPage";

// Header
import HeaderLayout from "@/layouts/MainLayout/HeaderLayout/HeaderLayout";
import ProtectedRoute from "@/routes/ProtectedRoute";

// Dashboard Admin
import { Home } from "@/pages/Admin/Dashboard/Home";
import KelolaAkunGuru from "@/pages/Admin/Dashboard/KelolaAkun/AkunGuru";
import { KelolaAkunSiswa } from "@/pages/Admin/Dashboard/KelolaAkun/AkunSiswa";
import TambahGuru from "@/pages/Admin/Dashboard/KelolaAkun/TambahAkun/TambahGuru";
import TambahSiswa from "@/pages/Admin/Dashboard/KelolaAkun/TambahAkun/TambahSiswa";
import EditAkunGuru from "@/pages/Admin/Dashboard/KelolaAkun/EditAkun/EditAkunGuru";
import EditAkunSiswa from "@/pages/Admin/Dashboard/KelolaAkun/EditAkun/EditAkunSiswa";

import MataPelajaran from "@/pages/Admin/Dashboard/DataMaster/MataPelajaran";
import DataKelas from "@/pages/Admin/Dashboard/DataMaster/DataKelas";
import RuangUjian from "@/pages/Admin/Dashboard/DataMaster/RuangUjian";
import DataSesi from "@/pages/Admin/Dashboard/DataMaster/DataSesi";
import TambahMataPelajaran from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahMapel";
import TambahKelas from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahKelas";
import TambahRuang from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahRuangUjian";
import TambahSesi from "@/pages/Admin/Dashboard/DataMaster/TambahDataMaster/TambahSesi";
import EditMapel from "@/pages/Admin/Dashboard/DataMaster/EditDataMaster/EditMapel";
import EditKelas from "@/pages/Admin/Dashboard/DataMaster/EditDataMaster/EditKelas";
import EditRuangUjian from "@/pages/Admin/Dashboard/DataMaster/EditDataMaster/EditRuangUjian";
import EditSesi from "@/pages/Admin/Dashboard/DataMaster/EditDataMaster/EditSesi";

import BankSoal from "@/pages/Admin/Dashboard/BankSoal/BankSoal";
import TambahBankSoal from "@/pages/Admin/Dashboard/BankSoal/TambahBankSoal";
import BuatBankSoal from "@/pages/Admin/Dashboard/BankSoal/BuatBankSoal";
import DetailBankSoal from "@/pages/Admin/Dashboard/BankSoal/DetailBankSoal";

import BuatUjian from "@/pages/Admin/Dashboard/Ujian/BuatUjian";
import JadwalUjian from "@/pages/Admin/Dashboard/Ujian/JadwalUjian";
import DetailUjian from "@/pages/Admin/Dashboard/Ujian/DetailUjian/DetailUjian";

import Cetak from "@/pages/Admin/Cetak/Cetak";

import TemplateLayout from "@/layouts/MainLayout/TemplateLayout/TemplateLayout";

import PengaturanProfil from "@/pages/Admin/Dashboard/Pengaturan/Pengaturan";
import HasilUjian from "@/pages/Admin/Dashboard/Ujian/HasilUjian";
import HasilUjianDetail from "@/pages/Admin/Dashboard/Ujian/DetailHasilUjian/HasilUjianDetail";
import StatistikSoalUjian from "@/pages/Admin/Dashboard/Ujian/StatistikSoalUjian/StatistikSoalUjian";
import KoreksiHasilUjian from "@/pages/Admin/Dashboard/Ujian/KoreksiHasilUjian/KoreksiHasilUjian";
import PengumumanManagement from "@/pages/Admin/Dashboard/Pengumuman/PengumumanManagement";
import TambahPengumuman from "@/pages/Admin/Dashboard/Pengumuman/TambahPengumuman";
import EditPengumuman from "@/pages/Admin/Dashboard/Pengumuman/EditPengumuman";
import KelolaSesi from "@/pages/Admin/Dashboard/KelolaSesi/KelolaSesi";
import HomeSiswa from "@/pages/Siswa/HomeSiswa";
import PengumumanSiswa from "@/pages/Siswa/Pengumuman/PengumumanSiswa";
import UjianSiswa from "@/pages/Siswa/Ujian/UjianSiswa/UjianSiswa";
import UjianTokenSiswa from "@/pages/Siswa/Ujian/UjianSiswa/UjianTokenSiswa";
import UjianMulaiSiswa from "@/pages/Siswa/Ujian/UjianSiswa/UjianMulaiSiswa";
import HasilUjianSiswa from "@/pages/Siswa/Ujian/HasilUjian/HasilUjianSiswa";
import HasilUjianSiswaDetailSiswa from "@/pages/Siswa/Ujian/HasilUjian/HasilUjianSiswaDetail";
import PublicOnlyRoute from "./PublicOnlyRoute";
import ProfilePage from "@/pages/Profile/ProfilePage";
import NotFound from "@/pages/NotFound/NotFound";

export const router = createBrowserRouter([
  // Login Page
  {
    path: paths.public.login,
    element: <PublicOnlyRoute/>,
    children : [
      {
        index:true,
        element:<LoginPage/>
      }
    ]
  },
  {
    path:"/",
    element:<Navigate to={paths.public.login} replace/>
  },
  {
    path: "/dashboard/administrator",
    element: <ProtectedRoute allowedRoles={["ADMIN"]} />,
    children: [
      {
        element: <TemplateLayout />,
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
                path: paths.dashboard.kelola_sesi,
                element: <KelolaSesi />,
              },
              {
                path: paths.dashboard.data_master_mapel,
                element: <MataPelajaran />,
              },
              {
                path: paths.dashboard.data_master_kelas,
                element: <DataKelas />,
              },
              {
                path: paths.dashboard.data_master_ruang,
                element: <RuangUjian />,
              },
              {
                path: paths.dashboard.pengaturan,
                element: <PengaturanProfil />,
              },
              { path: paths.dashboard.profil_admin, element: <ProfilePage /> },
              { path: paths.dashboard.data_master_sesi, element: <DataSesi /> },
              { path: paths.dashboard.bank_soal, element: <BankSoal /> },
              {
                path: paths.dashboard.preview_bank_soal,
                element: <DetailBankSoal />,
              },
              { path: paths.dashboard.buat_ujian, element: <BuatUjian /> },
              { path: paths.dashboard.jadwal_ujian, element: <JadwalUjian /> },
              { path: paths.dashboard.detail_ujian, element: <DetailUjian /> },
              { path: paths.dashboard.hasil_ujian, element: <HasilUjian /> },
              {
                path: paths.dashboard.hasil_ujian_detail,
                element: <HasilUjianDetail />,
              },
              {
                path: paths.dashboard.hasil_ujian_statistik_soal,
                element: <StatistikSoalUjian />,
              },
              {
                path: paths.dashboard.hasil_ujian_siswa_detail,
                element: <KoreksiHasilUjian />,
              },
              {
                path: paths.dashboard.pengumuman_admin,
                element: <PengumumanManagement />,
              },
              { path: paths.dashboard.cetak, element: <Cetak /> },
            ],
          },
          { path: paths.dashboard.tambah_guru, element: <TambahGuru /> },
          { path: paths.dashboard.tambah_siswa, element: <TambahSiswa /> },
          { path: paths.dashboard.edit_guru, element: <EditAkunGuru /> },
          { path: paths.dashboard.edit_siswa, element: <EditAkunSiswa /> },
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
            path: paths.dashboard.edit_data_master_mapel,
            element: <EditMapel />,
          },
          {
            path: paths.dashboard.edit_data_master_kelas,
            element: <EditKelas />,
          },
          {
            path: paths.dashboard.edit_data_master_ruang,
            element: <EditRuangUjian />,
          },
          {
            path: paths.dashboard.edit_data_master_sesi,
            element: <EditSesi />,
          },
          {
            path: paths.dashboard.tambah_bank_soal,
            element: <TambahBankSoal />,
          },
          {
            path: paths.dashboard.buat_bank_soal,
            element: <BuatBankSoal />,
          },
          {
            path: paths.dashboard.tambah_pengumuman_admin,
            element: <TambahPengumuman />,
          },
          {
            path: paths.dashboard.edit_pengumuman_admin,
            element: <EditPengumuman />,
          },
        ],
      },
    ],
  },
  {
    path: "/dashboard/guru",
    element: <ProtectedRoute allowedRoles={["GURU"]} />,
    children: [
      {
        element: <TemplateLayout />,
        children: [
          {
            element: <HeaderLayout />,
            children: [
              { index: true, element: <Home /> },
              { path: paths.dashboard.bank_soal_guru, element: <BankSoal /> },
              {
                path: paths.dashboard.preview_bank_soal_guru,
                element: <DetailBankSoal />,
              },
              { path: paths.dashboard.buat_ujian_guru, element: <BuatUjian /> },
              {
                path: paths.dashboard.jadwal_ujian_guru,
                element: <JadwalUjian />,
              },
              {
                path: paths.dashboard.detail_ujian_guru,
                element: <DetailUjian />,
              },
              {
                path: paths.dashboard.hasil_ujian_guru,
                element: <HasilUjian />,
              },
              {
                path: paths.dashboard.hasil_ujian_detail_guru,
                element: <HasilUjianDetail />,
              },
              {
                path: paths.dashboard.hasil_ujian_statistik_soal_guru,
                element: <StatistikSoalUjian />,
              },
              {
                path: paths.dashboard.hasil_ujian_siswa_detail_guru,
                element: <KoreksiHasilUjian />,
              },
              {
                path: paths.dashboard.pengumuman_guru,
                element: <PengumumanManagement />,
              },
              { path: paths.dashboard.cetak_guru, element: <Cetak /> },
              { path: paths.dashboard.profil_guru, element: <ProfilePage /> },
            ],
          },
          {
            path: paths.dashboard.tambah_pengumuman_guru,
            element: <TambahPengumuman />,
          },
          {
            path: paths.dashboard.edit_pengumuman_guru,
            element: <EditPengumuman />,
          },
          {
            path: paths.dashboard.tambah_bank_soal_guru,
            element: <TambahBankSoal />,
          },
          {
            path: paths.dashboard.buat_bank_soal_guru,
            element: <BuatBankSoal />,
          },
        ],
      },
    ],
  },
  {
    path: "/dashboard/siswa",
    element: <ProtectedRoute allowedRoles={["SISWA"]} />,
    children: [
      {
        element: <TemplateLayout />,
        children: [
          {
            element: <HeaderLayout />,
            children: [
              { index: true, element: <HomeSiswa /> },
              {
                path: paths.dashboard.pengumuman_siswa,
                element: <PengumumanSiswa />,
              },
              { path: paths.dashboard.ujian_siswa, element: <UjianSiswa /> },
              {
                path: paths.dashboard.hasil_ujian_siswa,
                element: <HasilUjianSiswa />,
              },
              {
                path: paths.dashboard.ujian_siswa_token,
                element: <UjianTokenSiswa />,
              },
              {
                path: paths.dashboard.ujian_siswa_mulai,
                element: <UjianMulaiSiswa />,
              },
              {
                path: paths.dashboard.hasil_ujian_detail_siswa,
                element: <HasilUjianSiswaDetailSiswa />,
              },
              { path: paths.dashboard.profil_siswa, element: <ProfilePage /> },
              
            ],
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <NotFound />,
  },
]);
