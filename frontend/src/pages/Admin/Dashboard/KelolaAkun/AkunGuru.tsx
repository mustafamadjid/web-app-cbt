// Component
import { Header } from "@/components/features/widget/Header/HeaderWidget";
import { AddButton } from "@/components/common/Button/AddButton";
import { AkunGuruTables } from "@/components/features/tables/AkunGuruTables/AkunGuruTables";

// React Link
import { useNavigate } from "react-router";

export const dummyHeaderUser = {
  title: "Dashboard",
  userName: "Administrator Sistem Pendidikan Nasional",
  roleLabel: "Admin",
  isOnline: true,
  avatarUrl: null,
};

// Nanti header ini akan diisi lewat context supaya state datanya global

export const KelolaAkunGuru = () => {
  const navigate = useNavigate();
  return (
    <>
      <div className="p-0">
        <div className="p-0 flex flex-col">
          {/* Header*/}
          <Header
            title={dummyHeaderUser.title}
            userName={dummyHeaderUser.userName}
            roleLabel={dummyHeaderUser.roleLabel}
            isOnline={dummyHeaderUser.isOnline}
            avatarUrl={dummyHeaderUser.avatarUrl}
          />

          {/* tombol Add User */}
          <div className=" flex justify-end px-4 py-6">
            <div>
              <AddButton
                label="Tambah Akun Guru"
                onClick={() =>
                  navigate(`/dashboard/administrator/kelola-akun/tambah-guru`)
                }
              />
            </div>
          </div>
          {/* User Tables */}
          <div className="px-4">
            <AkunGuruTables />
          </div>
        </div>
      </div>
    </>
  );
};
