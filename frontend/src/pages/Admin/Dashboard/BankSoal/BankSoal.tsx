"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router";
import {
  Search,

  Filter,
  Layers,
  LayoutGrid,
  BookOpen,
  XCircle,
  AlertCircle,
} from "lucide-react";

import { AddButton } from "@/components/common/Button/AddButton";

import { BankSoalLayout } from "@/layouts/BankSoalLayout/BankSoalLayout";

import type { BankSoalItem } from "@/types/DataMaster/BankSoal";
import type {
  KelasOption,
  MataPelajaranOption,
} from "@/types/DataMaster/MataPelajaran";

import {
  getBankSoalByKelas,
  getMataPelajaranOptions,
  getKelasOptions,
} from "@/services/Api/features-api/BankSoal/banksoal.service";

// Util: Debounce
function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

type ViewMode = "ALL" | "BY_KELAS";

export const BankSoal = () => {
  const navigate = useNavigate();

  // ----- UI State -----
  const [viewMode, setViewMode] = useState<ViewMode>("ALL");
  const [selectedKelasId, setSelectedKelasId] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebouncedValue(search, 400); // Sedikit diperlambat agar lebih smooth
  const [selectedMapelId, setSelectedMapelId] = useState<string>("");

  // ----- Data State -----
  const [kelasOptions, setKelasOptions] = useState<KelasOption[]>([]);
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
        const kelas = await getKelasOptions();
        if (!mounted) return;
        setKelasOptions(kelas);
        if (
          kelas.length > 0 &&
          selectedKelasId == null &&
          viewMode === "BY_KELAS"
        ) {
          setSelectedKelasId(kelas[0].id);
        }
      } catch (e) {
        if (mounted) setErrorMsg("Gagal memuat data kelas.");
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
        const kelasIdForMapel =
          viewMode === "BY_KELAS" ? selectedKelasId ?? undefined : undefined;
        const mapel = await getMataPelajaranOptions({
          kelasId: kelasIdForMapel,
        });
        if (!mounted) return;
        setMapelOptions(mapel);
        if (selectedMapelId && !mapel.some((m) => m.id === selectedMapelId)) {
          setSelectedMapelId("");
        }
      } catch (e) {
        if (mounted) setErrorMsg("Gagal memuat data mata pelajaran.");
      }
    })();
    return () => {
      mounted = false;
    };
  }, [viewMode, selectedKelasId]);

  // 3. Load Bank Soal Data
  useEffect(() => {
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const kelasId =
          viewMode === "BY_KELAS" ? selectedKelasId ?? undefined : undefined;
        const data = await getBankSoalByKelas({
          kelasId,
          mapelId: selectedMapelId || undefined,
          q: debouncedSearch?.trim() || undefined,
        });
        if (seq !== requestSeq.current) return;
        setItems(data);
      } catch (e) {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data bank soal. Silakan coba lagi.");
        setItems([]);
      } finally {
        if (seq === requestSeq.current) setLoading(false);
      }
    })();
  }, [viewMode, selectedKelasId, selectedMapelId, debouncedSearch]);

  const selectedKelasLabel = useMemo(() => {
    if (!selectedKelasId) return "";
    const k = kelasOptions.find((x) => x.id === selectedKelasId);
    if (!k) return "";
    return k.nama_kelas
      ? `Kelas ${k.tingkat_kelas} - ${k.nama_kelas}`
      : `Kelas ${k.tingkat_kelas}`;
  }, [kelasOptions, selectedKelasId]);

  return (
    <div className="min-h-screen bg-[#ecf1ed] pb-20">
      <div className="mx-auto max-w-[1600px] space-y-6 p-4 sm:p-6 lg:p-8">
        {/* === ADD BUTTON === */}
        <div className="flex justify-end">
          <AddButton label="Buat Bank Soal" onClick={() => navigate("/dashboard/administrator/bank-soal/tambah")} />
        </div>
      

        {/* === FILTER & CONTROL PANEL === */}
        <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm transition-all hover:shadow-md">
          {/* Top Bar: View Mode & Search */}
          <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
            {/* Left: View Mode Tabs */}
            <div className="flex flex-col gap-4 w-full lg:w-auto">
              <div>
                <label className="mb-2 block text-xs font-bold uppercase tracking-wider text-gray-500">
                  Mode Tampilan
                </label>
                <div className="inline-flex rounded-xl bg-gray-100 p-1">
                  <button
                    onClick={() => {
                      setViewMode("ALL");
                      setSelectedKelasId(null);
                    }}
                    className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-bold transition-all ${
                      viewMode === "ALL"
                        ? "bg-white text-[#397e50] shadow-sm"
                        : "text-gray-500 hover:text-gray-700 hover:bg-gray-200/50"
                    }`}
                  >
                    <Layers className="h-4 w-4" />
                    Semua Data
                  </button>
                  <button
                    onClick={() => {
                      setViewMode("BY_KELAS");
                      if (!selectedKelasId && kelasOptions.length > 0)
                        setSelectedKelasId(kelasOptions[0].id);
                    }}
                    className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-bold transition-all ${
                      viewMode === "BY_KELAS"
                        ? "bg-white text-[#397e50] shadow-sm"
                        : "text-gray-500 hover:text-gray-700 hover:bg-gray-200/50"
                    }`}
                  >
                    <LayoutGrid className="h-4 w-4" />
                    Filter Kelas
                  </button>
                </div>
              </div>

              {/* Kelas Chips (Only visible in BY_KELAS mode) */}
              {viewMode === "BY_KELAS" && (
                <div className="animate-in fade-in slide-in-from-top-2 duration-300">
                  <label className="mb-2 block text-xs font-bold uppercase tracking-wider text-gray-500">
                    Pilih Kelas
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {kelasOptions.length === 0 ? (
                      <span className="text-sm italic text-gray-400">
                        Data kelas tidak tersedia.
                      </span>
                    ) : (
                      kelasOptions.map((k) => {
                        const isActive = k.id === selectedKelasId;
                        const label = k.nama_kelas
                          ? `${k.tingkat_kelas} - ${k.nama_kelas}`
                          : `Kelas ${k.tingkat_kelas}`;
                        return (
                          <button
                            key={k.id}
                            onClick={() => setSelectedKelasId(k.id)}
                            className={`rounded-lg px-3 py-1.5 text-xs font-bold transition-all border ${
                              isActive
                                ? "bg-[#397e50] text-white border-[#397e50] shadow-md shadow-[#397e50]/20"
                                : "bg-white text-gray-600 border-gray-200 hover:border-[#397e50] hover:text-[#397e50]"
                            }`}
                          >
                            {label}
                          </button>
                        );
                      })
                    )}
                  </div>
                </div>
              )}
            </div>

            {/* Right: Search & Mapel Filter */}
            <div className="flex flex-col gap-4 w-full lg:w-[450px]">
              {/* Search Box */}
              <div>
                <label className="mb-2 block text-xs font-bold uppercase tracking-wider text-gray-500">
                  Pencarian
                </label>
                <div className="group relative">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                    <Search className="h-4 w-4" />
                  </div>
                  <input
                    type="text"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Cari judul, materi, atau kode..."
                    className="block w-full rounded-xl border border-gray-200 bg-gray-50 py-2.5 pl-10 pr-10 text-sm text-gray-800 outline-none transition-all placeholder:text-gray-400 focus:border-[#397e50] focus:bg-white focus:ring-2 focus:ring-[#397e50]/20"
                  />
                  {search && (
                    <button
                      onClick={() => setSearch("")}
                      className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-rose-500"
                    >
                      <XCircle className="h-4 w-4" />
                    </button>
                  )}
                </div>
              </div>

              {/* Mapel Dropdown */}
              <div>
                <label className="mb-2 block text-xs font-bold uppercase tracking-wider text-gray-500">
                  Mata Pelajaran
                </label>
                <div className="relative">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                    <BookOpen className="h-4 w-4" />
                  </div>
                  <select
                    value={selectedMapelId}
                    onChange={(e) => setSelectedMapelId(e.target.value)}
                    className="block w-full cursor-pointer appearance-none rounded-xl border border-gray-200 bg-gray-50 py-2.5 pl-10 pr-8 text-sm text-gray-800 outline-none transition-all focus:border-[#397e50] focus:bg-white focus:ring-2 focus:ring-[#397e50]/20"
                  >
                    <option value="">Semua Mata Pelajaran</option>
                    {mapelOptions.map((m) => (
                      <option key={m.id} value={m.id}>
                        {m.label}
                      </option>
                    ))}
                  </select>
                  <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400">
                    <Filter className="h-3.5 w-3.5" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Status Bar */}
          {(loading || items.length > 0 || errorMsg) && (
            <div className="mt-5 border-t border-gray-100 pt-4 flex items-center justify-between text-xs text-gray-500">
              <div className="flex items-center gap-2">
                {loading && (
                  <span className="flex items-center gap-2 text-[#397e50] font-medium animate-pulse">
                    <span className="h-2 w-2 rounded-full bg-[#397e50]" />
                    Sedang memuat data...
                  </span>
                )}
                {!loading && !errorMsg && (
                  <span>
                    Menampilkan{" "}
                    <strong className="text-gray-800">{items.length}</strong>{" "}
                    bank soal
                    {viewMode === "BY_KELAS" && selectedKelasLabel && (
                      <>
                        {" "}
                        untuk{" "}
                        <span className="font-bold text-[#397e50]">
                          {selectedKelasLabel}
                        </span>
                      </>
                    )}
                  </span>
                )}
                {errorMsg && (
                  <span className="flex items-center gap-1 text-rose-600 font-medium">
                    <AlertCircle className="h-3.5 w-3.5" />
                    {errorMsg}
                  </span>
                )}
              </div>
            </div>
          )}
        </div>

        {/* === CONTENT GRID === */}
        <div className="relative min-h-[300px]">
          {loading && items.length === 0 ? (
            // Loading State (Skeleton sederhana)
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-64 rounded-2xl bg-gray-200 animate-pulse"
                />
              ))}
            </div>
          ) : !loading && items.length === 0 ? (
            // Empty State
            <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 bg-white py-20 text-center">
              <div className="mb-4 rounded-full bg-gray-50 p-4">
                <Search className="h-8 w-8 text-gray-300" />
              </div>
              <h3 className="text-lg font-bold text-gray-800">
                Tidak ada data ditemukan
              </h3>
              <p className="max-w-md text-sm text-gray-500 mt-1">
                Coba ubah kata kunci pencarian, filter mata pelajaran, atau
                pilih kelas yang berbeda.
              </p>
              <button
                onClick={() => {
                  setSearch("");
                  setSelectedMapelId("");
                }}
                className="mt-6 text-sm font-bold text-[#397e50] hover:underline"
              >
                Reset Filter
              </button>
            </div>
          ) : (
            // Data Grid
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
