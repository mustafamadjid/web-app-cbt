import { useNavigate } from "react-router";


import { DataMataPelajaran } from "@/components/features/tables/DataMasterTables/DataMataPelajaranTables";
import { AddButton } from "@/components/common/Button/AddButton";
export const MataPelajaran = () => {
    const navigate = useNavigate();
  return (
    <>


      <div className="px-8">
        <div className="flex justify-end py-6">
          <AddButton label="Tambah Mata Pelajaran" onClick={() => {
            navigate('/dashboard/administrator/data-master/tambah-mapel')
          }} />
        </div>

      <DataMataPelajaran />
      </div>
    </>
  );
};
