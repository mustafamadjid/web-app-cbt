import { useMemo, useState } from "react";
import { Calendar, Layers, RefreshCw, RotateCcw } from "lucide-react";

import DatePicker from "@/components/common/DateInput/DatePicker";
import BoxHasilUjian from "@/components/features/Ujian/BoxHasilUjian";
import { useAuth } from "@/contexts/AuthContext";
import { tahunOption } from "@/helper/TahunOption/TahunOption";
import { paths } from "@/routes/paths";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import {
  useGetUjianEssayUngradedList,
} from "@/services/Api/features-api/Ujian/hasilUjian.service";
import { useGetJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { TingkatKelas } from "@/types/DataMaster/Kelas";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";

const MONTH_OPTIONS = [
  { value: 1, label: "Januari" },
  { value: 2, label: "Februari" },
  { value: 3, label: "Maret" },
  { value: 4, label: "April" },
  { value: 5, label: "Mei" },
  { value: 6, label: "Juni" },
  { value: 7, label: "Juli" },
  { value: 8, label: "Agustus" },
  { value: 9, label: "September" },
  { value: 10, label: "Oktober" },
  { value: 11, label: "November" },
  { value: 12, label: "Desember" },
] as const;

const HASIL_UJIAN_TABS = [
  { key: "selesai", label: "Ujian Selesai" },
  { key: "essay_belum_dinilai", label: "Essay Belum Dinilai" },
] as const;

type HasilUjianTabKey = (typeof HASIL_UJIAN_TABS)[number]["key"];

const getMonthFromTanggalUjian = (value?: string) => {
  if (!value) return null;

  const monthPart = value.split("-")[1];
  const month = Number(monthPart);

  if (!Number.isInteger(month) || month < 1 || month > 12) {
    return null;
  }

  return month;
};

const HasilUjian = () => {
  const { user } = useAuth();

  const [activeTab, setActiveTab] = useState<HasilUjianTabKey>("selesai");
  const [selectedTingkatId, setSelectedTingkatId] = useState<number | null>(
    null,
  );
  const [selectedTahun, setSelectedTahun] = useState<string | null>(null);
  const [selectedMonth, setSelectedMonth] = useState<number | null>(null);
  const [selectedDate, setSelectedDate] = useState("");

  const { data: kelasData } = useGetDataKelasFull();
  const tingkatKelasOptions: TingkatKelas[] = kelasData?.item_tingkat_kelas ?? [];

  const tahunOptions = useMemo(
    () => tahunOption().map((year) => String(year)),
    [],
  );

  const detailPathTemplate =
    user?.role === "GURU"
      ? paths.dashboard.hasil_ujian_detail_guru
      : paths.dashboard.hasil_ujian_detail;

  const {
    data: jadwalDataRaw,
    loading: loadingSelesai,
    error: errorSelesai,
    refetch: refetchSelesai,
  } = useGetJadwalUjian({
    tanggal: selectedDate || undefined,
    tingkatKelasId: selectedTingkatId ?? undefined,
    tahun: selectedTahun || undefined,
    kategoriUjian: "selesai",
  });

  const {
    data: essayUngradedData,
    loading: loadingEssayUngraded,
    error: errorEssayUngraded,
    refetch: refetchEssayUngraded,
  } = useGetUjianEssayUngradedList(
    {
      tanggal: selectedDate || undefined,
      tingkatKelasId: selectedTingkatId ?? undefined,
      tahun: selectedTahun || undefined,
      bulan: selectedMonth ?? undefined,
    },
    activeTab === "essay_belum_dinilai",
  );

  const daftarSelesai = useMemo(
    () =>
      (jadwalDataRaw ?? []).filter((ujian: JadwalUjianItem) => {
        if (selectedMonth == null) return true;

        return getMonthFromTanggalUjian(ujian.tanggal_ujian) === selectedMonth;
      }),
    [jadwalDataRaw, selectedMonth],
  );

  const daftarEssayBelumDinilai = essayUngradedData ?? [];
  const daftarAktif =
    activeTab === "essay_belum_dinilai"
      ? daftarEssayBelumDinilai
      : daftarSelesai;
  const loadingAktif =
    activeTab === "essay_belum_dinilai"
      ? loadingEssayUngraded
      : loadingSelesai;
  const errorAktif =
    activeTab === "essay_belum_dinilai"
      ? errorEssayUngraded
      : errorSelesai;
  const refetchAktif =
    activeTab === "essay_belum_dinilai"
      ? refetchEssayUngraded
      : refetchSelesai;

  const isEssayTab = activeTab === "essay_belum_dinilai";
  const titleDescription = isEssayTab
    ? "Daftar ujian dengan jawaban essay yang masih menunggu koreksi guru atau admin."
    : "Daftar ujian yang sudah selesai beserta akses hasilnya.";
  const totalLabel = isEssayTab ? "Perlu koreksi" : "Total selesai";

  const isFilterActive =
    selectedTingkatId != null ||
    selectedTahun != null ||
    selectedMonth != null ||
    selectedDate !== "";

  const handleResetFilter = () => {
    setSelectedTingkatId(null);
    setSelectedTahun(null);
    setSelectedMonth(null);
    setSelectedDate("");
  };

  const getEmptyStateMessage = () => {
    if (isEssayTab) {
      return isFilterActive
        ? "Tidak ada ujian essay yang cocok dengan filter yang dipilih."
        : "Belum ada ujian essay yang menunggu penilaian.";
    }

    return isFilterActive
      ? "Tidak ada ujian selesai yang cocok dengan filter yang dipilih."
      : "Belum ada ujian selesai yang memiliki hasil.";
  };

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="flex flex-col gap-6">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-slate-800">Hasil Ujian</h1>
              <p className="text-sm text-slate-500">{titleDescription}</p>
            </div>

            <div className="flex flex-wrap items-center gap-2">
              <span className="inline-flex items-center rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">
                {totalLabel}: {daftarAktif.length}
              </span>
              <button
                type="button"
                onClick={() => void refetchAktif()}
                disabled={loadingAktif}
                className="inline-flex items-center gap-2 rounded-xl border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-600 transition hover:border-[#397e50] hover:text-[#397e50] disabled:cursor-not-allowed disabled:opacity-60"
              >
                <RefreshCw
                  size={14}
                  className={loadingAktif ? "animate-spin" : ""}
                />
                Refresh
              </button>
              <button
                type="button"
                onClick={handleResetFilter}
                disabled={!isFilterActive}
                className="inline-flex items-center gap-2 rounded-xl border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-600 transition hover:border-[#397e50] hover:text-[#397e50] disabled:cursor-not-allowed disabled:opacity-60"
              >
                <RotateCcw size={14} />
                Reset Filter
              </button>
            </div>
          </div>

          <div className="mt-6 flex flex-wrap items-center gap-2 rounded-2xl bg-slate-100 p-1.5 sm:w-fit">
            {HASIL_UJIAN_TABS.map((tab) => {
              const isActive = activeTab === tab.key;

              return (
                <button
                  key={tab.key}
                  type="button"
                  onClick={() => setActiveTab(tab.key)}
                  className={`rounded-xl px-6 py-2.5 text-sm font-bold transition-all ${
                    isActive
                      ? "bg-white text-[#397e50] shadow-sm"
                      : "text-slate-500 hover:bg-slate-200 hover:text-slate-700"
                  }`}
                >
                  {tab.label}
                </button>
              );
            })}
          </div>

          <div className="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Layers size={16} /> Kelas
              </label>
              <select
                value={selectedTingkatId ?? ""}
                onChange={(event) =>
                  setSelectedTingkatId(
                    event.target.value === "" ? null : Number(event.target.value),
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Kelas</option>
                {tingkatKelasOptions.map((kelas) => (
                  <option key={kelas.id_tingkat_kelas} value={kelas.id_tingkat_kelas}>
                    Kelas {kelas.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Calendar size={16} /> Tahun
              </label>
              <select
                value={selectedTahun ?? ""}
                onChange={(event) =>
                  setSelectedTahun(
                    event.target.value === "" ? null : String(event.target.value),
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Tahun</option>
                {tahunOptions.map((tahun) => (
                  <option key={tahun} value={tahun}>
                    {tahun}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Calendar size={16} /> Bulan
              </label>
              <select
                value={selectedMonth ?? ""}
                onChange={(event) =>
                  setSelectedMonth(
                    event.target.value === "" ? null : Number(event.target.value),
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Bulan</option>
                {MONTH_OPTIONS.map((month) => (
                  <option key={month.value} value={month.value}>
                    {month.label}
                  </option>
                ))}
              </select>
            </div>

            <DatePicker
              id="hasil-ujian-tanggal"
              label="Tanggal"
              value={selectedDate}
              onChange={setSelectedDate}
            />
          </div>
        </div>

        {errorAktif ? (
          <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-semibold text-rose-700">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <span>{errorAktif}</span>
              <button
                type="button"
                onClick={() => void refetchAktif()}
                className="inline-flex items-center gap-2 self-start rounded-xl border border-rose-200 bg-white px-4 py-2 text-xs font-semibold text-rose-700 transition hover:border-rose-300 hover:bg-rose-100"
              >
                <RefreshCw size={14} />
                Coba Lagi
              </button>
            </div>
          </div>
        ) : null}

        {loadingAktif ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500">
            {isEssayTab
              ? "Memuat daftar ujian essay yang perlu dikoreksi..."
              : "Memuat daftar hasil ujian..."}
          </div>
        ) : null}

        {!loadingAktif && !errorAktif && daftarAktif.length === 0 ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
            {getEmptyStateMessage()}
          </div>
        ) : null}

        <div className="grid gap-6 lg:grid-cols-2">
          {daftarAktif.map((ujian) => {
            const hasilDetailId = ujian.id_ujian ?? ujian.id;

            return (
              <BoxHasilUjian
                key={`${activeTab}-${hasilDetailId}`}
                {...ujian}
                linkHasil={detailPathTemplate.replace(":id", String(hasilDetailId))}
              />
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default HasilUjian;
