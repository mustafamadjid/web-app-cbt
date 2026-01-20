import React from "react";
import { useNavigate, useParams } from "react-router";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import type { SoalUjianItem } from "@/types/BankSoal/BankSoal";
import { paths } from "@/routes/paths";

const DUMMY_SOAL: SoalUjianItem[] = [
  {
    id_soal: 1,
    nomor_urut_soal: 1,
    tipe_soal: "PILIHAN_GANDA",
    pertanyaan: "Hasil dari 12 + 8 adalah ...",
    opsi_a: "18",
    opsi_b: "20",
    opsi_c: "22",
    opsi_d: "24",
  },
  {
    id_soal: 2,
    nomor_urut_soal: 2,
    tipe_soal: "PILIHAN_GANDA",
    pertanyaan: "Ibu kota Indonesia adalah ...",
    opsi_a: "Bandung",
    opsi_b: "Jakarta",
    opsi_c: "Surabaya",
    opsi_d: "Yogyakarta",
  },
  {
    id_soal: 3,
    nomor_urut_soal: 3,
    tipe_soal: "ESSAY",
    pertanyaan: "Jelaskan alasan pentingnya menjaga kelestarian lingkungan.",
  },
];

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [selectedOptions, setSelectedOptions] = React.useState<
    Record<number, string>
  >({});

  const totalSoal = DUMMY_SOAL.length;
  const currentSoal = DUMMY_SOAL[currentIndex];

  const handleSelectOption = (soalId: number, value: string) => {
    setSelectedOptions((prev) => ({ ...prev, [soalId]: value }));
  };

  const questionNavigator = (
    <div className="space-y-3">
      <p className="text-xs font-semibold text-slate-400">Nomor Soal</p>
      <div className="flex flex-wrap gap-2">
        {DUMMY_SOAL.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id_soal}
              type="button"
              onClick={() => setCurrentIndex(index)}
              className={[
                "flex h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold transition",
                isActive
                  ? "border-[#397e50] bg-[#397e50] text-white"
                  : "border-slate-200 text-slate-500 hover:border-[#397e50] hover:text-[#397e50]",
              ].join(" ")}
              aria-label={`Soal nomor ${soal.nomor_urut_soal}`}
            >
              {soal.nomor_urut_soal}
            </button>
          );
        })}
      </div>
    </div>
  );

  return (
    <SoalLayout
      title={`Ujian ${id ?? ""}`}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu="01:30:00"
      soal={currentSoal}
      questionNavigator={questionNavigator}
      selectedOption={selectedOptions[currentSoal.id_soal]}
      onSelectOption={(value) => handleSelectOption(currentSoal.id_soal, value)}
      onPrev={
        currentIndex > 0 ? () => setCurrentIndex((prev) => prev - 1) : undefined
      }
      onNext={
        currentIndex < totalSoal - 1
          ? () => setCurrentIndex((prev) => prev + 1)
          : undefined
      }
      onBack={() => navigate(paths.dashboard.ujian_siswa)}
    />
  );
};

export default UjianMulaiSiswa;
