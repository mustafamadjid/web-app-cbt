import { AkunSiswaForm } from "@/layouts/Form/Admin/KelolaAkun/AkunSiswaForm"
import { Header } from "@/components/features/widget/Header/HeaderWidget"
export const TambahSiswa = () => {
    return (
        <>
            <div className="p-0">
                      <div className="p-0 flex flex-col">
                        {/* Header*/}
                          <Header />
            
                        {/* Konten form tambah guru*/}
                        <AkunSiswaForm />
                      </div>
                    </div>
        </>
    )
}