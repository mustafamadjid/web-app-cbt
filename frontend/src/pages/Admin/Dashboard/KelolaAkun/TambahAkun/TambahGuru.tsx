import { AkunGuruForm } from "@/layouts/Form/Admin/KelolaAkun/AkunGuruForm";
import { Header } from "@/components/features/widget/Header/HeaderWidget";
export const TambahGuru = () => {
    return (
      <>
        <div className="p-0">
          <div className="p-0 flex flex-col">
            {/* Header*/}
              <Header />

            {/* Konten form tambah guru*/}
            <AkunGuruForm />
          </div>
        </div>
      </>
    );
}