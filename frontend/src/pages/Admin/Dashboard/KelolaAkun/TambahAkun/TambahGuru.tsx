import { AkunGuruForm } from "@/layouts/Form/Admin/KelolaAkun/AkunGuruForm";

export const TambahGuru = () => {
    return (
      <>
        <div className="p-0">
          <div className="p-0 flex flex-col">
            {/* Konten form tambah guru*/}
            <AkunGuruForm />
          </div>
        </div>
      </>
    );
}