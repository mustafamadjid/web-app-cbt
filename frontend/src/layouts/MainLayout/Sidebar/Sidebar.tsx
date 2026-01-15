import { useEffect, useMemo, useRef, useState } from "react";
import { NavLink, useLocation, matchPath } from "react-router";

import { filterMenuByRole } from "@/helper/FilterMenuSidebar/sidebarMenuFilter";
import type { Role} from "../../../types/Sidebar/SidebarMenu";
import { useAuth } from "@/contexts/AuthContext";

import SvgIcons from "@/assets/SvgIcons/svgIcons";

import type {
  SidebarMenuGroupItem,
  SidebarMenuItem,
  SidebarSubMenuItem,
} from "../../../types/Sidebar/SidebarMenu";
import { footerMenuItems, mainMenuItems } from "./sidebarMenuItems";

type SidebarProps = {
  isOpen: boolean;
  onToggle: () => void;
  onClose: () => void;
  onOpen: () => void;
};


const Sidebar = ({ isOpen, onToggle, onClose }: SidebarProps) => {
  const { user } = useAuth();
  const role = (user?.role ?? "SISWA") as Role;


  const [openGroups, setOpenGroups] = useState<Record<number, boolean>>({});
  const sidebarRef = useRef<HTMLElement | null>(null);

  const { pathname } = useLocation();

  
const filteredMainMenuItems = useMemo(
  () => filterMenuByRole(mainMenuItems, role),
  [role]
);

const filteredFooterMenuItems = useMemo(
  () => filterMenuByRole(footerMenuItems, role),
  [role]
);

 const activeGroupMap = useMemo(() => {
   const activeGroups: Record<number, boolean> = {};

   filteredMainMenuItems.forEach((item) => {
     if (item.type !== "group") return;

     const hasActiveChild = item.children.some((child) =>
       matchPath({ path: child.to, end: false }, pathname)
     );

     if (hasActiveChild) activeGroups[item.id] = true;
   });

   return activeGroups;
 }, [pathname, filteredMainMenuItems]);


  useEffect(() => {
    if (Object.keys(activeGroupMap).length === 0) return;
    setOpenGroups((prev) => ({ ...prev, ...activeGroupMap }));
  }, [activeGroupMap]);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [onClose]);

  useEffect(() => {
    if (!isOpen) return;

    const mql = window.matchMedia("(max-width: 639px)");

    const onMouseDown = (e: MouseEvent) => {
      // Desktop: abaikan
      if (!mql.matches) return;

      const el = sidebarRef.current;
      if (!el) return;

      if (e.target instanceof Node && !el.contains(e.target)) {
        onClose();
      }
    };

    document.addEventListener("mousedown", onMouseDown);
    return () => document.removeEventListener("mousedown", onMouseDown);
  }, [isOpen, onClose]);

  const navItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "flex items-center px-2 py-2 rounded-base group transition-colors cursor-pointer",
      "text-white/80",
      "hover:bg-white/10 hover:text-white",
      isActive ? "bg-white/15 text-white font-semibold" : "",
    ].join(" ");

  const subNavItemClass = ({ isActive }: { isActive: boolean }) =>
    [
      "pl-10 flex items-center px-2 py-2 rounded-base group transition-colors cursor-pointer",
      "text-white/70 text-sm",
      "hover:bg-white/10 hover:text-white",
      isActive ? "bg-white/15 text-white font-semibold" : "",
    ].join(" ");

  const iconClass = (active?: boolean) =>
    [
      "shrink-0 w-5 h-5 transition-colors duration-150 cursor-pointer",
      active ? "text-white" : "text-white/75",
      "group-hover:text-white",
    ].join(" ");

  const subIconClass = (active?: boolean) =>
    [
      "shrink-0 w-5 h-5 transition-colors duration-150 cursor-pointer",
      active ? "text-white" : "text-white/70",
      "group-hover:text-white",
    ].join(" ");

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
            "flex items-center w-full justify-between px-2 py-2 rounded-base group transition-colors cursor-pointer",
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

          {SvgIcons.chevronDown(
            [
              "w-5 h-5 transition-transform duration-200",
              "text-white/80 group-hover:text-white",
              isGroupOpen ? "rotate-180" : "rotate-0",
            ].join(" ")
          )}
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

  return (
    <>
      {!isOpen && (
        <button
          type="button"
          onClick={onToggle}
          aria-controls="separator-sidebar"
          aria-expanded={isOpen}
          className={[
            "z-50 top-3 left-0 cursor-pointer",
            "text-heading bg-transparent border border-transparent",
            "hover:bg-neutral-secondary-medium focus:ring-4 focus:ring-neutral-tertiary",
            "font-medium rounded-base text-sm p-2 focus:outline-none inline-flex",
          ].join(" ")}
        >
          <span className="sr-only">Open sidebar</span>
          {SvgIcons.menu("w-6 h-6")}
        </button>
      )}

      {/* Overlay (mobile only) */}
      <div
        className={[
          "fixed inset-0 z-30 bg-black/40 transition-opacity sm:hidden ",
          isOpen
            ? "opacity-100 pointer-events-auto"
            : "opacity-0 pointer-events-none",
        ].join(" ")}
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Sidebar */}
      <aside
        id="separator-sidebar"
        aria-label="Sidebar"
        ref={sidebarRef}
        className={[
          "fixed top-0 left-0 z-40 w-64 h-screen rounded-r-sm",
          "transition-transform duration-200 ease-out",
          isOpen ? "translate-x-0" : "-translate-x-full",
          "flex flex-col overflow-hidden",
        ].join(" ")}
      >
        <div className="flex-1 min-h-0 px-3 py-5 bg-[#37513d] border-e border-black/10 flex flex-col">
          {/* Header mobile (boleh tetap, tapi sekarang juga aman untuk desktop) */}
          <div className="flex items-center justify-between  shrink-0">
            <span className=""></span>
            <button
              type="button"
              onClick={onClose}
              className="p-2 rounded-base cursor-pointer hover:bg-white/10 focus:outline-none focus:ring-4 focus:ring-white/20 text-white"
            >
              <span className="sr-only">Close sidebar</span>
              {SvgIcons.close("w-5 h-5")}
            </button>
          </div>

          <div className="flex-1 min-h-0 overflow-y-auto pb-4">
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
                  <h2 className="text-white font-semibold text-sm leading-tight wrap-break-word">
                    SMA IT Fitrah Insani
                  </h2>
                  <p className="text-white/70 text-xs mt-0.5">Panel</p>
                </div>
              </div>

              <div className="mt-3 h-px bg-white/10" />
            </div>

            <ul className="space-y-2 font-medium flex flex-col gap-3 py-4">
              {filteredMainMenuItems.map((item) => renderMenuItem(item))}
            </ul>
          </div>

          <div className="shrink-0 mt-4 border-t border-white/10 pt-4">
            <ul className="space-y-2 font-medium">
              {filteredFooterMenuItems.map((item) => renderMenuItem(item))}
            </ul>
          </div>
        </div>
      </aside>
    </>
  );
};

export default Sidebar;
