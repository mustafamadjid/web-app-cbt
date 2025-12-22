import { createBrowserRouter } from "react-router";
import { paths } from "./paths";

// Components

// Login Page
import { LoginPage } from "../pages/Auth/LoginPage";

// Dashboard Admin
import { Home } from "@/pages/Admin/Dashboard/Home";
import { KelolaAkunGuru } from "@/pages/Admin/Dashboard/KelolaAkun/AkunGuru";
import { AdminLayout } from "@/layouts/MainLayout/AdminLayout/AdminLayout";

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
            {path:paths.dashboard.kelola_akun_guru,element:<KelolaAkunGuru/>}
        ]
    }
])