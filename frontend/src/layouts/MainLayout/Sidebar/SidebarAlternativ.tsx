// import { useEffect, useMemo, useRef, useState } from "react";
// import { NavLink, useLocation, matchPath } from "react-router";

// export const SidebarAdmin = () => {
//   const [isSidebarOpen, setIsSidebarOpen] = useState(false);
//   const [isEcomOpen, setIsEcomOpen] = useState(false);
//   const [isDataMasterOpen, setIsDataMasterOpen] = useState(false);
//   const sidebarRef = useRef<HTMLElement | null>(null);

//   const { pathname } = useLocation();

//   const kelolaPenggunaRoutes = useMemo(
//     () => ({
//       guru: "/dashboard/administrator/kelola-akun/guru",
//       siswa: "/dashboard/administrator/kelola-akun/siswa",
//     }),
//     []
//   );

//   const isKelolaPenggunaActive =
//     !!matchPath({ path: kelolaPenggunaRoutes.guru, end: false }, pathname) ||
//     !!matchPath({ path: kelolaPenggunaRoutes.siswa, end: false }, pathname);

//   useEffect(() => {
//     if (isKelolaPenggunaActive) setIsEcomOpen(true);
//   }, [isKelolaPenggunaActive]);

//   useEffect(() => {
//     const onKeyDown = (e: KeyboardEvent) => {
//       if (e.key === "Escape") setIsSidebarOpen(false);
//     };
//     window.addEventListener("keydown", onKeyDown);
//     return () => window.removeEventListener("keydown", onKeyDown);
//   }, []);

//   useEffect(() => {
//     if (!isSidebarOpen) return;
//     const onMouseDown = (e: MouseEvent) => {
//       const el = sidebarRef.current;
//       if (!el) return;
//       if (e.target instanceof Node && !el.contains(e.target))
//         setIsSidebarOpen(false);
//     };
//     document.addEventListener("mousedown", onMouseDown);
//     return () => document.removeEventListener("mousedown", onMouseDown);
//   }, [isSidebarOpen]);

//   // Soft sidebar style (mirip desain gambar)
//   const navItemClass = ({ isActive }: { isActive: boolean }) =>
//     [
//       "group flex items-center gap-3 rounded-xl px-3 py-2 transition",
//       "text-[#929f98] hover:text-slate-900",
//       "hover:bg-slate-100/80",
//       isActive ? "bg-white shadow-sm ring-1 ring-slate-200 text-slate-900" : "",
//     ].join(" ");

//   const subNavItemClass = ({ isActive }: { isActive: boolean }) =>
//     [
//       "group flex items-center gap-3 rounded-xl px-3 py-2 transition",
//       "text-[#929f98] hover:text-slate-900",
//       "hover:bg-slate-100/80",
//       isActive ? "bg-white shadow-sm ring-1 ring-slate-200 text-slate-900" : "",
//     ].join(" ");

//   const sectionLabelClass =
//     "px-2 pt-4 pb-2 text-xs font-semibold text-slate-400";

//   return (
//     <>
//       {/* Mobile: open button */}
//       <button
//         type="button"
//         onClick={() => setIsSidebarOpen(true)}
//         aria-controls="admin-sidebar"
//         aria-expanded={isSidebarOpen}
//         className="sm:hidden inline-flex items-center gap-2 rounded-xl px-3 py-2 text-sm font-medium
//                    text-slate-700 ring-1 ring-slate-200 bg-white hover:bg-slate-50"
//       >
//         <span className="sr-only">Open sidebar</span>
//         <svg
//           className="h-5 w-5"
//           viewBox="0 0 24 24"
//           fill="none"
//           aria-hidden="true"
//         >
//           <path
//             d="M5 7h14M5 12h14M5 17h10"
//             stroke="currentColor"
//             strokeWidth="2"
//             strokeLinecap="round"
//           />
//         </svg>
//         Menu
//       </button>

