import { useEffect, useMemo, useRef, useState } from "react";
import { NavLink, useLocation, matchPath } from "react-router";

export const SidebarAdmin = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [isEcomOpen, setIsEcomOpen] = useState(false);
  const [isDataMasterOpen, setIsDataMasterOpen] = useState(false);
  const sidebarRef = useRef<HTMLElement | null>(null);

  const { pathname } = useLocation();

  // Definisikan route submenu
  const kelolaPenggunaRoutes = useMemo(
    () => ({
      guru: "/dashboard/administrator/kelola-akun/guru",
      siswa: "/dashboard/administrator/kelola-akun/siswa",
    }),
    []
  );

  // Cek apakah user sedang ada di salah satu halaman Kelola Pengguna
  const isKelolaPenggunaActive =
    !!matchPath({ path: kelolaPenggunaRoutes.guru, end: false }, pathname) ||
    !!matchPath({ path: kelolaPenggunaRoutes.siswa, end: false }, pathname);

  // Auto open dropdown jika sedang berada di dalam section itu
  useEffect(() => {
    if (isKelolaPenggunaActive) setIsEcomOpen(true);
  }, [isKelolaPenggunaActive]);

  // Tutup sidebar via ESC
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setIsSidebarOpen(false);
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  // Tutup sidebar ketika klik di luar (mobile)
  useEffect(() => {
    if (!isSidebarOpen) return;

    const onMouseDown = (e: MouseEvent) => {
      const el = sidebarRef.current;
      if (!el) return;

      if (e.target instanceof Node && !el.contains(e.target)) {
        setIsSidebarOpen(false);
      }
    };

    document.addEventListener("mousedown", onMouseDown);
    return () => document.removeEventListener("mousedown", onMouseDown);
  }, [isSidebarOpen]);

  // Helper class untuk NavLink
  const navItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "flex items-center px-2 py-1.5 rounded-base group transition-colors",
      "hover:bg-[#e3f4ea] hover:text-[#2f6f45]",
      isActive
        ? "bg-[#d5f0df] text-[#2f6f45] font-semibold shadow-sm"
        : "text-body",
    ].join(" ");

  const subNavItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "pl-10 flex items-center px-2 py-1.5 rounded-base group transition-colors",
      "hover:bg-[#e3f4ea] hover:text-[#2f6f45]",
      isActive
        ? "bg-[#d5f0df] text-[#2f6f45] font-semibold shadow-sm"
        : "text-body",
    ].join(" ");

  return (
    <>
      {/* Mobile: tombol buka sidebar */}
      <button
        type="button"
        onClick={() => setIsSidebarOpen(true)}
        aria-controls="separator-sidebar"
        aria-expanded={isSidebarOpen}
        className="text-heading bg-transparent box-border border border-transparent hover:bg-[#e3f4ea] focus:ring-4 focus:ring-[#cbead7] font-medium leading-5 rounded-base ms-3 mt-3 text-sm p-2 focus:outline-none inline-flex sm:hidden"
      >
        <span className="sr-only">Open sidebar</span>
        <svg
          className="w-6 h-6"
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
            strokeWidth="2"
            d="M5 7h14M5 12h14M5 17h10"
          />
        </svg>
      </button>

      {/* Overlay (mobile) */}
      <div
        className={[
          "fixed inset-0 z-30 bg-[#1a3d2a]/40 transition-opacity sm:hidden",
          isSidebarOpen
            ? "opacity-100 pointer-events-auto"
            : "opacity-0 pointer-events-none",
        ].join(" ")}
        onClick={() => setIsSidebarOpen(false)}
        aria-hidden="true"
      />

      {/* Sidebar */}
      <aside
        id="separator-sidebar"
        aria-label="Sidebar"
        ref={sidebarRef}
        className={[
          "fixed top-0 left-0 z-40 w-64 h-screen",
          "transition-transform duration-200 ease-out",
          "sm:translate-x-0",
          isSidebarOpen ? "translate-x-0" : "-translate-x-full",
          // penting:
          "flex flex-col overflow-hidden",
        ].join(" ")}
      >
        {/* BODY (mengisi sisa tinggi layar) */}
        <div className="flex-1 min-h-0 px-3 py-5 bg-gradient-to-b from-[#f0faf3] via-[#e6f5ec] to-[#dbf0e3] border-e border-[#cbead7] flex flex-col">
          {/* Header mobile: stated (tidak ikut scroll) */}
          <div className="flex items-center justify-between sm:hidden mb-3 shrink-0">
            <span className="text-sm font-medium text-[#2f6f45]">Menu</span>
            <button
              type="button"
              onClick={() => setIsSidebarOpen(false)}
              className="p-2 rounded-base hover:bg-[#e3f4ea] focus:outline-none focus:ring-4 focus:ring-[#cbead7] text-[#2f6f45]"
            >
              <span className="sr-only">Close sidebar</span>
              <svg
                className="w-5 h-5"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
              >
                <path
                  d="M6 6l12 12M18 6L6 18"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                />
              </svg>
            </button>
          </div>

          {/* MENU UTAMA (scroll di sini) */}
          <div className="flex-1 min-h-0 overflow-y-auto pb-4">
            {/* HEADER / LOGO (tidak ikut scroll) */}
            <div className="p-2 flex items-center shrink-0 rounded-base bg-white/70 border border-[#cbead7] shadow-sm">
              <img
                src="/Images/LoginPageImg/logo-fi.png"
                alt=""
                className="w-20"
              />
              <h2 className="text-sm font-semibold text-[#2f6f45]">
                SMA IT Fitrah Insani
              </h2>
            </div>
            <ul className="space-y-2 font-medium flex flex-col gap-3 py-4">
              {/* Dashboard */}
              <li>
                <NavLink
                  to="/dashboard/administrator"
                  end
                  className={navItemClass}
                >
                  <svg
                    className="w-5 h-5 transition duration-75 group-hover:text-[#397e50]"
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
                  <span className="ms-3">Dashboard</span>
                </NavLink>
              </li>

              {/* Dropdown Kelola Pengguna */}
              <li>
                <button
                  type="button"
                  onClick={() => setIsEcomOpen((v) => !v)}
                  className={[
                    "flex items-center w-full justify-between px-2 py-1.5 rounded-base group transition-colors",
                    "hover:bg-[#e3f4ea] hover:text-[#2f6f45]",
                    isKelolaPenggunaActive
                      ? "bg-[#d5f0df] text-[#2f6f45] font-semibold shadow-sm"
                      : "text-body",
                  ].join(" ")}
                  aria-controls="dropdown-kelola-pengguna"
                  aria-expanded={isEcomOpen}
                >
                  <svg
                    className="shrink-0 w-5 h-5 transition duration-75 group-hover:text-[#397e50]"
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

                  <span className="flex-1 ms-3 text-left rtl:text-right whitespace-nowrap">
                    Kelola Pengguna
                  </span>

                  <svg
                    className={[
                      "w-5 h-5 transition-transform duration-200",
                      isEcomOpen ? "rotate-180" : "rotate-0",
                    ].join(" ")}
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
                      d="m19 9-7 7-7-7"
                    />
                  </svg>
                </button>

                <ul
                  id="dropdown-kelola-pengguna"
                  className={[
                    "py-2 space-y-2",
                    isEcomOpen ? "block" : "hidden",
                  ].join(" ")}
                >
                  <li>
                    <NavLink
                      to={kelolaPenggunaRoutes.guru}
                      className={subNavItemClass}
                    >
                      <svg
                        className="w-6 h-6 group-hover:text-[#397e50]"
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
                      <span className="ms-2">Akun Guru</span>
                    </NavLink>
                  </li>

                  <li>
                    <NavLink
                      to={kelolaPenggunaRoutes.siswa}
                      className={subNavItemClass}
                    >
                      <svg
                        className="w-6 h-6 group-hover:text-[#397e50]"
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
                      <span className="ms-2">Akun Siswa</span>
                    </NavLink>
                  </li>
                </ul>
              </li>

              {/* Data Master */}
              <li>
                <button
                  type="button"
                  onClick={() => setIsDataMasterOpen((v) => !v)}
                    className="group flex w-full items-center justify-between rounded-base px-2 py-1.5 text-body hover:bg-[#e3f4ea] hover:text-[#2f6f45]"
                  aria-controls="dropdown-data-master"
                  aria-expanded={isDataMasterOpen}
                >
                  <span className="flex items-center gap-3">
                    <svg
                      className="h-5 w-5 shrink-0 transition duration-75 group-hover:text-[#397e50]"
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
                    <span className="whitespace-nowrap">Data Master</span>
                  </span>

                  <svg
                    className={[
                      "h-5 w-5 transition-transform duration-200",
                      isDataMasterOpen ? "rotate-180" : "rotate-0",
                    ].join(" ")}
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
                      d="m19 9-7 7-7-7"
                    />
                  </svg>
                </button>
              </li>

              {/* Ujian */}
              <li>
                <NavLink
                  to="/dashboard/administrator/ujian"
                  className={navItemClass}
                >
                  <svg
                    className="w-6 h-6"
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
                      d="M12 6.03v13m0-13c-2.819-.831-4.715-1.076-8.029-1.023A.99.99 0 0 0 3 6v11c0 .563.466 1.014 1.03 1.007 3.122-.043 5.018.212 7.97 1.023m0-13c2.819-.831 4.715-1.076 8.029-1.023A.99.99 0 0 1 21 6v11c0 .563-.466 1.014-1.03 1.007-3.122-.043-5.018.212-7.97 1.023"
                    />
                  </svg>
                  <span className="ms-3">Ujian</span>
                </NavLink>
              </li>
            </ul>
          </div>

          {/* Divider section (tidak ikut scroll) */}
          <div className="shrink-0 mt-4 border-t border-[#cbead7] pt-4">
            <ul className="space-y-2 font-medium">
              <li>
                <NavLink
                  to="/dashboard/administrator/pengaturan"
                  className={navItemClass}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 24 24"
                    width="24"
                    height="24"
                  >
                    <path
                      fill="currentColor"
                      d="M19.14 12.94c.04-.3.06-.61.06-.94c0-.32-.02-.64-.07-.94l2.03-1.58a.49.49 0 0 0 .12-.61l-1.92-3.32a.49.49 0 0 0-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.484.484 0 0 0-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58a.49.49 0 0 0-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6s3.6 1.62 3.6 3.6s-1.62 3.6-3.6 3.6"
                    ></path>
                  </svg>
                  <span className="flex-1 ms-3 whitespace-nowrap">
                    Pengaturan
                  </span>
                </NavLink>
              </li>

              {/* <li>
                <NavLink
                  to="/dashboard/administrator/support"
                  className={navItemClass}
                >
                  <svg
                    className="shrink-0 w-5 h-5 transition duration-75 group-hover:text-[#397e50]"
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
                      d="m13.46 8.291 3.849-3.849a1.5 1.5 0 0 1 2.122 0l.127.127a1.5 1.5 0 0 1 0 2.122l-3.84 3.838a4 4 0 0 0-2.258-2.238Zm0 0a4 4 0 0 1 2.263 2.238l3.662-3.662a8.961 8.961 0 0 1 0 10.27l-3.676-3.676m-2.25-5.17 3.678-3.676a8.961 8.961 0 0 0-10.27 0l3.662 3.662a4 4 0 0 0-2.238 2.258L4.615 6.863a8.96 8.96 0 0 0 0 10.27l3.662-3.662a4 4 0 0 0 2.258 2.238l-3.672 3.676a8.96 8.96 0 0 0 10.27 0l-3.662-3.662a4.001 4.001 0 0 0 2.238-2.262m0 0 3.849 3.848a1.5 1.5 0 0 1 0 2.122l-.127.126a1.499 1.499 0 0 1-2.122 0l-3.838-3.838a4 4 0 0 0 2.238-2.258Zm.29-1.461a4 4 0 1 1-8 0 4 4 0 0 1 8 0Zm-7.718 1.471-3.84 3.838a1.5 1.5 0 0 0 0 2.122l.128.126a1.5 1.5 0 0 0 2.122 0l3.848-3.848a4 4 0 0 1-2.258-2.238Zm2.248-5.19L6.69 4.442a1.5 1.5 0 0 0-2.122 0l-.127.127a1.5 1.5 0 0 0 0 2.122l3.849 3.848a4 4 0 0 1 2.238-2.258Z"
                    />
                  </svg>
                  <span className="flex-1 ms-3 whitespace-nowrap">Support</span>
                </NavLink>
              </li>

              <li>
                <NavLink
                  to="/dashboard/administrator/pro"
                  className={navItemClass}
                >
                  <svg
                    className="shrink-0 w-5 h-5 transition duration-75 group-hover:text-[#397e50]"
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
                      d="m10.051 8.102-3.778.322-1.994 1.994a.94.94 0 0 0 .533 1.6l2.698.316m8.39 1.617-.322 3.78-1.994 1.994a.94.94 0 0 1-1.595-.533l-.4-2.652m8.166-11.174a1.366 1.366 0 0 0-1.12-1.12c-1.616-.279-4.906-.623-6.38.853-1.671 1.672-5.211 8.015-6.31 10.023a.932.932 0 0 0 .162 1.111l.828.835.833.832a.932.932 0 0 0 1.111.163c2.008-1.102 8.35-4.642 10.021-6.312 1.475-1.478 1.133-4.77.855-6.385Zm-2.961 3.722a1.88 1.88 0 1 1-3.76 0 1.88 1.88 0 0 1 3.76 0Z"
                    />
                  </svg>
                  <span className="flex-1 ms-3 whitespace-nowrap">
                    PRO version
                  </span>
                </NavLink>
              </li> */}
            </ul>
          </div>
        </div>
      </aside>
    </>
  );
};
