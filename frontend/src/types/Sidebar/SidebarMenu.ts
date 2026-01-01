export type SidebarIconRenderer = (className: string) => JSX.Element;

export type SidebarMenuItem =
  | SidebarMenuLinkItem
  | SidebarMenuGroupItem;

export interface SidebarMenuItemBase {
  id: string;
  label: string;
  icon: SidebarIconRenderer;
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
  id: string;
  label: string;
  to: string;
  icon: SidebarIconRenderer;
}
