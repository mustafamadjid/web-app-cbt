import type React from "react";
export type SidebarIconRenderer = (className: string) => React.JSX.Element;

export type Role = "ADMIN" | "GURU" | "SISWA";

export type SidebarMenuItem =
  | SidebarMenuLinkItem
  | SidebarMenuGroupItem;

export interface SidebarMenuItemBase {
  id: number;
  label: string;
  icon: SidebarIconRenderer;
  roles?: Role[];
}

export interface SidebarMenuLinkItem extends SidebarMenuItemBase {
  type: "link";
  to: string;
  end?: boolean;
}

export interface SidebarMenuGroupItem extends SidebarMenuItemBase {
  type: "group";
  children: SidebarSubMenuItem[];
}

export interface SidebarSubMenuItem {
  id: number;
  label: string;
  to: string;
  icon: SidebarIconRenderer;
  roles?: Role[];
}