//       {/* Overlay (mobile) */}
//       <div
//         className={[
//           "fixed inset-0 z-30 bg-black/30 backdrop-blur-[1px] transition-opacity sm:hidden",
//           isSidebarOpen
//             ? "opacity-100 pointer-events-auto"
//             : "opacity-0 pointer-events-none",
//         ].join(" ")}
//         onClick={() => setIsSidebarOpen(false)}
//         aria-hidden="true"
//       />

//       {/* Sidebar */}
//       <aside
//         id="admin-sidebar"
//         aria-label="Sidebar"
//         ref={sidebarRef}
//         className={[
//           "fixed top-0 left-0 z-40 h-screen w-[280px]",
//           "transition-transform duration-200 ease-out",
//           "sm:translate-x-0",
//           isSidebarOpen ? "translate-x-0" : "-translate-x-full",
//           "p-4",
//         ].join(" ")}
//       >
//         {/* Outer card container (rounded + shadow) */}
//         <div
//           className={[
//             "h-full w-full rounded-3xl",
//             "bg-white ring-1 ring-slate-200",
//             "shadow-[0_10px_30px_-18px_rgba(0,0,0,0.25)]",
//             "flex flex-col overflow-hidden",
//           ].join(" ")}
//         >
//           {/* Header */}
//           <div className="px-5 pt-5 pb-4">
//             <div className="flex items-center gap-3">
//               <div className="h-10 w-10 rounded-2xl  flex items-center justify-center shadow-sm">
//                 {/* App icon */}
//                 <img
//                   src="/Images/assetUpload/logo-fi-bg.webp"
//                   className="object-cover w-full h-full"
//                   alt=""
//                 />
//               </div>
//               <div className="min-w-0">
//                 <div className="text-slate-900 font-semibold leading-tight">
//                   SMA IT Fitrah Insani
//                 </div>
//                 <div className="text-xs text-slate-500">Panel</div>
//               </div>
//               <button
//                 type="button"
//                 onClick={() => setIsSidebarOpen(false)}
//                 className="sm:hidden ml-auto rounded-xl p-2 hover:bg-slate-200/60 text-slate-600"
//               >
//                 <span className="sr-only">Close sidebar</span>
//                 <svg
//                   className="h-5 w-5"
//                   viewBox="0 0 24 24"
//                   fill="none"
//                   aria-hidden="true"
//                 >
//                   <path
//                     d="M6 6l12 12M18 6L6 18"
//                     stroke="currentColor"
//                     strokeWidth="2"
//                     strokeLinecap="round"
//                   />
//                 </svg>
//               </button>
//             </div>

//             {/* Search (mirip top search feel) */}
//             <div className="mt-4">
//               <div className="flex items-center gap-2 rounded-2xl bg-white ring-1 ring-slate-200 px-3 py-2">
//                 <svg
//                   className="h-4 w-4 text-slate-400"
//                   viewBox="0 0 24 24"
//                   fill="none"
//                   aria-hidden="true"
//                 >
//                   <path
//                     d="M21 21l-4.3-4.3m1.8-5.2a7 7 0 11-14 0 7 7 0 0114 0z"
//                     stroke="currentColor"
//                     strokeWidth="2"
//                     strokeLinecap="round"
//                   />
//                 </svg>
//                 <input
//                   className="w-full bg-transparent text-sm text-slate-700 placeholder:text-slate-400 outline-none"
//                   placeholder="Search here..."
//                 />
//               </div>
//             </div>
//           </div>

//           {/* Menu scroll area */}
//           <div className="flex-1 overflow-y-auto px-3 pb-3">
//             <div className={sectionLabelClass}>MENU</div>

