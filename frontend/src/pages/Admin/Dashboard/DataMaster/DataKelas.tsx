import { useNavigate } from "react-router";

import { DataKelasTables } from "@/components/features/tables/DataMasterTables/DataKelasTables";
import { AddButton } from "@/components/common/Button/AddButton";

export const DataKelas = () => {
  const navigate = useNavigate();
  return (
    <>
      <div className="px-8">
        <div className="flex justify-end py-6 ">
          <AddButton
            label="Tambah Kelas"
            onClick={() => {
              navigate("/dashboard/administrator/data-master/tambah-kelas");
            }}
          />
        </div>
        <DataKelasTables />
      </div>
    </>
  );
};
