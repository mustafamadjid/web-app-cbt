import { Header } from "@/components/features/widget/Header/HeaderWidget";
import { Outlet } from "react-router";

// Nanti header ini akan diisi lewat context supaya state datanya global
export const HeaderLayout = () => {
    return (
      <div>
        <Header
          title="Dashboard"
          userName="Ruben Onsu Marpaung"
          roleLabel="Administrator"
          isOnline={true}
          avatarUrl={null}
        />
        <Outlet />
      </div>
    );
}