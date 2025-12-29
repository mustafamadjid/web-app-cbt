// Component
import { Header } from "@/components/features/widget/Header/HeaderWidget";
import { AddButton } from "@/components/common/Button/AddButton";
import { AkunGuruTables } from "@/components/features/tables/AkunGuruTables/AkunGuruTables";

// React Link
import { useNavigate } from "react-router";



export const KelolaAkunGuru = () => {
    const navigate = useNavigate();
  return (
    <>
      <div className="p-0">
        <div className="p-0 flex flex-col">
          {/* Header*/}
            <Header />

          {/* tombol Add User */}
          <div className=" flex justify-end px-4 py-6">
            <div>
                <AddButton label="Tambah Akun Guru" onClick={() => navigate(`/dashboard/administrator/kelola-akun/tambah-guru`)} />
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