//             <ul className="space-y-1">
//               <li>
//                 <NavLink
//                   to="/dashboard/administrator"
//                   end
//                   className={navItemClass}
//                 >
//                   <svg
//                     className="h-5 w-5"
//                     viewBox="0 0 24 24"
//                     fill="none"
//                     aria-hidden="true"
//                   >
//                     <path
//                       d="M10 6.025A7.5 7.5 0 1 0 17.975 14H10V6.025Z"
//                       stroke="currentColor"
//                       strokeWidth="2"
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                     />
//                     <path
//                       d="M13.5 3c-.169 0-.334.014-.5.025V11h7.975c.011-.166.025-.331.025-.5A7.5 7.5 0 0 0 13.5 3Z"
//                       stroke="currentColor"
//                       strokeWidth="2"
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                     />
//                   </svg>
//                   <span className="text-sm font-medium">Dashboard</span>
//                 </NavLink>
//               </li>

//               {/* Kelola Pengguna */}
//               <li>
//                 <button
//                   type="button"
//                   onClick={() => setIsEcomOpen((v) => !v)}
//                   className={[
//                     "w-full group flex items-center justify-between gap-3 rounded-xl px-3 py-2 transition",
//                     "text-slate-600 hover:text-slate-900 hover:bg-slate-100/80",
//                     isKelolaPenggunaActive
//                       ? "bg-white shadow-sm ring-1 ring-slate-200 text-slate-900"
//                       : "",
//                   ].join(" ")}
//                   aria-controls="dropdown-kelola-pengguna"
//                   aria-expanded={isEcomOpen}
//                 >
//                   <span className="flex items-center gap-3">
//                     <svg
//                       className="h-5 w-5"
//                       viewBox="0 0 24 24"
//                       fill="none"
//                       aria-hidden="true"
//                     >
//                       <path
//                         d="M16 19h4a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-2m-2.236-4a3 3 0 1 0 0-4M3 18v-1a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v1a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1Zm8-10a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
//                         stroke="currentColor"
//                         strokeWidth="2"
//                         strokeLinecap="round"
//                         strokeLinejoin="round"
//                       />
//                     </svg>
//                     <span className="text-sm font-medium">Kelola Pengguna</span>
//                   </span>

//                   <svg
//                     className={[
//                       "h-5 w-5 transition-transform",
//                       isEcomOpen ? "rotate-180" : "rotate-0",
//                     ].join(" ")}
//                     viewBox="0 0 24 24"
//                     fill="none"
//                     aria-hidden="true"
//                   >
//                     <path
//                       d="m19 9-7 7-7-7"
//                       stroke="currentColor"
//                       strokeWidth="2"
//                       strokeLinecap="round"
//                     />
//                   </svg>
//                 </button>

//                 <ul
//                   id="dropdown-kelola-pengguna"
//                   className={[
//                     "mt-1 space-y-1 pl-7 border-l border-slate-200 ml-3",
//                     isEcomOpen ? "block" : "hidden",
//                   ].join(" ")}
//                 >
//                   <li>
//                     <NavLink
//                       to={kelolaPenggunaRoutes.guru}
//                       className={subNavItemClass}
//                     >
//                       <span className="h-2 w-2 rounded-full bg-slate-300 group-hover:bg-slate-500" />
//                       <span className="text-sm font-medium">Akun Guru</span>
//                     </NavLink>
//                   </li>
//                   <li>
//                     <NavLink
//                       to={kelolaPenggunaRoutes.siswa}
//                       className={subNavItemClass}
//                     >
//                       <span className="h-2 w-2 rounded-full bg-slate-300 group-hover:bg-slate-500" />
//                       <span className="text-sm font-medium">Akun Siswa</span>
//                     </NavLink>
//                   </li>
//                 </ul>
//               </li>

//               {/* Data Master */}
//               <li>
//                 <button
//                   type="button"
//                   onClick={() => setIsDataMasterOpen((v) => !v)}
//                   className={[
//                     "w-full group flex items-center justify-between gap-3 rounded-xl px-3 py-2 transition",
//                     "text-slate-600 hover:text-slate-900 hover:bg-slate-100/80",
//                     isDataMasterOpen
//                       ? "bg-white shadow-sm ring-1 ring-slate-200 text-slate-900"
//                       : "",
//                   ].join(" ")}
//                   aria-controls="dropdown-data-master"
//                   aria-expanded={isDataMasterOpen}
//                 >
//                   <span className="flex items-center gap-3">
//                     <svg
//                       className="h-5 w-5"
//                       viewBox="0 0 24 24"
//                       fill="none"
//                       aria-hidden="true"
//                     >
//                       <path
//                         d="M15 5v14M9 5v14M4 5h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1Z"
//                         stroke="currentColor"
//                         strokeWidth="2"
//                         strokeLinecap="round"
//                         strokeLinejoin="round"
//                       />
//                     </svg>
//                     <span className="text-sm font-medium">Data Master</span>
//                   </span>

