import { useMemo, useState } from "react";
import { Calendar, Layers, RefreshCw, RotateCcw } from "lucide-react";

import DatePicker from "@/components/common/DateInput/DatePicker";
import BoxHasilUjian from "@/components/features/Ujian/BoxHasilUjian";
import { useAuth } from "@/contexts/AuthContext";
import { tahunOption } from "@/helper/TahunOption/TahunOption";
import { paths } from "@/routes/paths";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
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
];

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
    loading,
    error: errorMsg,
    refetch,
  } = useGetJadwalUjian({
    tanggal: selectedDate || undefined,
    tingkatKelasId: selectedTingkatId ?? undefined,
    tahun: selectedTahun || undefined,
    kategoriUjian: "selesai",
  });

  console.log("jadwalDataRaw", jadwalDataRaw);

  const daftarSelesai = useMemo(
    () =>
      (jadwalDataRaw ?? []).filter((ujian: JadwalUjianItem) => {
        if (selectedMonth == null) return true;

        return getMonthFromTanggalUjian(ujian.tanggal_ujian) === selectedMonth;
      }),
    [jadwalDataRaw, selectedMonth],
  );

  console.log("daftarSelesai", daftarSelesai);

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

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="flex flex-col gap-6">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-slate-800">Hasil Ujian</h1>
              <p className="text-sm text-slate-500">
                Daftar ujian yang sudah selesai beserta akses hasilnya.
              </p>
            </div>

            <div className="flex flex-wrap items-center gap-2">
              <span className="inline-flex items-center rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">
                Total selesai: {daftarSelesai.length}
              </span>
              <button
                type="button"
                onClick={() => void refetch()}
                disabled={loading}
                className="inline-flex items-center gap-2 rounded-xl border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-600 transition hover:border-[#397e50] hover:text-[#397e50] disabled:cursor-not-allowed disabled:opacity-60"
              >
                <RefreshCw size={14} className={loading ? "animate-spin" : ""} />
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

        {errorMsg ? (
          <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-semibold text-rose-700">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <span>{errorMsg}</span>
              <button
                type="button"
                onClick={() => void refetch()}
                className="inline-flex items-center gap-2 self-start rounded-xl border border-rose-200 bg-white px-4 py-2 text-xs font-semibold text-rose-700 transition hover:border-rose-300 hover:bg-rose-100"
              >
                <RefreshCw size={14} />
                Coba Lagi
              </button>
            </div>
          </div>
        ) : null}

        {loading ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500">
            Memuat daftar hasil ujian...
          </div>
        ) : null}

        {!loading && !errorMsg && daftarSelesai.length === 0 ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
            {isFilterActive
              ? "Tidak ada ujian selesai yang cocok dengan filter yang dipilih."
              : "Belum ada ujian selesai yang memiliki hasil."}
          </div>
        ) : null}

        <div className="grid gap-6 lg:grid-cols-2">
          {daftarSelesai.map((ujian) => {
            const hasilDetailId = ujian.id_ujian ?? ujian.id;

            return (
              <BoxHasilUjian
                key={hasilDetailId}
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
