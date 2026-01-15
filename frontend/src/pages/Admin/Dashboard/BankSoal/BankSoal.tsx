import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router";
import {
  Search,
  Filter,
  Layers,
  LayoutGrid,
  BookOpen,
  X,
  ChevronDown,
} from "lucide-react";

import AddButton from "@/components/common/Button/AddButton";
import BankSoalLayout from "@/layouts/BankSoalLayout/BankSoalLayout";

import type { BankSoalItem } from "@/types/DataMaster/BankSoal";
import type { TingkatKelasOption } from "@/types/DataMaster/Kelas";
import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";

import {
  getBankSoalByKelas,
  getMataPelajaranOptions,
} from "@/services/Api/features-api/BankSoal/banksoal.service";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import { paths } from "@/routes/paths";

// Util: Debounce (Tetap sama)
function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

type ViewMode = "ALL" | "BY_KELAS";

const BankSoal = () => {
  const navigate = useNavigate();

  // ----- UI State -----
  const [viewMode, setViewMode] = useState<ViewMode>("ALL");
  const [selectedTingkatId, setSelectedTingkatId] = useState<number | null>(
    null
  );
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebouncedValue(search, 400);
  const [selectedMapelId, setSelectedMapelId] = useState<number | null>(null);

  // ----- Data State -----
  const [tingkatKelasOptions, setTingkatKelasOptions] = useState<
    TingkatKelasOption[]
  >([]);
  const [mapelOptions, setMapelOptions] = useState<MataPelajaranOption[]>([]);
  const [items, setItems] = useState<BankSoalItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string>("");

  const requestSeq = useRef(0);

  // 1. Load Kelas Options
  useEffect(() => {
    let mounted = true;
    (async () => {
      try {
        const kelas = await getTingkatKelasOptions();
        if (!mounted) return;
        setTingkatKelasOptions(kelas);
      } catch (e) {
        if (mounted) setErrorMsg("Gagal memuat data tingkat kelas.");
      }
    })();
    return () => {
      mounted = false;
    };
  }, []);

  // 2. Load Mapel Options
  useEffect(() => {
    let mounted = true;
    (async () => {
      try {
        const mapel = await getMataPelajaranOptions({
          tingkatKelasId:
            viewMode === "BY_KELAS"
              ? selectedTingkatId ?? undefined
              : undefined,
        });
        if (!mounted) return;
        setMapelOptions(mapel);
        if (
          selectedMapelId != null &&
          !mapel.some((m) => m.id === selectedMapelId)
        ) {
          setSelectedMapelId(null);
        }
      } catch (e) {
        if (mounted) setErrorMsg("Gagal memuat data mata pelajaran.");
      }
    })();
    return () => {
      mounted = false;
    };
  }, [viewMode, selectedTingkatId]);

  // 3. Load Bank Soal Data
  useEffect(() => {
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getBankSoalByKelas({
          tingkatKelasId:
            viewMode === "BY_KELAS"
              ? selectedTingkatId ?? undefined
              : undefined,
          mapelId: selectedMapelId ?? undefined,
          q: debouncedSearch?.trim() || undefined,
        });
        if (seq !== requestSeq.current) return;
        setItems(data);
      } catch (e) {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data bank soal.");
        setItems([]);
      } finally {
        if (seq === requestSeq.current) setLoading(false);
      }
    })();
  }, [viewMode, selectedTingkatId, selectedMapelId, debouncedSearch]);

  // Helper untuk mengubah kelas (Logic penggabungan ViewMode & TingkatID)
  const handleTingkatChange = (value: string) => {
    if (value === "ALL") {
      setViewMode("ALL");
      setSelectedTingkatId(null);
    } else {
      setViewMode("BY_KELAS");
      setSelectedTingkatId(Number(value));
    }
  };

  const isFiltering =
    search !== "" || selectedMapelId !== null || viewMode === "BY_KELAS";

  return (
    <div className="min-h-screen bg-[#F8F9FA] pb-20">
      <div className="mx-auto max-w-[1600px] space-y-8 p-6 lg:p-8">
        {/* === HEADER SECTION === */}
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Bank Soal</h1>
            <p className="text-sm text-gray-500 mt-1">
              Kelola dan atur kumpulan soal ujian sekolah.
            </p>
          </div>
          <AddButton
            label="Buat Bank Soal"
            onClick={() => navigate(`${paths.dashboard.tambah_bank_soal}`)}
          />
        </div>

        {/* === MINIMALIST TOOLBAR === */}
        {/* Menggunakan sticky agar tetap terlihat saat scroll jika konten panjang */}
        <div className="sticky top-4 z-20 flex flex-col gap-3 rounded-2xl bg-white p-2 shadow-sm ring-1 ring-gray-200/70 sm:flex-row sm:items-center">
          {/* 1. Search Bar (Flexible width) */}
          <div className="relative flex-1 group">
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-gray-400 group-focus-within:text-[#397e50] transition-colors">
              <Search className="h-5 w-5" />
            </div>
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Cari judul, materi, atau kode..."
              className="block w-full rounded-xl border-none bg-transparent py-3 pl-10 pr-10 text-sm font-medium text-gray-800 placeholder:text-gray-400 focus:ring-0"
            />
            {search && (
              <button
                onClick={() => setSearch("")}
                className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-rose-500 transition-colors"
              >
                <X className="h-4 w-4" />
              </button>
            )}
          </div>

          {/* Divider Vertical (Desktop Only) */}
          <div className="hidden h-8 w-px bg-gray-200 sm:block" />

          {/* 2. Filter Controls Wrapper */}
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            {/* Filter: Tingkat Kelas (Gabungan ViewMode & Dropdown) */}
            <div className="relative min-w-[180px]">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
                {viewMode === "ALL" ? (
                  <Layers className="h-4 w-4" />
                ) : (
                  <LayoutGrid className="h-4 w-4" />
                )}
              </div>
              <select
                value={
                  viewMode === "ALL"
                    ? "ALL"
                    : selectedTingkatId?.toString() ?? ""
                }
                onChange={(e) => handleTingkatChange(e.target.value)}
                className="block w-full cursor-pointer appearance-none rounded-xl border-none bg-gray-50 py-2.5 pl-9 pr-8 text-sm font-medium text-gray-700 hover:bg-gray-100 focus:ring-2 focus:ring-[#397e50]/20 transition-all"
              >
                <option value="ALL">Semua Kelas</option>
                <optgroup label="Pilih Tingkat">
                  {tingkatKelasOptions.map((t) => (
                    <option key={t.id_tingkat_kelas} value={t.id_tingkat_kelas}>
                      Kelas {t.tingkat_kelas}
                    </option>
                  ))}
                </optgroup>
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400">
                <ChevronDown className="h-3.5 w-3.5" />
              </div>
            </div>

            {/* Filter: Mata Pelajaran */}
            <div className="relative min-w-[200px]">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
                <BookOpen className="h-4 w-4" />
              </div>
              <select
                value={selectedMapelId ?? ""}
                onChange={(e) =>
                  setSelectedMapelId(
                    e.target.value === "" ? null : Number(e.target.value)
                  )
                }
                className="block w-full cursor-pointer appearance-none rounded-xl border-none bg-gray-50 py-2.5 pl-9 pr-8 text-sm font-medium text-gray-700 hover:bg-gray-100 focus:ring-2 focus:ring-[#397e50]/20 transition-all"
              >
                <option value="">Semua Mapel</option>
                {mapelOptions.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.label}
                  </option>
                ))}
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400">
                <Filter className="h-3 w-3" />
              </div>
            </div>
          </div>
        </div>

        {/* === STATUS & COUNT INDICATOR === */}
        <div className="flex items-center justify-between px-1">
          <div className="flex items-center gap-2 text-sm text-gray-500">
            {loading ? (
              <span className="flex items-center gap-2 text-[#397e50] animate-pulse">
                <span className="h-2 w-2 rounded-full bg-[#397e50]" />
                Memuat data...
              </span>
            ) : (
              <>
                <span className="font-medium text-gray-900">
                  {items.length}
                </span>{" "}
                Soal ditemukan
                {errorMsg && (
                  <span className="text-rose-500 ml-2">• {errorMsg}</span>
                )}
              </>
            )}
          </div>

          {/* Reset Button (Only visible if filtering) */}
          {isFiltering && !loading && (
            <button
              onClick={() => {
                setSearch("");
                handleTingkatChange("ALL");
                setSelectedMapelId(null);
              }}
              className="text-xs font-semibold text-[#397e50] hover:text-[#2c633f] hover:underline"
            >
              Reset Filter
            </button>
          )}
        </div>

        {/* === CONTENT GRID === */}
        <div className="relative min-h-[300px]">
          {loading && items.length === 0 ? (
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {[1, 2, 3, 4, 5, 6].map((i) => (
                <div
                  key={i}
                  className="h-64 rounded-2xl bg-gray-200 animate-pulse"
                />
              ))}
            </div>
          ) : !loading && items.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 bg-gray-50/50 py-24 text-center">
              <div className="mb-4 rounded-full bg-white p-4 shadow-sm ring-1 ring-gray-100">
                <Search className="h-8 w-8 text-gray-300" />
              </div>
              <h3 className="text-lg font-bold text-gray-900">
                Data tidak ditemukan
              </h3>
              <p className="max-w-md text-sm text-gray-500 mt-2">
                Kami tidak dapat menemukan bank soal dengan filter tersebut.
              </p>
            </div>
          ) : (
            <BankSoalLayout
              items={items}
              onPreview={(item) => navigate(`/banksoal/preview/${item.id}`)}
              onKelola={(item) => navigate(`/banksoal/kelola/${item.id}`)}
              onHapus={(item) => console.log("Hapus", item.id)}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default BankSoal;