//                   <svg
//                     className={[
//                       "h-5 w-5 transition-transform",
//                       isDataMasterOpen ? "rotate-180" : "rotate-0",
//                     ].join(" ")}
//                     viewBox="0 0 24 24"
//                     fill="none"
//                     aria-hidden="true"
//                   >
//                     <path
//                       d="m19 9-7 7-7-7"
//                       stroke="currentColor"
//                       strokeWidth="2"
//                       strokeLinecap="round"
//                     />
//                   </svg>
//                 </button>

//                 {/* Submenu Data Master kalau sudah ada tinggal ditaruh di sini (modelnya sama kayak Kelola Pengguna) */}
//               </li>

//               {/* Ujian */}
//               <li>
//                 <NavLink
//                   to="/dashboard/administrator/ujian"
//                   className={navItemClass}
//                 >
//                   <svg
//                     className="h-5 w-5"
//                     viewBox="0 0 24 24"
//                     fill="none"
//                     aria-hidden="true"
//                   >
//                     <path
//                       d="M12 6.03v13m0-13c-2.819-.831-4.715-1.076-8.029-1.023A.99.99 0 0 0 3 6v11c0 .563.466 1.014 1.03 1.007 3.122-.043 5.018.212 7.97 1.023m0-13c2.819-.831 4.715-1.076 8.029-1.023A.99.99 0 0 1 21 6v11c0 .563-.466 1.014-1.03 1.007-3.122-.043-5.018.212-7.97 1.023"
//                       stroke="currentColor"
//                       strokeWidth="2"
//                       strokeLinecap="round"
//                       strokeLinejoin="round"
//                     />
//                   </svg>
//                   <span className="text-sm font-medium">Ujian</span>
//                 </NavLink>
//               </li>
//             </ul>
//           </div>

//           {/* Bottom section */}
//           <div className="px-3 pb-4 pt-3 border-t border-slate-200">
//             <NavLink
//               to="/dashboard/administrator/pengaturan"
//               className={navItemClass}
//             >
//               <svg
//                 className="h-5 w-5"
//                 viewBox="0 0 24 24"
//                 fill="none"
//                 aria-hidden="true"
//               >
//                 <path
//                   d="M19.14 12.94c.04-.3.06-.61.06-.94c0-.32-.02-.64-.07-.94l2.03-1.58a.49.49 0 0 0 .12-.61l-1.92-3.32a.49.49 0 0 0-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.484.484 0 0 0-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58a.49.49 0 0 0-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6s3.6 1.62 3.6 3.6s-1.62 3.6-3.6 3.6"
//                   fill="currentColor"
//                 />
//               </svg>
//               <span className="text-sm font-medium">Pengaturan</span>
//             </NavLink>
//           </div>
//         </div>
//       </aside>
//     </>
//   );
// };





import { useEffect, useMemo, useRef, useState } from "react";
import type { ReactNode } from "react";
import { NavLink, useLocation, matchPath } from "react-router";

type SidebarMenuItem = {
  id: string;
  label: string;
  icon: ReactNode;
  to?: string;
  end?: boolean;
  children?: SidebarMenuItem[];
};

