import { useEffect, useMemo, useRef, useState } from "react";
import { NavLink, useLocation, matchPath } from "react-router";

import type {
  SidebarMenuGroupItem,
  SidebarMenuItem,
  SidebarSubMenuItem,
} from "../../../types/Sidebar/SidebarMenu";
import {
  footerMenuItems,
  mainMenuItems,
} from "./sidebarMenuItems";

export const SidebarAdmin = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>({});
  const sidebarRef = useRef<HTMLElement | null>(null);

  const { pathname } = useLocation();

  const activeGroupMap = useMemo(() => {
    const activeGroups: Record<string, boolean> = {};

    mainMenuItems.forEach((item) => {
      if (item.type !== "group") return;

      const hasActiveChild = item.children.some((child) =>
        matchPath({ path: child.to, end: false }, pathname)
      );

      if (hasActiveChild) activeGroups[item.id] = true;
    });

    return activeGroups;
  }, [pathname]);

  useEffect(() => {
    if (Object.keys(activeGroupMap).length === 0) return;
    setOpenGroups((prev) => ({ ...prev, ...activeGroupMap }));
  }, [activeGroupMap]);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setIsSidebarOpen(false);
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

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

  /**
   * Warna sidebar gelap -> semua menu harus punya text putih konsisten.
   * State:
   * - default: text-white/80
   * - hover: bg putih tipis, text putih full
   * - active: bg putih lebih kuat, text putih full, font-semibold
   */
  const navItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "flex items-center px-2 py-2 rounded-base group transition-colors",
      "text-white/80",
      "hover:bg-white/10 hover:text-white",
      isActive ? "bg-white/15 text-white font-semibold" : "",
    ].join(" ");

  const subNavItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "pl-10 flex items-center px-2 py-2 rounded-base group transition-colors",
      "text-white/70 text-sm",
      "hover:bg-white/10 hover:text-white",
      isActive ? "bg-white/15 text-white font-semibold" : "",
    ].join(" ");

  // Icon class helper supaya ikon ikut konsisten warnanya
  const iconClass = (active?: boolean) =>
    [
      "shrink-0 w-5 h-5 transition-colors duration-150",
      active ? "text-white" : "text-white/75",
      "group-hover:text-white",
    ].join(" ");

  const subIconClass = (active?: boolean) =>
    [
      "shrink-0 w-5 h-5 transition-colors duration-150",
      active ? "text-white" : "text-white/70",
      "group-hover:text-white",
    ].join(" ");

  const renderMenuItem = (item: SidebarMenuItem) => {
    if (item.type === "link") {
      return (
        <li key={item.id}>
          <NavLink to={item.to} end={item.end} className={navItemClass}>
            {({ isActive }) => (
              <>
                {item.icon(iconClass(isActive))}
                <span className="ms-3">{item.label}</span>
              </>
            )}
          </NavLink>
        </li>
      );
    }

    const groupItem = item as SidebarMenuGroupItem;
    const isGroupOpen = !!openGroups[groupItem.id];
    const isGroupActive = !!activeGroupMap[groupItem.id];
    const groupIconActive = isGroupActive || isGroupOpen;

    return (
      <li key={groupItem.id}>
        <button
          type="button"
          onClick={() =>
            setOpenGroups((prev) => ({
              ...prev,
              [groupItem.id]: !isGroupOpen,
            }))
          }
          className={[
            "flex items-center w-full justify-between px-2 py-2 rounded-base group transition-colors",
            "text-white/80",
            "hover:bg-white/10 hover:text-white",
            groupIconActive ? "bg-white/15 text-white font-semibold" : "",
          ].join(" ")}
          aria-controls={`dropdown-${groupItem.id}`}
          aria-expanded={isGroupOpen}
        >
          <span className="flex items-center gap-3">
            {groupItem.icon(iconClass(groupIconActive))}
            <span className="whitespace-nowrap">{groupItem.label}</span>
          </span>

          <svg
            className={[
              "w-5 h-5 transition-transform duration-200",
              "text-white/80 group-hover:text-white",
              isGroupOpen ? "rotate-180" : "rotate-0",
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
          id={`dropdown-${groupItem.id}`}
          className={["py-2 space-y-2", isGroupOpen ? "block" : "hidden"].join(
            " "
          )}
        >
          {groupItem.children.map((child) =>
            renderSubMenuItem(groupItem, child)
          )}
        </ul>
      </li>
    );
  };

  const renderSubMenuItem = (
    groupItem: SidebarMenuGroupItem,
    child: SidebarSubMenuItem
  ) => (
    <li key={`${groupItem.id}-${child.id}`}>
      <NavLink to={child.to} className={subNavItemClass}>
        {({ isActive }) => (
          <>
            {child.icon(subIconClass(isActive))}
            <span className="ms-2">{child.label}</span>
          </>
        )}
      </NavLink>
    </li>
  );

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

          {/* MENU UTAMA (scroll) */}
          <div className="flex-1 min-h-0 overflow-y-auto pb-4">
            {/* HEADER / LOGO */}
            <div className="px-2 pt-2 pb-3">
              <div className="flex items-center gap-3">
                <div className="h-12 w-12 rounded-xl ring-1 ring-white/15 flex items-center justify-center shrink-0 overflow-hidden">
                  <img
                    src="/Images/assetUpload/logo-fi-bg.webp"
                    alt="Logo sekolah"
                    className="h-12 w-12 object-cover"
                  />
                </div>

                <div className="min-w-0">
                  <h2 className="text-white font-semibold text-sm leading-tight break-words">
                    SMA IT Fitrah Insani
                  </h2>
                  <p className="text-white/70 text-xs mt-0.5">Panel</p>
                </div>
              </div>

              <div className="mt-3 h-px bg-white/10" />
            </div>

            <ul className="space-y-2 font-medium flex flex-col gap-3 py-4">
              {mainMenuItems.map((item) => renderMenuItem(item))}
            </ul>
          </div>

          {/* Divider section */}
          <div className="shrink-0 mt-4 border-t border-white/10 pt-4">
            <ul className="space-y-2 font-medium">
              {footerMenuItems.map((item) => renderMenuItem(item))}
            </ul>
          </div>
        </div>
      </aside>
    </>
  );
};
