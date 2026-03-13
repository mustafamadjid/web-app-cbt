import React from "react";
import { useNavigate, useParams } from "react-router";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import { paths } from "@/routes/paths";
import { useGetSoalUjianForSiswa } from "@/services/Api/features-api/Ujian/soalUjian.service";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";

const OPTION_LABELS = ["A", "B", "C", "D", "E", "F", "G", "H"];

type SiswaSoalPreviewContentProps = {
  title: string;
  sisaWaktu: string;
  soalPreview: SoalPreviewItem[];
  onBack: () => void;
};

const SiswaSoalPreviewContent: React.FC<SiswaSoalPreviewContentProps> = ({
  title,
  sisaWaktu,
  soalPreview,
  onBack,
}) => {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [selectedOptions, setSelectedOptions] = React.useState<
    Record<number, number>
  >({});

  const totalSoal = soalPreview.length;
  const currentSoal = soalPreview[currentIndex];

  const handleSelectOption = (soalId: number, optionId: number) => {
    setSelectedOptions((prev) => ({ ...prev, [soalId]: optionId }));
  };

  const questionNavigator = (
    <div className="space-y-3">
      <p className="text-xs font-semibold text-slate-400">Nomor Soal</p>
      <div className="flex flex-wrap gap-2">
        {soalPreview.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id}
              type="button"
              onClick={() => setCurrentIndex(index)}
              className={[
                "flex h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold transition",
                isActive
                  ? "border-[#397e50] bg-[#397e50] text-white"
                  : "border-slate-200 text-slate-500 hover:border-[#397e50] hover:text-[#397e50]",
              ].join(" ")}
              aria-label={`Soal nomor ${soal.nomor}`}
            >
              {soal.nomor}
            </button>
          );
        })}
      </div>
    </div>
  );

  return (
    <SoalLayout
      title={title}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu={sisaWaktu}
      soal={currentSoal}
      questionNavigator={questionNavigator}
      selectedOptionId={selectedOptions[currentSoal.id]}
      onSelectOption={(optionId) => handleSelectOption(currentSoal.id, optionId)}
      onPrev={
        currentIndex > 0 ? () => setCurrentIndex((prev) => prev - 1) : undefined
      }
      onNext={
        currentIndex < totalSoal - 1
          ? () => setCurrentIndex((prev) => prev + 1)
          : undefined
      }
      onBack={onBack}
    />
  );
};

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { idJadwalUjian } = useParams();

  const parsedIdJadwalUjian = Number(idJadwalUjian);
  const isIdJadwalUjianValid =
    Number.isInteger(parsedIdJadwalUjian) && parsedIdJadwalUjian > 0;
  const {
    data: soalRows,
    loading,
    error,
  } = useGetSoalUjianForSiswa(parsedIdJadwalUjian, isIdJadwalUjianValid);

  const soalPreview = React.useMemo<SoalPreviewItem[]>(
    () =>
      (soalRows ?? []).map((soal) => ({
        id: soal.id_soal,
        nomor: soal.no_urut_soal,
        tipe: soal.tipe_soal,
        pertanyaan: soal.pertanyaan,
        gambar_url: resolveImageUrl(soal.gambar) || undefined,
        opsi: soal.opsi_jawaban.map((opsi, index) => ({
          id: opsi.id_pilihan_ganda,
          label: OPTION_LABELS[index] ?? String(index + 1),
          text: opsi.isi_pilihan,
        })),
      })),
    [soalRows],
  );

  const errorMessage =
    !idJadwalUjian || !isIdJadwalUjianValid ? "Jadwal ujian tidak ditemukan." : error;
  const title = isIdJadwalUjianValid ? `Ujian #${parsedIdJadwalUjian}` : "Ujian";

  if (loading) {
    return (
      <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
        Memuat soal ujian...
      </div>
    );
  }

  if (errorMessage || soalPreview.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-red-200 bg-white p-6 text-center text-sm text-red-500">
        {errorMessage ?? "Soal ujian tidak tersedia."}
      </div>
    );
  }

  return (
    <SiswaSoalPreviewContent
      key={parsedIdJadwalUjian}
      title={title}
      sisaWaktu="00:00:00"
      soalPreview={soalPreview}
      onBack={() => navigate(paths.dashboard.ujian_siswa)}
    />
  );
};

export default UjianMulaiSiswa;
