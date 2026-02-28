import React, { useEffect, useMemo, useState } from "react";
import {
  Archive,
  BookOpen,
  ChevronDown,
  Edit3,
  Search,
  Trash2,
} from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import type {
  MataPelajaranOption,
} from "@/types/DataMaster/MataPelajaran";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import {
  deleteMataPelajaran,
  useGetMapel,
} from "@/services/Api/features-api/DataMaster/mapel.service";
import toast from "react-hot-toast";

import { paths } from "@/routes/paths";

function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

const DataMataPelajaran: React.FC = () => {
  const navigate = useNavigate();

  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");

  const [tingkatTerpilih, setTingkatTerpilih] = useState<number | null>(null);
  const [mapelTerpilih, setMapelTerpilih] = useState<string>("");

  const [idTerpilih, setIdTerpilih] = useState<Set<number>>(new Set());
  const [batasData, setBatasData] = useState(12);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);
  const [targetDeleteId, setTargetDeleteId] = useState<number | null>(null);
  const [isBulkDeleteConfirmOpen, setIsBulkDeleteConfirmOpen] = useState(false);
  const [isDeleteLoading, setIsDeleteLoading] = useState(false);

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);

  // Hook: fetch kelas options on mount
  const { data: kelasData } = useGetDataKelasFull();
  const opsiTingkatKelas = kelasData?.item_tingkat_kelas ?? [];

  // Hook: fetch mapel list with filters (auto re-fetch on dep change)
  const {
    data: daftarMapel,
    loading,
    error: errorMsg,
    refetch: refetchMapel,
  } = useGetMapel({
    search: debouncedKataKunci.trim() || undefined,
    tingkatKelas: tingkatTerpilih ?? undefined,
    namaMapel: mapelTerpilih || undefined,
    limit: batasData,
    offset: (halamanSaatIni - 1) * batasData,
  });

  // Build mapel filter options from the full list
  const { data: mapelFullForOptions } = useGetMapel();
  const opsiMapel = useMemo<MataPelajaranOption[]>(() => {
    if (!mapelFullForOptions) return [];
    const uniqNamaMapel = Array.from(
      new Set(mapelFullForOptions.map((item) => item.namaMapel)),
    ).sort((a, b) => a.localeCompare(b, "id", { sensitivity: "base" }));
    return uniqNamaMapel.map((nama, index) => ({ id: index + 1, label: nama }));
  }, [mapelFullForOptions]);

  const kelasLabelById = useMemo(() => {
    return opsiTingkatKelas.reduce<Record<number, string>>((acc, tingkat) => {
      acc[tingkat.id_tingkat_kelas] = `Kelas ${tingkat.tingkat_kelas}`;
      return acc;
    }, {});
  }, [opsiTingkatKelas]);

  const mapelList = daftarMapel ?? [];
  const totalData = mapelList.length;
  const dataTerlihat = mapelList;

  const semuaTerlihatTerpilih =
    dataTerlihat.length > 0 &&
    dataTerlihat.every((mapel) => idTerpilih.has(mapel.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        dataTerlihat.forEach((mapel) => next.delete(mapel.id));
      } else {
        dataTerlihat.forEach((mapel) => next.add(mapel.id));
      }
      return next;
    });
  };

  const togglePilihBaris = (id: number) => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const executeDeleteMapel = async (idMapel: number) => {
    try {
      setIsDeleteLoading(true);
      await deleteMataPelajaran(idMapel);
      await refetchMapel();
      setIdTerpilih((prev) => {
        const next = new Set(prev);
        next.delete(idMapel);
        return next;
      });
      toast.success("Berhasil menghapus data mata pelajaran");
    } catch {
      toast.error("Gagal menghapus data mata pelajaran");
    } finally {
      setIsDeleteLoading(false);
      setTargetDeleteId(null);
    }
  };

  const executeBulkDelete = async () => {
    if (idTerpilih.size === 0) return;
    const ids = Array.from(idTerpilih);
    try {
      setIsDeleteLoading(true);
      await Promise.all(ids.map((id) => deleteMataPelajaran(id)));
      await refetchMapel();
      setIdTerpilih(new Set());
      setDropdownAksiTerbuka(false);
      toast.success("Berhasil menghapus data mata pelajaran terpilih");
    } catch {
      toast.error("Gagal menghapus data mata pelajaran terpilih");
    } finally {
      setIsDeleteLoading(false);
      setIsBulkDeleteConfirmOpen(false);
    }
  };

  const handleDeleteConfirm = async () => {
    if (targetDeleteId !== null) {
      await executeDeleteMapel(targetDeleteId);
      return;
    }

    if (isBulkDeleteConfirmOpen) {
      await executeBulkDelete();
    }
  };

  const jumlahTerpilih = idTerpilih.size;

  const resetFilter = () => {
    setKataKunci("");
    setHalamanSaatIni(1);
    setTingkatTerpilih(null);
    setMapelTerpilih("");
  };

  useEffect(() => {
    setHalamanSaatIni(1);
  }, [debouncedKataKunci, batasData]);

  const awalData =
    totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + 1;
  const akhirData =
    totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + totalData;
  const bisaSebelumnya = halamanSaatIni > 1;
  const bisaSelanjutnya = totalData === batasData;

  return (
    <div className="w-full space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Data Mata Pelajaran
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola daftar mata pelajaran, kode mapel, serta tingkat kelas
            terkait.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <AddButton
            label="Tambah Mata Pelajaran"
            onClick={() =>
              navigate(`${paths.dashboard.tambah_data_master_mapel}`)
            }
          />
        </div>
      </div>

      <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <div className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
          <div className="w-full lg:w-[340px]">
            <label className="text-xs font-medium text-slate-600">
              Pencarian
            </label>
            <div className="relative mt-1">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <Search className="h-4 w-4 text-slate-400" />
              </div>
              <input
                type="text"
                value={kataKunci}
                onChange={(e) => setKataKunci(e.target.value)}
                className="block w-full cursor-pointer rounded-lg border border-slate-200 bg-slate-50 py-2 pl-10 pr-3 text-sm text-slate-900 placeholder:text-slate-400 focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-1 focus:ring-[#397e50]"
                placeholder="Cari kode, nama mapel, atau deskripsi..."
              />
            </div>
          </div>

          <div className="grid w-full grid-cols-1 gap-3 sm:grid-cols-3">
            <div>
              <label className="text-xs font-medium text-slate-600">
                Tingkat Kelas
              </label>
              <select
                value={tingkatTerpilih ?? ""}
                onChange={(e) =>
                  setTingkatTerpilih(
                    e.target.value === "" ? null : Number(e.target.value),
                  )
                }
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {opsiTingkatKelas.map((tingkat) => (
                  <option
                    key={tingkat.id_tingkat_kelas}
                    value={tingkat.id_tingkat_kelas}
                  >
                    {tingkat.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-xs font-medium text-slate-600">
                Mapel
              </label>
              <select
                value={mapelTerpilih}
                onChange={(e) => setMapelTerpilih(e.target.value)}
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {opsiMapel.map((mapel) => (
                  <option key={mapel.label} value={mapel.label}>
                    {mapel.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="flex items-center justify-between gap-2 lg:justify-end">
            <button
              type="button"
              onClick={resetFilter}
              className="min-w-[140px] cursor-pointer rounded-lg border border-slate-200 bg-white px-5 py-2 text-sm text-slate-700 hover:bg-slate-50"
            >
              Reset Filter
            </button>

            <div className="relative">
              <button
                type="button"
                onClick={() => setDropdownAksiTerbuka((v) => !v)}
                className="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
                aria-haspopup="menu"
                aria-expanded={dropdownAksiTerbuka}
              >
                Aksi
                <ChevronDown className="h-4 w-4 text-slate-400" />
              </button>

              {dropdownAksiTerbuka && (
                <div
                  role="menu"
                  className="absolute right-0 z-20 mt-2 w-44 rounded-lg border border-slate-200 bg-white shadow-lg"
                  onMouseLeave={() => setDropdownAksiTerbuka(false)}
                >
                  <ul className="p-2 text-sm text-slate-700">
                    <li>
                      <button
                        type="button"
                        className="flex w-full items-center gap-2 rounded-lg px-3 py-2 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
                        onClick={() => setDropdownAksiTerbuka(false)}
                        disabled={jumlahTerpilih === 0}
                      >
                        <Archive className="h-4 w-4 text-slate-500" />
                        Arsipkan
                      </button>
                    </li>
                    <li>
                      <button
                        type="button"
                        className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-rose-600 hover:bg-rose-50 disabled:cursor-not-allowed disabled:opacity-50"
                        onClick={() => setIsBulkDeleteConfirmOpen(true)}
                        disabled={jumlahTerpilih === 0}
                      >
                        <Trash2 className="h-4 w-4" />
                        Hapus
                      </button>
                    </li>
                  </ul>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="text-sm text-slate-600">
          {loading ? (
            <span>Memuat data...</span>
          ) : errorMsg ? (
            <span className="text-rose-600">{errorMsg}</span>
          ) : (
            <span>
              Menampilkan{" "}
              <span className="font-medium">{mapelList.length}</span> hasil.
            </span>
          )}
        </div>
      </div>

      <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="border-b border-slate-200 bg-slate-50 text-xs uppercase text-slate-500">
              <tr>
                <th scope="col" className="w-4 p-4">
                  <input
                    type="checkbox"
                    checked={semuaTerlihatTerpilih}
                    onChange={togglePilihSemuaTerlihat}
                    className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                  />
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Kode Mapel
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Mata Pelajaran
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Tingkat Kelas
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Deskripsi
                </th>
                <th scope="col" className="px-6 py-3 text-right font-semibold">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {dataTerlihat.length > 0 ? (
                dataTerlihat.map((mapel) => {
                  const kelasLabel = kelasLabelById[mapel.kelasId];
                  return (
                    <tr
                      key={mapel.id}
                      className={`transition-colors hover:bg-slate-50 ${
                        idTerpilih.has(mapel.id) ? "bg-indigo-50/30" : ""
                      }`}
                    >
                      <td className="p-4">
                        <input
                          type="checkbox"
                          checked={idTerpilih.has(mapel.id)}
                          onChange={() => togglePilihBaris(mapel.id)}
                          className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                        />
                      </td>
                      <td className="px-6 py-4 font-semibold text-slate-900">
                        {mapel.kodeMapel}
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex flex-col">
                          <span className="font-semibold text-slate-900">
                            {mapel.namaMapel}
                          </span>
                          <span className="text-xs text-slate-500">
                            ID: {mapel.id}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <span className="text-slate-700">
                          {kelasLabel ?? "-"}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <p className="text-sm text-slate-600">
                          {mapel.deskripsiMapel}
                        </p>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600"
                            title="Edit"
                            onClick={() =>
                              navigate(
                                paths.dashboard.edit_data_master_mapel.replace(
                                  ":id",
                                  String(mapel.id),
                                ),
                              )
                            }
                          >
                            <Edit3 className="h-4 w-4" />
                          </button>
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                            title="Hapus"
                            onClick={() => setTargetDeleteId(mapel.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center">
                    <div className="flex flex-col items-center justify-center gap-2">
                      <BookOpen className="h-10 w-10 text-slate-300" />
                      <p className="text-base font-medium text-slate-900">
                        Tidak ada mata pelajaran ditemukan
                      </p>
                      <p className="text-sm text-slate-500">
                        Coba sesuaikan kata kunci atau filter Anda.
                      </p>
                    </div>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
        <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
          <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm text-slate-700">
                Menampilkan <span className="font-medium">{awalData}</span>{" "}
                sampai <span className="font-medium">{akhirData}</span> dari{" "}
                <span className="font-medium">{totalData}</span> hasil
              </p>
            </div>
            <div className="flex items-center gap-3">
              <div className="flex cursor-pointer items-center gap-2 text-sm text-slate-600">
                <span>Tampilkan</span>
                <select
                  value={batasData}
                  onChange={(event) => setBatasData(Number(event.target.value))}
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
                    setHalamanSaatIni((sebelumnya) =>
                      Math.max(1, sebelumnya - 1),
                    )
                  }
                  disabled={!bisaSebelumnya}
                  className="rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Sebelumnya
                </button>
                <span className="text-sm text-slate-600">
                  Halaman {halamanSaatIni}
                </span>
                <button
                  type="button"
                  onClick={() =>
                    setHalamanSaatIni((sebelumnya) => sebelumnya + 1)
                  }
                  disabled={!bisaSelanjutnya}
                  className="rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Selanjutnya
                </button>
              </div>
            </div>
          </div>
        </div>{" "}
      </div>

      <ConfirmAlert
        isOpen={targetDeleteId !== null || isBulkDeleteConfirmOpen}
        title="Konfirmasi Hapus Mata Pelajaran"
        message={
          targetDeleteId !== null
            ? "Data mata pelajaran ini akan dihapus permanen. Lanjutkan?"
            : `Anda akan menghapus ${jumlahTerpilih} data mata pelajaran terpilih. Lanjutkan?`
        }
        onClose={() => {
          if (isDeleteLoading) return;
          setTargetDeleteId(null);
          setIsBulkDeleteConfirmOpen(false);
        }}
        onConfirm={() => void handleDeleteConfirm()}
        isLoading={isDeleteLoading}
        confirmLabel="Ya, Hapus"
        loadingLabel="Menghapus..."
      />
    </div>
  );
};

export default DataMataPelajaran;
