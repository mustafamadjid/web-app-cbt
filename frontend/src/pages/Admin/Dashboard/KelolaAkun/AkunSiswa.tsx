
import AkunSiswaTables from "@/components/features/tables/AkunSiswaTables/AkunSiswaTables";


export const dummyHeaderUser = {
  title: "Dashboard",
  userName: "Administrator Sistem Pendidikan Nasional",
  roleLabel: "Admin",
  isOnline: true,
  avatarUrl: null,
};

// Nanti header ini akan diisi lewat context supaya state datanya global

export const KelolaAkunSiswa = () => {


  return (
    <div className="px-8 py-13">
      <AkunSiswaTables />
    </div>
  );
};
