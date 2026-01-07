import { paths } from "@/routes/paths";
import type { SidebarMenuItem } from "../../../types/Sidebar/SidebarMenu";

const chevronDoubleIcon = (className: string) => (
  <svg
    className={className}
    aria-hidden="true"
    xmlns="http://www.w3.org/2000/svg"
    width="24"
    height="24"
    fill="none"
    viewBox="0 0 24 24"
  >
    <path
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      d="m7 16 4-4-4-4m6 8 4-4-4-4"
    />
  </svg>
);

export const mainMenuItems: SidebarMenuItem[] = [
  {
    id: "dashboard",
    type: "link",
    label: "Dashboard",
    to: paths.dashboard.home_admin,
    end: true,
    icon: (className) => (
      <svg
        className={className}
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M10 6.025A7.5 7.5 0 1 0 17.975 14H10V6.025Z"
        />
        <path
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M13.5 3c-.169 0-.334.014-.5.025V11h7.975c.011-.166.025-.331.025-.5A7.5 7.5 0 0 0 13.5 3Z"
        />
      </svg>
    ),
  },
  {
    id: "kelola-pengguna",
    type: "group",
    label: "Kelola Pengguna",
    icon: (className) => (
      <svg
        className={className}
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M16 19h4a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-2m-2.236-4a3 3 0 1 0 0-4M3 18v-1a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v1a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1Zm8-10a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
        />
      </svg>
    ),
    children: [
      {
        id: "akun-guru",
        label: "Akun Guru",
        to: paths.dashboard.kelola_akun_guru,
        icon: (className) => (
          <svg
            className={className}
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              strokeWidth="2"
              d="M7 17v1a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-4a3 3 0 0 0-3 3Zm8-9a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
            />
          </svg>
        ),
      },
      {
        id: "akun-siswa",
        label: "Akun Siswa",
        to: paths.dashboard.kelola_akun_siswa,
        icon: (className) => (
          <svg
            className={className}
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              strokeWidth="2"
              d="M7 17v1a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-4a3 3 0 0 0-3 3Zm8-9a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
            />
          </svg>
        ),
      },
    ],
  },
  {
    id: "data-master",
    type: "group",
    label: "Data Master",
    icon: (className) => (
      <svg
        className={className}
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M15 5v14M9 5v14M4 5h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1Z"
        />
      </svg>
    ),
    children: [
      {
        id: "data-mata-pelajaran",
        label: "Data Mata Pelajaran",
        to: paths.dashboard.data_master_mapel,
        icon: chevronDoubleIcon,
      },
      {
        id: "data-kelas",
        label: "Data Kelas",
        to: paths.dashboard.data_master_kelas,
        icon: chevronDoubleIcon,
      },
      {
        id: "data-ruang-ujian",
        label: "Data Ruang Ujian",
        to: paths.dashboard.data_master_ruang,
        icon: chevronDoubleIcon,
      },
      {
        id: "data-sesi",
        label: "Data Sesi",
        to: paths.dashboard.data_master_sesi,
        icon: chevronDoubleIcon,
      },
    ],
  },
  {
    id: "bank-soal",
    type: "link",
    label: "Bank Soal",
    to: paths.dashboard.bank_soal,
    end: true,
    icon: (className) => (
      <svg
        className={className}
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M11 18h2M5.875 3h12.25c.483 0 .875.448.875 1v16c0 .552-.392 1-.875 1H5.875C5.392 21 5 20.552 5 20V4c0-.552.392-1 .875-1Z"
        />
      </svg>
    ),
  },
  {
    id: "ujian",
    type: "group",
    label: "Ujian",
    icon: (className) => (
      <svg
        className={className}
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M15 5v14M9 5v14M4 5h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1Z"
        />
      </svg>
    ),
    children: [
      {
        id: "jadwal-ujian",
        label: "Jadwal Ujian",
        to: paths.dashboard.jadwal_ujian,
        icon: chevronDoubleIcon,
      },
      {
        id: "buat-ujian",
        label: "Buat Ujian",
        to: paths.dashboard.buat_ujian,
        icon: chevronDoubleIcon,
      },
    ],
  },
];

export const footerMenuItems: SidebarMenuItem[] = [
  {
    id: "pengaturan",
    type: "link",
    label: "Pengaturan",
    to: paths.dashboard.pengaturan,
    icon: (className) => (
      <svg
        className={className}
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        width="24"
        height="24"
      >
        <path
          fill="currentColor"
          d="M19.14 12.94c.04-.3.06-.61.06-.94c0-.32-.02-.64-.07-.94l2.03-1.58a.49.49 0 0 0 .12-.61l-1.92-3.32a.49.49 0 0 0-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.484.484 0 0 0-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58a.49.49 0 0 0-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6s3.6 1.62 3.6 3.6s-1.62 3.6-3.6 3.6"
        />
      </svg>
    ),
  },
];
