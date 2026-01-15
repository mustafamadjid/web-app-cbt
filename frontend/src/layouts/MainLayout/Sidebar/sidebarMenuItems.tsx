import { paths } from "@/routes/paths";
import type { SidebarMenuItem, Role } from "../../../types/Sidebar/SidebarMenu";

import SvgIcons from "@/assets/SvgIcons/svgIcons";

const ADMIN_ONLY: Role[] = ["ADMIN"];
const GURU_ONLY: Role[] = ["GURU"];
const SISWA_ONLY: Role[] = ["SISWA"];


export const mainMenuItems: SidebarMenuItem[] = [
  // ADMIN MENU
  {
    id: 1,
    type: "link",
    label: "Dashboard",
    to: paths.dashboard.home_admin,
    end: true,
    icon: SvgIcons.dashboard,
    roles: ADMIN_ONLY,
  },
  {
    id: 2,
    type: "group",
    label: "Kelola Pengguna",
    icon: SvgIcons.users,
    roles: ADMIN_ONLY,
    children: [
      {
        id: 201,
        label: "Akun Guru",
        to: paths.dashboard.kelola_akun_guru,
        icon: SvgIcons.userSingle,
        roles: ADMIN_ONLY,
      },
      {
        id: 202,
        label: "Akun Siswa",
        to: paths.dashboard.kelola_akun_siswa,
        icon: SvgIcons.userSingle,
        roles: ADMIN_ONLY,
      },
    ],
  },
  {
    id: 3,
    type: "group",
    label: "Data Master",
    icon: SvgIcons.book,
    roles: ADMIN_ONLY,
    children: [
      {
        id: 301,
        label: "Data Mata Pelajaran",
        to: paths.dashboard.data_master_mapel,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
      },
      {
        id: 302,
        label: "Data Kelas",
        to: paths.dashboard.data_master_kelas,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
      },
      {
        id: 303,
        label: "Data Ruang Ujian",
        to: paths.dashboard.data_master_ruang,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
      },
      {
        id: 304,
        label: "Data Sesi",
        to: paths.dashboard.data_master_sesi,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
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
    roles: ADMIN_ONLY,
  },
  {
    id: 5,
    type: "group",
    label: "Ujian",
    icon: SvgIcons.book,
    roles: ADMIN_ONLY,
    children: [
      {
        id: 501,
        label: "Buat Ujian",
        to: paths.dashboard.buat_ujian,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
      },
      {
        id: 502,
        label: "Jadwal Ujian",
        to: paths.dashboard.jadwal_ujian,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
      },
      {
        id: 503,
        label: "Hasil Ujian",
        to: paths.dashboard.hasil_ujian,
        icon: SvgIcons.chevronDouble,
        roles: ADMIN_ONLY,
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
    roles: ADMIN_ONLY,
  },

  // GURU Menu
  {
    id: 101,
    type: "link",
    label: "Dashboard",
    to: paths.dashboard.home_guru,
    end: true,
    icon: SvgIcons.dashboard,
    roles: GURU_ONLY,
  },
  {
    id: 4,
    type: "link",
    label: "Bank Soal",
    to: paths.dashboard.bank_soal_guru,
    end: true,
    icon: SvgIcons.bankSoal,
    roles: GURU_ONLY,
  },
  {
    id: 5,
    type: "group",
    label: "Ujian",
    icon: SvgIcons.book,
    roles: GURU_ONLY,
    children: [
      {
        id: 501,
        label: "Buat Ujian",
        to: paths.dashboard.buat_ujian_guru,
        icon: SvgIcons.chevronDouble,
        roles: GURU_ONLY,
      },
      {
        id: 502,
        label: "Jadwal Ujian",
        to: paths.dashboard.jadwal_ujian_guru,
        icon: SvgIcons.chevronDouble,
        roles: GURU_ONLY,
      },
      {
        id: 503,
        label: "Hasil Ujian",
        to: paths.dashboard.hasil_ujian_guru,
        icon: SvgIcons.chevronDouble,
        roles: GURU_ONLY,
      },
    ],
  },
];

export const footerMenuItems: SidebarMenuItem[] = [
  {
    id: 101,
    type: "link",
    label: "Pengaturan",
    to: paths.dashboard.pengaturan,
    icon: SvgIcons.settings,
    roles: ADMIN_ONLY,
  },
];
