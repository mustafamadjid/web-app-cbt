import { createBrowserRouter } from "react-router";
import { paths } from "./paths";

// Components

// Login Page
import { LoginPage } from "../pages/Auth/LoginPage";

export const router = createBrowserRouter([
    // Login Page
    {
        path: paths.public.login,
        element: <LoginPage />
    }
])