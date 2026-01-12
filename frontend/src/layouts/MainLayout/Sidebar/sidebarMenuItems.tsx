import { paths } from "@/routes/paths";
import type { SidebarMenuItem } from "../../../types/Sidebar/SidebarMenu";

import SvgIcons from "@/assets/SvgIcons/svgIcons";

export const mainMenuItems: SidebarMenuItem[] = [
  {
    id: 1,
    type: "link",
    label: "Dashboard",
    to: paths.dashboard.home_admin,
    end: true,
    icon: SvgIcons.dashboard,
  },
  {
    id: 2,
    type: "group",
    label: "Kelola Pengguna",
    icon: SvgIcons.users,
    children: [
      {
        id: 201,
        label: "Akun Guru",
        to: paths.dashboard.kelola_akun_guru,
        icon: SvgIcons.userSingle,
      },
      {
        id: 202,
        label: "Akun Siswa",
        to: paths.dashboard.kelola_akun_siswa,
        icon: SvgIcons.userSingle,
      },
    ],
  },
  {
    id: 3,
    type: "group",
    label: "Data Master",
    icon: SvgIcons.book,
    children: [
      {
        id: 301,
        label: "Data Mata Pelajaran",
        to: paths.dashboard.data_master_mapel,
        icon: SvgIcons.chevronDouble,
      },
      {
        id: 302,
        label: "Data Kelas",
        to: paths.dashboard.data_master_kelas,
        icon: SvgIcons.chevronDouble,
      },
      {
        id: 303,
        label: "Data Ruang Ujian",
        to: paths.dashboard.data_master_ruang,
        icon: SvgIcons.chevronDouble,
      },
      {
        id: 304,
        label: "Data Sesi",
        to: paths.dashboard.data_master_sesi,
        icon: SvgIcons.chevronDouble,
      },
    ],
  },
  {
    id: 4,
    type: "link",
    label: "Bank Soal",
    to: paths.dashboard.bank_soal,
    end: true,
    icon: SvgIcons.bankSoal,
  },
  {
    id: 5,
    type: "group",
    label: "Ujian",
    icon: SvgIcons.book,
    children: [
      {
        id: 501,
        label: "Buat Ujian",
        to: paths.dashboard.buat_ujian,
        icon: SvgIcons.chevronDouble,
      },
      {
        id: 502,
        label: "Jadwal Ujian",
        to: paths.dashboard.jadwal_ujian,
        icon: SvgIcons.chevronDouble,
      },
      {
        id: 503,
        label: "Hasil Ujian",
        to: paths.dashboard.hasil_ujian,
        icon: SvgIcons.chevronDouble,
      },
    ],
  },
  {
    id: 6,
    type: "link",
    label: "Cetak",
    to: paths.dashboard.cetak,
    end: true,
    icon: SvgIcons.print,
  },
];

export const footerMenuItems: SidebarMenuItem[] = [
  {
    id: 101, 
    type: "link",
    label: "Pengaturan",
    to: paths.dashboard.pengaturan,
    icon: SvgIcons.settings,
  },
];
