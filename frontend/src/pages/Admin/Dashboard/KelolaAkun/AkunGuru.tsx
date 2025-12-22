

// Component
import { Header } from "@/components/features/widget/Header/HeaderWidget"
import { AddButton } from "@/components/common/Button/AddButton"
import { UserTables } from "@/components/features/tables/UserTables/UserTables"

export const KelolaAkunGuru = () => {


    return (
        <>
        <div className="p-0h-screen ">
            <div className="p-0 flex flex-col">
                {/* Header*/}
                <header className="bg-white p-5 ">
                    <Header />
                </header>
        
                {/* tombol Add User */}
                <div className=" flex justify-end px-4 py-6">
                    <div>
                        <AddButton label="Tambah Akun Guru" onClick={() => {}} />
                    </div>
                </div>
                {/* User Tables */}
                <div className="px-4">
                    <UserTables />
                </div>
                
            </div>

        </div>
        </>
    )
}