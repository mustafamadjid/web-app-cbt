import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router";
import { BookOpen, ChevronDown, Filter, Search, X } from "lucide-react";
import toast from "react-hot-toast";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import BankSoalLayout from "@/layouts/BankSoalLayout/BankSoalLayout";

import type { BankSoalItem } from "@/types/BankSoal/BankSoal";
import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";

import {
  useDeleteBankSoal,
  useGetBankSoal,
  useGetBankSoalByGuru,
} from "@/services/Api/features-api/BankSoal/banksoal.service";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import { useGetMapel } from "@/services/Api/features-api/DataMaster/mapel.service";
import { paths } from "@/routes/paths";
import { useAuth } from "@/contexts/AuthContext";

function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

type ActiveSection = "ALL" | "MY";

const BankSoal = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const isGuru = user?.role === "GURU";

  const [activeSection, setActiveSection] = useState<ActiveSection>("ALL");
  const [selectedKelasId, setSelectedKelasId] = useState<number | null>(null);
  const [selectedMapelId, setSelectedMapelId] = useState<number | null>(null);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebouncedValue(search, 400);

  const [batasData, setBatasData] = useState(12);
  const [halamanSemua, setHalamanSemua] = useState(1);
  const [halamanSaya, setHalamanSaya] = useState(1);

  const [targetDeleteId, setTargetDeleteId] = useState<number | null>(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  const { execute: executeDeleteBankSoal } = useDeleteBankSoal();

  const { data: kelasData } = useGetDataKelasFull();
  const opsiTingkatKelas = useMemo(
    () => kelasData?.item_tingkat_kelas ?? [],
    [kelasData],
  );

  const kelasLabelById = useMemo(() => {
    return opsiTingkatKelas.reduce<Record<number, string>>((acc, tingkat) => {
      acc[tingkat.id_tingkat_kelas] = `Kelas ${tingkat.tingkat_kelas}`;
      return acc;
    }, {});
  }, [opsiTingkatKelas]);

  const { data: mapelRowsData } = useGetMapel({ limit: 50, offset: 0 });
  const mapelRows = useMemo(() => mapelRowsData ?? [], [mapelRowsData]);

  const mapelOptions = useMemo<MataPelajaranOption[]>(() => {
    const rows =
      selectedKelasId == null
        ? mapelRows
        : mapelRows.filter((mapel) => mapel.kelasId === selectedKelasId);
    return rows.map((mapel) => ({
      id: mapel.id,
      label: mapel.namaMapel,
    }));
  }, [mapelRows, selectedKelasId]);

  const mapelLabelById = useMemo(() => {
    return mapelRows.reduce<Record<number, string>>((acc, mapel) => {
      acc[mapel.id] = mapel.namaMapel;
      return acc;
    }, {});
  }, [mapelRows]);

  const {
    data: itemsData,
    loading,
    error,
    refetch: refetchAll,
  } = useGetBankSoal({
    search: debouncedSearch.trim() || undefined,
    id_kelas: selectedKelasId ?? undefined,
    id_mapel: selectedMapelId ?? undefined,
    limit: batasData,
    offset: (halamanSemua - 1) * batasData,
  });

  const idPengguna = user?.id_pengguna ?? 0;
  const {
    data: myRawItemsData,
    loading: loadingMyItemsRaw,
    error: myErrorMsg,
    refetch: refetchMyItems,
  } = useGetBankSoalByGuru(idPengguna);
  const myRawItems = useMemo(() => myRawItemsData ?? [], [myRawItemsData]);

  const myFilteredItems = useMemo(() => {
    if (idPengguna <= 0) return [];

    let rows = [...myRawItems];
    if (selectedKelasId != null) {
      rows = rows.filter((item) => item.id_kelas === selectedKelasId);
    }
    if (selectedMapelId != null) {
      rows = rows.filter((item) => item.id_mapel === selectedMapelId);
    }

    const q = debouncedSearch.trim().toLowerCase();
    if (q.length > 0) {
      rows = rows.filter((item) => {
        const haystack = [
          item.nama_bank_soal,
          item.materi,
          item.deskripsi,
          mapelLabelById[item.id_mapel],
          kelasLabelById[item.id_kelas],
        ]
          .filter(Boolean)
          .join(" ")
          .toLowerCase();

        return haystack.includes(q);
      });
    }

    return rows;
  }, [
    debouncedSearch,
    idPengguna,
    kelasLabelById,
    mapelLabelById,
    myRawItems,
    selectedKelasId,
    selectedMapelId,
  ]);

  const myOffset = (halamanSaya - 1) * batasData;
  const myItems = myFilteredItems.slice(myOffset, myOffset + batasData);

  const items: BankSoalItem[] = itemsData ?? [];
  const errorMsg = error ?? "";

  useEffect(() => {
    if (
      selectedMapelId != null &&
      !mapelOptions.some((mapel) => mapel.id === selectedMapelId)
    ) {
      setSelectedMapelId(null);
    }
  }, [mapelOptions, selectedMapelId]);

  useEffect(() => {
    setHalamanSemua(1);
    setHalamanSaya(1);
  }, [debouncedSearch, selectedKelasId, selectedMapelId, batasData]);

  const isFiltering =
    search.trim() !== "" || selectedMapelId !== null || selectedKelasId !== null;

  const createPath = isGuru
    ? paths.dashboard.buat_bank_soal_guru
    : paths.dashboard.buat_bank_soal;

  const detailPathTemplate = isGuru
    ? paths.dashboard.preview_bank_soal_guru
    : paths.dashboard.preview_bank_soal;

  const uploadPathTemplate = isGuru
    ? paths.dashboard.tambah_bank_soal_guru
    : paths.dashboard.tambah_bank_soal;

  const getDetailPath = (idBankSoal: number) =>
    detailPathTemplate.replace(":id", String(idBankSoal));

  const getUploadPath = (idBankSoal: number) =>
    uploadPathTemplate.replace(":idBankSoal", String(idBankSoal));

  const resetFilter = () => {
    setSearch("");
    setSelectedKelasId(null);
    setSelectedMapelId(null);
  };

  const resolveGuruLabel = (item: BankSoalItem) => {
    if (item.id_pengguna === idPengguna) {
      return user?.username ?? `Pengguna #${item.id_pengguna}`;
    }
    return `Pengguna #${item.id_pengguna}`;
  };

  const resolveMapelLabel = (item: BankSoalItem) =>
    mapelLabelById[item.id_mapel] ?? `Mapel #${item.id_mapel}`;

  const resolveKelasLabel = (item: BankSoalItem) =>
    kelasLabelById[item.id_kelas] ?? `Kelas #${item.id_kelas}`;

  const handleDelete = async () => {
    if (targetDeleteId == null) return;

    setDeleteLoading(true);
    try {
      await executeDeleteBankSoal(targetDeleteId);
      await Promise.all([refetchAll(), refetchMyItems()]);
      toast.success("Bank soal berhasil dihapus.");
      setTargetDeleteId(null);
    } catch (e) {
      const message =
        e instanceof Error ? e.message : "Gagal menghapus bank soal.";
      toast.error(message);
    } finally {
      setDeleteLoading(false);
    }
  };

  const semuaAwal =
    items.length === 0 ? 0 : (halamanSemua - 1) * batasData + 1;
  const semuaAkhir =
    items.length === 0 ? 0 : (halamanSemua - 1) * batasData + items.length;
  const bisaSemuaSebelumnya = halamanSemua > 1;
  const bisaSemuaSelanjutnya = items.length === batasData;

  const totalSaya = myFilteredItems.length;
  const sayaAwal = myItems.length === 0 ? 0 : (halamanSaya - 1) * batasData + 1;
  const sayaAkhir =
    myItems.length === 0 ? 0 : (halamanSaya - 1) * batasData + myItems.length;
  const bisaSayaSebelumnya = halamanSaya > 1;
  const bisaSayaSelanjutnya = myOffset + myItems.length < totalSaya;

  return (
    <div className="min-h-screen bg-[#F8F9FA] pb-20">
      <div className="mx-auto max-w-[1600px] space-y-8 p-6 lg:p-8">
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Bank Soal</h1>
            <p className="mt-1 text-sm text-gray-500">
              Kelola dan atur kumpulan soal ujian sekolah.
            </p>
          </div>
          <AddButton
            label="Buat Bank Soal"
            onClick={() => navigate(createPath)}
          />
        </div>

        <div className="flex flex-wrap items-center gap-3">
          <button
            type="button"
            onClick={() => setActiveSection("ALL")}
            className={`rounded-full px-4 py-2 text-sm font-semibold transition cursor-pointer ${
              activeSection === "ALL"
                ? "bg-[#397e50] text-white shadow-sm"
                : "bg-white text-gray-700 ring-1 ring-gray-200 hover:bg-gray-50"
            }`}
          >
            Semua Bank Soal
          </button>
          <button
            type="button"
            onClick={() => setActiveSection("MY")}
            className={`rounded-full px-4 py-2 text-sm font-semibold transition cursor-pointer ${
              activeSection === "MY"
                ? "bg-[#397e50] text-white shadow-sm"
                : "bg-white text-gray-700 ring-1 ring-gray-200 hover:bg-gray-50"
            }`}
          >
            Soal Saya
            <span
              className={`ml-2 rounded-full px-2 py-0.5 text-xs font-semibold ${
                activeSection === "MY"
                  ? "bg-white/20 text-white"
                  : "bg-gray-100 text-gray-600"
              }`}
            >
              {totalSaya}
            </span>
          </button>
        </div>

        <div className="sticky top-4 z-20 flex flex-col gap-3 rounded-2xl bg-white p-2 shadow-sm ring-1 ring-gray-200/70 sm:flex-row sm:items-center">
          <div className="relative flex-1 group">
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-gray-400 group-focus-within:text-[#397e50] transition-colors">
              <Search className="h-5 w-5" />
            </div>
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Cari nama bank soal, materi, atau deskripsi..."
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

          <div className="hidden h-8 w-px bg-gray-200 sm:block" />

          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <div className="relative min-w-[180px]">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
                <Filter className="h-4 w-4" />
              </div>
              <select
                value={selectedKelasId ?? ""}
                onChange={(e) =>
                  setSelectedKelasId(
                    e.target.value === "" ? null : Number(e.target.value),
                  )
                }
                className="block w-full cursor-pointer appearance-none rounded-xl border-none bg-gray-50 py-2.5 pl-9 pr-8 text-sm font-medium text-gray-700 hover:bg-gray-100 focus:ring-2 focus:ring-[#397e50]/20 transition-all"
              >
                <option value="">Semua Kelas</option>
                {opsiTingkatKelas.map((tingkat) => (
                  <option
                    key={tingkat.id_tingkat_kelas}
                    value={tingkat.id_tingkat_kelas}
                  >
                    Kelas {tingkat.tingkat_kelas}
                  </option>
                ))}
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400">
                <ChevronDown className="h-3.5 w-3.5" />
              </div>
            </div>

            <div className="relative min-w-[200px]">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
                <BookOpen className="h-4 w-4" />
              </div>
              <select
                value={selectedMapelId ?? ""}
                onChange={(e) =>
                  setSelectedMapelId(
                    e.target.value === "" ? null : Number(e.target.value),
                  )
                }
                className="block w-full cursor-pointer appearance-none rounded-xl border-none bg-gray-50 py-2.5 pl-9 pr-8 text-sm font-medium text-gray-700 hover:bg-gray-100 focus:ring-2 focus:ring-[#397e50]/20 transition-all"
              >
                <option value="">Semua Mapel</option>
                {mapelOptions.map((mapel) => (
                  <option key={mapel.id} value={mapel.id}>
                    {mapel.label}
                  </option>
                ))}
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400">
                <Filter className="h-3 w-3" />
              </div>
            </div>
          </div>
        </div>

        {activeSection === "MY" ? (
          <div className="space-y-4">
            <div className="flex items-center justify-between px-1">
              <h2 className="text-lg font-bold text-gray-900">Soal Saya</h2>
              {isFiltering && !loadingMyItemsRaw && (
                <button
                  onClick={resetFilter}
                  className="text-xs font-semibold text-[#397e50] hover:text-[#2c633f] hover:underline"
                >
                  Reset Filter
                </button>
              )}
            </div>
            <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6 rounded-lg">
              <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-sm text-slate-700">
                  Menampilkan <span className="font-medium">{sayaAwal}</span>{" "}
                  sampai <span className="font-medium">{sayaAkhir}</span> dari{" "}
                  <span className="font-medium">{totalSaya}</span> hasil
                </p>
                <div className="flex items-center gap-3">
                  <div className="flex cursor-pointer items-center gap-2 text-sm text-slate-600">
                    <span>Tampilkan</span>
                    <select
                      value={batasData}
                      onChange={(event) =>
                        setBatasData(Number(event.target.value))
                      }
                      className="cursor-pointer appearance-none rounded-lg border border-slate-200 bg-white px-7 py-1 text-sm text-slate-700 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
                    >
                      {[12, 20, 30, 40, 50].map((opsi) => (
                        <option key={opsi} value={opsi}>
                          {opsi}
                        </option>
                      ))}
                    </select>
                    <span>baris</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() =>
                        setHalamanSaya((prev) => Math.max(1, prev - 1))
                      }
                      disabled={!bisaSayaSebelumnya}
                      className="rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      Sebelumnya
                    </button>
                    <span className="text-sm text-slate-600">
                      Halaman {halamanSaya}
                    </span>
                    <button
                      type="button"
                      onClick={() => setHalamanSaya((prev) => prev + 1)}
                      disabled={!bisaSayaSelanjutnya}
                      className="rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      Selanjutnya
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between px-1">
              <div className="flex items-center gap-2 text-sm text-gray-500">
                {loadingMyItemsRaw ? (
                  <span className="flex items-center gap-2 text-[#397e50] animate-pulse">
                    <span className="h-2 w-2 rounded-full bg-[#397e50]" />
                    Memuat data soal saya...
                  </span>
                ) : (
                  <>
                    <span className="font-medium text-gray-900">
                      {myItems.length}
                    </span>{" "}
                    dari{" "}
                    <span className="font-medium text-gray-900">
                      {totalSaya}
                    </span>{" "}
                    soal saya
                    {myErrorMsg && (
                      <span className="text-rose-500 ml-2">• {myErrorMsg}</span>
                    )}
                  </>
                )}
              </div>
            </div>

            <div className="relative min-h-[220px]">
              {loadingMyItemsRaw && myItems.length === 0 ? (
                <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className="h-64 rounded-2xl bg-gray-200 animate-pulse"
                    />
                  ))}
                </div>
              ) : !loadingMyItemsRaw && myItems.length === 0 ? (
                <div className="flex flex-col items-center justify-center rounded-3xl border border-dashed border-gray-300 bg-gray-50/50 py-16 text-center">
                  <div className="mb-4 rounded-full bg-white p-4 shadow-sm ring-1 ring-gray-100">
                    <Search className="h-8 w-8 text-gray-300" />
                  </div>
                  <h3 className="text-base font-bold text-gray-900">
                    Soal saya belum tersedia
                  </h3>
                  <p className="max-w-md text-sm text-gray-500 mt-2">
                    Bank soal yang kamu buat akan tampil di sini.
                  </p>
                </div>
              ) : (
                <BankSoalLayout
                  items={myItems}
                  startIndex={(halamanSaya - 1) * batasData + 1}
                  resolveGuruLabel={resolveGuruLabel}
                  resolveMapelLabel={resolveMapelLabel}
                  resolveKelasLabel={resolveKelasLabel}
                  onKelola={(item) =>
                    navigate(getDetailPath(item.id_bank_soal))
                  }
                  onPreview={(item) =>
                    navigate(getDetailPath(item.id_bank_soal))
                  }
                  onUpload={(item) =>
                    navigate(getUploadPath(item.id_bank_soal))
                  }
                  onHapus={(item) => setTargetDeleteId(item.id_bank_soal)}
                />
              )}
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center justify-between px-1">
              <h2 className="text-lg font-bold text-gray-900">
                Semua Bank Soal
              </h2>
              {isFiltering && !loading && (
                <button
                  onClick={resetFilter}
                  className="text-xs font-semibold text-[#397e50] hover:text-[#2c633f] hover:underline"
                >
                  Reset Filter
                </button>
              )}
            </div>

            <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6 rounded-lg">
              <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-sm text-slate-700">
                  Menampilkan <span className="font-medium">{semuaAwal}</span>{" "}
                  sampai <span className="font-medium">{semuaAkhir}</span> dari{" "}
                  <span className="font-medium">{items.length}</span> hasil
                </p>
                <div className="flex items-center gap-3">
                  <div className="flex cursor-pointer items-center gap-2 text-sm text-slate-600">
                    <span>Tampilkan</span>
                    <select
                      value={batasData}
                      onChange={(event) =>
                        setBatasData(Number(event.target.value))
                      }
                      className="cursor-pointer appearance-none rounded-lg border border-slate-200 bg-white px-7 py-1 text-sm text-slate-700 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
                    >
                      {[12, 20, 30, 40, 50].map((opsi) => (
                        <option key={opsi} value={opsi}>
                          {opsi}
                        </option>
                      ))}
                    </select>
                    <span>baris</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() =>
                        setHalamanSemua((prev) => Math.max(1, prev - 1))
                      }
                      disabled={!bisaSemuaSebelumnya}
                      className="rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      Sebelumnya
                    </button>
                    <span className="text-sm text-slate-600">
                      Halaman {halamanSemua}
                    </span>
                    <button
                      type="button"
                      onClick={() => setHalamanSemua((prev) => prev + 1)}
                      disabled={!bisaSemuaSelanjutnya}
                      className="rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      Selanjutnya
                    </button>
                  </div>
                </div>
              </div>
            </div>

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
                    data halaman ini
                    {errorMsg && (
                      <span className="text-rose-500 ml-2">• {errorMsg}</span>
                    )}
                  </>
                )}
              </div>
            </div>

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
                  startIndex={(halamanSemua - 1) * batasData + 1}
                  resolveGuruLabel={resolveGuruLabel}
                  resolveMapelLabel={resolveMapelLabel}
                  resolveKelasLabel={resolveKelasLabel}
                  onKelola={(item) =>
                    navigate(getDetailPath(item.id_bank_soal))
                  }
                  onPreview={(item) =>
                    navigate(getDetailPath(item.id_bank_soal))
                  }
                  onUpload={(item) =>
                    navigate(getUploadPath(item.id_bank_soal))
                  }
                  onHapus={(item) => setTargetDeleteId(item.id_bank_soal)}
                />
              )}
            </div>
          </div>
        )}
      </div>

      <ConfirmAlert
        isOpen={targetDeleteId != null}
        title="Konfirmasi Hapus Bank Soal"
        message="Data bank soal ini akan dihapus permanen. Lanjutkan?"
        onClose={() => {
          if (deleteLoading) return;
          setTargetDeleteId(null);
        }}
        onConfirm={() => void handleDelete()}
        isLoading={deleteLoading}
        confirmLabel="Ya, Hapus"
        loadingLabel="Menghapus..."
      />
    </div>
  );
};

export default BankSoal;
