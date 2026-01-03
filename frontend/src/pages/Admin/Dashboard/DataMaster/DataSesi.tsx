import { useNavigate } from "react-router";

import { DataSesiTables } from "@/components/features/tables/DataMasterTables/DataSesiTables";
import { AddButton } from "@/components/common/Button/AddButton";


export const DataSesi = () => {
  const navigate = useNavigate();
  return (
    <>
    

      <div className="px-8">
        <div className="flex justify-end py-6">
          <AddButton
            label="Tambah Sesi"
            onClick={() => {
              navigate("/dashboard/administrator/data-master/tambah-sesi");
            }}
          />
        </div>

        <DataSesiTables />
      </div>
    </>
  );
};
