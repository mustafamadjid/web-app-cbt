import type { Role,SidebarMenuItem } from "@/types/Sidebar/SidebarMenu";

const hasAccess = (roles:Role[] | undefined, role:Role) => {
    if (!roles || roles.length === 0) return true;
    return roles.includes(role);
}

export function filterMenuByRole(
    items: SidebarMenuItem[],
    role: Role
): SidebarMenuItem[] {
    return items.
    map((item) => {
        if(item.type === "group") {
            const children = item.children.filter((child) => hasAccess(child.roles, role));
            if(!hasAccess(item.roles, role) || children.length === 0) return null;
            return {
                ...item,
                children
            }
        }
        if(!hasAccess(item.roles, role)) return null;
        return item;
    }).filter(Boolean) as SidebarMenuItem[];
}

