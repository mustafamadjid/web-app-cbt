import { createBrowserRouter } from "react-router";
import { paths } from "./paths";

// Components

// Login Page
import { LoginPage } from "../pages/Auth/LoginPage";

// Dashboard Admin
import { Home } from "@/pages/Admin/Dashboard/Home";
import { KelolaAkunGuru } from "@/pages/Admin/Dashboard/KelolaAkun/AkunGuru";
import { KelolaAkunSiswa } from "@/pages/Admin/Dashboard/KelolaAkun/AkunSiswa";
import { TambahGuru } from "@/pages/Admin/Dashboard/KelolaAkun/TambahAkun/TambahGuru";
import { TambahSiswa } from "@/pages/Admin/Dashboard/KelolaAkun/TambahAkun/TambahSiswa";
import { AdminLayout } from "@/layouts/MainLayout/AdminLayout/AdminLayout";
import { PengaturanProfil } from "@/pages/Admin/Dashboard/Pengaturan/Pengaturan";

export const router = createBrowserRouter([
    // Login Page
    {
        path: paths.public.login,
        element: <LoginPage />
    },
    {
        path:"/dashboard/administrator",
        element:<AdminLayout/>,
        children:[
            {index:true,element:<Home/>},
            {path:paths.dashboard.kelola_akun_guru,element:<KelolaAkunGuru/>},
            {path:paths.dashboard.kelola_akun_siswa,element:<KelolaAkunSiswa/>},
            {path:paths.dashboard.tambah_guru,element:<TambahGuru/>},
            {path:paths.dashboard.tambah_siswa,element:<TambahSiswa/>},
            {path:paths.dashboard.pengaturan,element:<PengaturanProfil/>},
        ]
    }
])