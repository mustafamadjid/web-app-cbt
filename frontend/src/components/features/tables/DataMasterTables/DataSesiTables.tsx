import React, { useCallback, useEffect, useRef, useState } from "react";
import { Clock3, Edit3, Search, Trash2 } from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  deleteSesi,
  getSesi,
} from "@/services/Api/features-api/DataMaster/sesi.service";
import type { SesiRow } from "@/types/DataMaster/Sesi";

function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

const DataSesiTables: React.FC = () => {
  const navigate = useNavigate();

  const [kataKunci, setKataKunci] = useState("");
  const [daftarSesi, setDaftarSesi] = useState<SesiRow[]>([]);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [batasData, setBatasData] = useState(12);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);

  const [modalKonfirmasiTerbuka, setModalKonfirmasiTerbuka] = useState(false);
  const [idSesiAkanDihapus, setIdSesiAkanDihapus] = useState<number | null>(null);
  const [sedangHapus, setSedangHapus] = useState(false);

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);
  const requestSeq = useRef(0);

  const fetchSesi = useCallback(async () => {
    return getSesi({
      q: debouncedKataKunci.trim() || undefined,
      limit: batasData,
      offset: (halamanSaatIni - 1) * batasData,
    });
  }, [debouncedKataKunci, batasData, halamanSaatIni]);

  useEffect(() => {
    const seq = ++requestSeq.current;

    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");

        const data = await fetchSesi();
        if (seq !== requestSeq.current) return;

        setDaftarSesi(data);
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data sesi.");
        setDaftarSesi([]);
      } finally {
        if (seq === requestSeq.current) {
          setLoading(false);
        }
      }
    })();
  }, [fetchSesi]);

  useEffect(() => {
    setHalamanSaatIni(1);
  }, [debouncedKataKunci, batasData]);

  const resetFilter = () => {
    setKataKunci("");
    setHalamanSaatIni(1);
  };

  const handleOpenDeleteConfirm = (id_sesi: number) => {
    setIdSesiAkanDihapus(id_sesi);
    setModalKonfirmasiTerbuka(true);
  };

  const handleConfirmDelete = async () => {
    if (!idSesiAkanDihapus) return;

    setSedangHapus(true);
    setErrorMsg("");

    try {
      await deleteSesi(idSesiAkanDihapus);
      setModalKonfirmasiTerbuka(false);
      setIdSesiAkanDihapus(null);
      const data = await fetchSesi();
      setDaftarSesi(data);
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.message === "data not found"
            ? "Data sesi tidak ditemukan."
            : e.message === "delete restricted : constraint violation"
              ? "Sesi tidak bisa dihapus karena masih dipakai."
              : "Data sesi gagal dihapus."
          : "Data sesi gagal dihapus.";
      setErrorMsg(message);
    } finally {
      setSedangHapus(false);
    }
  };

  const totalData = daftarSesi.length;
  const awalData = totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + 1;
  const akhirData =
    totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + totalData;
  const bisaSebelumnya = halamanSaatIni > 1;
  const bisaSelanjutnya = totalData === batasData;

  return (
    <div className="w-full space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">Data Sesi</h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola daftar sesi ujian sesuai kebutuhan jadwal.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <AddButton
            label="Tambah Sesi"
            onClick={() => navigate(`${paths.dashboard.tambah_data_master_sesi}`)}
          />
        </div>
      </div>

      <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <div className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
          <div className="w-full lg:w-[340px]">
            <label className="text-xs font-medium text-slate-600">Pencarian</label>
            <div className="relative mt-1">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <Search className="h-4 w-4 text-slate-400" />
              </div>
              <input
                type="text"
                value={kataKunci}
                onChange={(e) => setKataKunci(e.target.value)}
                className="block w-full cursor-pointer rounded-lg border border-slate-200 bg-slate-50 py-2 pl-10 pr-3 text-sm text-slate-900 placeholder:text-slate-400 focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-1 focus:ring-[#397e50]"
                placeholder="Cari kode atau nama sesi..."
              />
            </div>
          </div>

          <button
            type="button"
            onClick={resetFilter}
            className="min-w-[140px] cursor-pointer rounded-lg border border-slate-200 bg-white px-5 py-2 text-sm text-slate-700 hover:bg-slate-50"
          >
            Reset Filter
          </button>
        </div>

        {errorMsg && (
          <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-2 text-sm text-rose-600">
            {errorMsg}
          </div>
        )}
      </div>

      <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        {loading ? (
          <div className="px-6 py-12 text-center text-sm text-slate-500">Memuat data...</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="border-b border-slate-200 bg-slate-50 text-xs uppercase text-slate-500">
                <tr>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">Kode Sesi</th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">Nama Sesi</th>
                  <th className="px-6 py-3 text-right font-semibold text-slate-600">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {daftarSesi.length > 0 ? (
                  daftarSesi.map((sesi) => (
                    <tr key={sesi.id_sesi} className="transition-colors hover:bg-slate-50">
                      <td className="px-6 py-4 font-semibold text-slate-900">{sesi.kode_sesi}</td>
                      <td className="px-6 py-4">
                        <div className="flex flex-col">
                          <span className="font-semibold text-slate-900">{sesi.nama_sesi}</span>
                          <span className="text-xs text-slate-500">ID: {sesi.id_sesi}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600"
                            title="Edit"
                            onClick={() =>
                              navigate(
                                paths.dashboard.edit_data_master_sesi.replace(
                                  ":id",
                                  String(sesi.id_sesi),
                                ),
                              )
                            }
                          >
                            <Edit3 className="h-4 w-4" />
                          </button>
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                            title="Hapus"
                            onClick={() => handleOpenDeleteConfirm(sesi.id_sesi)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={3} className="px-6 py-12 text-center">
                      <div className="flex flex-col items-center justify-center gap-2">
                        <Clock3 className="h-10 w-10 text-slate-300" />
                        <p className="text-base font-medium text-slate-900">Tidak ada sesi ditemukan</p>
                      </div>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}

        <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
          <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <p className="text-sm text-slate-700">
              Menampilkan <span className="font-medium">{awalData}</span> sampai{" "}
              <span className="font-medium">{akhirData}</span> dari{" "}
              <span className="font-medium">{totalData}</span> hasil
            </p>

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
                    setHalamanSaatIni((sebelumnya) => Math.max(1, sebelumnya - 1))
                  }
                  disabled={!bisaSebelumnya}
                  className="cursor-pointer rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Sebelumnya
                </button>
                <span className="text-sm text-slate-600">Halaman {halamanSaatIni}</span>
                <button
                  type="button"
                  onClick={() => setHalamanSaatIni((sebelumnya) => sebelumnya + 1)}
                  disabled={!bisaSelanjutnya}
                  className="cursor-pointer rounded-lg border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 transition hover:border-slate-300 hover:text-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Selanjutnya
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <ConfirmAlert
        isOpen={modalKonfirmasiTerbuka}
        title="Konfirmasi Hapus Sesi"
        message="Data sesi yang dihapus tidak bisa dikembalikan. Lanjutkan?"
        onClose={() => {
          if (sedangHapus) return;
          setModalKonfirmasiTerbuka(false);
          setIdSesiAkanDihapus(null);
        }}
        onConfirm={handleConfirmDelete}
        isLoading={sedangHapus}
        confirmLabel="Ya, Hapus"
        loadingLabel="Menghapus..."
      />
    </div>
  );
};

export default DataSesiTables;
