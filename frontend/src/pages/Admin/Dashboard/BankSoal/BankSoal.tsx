import { useNavigate } from "react-router";

import { BankSoalLayout } from "@/layouts/BankSoalLayout/BankSoalLayout";
import type { BankSoalItem } from "@/types/DataMaster/BankSoal";
import { AddButton } from "@/components/common/Button/AddButton";

export const BankSoal= () => {
  const navigate = useNavigate();

  const data: BankSoalItem[] = [
    {
      id: "bs-001",
      nama_banksoal: "Bank Soal Ujian Bahasa",
      kelas: 11,
      mata_pelajaran: "Bahasa Indonesia",
      materi: "Ejaan Bahasa",
      jumlah_soal_pg: 10,
      jumlah_soal_essay: 5,
      deskripsi:
        "Bank Soal Ujian Bahasa ini disusun untuk peserta didik kelas XI...",
      tgl_buat: "2025-03-10",
    },
    {
      id: "bs-002",
      nama_banksoal: "Bank Soal Ujian Fisika",
      kelas: 10,
      mata_pelajaran: "Fisika",
      materi: "Hukum Newton",
      jumlah_soal_pg: 10,
      jumlah_soal_essay: 5,
      deskripsi: "Bank soal untuk latihan ujian fisika...",
      tgl_buat: "2025-03-12",
    },
    {
      id: "bs-003",
      nama_banksoal: "Bank Soal Ujian Fisika",
      kelas: 10,
      mata_pelajaran: "Fisika",
      materi: "Hukum Newton",
      jumlah_soal_pg: 10,
      jumlah_soal_essay: 5,
      deskripsi: "Bank soal untuk latihan ujian fisika...",
      tgl_buat: "2025-03-12",
    },
  ];

  return (
    <div className="w-full mt-4 px-8 flex flex-col gap-5">
      {/* tombol Add Bank Soal */}
      <div className="flex justify-end">
        <AddButton
        
          label="Tambah Bank Soal"
          onClick={() => navigate(`/dashboard/administrator/bank-soal/tambah`)}
        />
      </div>
      <BankSoalLayout
        items={data}
        onPreview={(item) => navigate(`/banksoal/preview/${item.id}`)}
        onKelola={(item) => navigate(`/banksoal/kelola/${item.id}`)}
        onHapus={(item) => console.log("hapus", item.id)}
      />
    </div>
  );
};
