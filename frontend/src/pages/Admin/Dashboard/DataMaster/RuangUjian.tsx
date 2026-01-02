import { useNavigate } from "react-router";

import { RuangUjianTables } from "@/components/features/tables/DataMasterTables/RuangUjianTables";
import { AddButton } from "@/components/common/Button/AddButton";
import { Header } from "@/components/features/widget/Header/HeaderWidget";

export const RuangUjian = () => {
  const navigate = useNavigate();
  return (
    <>
      <Header className="px-4" />

      <div className="px-4">
        <div className="flex justify-end py-6">
          <AddButton
            label="Tambah Ruang"
            onClick={() => {
              navigate("/dashboard/administrator/data-master/tambah-ruang");
            }}
          />
        </div>

        <RuangUjianTables />
      </div>
    </>
  );
};
