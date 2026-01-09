import { BoxJadwalUjian } from "@/components/features/Ujian/BoxJadwalUjian";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { getJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { TingkatKelasOption } from "@/types/DataMaster/Kelas";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { useEffect, useMemo, useRef, useState } from "react";
import { Search, Calendar, Layers, MapPin } from "lucide-react"; // Tambahan icon untuk filter

const STATUS_SECTIONS = [
  { key: "berlangsung", label: "Berlangsung", color: "bg-emerald-500" },
  { key: "belum_dimulai", label: "Belum Mulai", color: "bg-amber-500" },
  { key: "selesai", label: "Selesai", color: "bg-slate-500" },
] as const;

type StatusKey = (typeof STATUS_SECTIONS)[number]["key"];

function useDebouncedValue<T>(value: T, delayMs: number) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = window.setTimeout(() => setDebounced(value), delayMs);
    return () => window.clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

export const JadwalUjian = () => {
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [activeTab, setActiveTab] = useState<StatusKey>("berlangsung"); // Default ke 'berlangsung'

  const [daftarJadwalUjian, setDaftarJadwalUjian] = useState<JadwalUjianItem[]>(
    []
  );
  const [tingkatKelasOptions, setTingkatKelasOptions] = useState<
    TingkatKelasOption[]
  >([]);
  const [ruangOptions, setRuangOptions] = useState<RuangUjianRow[]>([]);

  const [searchTerm, setSearchTerm] = useState("");
  const debouncedSearchTerm = useDebouncedValue(searchTerm, 300);

  const [selectedDate, setSelectedDate] = useState("");
  const [selectedTingkatId, setSelectedTingkatId] = useState("");
  const [selectedRuang, setSelectedRuang] = useState("");

  const requestSeq = useRef(0);

  // Load options
  useEffect(() => {
    (async () => {
      try {
        const [tingkatOptions, ruangUjianOptions] = await Promise.all([
          getTingkatKelasOptions(),
          getRuangUjian(),
        ]);
        setTingkatKelasOptions(tingkatOptions);
        setRuangOptions(ruangUjianOptions);
      } catch {
        setTingkatKelasOptions([]);
        setRuangOptions([]);
      }
    })();
  }, []);

  // Fetch data
  useEffect(() => {
    const seq = ++requestSeq.current;
    const tingkat =
      selectedTingkatId.trim() === "" ? undefined : Number(selectedTingkatId);
    const tingkatKelas = Number.isFinite(tingkat) ? tingkat : undefined;
    const ruangUjianId =
      selectedRuang.trim() === "" ? undefined : Number(selectedRuang);

    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getJadwalUjian({
          q: debouncedSearchTerm.trim() || undefined,
          tanggal: selectedDate || undefined,
          ruangUjianId,
          tingkatKelasId: tingkatKelas,
        });
        if (seq !== requestSeq.current) return;
        setDaftarJadwalUjian(data);
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data jadwal ujian.");
        setDaftarJadwalUjian([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, [debouncedSearchTerm, selectedDate, selectedRuang, selectedTingkatId]);

  const groupedJadwal = useMemo(() => {
    const grouped: Record<StatusKey, JadwalUjianItem[]> = {
      belum_dimulai: [],
      berlangsung: [],
      selesai: [],
    };
    for (const ujian of daftarJadwalUjian) {
      const status = ujian.status_ujian as StatusKey;
      if (grouped[status]) grouped[status].push(ujian);
    }
    return grouped;
  }, [daftarJadwalUjian]);

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="flex flex-col gap-8">
        {/* Header & Filter Card */}
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <h1 className="mb-6 text-2xl font-bold text-slate-800">
            Jadwal Ujian
          </h1>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Search size={16} /> Cari Jadwal
              </label>
              <input
                type="text"
                placeholder="Cari ujian, pengawas..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              />
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Calendar size={16} /> Tanggal
              </label>
              <input
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              />
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Layers size={16} /> Tingkat
              </label>
              <select
                value={selectedTingkatId}
                onChange={(e) => setSelectedTingkatId(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Tingkat</option>
                {tingkatKelasOptions.map((t) => (
                  <option key={t.id_tingkat_kelas} value={t.id_tingkat_kelas}>
                    Kelas {t.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <MapPin size={16} /> Ruangan
              </label>
              <select
                value={selectedRuang}
                onChange={(e) => setSelectedRuang(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Ruangan</option>
                {ruangOptions.map((r) => (
                  <option key={r.id} value={String(r.id)}>
                    {r.namaRuangan}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>

        {/* Status Tabs Navigation */}
        <div className="flex flex-wrap items-center gap-2 rounded-2xl bg-slate-100 p-1.5 sm:w-fit">
          {STATUS_SECTIONS.map((section) => {
            const isActive = activeTab === section.key;
            const count = groupedJadwal[section.key].length;
            return (
              <button
                key={section.key}
                onClick={() => setActiveTab(section.key)}
                className={`
                  flex items-center cursor-pointer gap-2 rounded-xl px-6 py-2.5 text-sm font-bold transition-all
                  ${
                    isActive
                      ? "bg-white text-[#397e50] shadow-sm"
                      : "text-slate-500 hover:bg-slate-200 hover:text-slate-700"
                  }
                `}
              >
                {section.label}
                <span
                  className={`
                  ml-1 rounded-md px-1.5 py-0.5 text-2xs
                  ${
                    isActive
                      ? "bg-[#397e50] text-white"
                      : "bg-slate-200 text-slate-600"
                  }
                `}
                >
                  {count}
                </span>
              </button>
            );
          })}
        </div>

        {/* Content Area */}
        <div className="min-h-[300px]">
          {errorMsg && (
            <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
              {errorMsg}
            </div>
          )}

          {loading ? (
            <div className="flex h-40 items-center justify-center rounded-2xl border border-dashed border-slate-300 text-slate-500">
              <div className="flex flex-col items-center gap-2">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-[#397e50] border-t-transparent" />
                <span>Memuat data...</span>
              </div>
            </div>
          ) : (
            <div className="flex flex-col gap-5">
              {groupedJadwal[activeTab].length > 0 ? (
                groupedJadwal[activeTab].map((ujian) => (
                  <BoxJadwalUjian key={ujian.id} {...ujian} />
                ))
              ) : (
                <div className="flex h-40 flex-col items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white text-slate-500">
                  <p className="font-medium">Tidak ada jadwal ujian</p>
                  <p className="text-xs">
                    Ujian dengan status "{activeTab.replace("_", " ")}" tidak
                    ditemukan.
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
