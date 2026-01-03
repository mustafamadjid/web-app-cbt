// Component
import { AddButton } from "@/components/common/Button/AddButton";
import { AkunSiswaTables } from "@/components/features/tables/AkunSiswaTables/AkunSiswaTables";

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

export const KelolaAkunSiswa = () => {
  const navigate = useNavigate();

  return (
    <div className="w-full min-w-0">
      <div className="flex flex-col w-full min-w-0">
        {/* tombol Add User */}
        <div className="flex justify-end px-8 py-6 shrink-0">
          <AddButton
            label="Tambah Akun Siswa"
            onClick={() =>
              navigate(`/dashboard/administrator/kelola-akun/tambah-siswa`)
            }
          />
        </div>

        {/* User Tables */}
        <div className="px-8 pb-6 flex-1 min-w-0 w-full">
          {/* wrapper scroll horizontal */}
          <div className="w-full min-w-0 overflow-x-auto">
            <AkunSiswaTables />
          </div>
        </div>
      </div>
    </div>
  );
};
