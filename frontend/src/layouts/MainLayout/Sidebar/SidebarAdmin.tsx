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
      "hover:bg-neutral-tertiary hover:text-[#397e50]",
      isActive
        ? "bg-neutral-tertiary text-[#397e50] font-semibold"
        : "text-body",
    ].join(" ");

  const subNavItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "pl-10 flex items-center px-2 py-1.5 rounded-base group transition-colors",
      "hover:bg-neutral-tertiary hover:text-[#397e50]",
      isActive
        ? "bg-neutral-tertiary text-[#397e50] font-semibold"
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
        className="text-heading bg-transparent box-border border border-transparent hover:bg-neutral-secondary-medium focus:ring-4 focus:ring-neutral-tertiary font-medium leading-5 rounded-base ms-3 mt-3 text-sm p-2 focus:outline-none inline-flex sm:hidden"
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
          "fixed inset-0 z-30 bg-black/40 transition-opacity sm:hidden",
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
        ].join(" ")}
      >
        <div className="h-full px-3 py-25 bg-neutral-primary-soft border-e border-default flex flex-col">
          {/* Header mobile: tombol close */}
          <div className="flex items-center justify-between sm:hidden mb-3">
            <span className="text-sm font-medium text-heading">Menu</span>
            <button
              type="button"
              onClick={() => setIsSidebarOpen(false)}
              className="p-2 rounded-base hover:bg-neutral-secondary-medium focus:outline-none focus:ring-4 focus:ring-neutral-tertiary"
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

          <ul className="space-y-2 font-medium flex flex-col gap-3">
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

            {/* Dropdown Kelola Akun */}
            <li>
              <button
                type="button"
                onClick={() => setIsEcomOpen((v) => !v)}
                className={[
                  "flex items-center w-full justify-between px-2 py-1.5 rounded-base group transition-colors",
                  "hover:bg-neutral-tertiary hover:text-[#397e50]",
                  isKelolaPenggunaActive
                    ? "bg-neutral-tertiary text-[#397e50] font-semibold"
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

            {/* Data Master (saya biarkan tombol dulu; nanti submenu pakai NavLink juga) */}
            <li>
              <button
                type="button"
                onClick={() => setIsDataMasterOpen((v) => !v)}
                className="group flex w-full items-center justify-between rounded-base px-2 py-1.5 text-body hover:bg-neutral-tertiary hover:text-[#397e50]"
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
                  className={`h-5 w-5 transition-transform duration-200 ${
                    isDataMasterOpen ? "rotate-180" : "rotate-0"
                  }`}
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

            {/* Ujian (ganti jadi NavLink agar bisa active highlight) */}
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

          {/* Divider section */}
          <ul className="mt-auto space-y-2 font-medium border-t border-default pt-4">
            {/* Ini sebaiknya juga NavLink jika memang ada route nyata */}
            <li>
              <NavLink
                to="/dashboard/administrator/docs"
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
                    d="M5 19V4a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1v13H7a2 2 0 0 0-2 2Zm0 0a2 2 0 0 0 2 2h12M9 3v14m7 0v4"
                  />
                </svg>
                <span className="flex-1 ms-3 whitespace-nowrap">
                  Documentation
                </span>
              </NavLink>
            </li>

            <li>
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
            </li>
          </ul>
        </div>
      </aside>
    </>
  );
};
