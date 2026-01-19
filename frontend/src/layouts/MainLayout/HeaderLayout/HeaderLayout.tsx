import Header from "@/components/features/widget/Header/HeaderWidget";
import { useAuth } from "@/contexts/AuthContext";
import type { Role } from "@/types/Sidebar/SidebarMenu";
import { Outlet } from "react-router";

// Nanti header ini akan diisi lewat context supaya state datanya global
const HeaderLayout = () => {
  const {user} = useAuth();
    return (
      <div>
        <Header
          title="Dashboard"
          userName={user?.username ?? ""}
          roleLabel={user?.role as Role}
          isOnline={true}
          avatarUrl={null}
        />
        <Outlet />
      </div>
    );
}

export default HeaderLayout;