const MENU_ITEMS: SidebarMenuItem[] = [
  {
    id: "dashboard",
    label: "Dashboard",
    to: "/dashboard/administrator",
    end: true,
    icon: (
      <svg
        className="w-5 h-5 transition duration-75"
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
    label: "Kelola Pengguna",
    icon: (
      <svg
        className="shrink-0 w-5 h-5 transition duration-75"
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
        to: "/dashboard/administrator/kelola-akun/guru",
        icon: (
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
              strokeWidth="2"
              d="M7 17v1a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1a3 3 0 0 0-3-3h-4a3 3 0 0 0-3 3Zm8-9a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
            />
          </svg>
        ),
      },
      {
        id: "akun-siswa",
        label: "Akun Siswa",
        to: "/dashboard/administrator/kelola-akun/siswa",
        icon: (
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
    label: "Data Master",
    icon: (
      <svg
        className="h-5 w-5 shrink-0 transition duration-75"
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
    children: [],
  },
  {
    id: "ujian",
    label: "Ujian",
    to: "/dashboard/administrator/ujian",
    icon: (
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
    ),
  },
];

const BOTTOM_MENU_ITEMS: SidebarMenuItem[] = [
  {
    id: "pengaturan",
    label: "Pengaturan",
    to: "/dashboard/administrator/pengaturan",
    icon: (
      <svg
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

export const SidebarAdmin = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [openMenus, setOpenMenus] = useState<Record<string, boolean>>({});
  const sidebarRef = useRef<HTMLElement | null>(null);

  const { pathname } = useLocation();

  const menuItems = MENU_ITEMS;
  const bottomMenuItems = BOTTOM_MENU_ITEMS;

  const activeMenuIds = useMemo(() => {
    const activeIds = new Set<string>();

    const isRouteActive = (route?: string, end = false) =>
      route ? !!matchPath({ path: route, end }, pathname) : false;

    const checkItem = (item: SidebarMenuItem) => {
      if (isRouteActive(item.to, item.end)) {
        activeIds.add(item.id);
      }

      item.children?.forEach((child) => {
        if (isRouteActive(child.to, child.end)) {
          activeIds.add(item.id);
          activeIds.add(child.id);
        }
      });
    };

    menuItems.forEach(checkItem);
    bottomMenuItems.forEach(checkItem);

    return activeIds;
  }, [menuItems, bottomMenuItems, pathname]);

  // Auto open dropdown jika sedang berada di dalam section itu
  useEffect(() => {
    setOpenMenus((prev) => {
      const next = { ...prev };
      menuItems.forEach((item) => {
        if (item.children && activeMenuIds.has(item.id)) {
          next[item.id] = true;
        }
      });
      return next;
    });
  }, [activeMenuIds, menuItems]);

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

  // === Sidebar theme tokens (seragam) ===
  const ITEM_BASE =
    "flex items-center w-full px-2 py-1.5 rounded-base group transition-colors focus:outline-none focus:ring-4 focus:ring-white/10";
  const ITEM_TEXT = "text-white/80";
  const ITEM_HOVER = "hover:bg-white/10 hover:text-white";
  const ITEM_ACTIVE = "bg-white/15 text-white font-semibold";

  const SUB_BASE =
    "pl-8 flex items-center gap-3 px-2 py-1.5 rounded-base transition-colors focus:outline-none focus:ring-4 focus:ring-white/10";
  const SUB_TEXT = "text-white/70";
  const SUB_HOVER = "hover:bg-white/10 hover:text-white";
  const SUB_ACTIVE = "bg-white/10 text-white font-semibold";

  const navItemClass = ({ isActive }: { isActive: boolean }) =>
    [ITEM_BASE, ITEM_TEXT, ITEM_HOVER, isActive ? ITEM_ACTIVE : ""].join(" ");

  const subNavItemClass = ({ isActive }: { isActive: boolean }) =>
    [SUB_BASE, SUB_TEXT, SUB_HOVER, isActive ? SUB_ACTIVE : ""].join(" ");

  // Optional: samakan ukuran icon via wrapper (biar rapi)
  const IconWrap = ({ children }: { children: ReactNode }) => (
    <span className="shrink-0 w-5 h-5 text-current [&_svg]:w-5 [&_svg]:h-5">
      {children}
    </span>
  );

  const renderMenuItem = (item: SidebarMenuItem) => {
    const isActive = activeMenuIds.has(item.id);
    const hasChildren = (item.children?.length ?? 0) > 0;
    const isOpen = !!openMenus[item.id];

    if (hasChildren) {
      return (
        <li key={item.id}>
          <button
            type="button"
            onClick={() =>
              setOpenMenus((prev) => ({
                ...prev,
                [item.id]: !prev[item.id],
              }))
            }
            className={[
              ITEM_BASE,
              ITEM_TEXT,
              ITEM_HOVER,
              isActive ? ITEM_ACTIVE : "",
              "justify-between",
            ].join(" ")}
            aria-controls={`dropdown-${item.id}`}
            aria-expanded={isOpen}
          >
            <IconWrap>{item.icon}</IconWrap>
            <span className="flex-1 ms-3 text-left rtl:text-right whitespace-nowrap">
              {item.label}
            </span>

            <svg
              className={[
                "w-5 h-5 transition-transform duration-200",
                isOpen ? "rotate-180" : "rotate-0",
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

          {item.children && (
            <ul
              id={`dropdown-${item.id}`}
              className={["py-2 space-y-2", isOpen ? "block" : "hidden"].join(
                " "
              )}
            >
              {item.children.map((child) => (
                <li key={child.id}>
                  <NavLink to={child.to ?? ""} className={subNavItemClass}>
                    <IconWrap>{child.icon}</IconWrap>
                    <span className="ms-2">{child.label}</span>
                  </NavLink>
                </li>
              ))}
            </ul>
          )}
        </li>
      );
    }

    return (
      <li key={item.id}>
        <NavLink to={item.to ?? ""} end={item.end} className={navItemClass}>
          <IconWrap>{item.icon}</IconWrap>
          <span className="ms-3">{item.label}</span>
        </NavLink>
      </li>
    );
  };

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
          "fixed top-0 left-0 z-40 w-64 h-screen rounded-r-",
          "transition-transform duration-200 ease-out",
          "sm:translate-x-0",
          isSidebarOpen ? "translate-x-0" : "-translate-x-full",
          "flex flex-col overflow-hidden",
        ].join(" ")}
      >
        {/* BODY */}
        <div className="flex-1 min-h-0 px-3 py-5 bg-[#37513d] border-e border-black/10 flex flex-col">
          {/* Header mobile */}
          <div className="flex items-center justify-between sm:hidden mb-3 shrink-0">
            <span className="text-sm font-medium text-white">Menu</span>
            <button
              type="button"
              onClick={() => setIsSidebarOpen(false)}
              className="p-2 rounded-base hover:bg-white/10 focus:outline-none focus:ring-4 focus:ring-white/20 text-white"
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
            {/* HEADER / LOGO */}
            <div className="px-2 pt-2 pb-3">
              <div className="flex items-center gap-3">
                <div className="h-12 w-12 rounded-xl ring-1 ring-white/15 flex items-center justify-center shrink-0">
                  <img
                    src="/Images/assetUpload/logo-fi-bg.webp"
                    alt="Logo sekolah"
                    className="h-12 w-12 object-cover rounded-xl"
                  />
                </div>

                <div className="min-w-0">
                  <h2 className="text-white font-semibold text-sm leading-tight wrap-break-word">
                    SMA IT Fitrah Insani
                  </h2>
                  <p className="text-white/70 text-xs mt-0.5">Panel</p>
                </div>
              </div>

              <div className="mt-3 h-px bg-white/10" />
            </div>

            <ul className="space-y-2 font-medium flex flex-col gap-3 py-4">
              {menuItems.map(renderMenuItem)}
            </ul>
          </div>

          {/* Divider section */}
          <div className="shrink-0 mt-4 border-t border-white/15 pt-4">
            <ul className="space-y-2 font-medium">
              {bottomMenuItems.map(renderMenuItem)}
            </ul>
          </div>
        </div>
      </aside>
    </>
  );
};
