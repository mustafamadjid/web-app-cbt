import { useNavigate } from "react-router";

import { DataKelasTables } from "@/components/features/tables/DataMasterTables/DataKelasTables";
import { AddButton } from "@/components/common/Button/AddButton";
import { Header } from "@/components/features/widget/Header/HeaderWidget";

export const DataKelas = () => {
  const navigate = useNavigate();
  return (
    <>
      <Header />
      <div className="flex justify-end py-6 px-4">
        <AddButton
          label="Tambah Kelas"
          onClick={() => {
            navigate("/dashboard/administrator/data-master/tambah-kelas");
          }}
        />
      </div>
      <div className="flex flex-col px-5">
        <DataKelasTables />
      </div>
    </>
  );
};
