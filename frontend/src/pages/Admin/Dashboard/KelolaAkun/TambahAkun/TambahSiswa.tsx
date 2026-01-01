import { AkunSiswaForm } from "@/layouts/Form/Admin/KelolaAkun/AkunSiswaForm"

export const TambahSiswa = () => {
    return (
      <>
        <div className="p-0">
          <div className="p-0 flex flex-col">
            {/* Konten form tambah guru*/}
            <AkunSiswaForm />
          </div>
        </div>
      </>
    );
}