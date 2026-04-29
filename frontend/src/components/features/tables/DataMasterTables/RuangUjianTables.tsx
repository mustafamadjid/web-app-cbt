import React, { useEffect, useState } from "react";
import { Building2, Edit3, Search, Trash2 } from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import { paths } from "@/routes/paths";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import {
  deleteRuangUjian,
  useGetRuangUjian,
} from "@/services/Api/features-api/DataMaster/ruang-ujian.service";

function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

const RuangUjianTables: React.FC = () => {
  const navigate = useNavigate();

  const [kataKunci, setKataKunci] = useState("");
  const [batasData, setBatasData] = useState(12);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);

  const [modalKonfirmasiTerbuka, setModalKonfirmasiTerbuka] = useState(false);
  const [idRuangAkanDihapus, setIdRuangAkanDihapus] = useState<number | null>(null);
  const [sedangHapus, setSedangHapus] = useState(false);
  const [deleteErrorMsg, setDeleteErrorMsg] = useState("");

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);

  // Hook: fetch ruang ujian with filters
  const {
    data: rawData,
    loading,
    error: errorMsg,
    refetch: refetchRuangUjian,
  } = useGetRuangUjian({
    q: debouncedKataKunci.trim() || undefined,
    limit: batasData,
    offset: (halamanSaatIni - 1) * batasData,
  });

  const daftarRuangUjian = rawData ?? [];

  useEffect(() => {
    setHalamanSaatIni(1);
  }, [debouncedKataKunci, batasData]);

  const resetFilter = () => {
    setKataKunci("");
    setHalamanSaatIni(1);
  };

  const handleOpenDeleteConfirm = (id_ruangan: number) => {
    setIdRuangAkanDihapus(id_ruangan);
    setModalKonfirmasiTerbuka(true);
  };

  const handleConfirmDelete = async () => {
    if (!idRuangAkanDihapus) return;

    setSedangHapus(true);
    setDeleteErrorMsg("");

    try {
      await deleteRuangUjian(idRuangAkanDihapus);
      setModalKonfirmasiTerbuka(false);
      setIdRuangAkanDihapus(null);
      await refetchRuangUjian();
    } catch (e) {
      setDeleteErrorMsg(
        getUserFriendlyErrorMessage(e, {
          action: "delete",
          entity: "ruang ujian",
        }),
      );
    } finally {
      setSedangHapus(false);
    }
  };

  const totalData = daftarRuangUjian.length;
  const awalData = totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + 1;
  const akhirData =
    totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + totalData;
  const bisaSebelumnya = halamanSaatIni > 1;
  const bisaSelanjutnya = totalData === batasData;

  return (
    <div className="w-full space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Data Ruang Ujian
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola daftar ruangan yang digunakan untuk ujian.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <AddButton
            label="Tambah Ruang"
            onClick={() => navigate(`${paths.dashboard.tambah_data_master_ruang}`)}
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
                placeholder="Cari nama ruang / kode ruang..."
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

        {(errorMsg || deleteErrorMsg) && (
          <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
            {errorMsg || deleteErrorMsg}
          </div>
        )}

        {loading ? (
          <div className="py-12 text-center text-sm text-slate-500">Memuat data...</div>
        ) : (
          <div className="overflow-x-auto rounded-xl border border-slate-200">
            <table className="min-w-full divide-y divide-slate-200 text-sm">
              <thead className="bg-slate-50">
                <tr>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">Nama Ruangan</th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">Kode Ruang</th>
                  <th className="px-6 py-3 text-right font-semibold text-slate-600">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {daftarRuangUjian.length > 0 ? (
                  daftarRuangUjian.map((ruang) => (
                    <tr key={ruang.id_ruangan} className="transition-colors hover:bg-slate-50">
                      <td className="px-6 py-4">
                        <div className="flex flex-col">
                          <span className="font-semibold text-slate-900">{ruang.nama_ruangan}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-slate-700">{ruang.kode_ruang}</td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600"
                            title="Edit"
                            onClick={() =>
                              navigate(
                                paths.dashboard.edit_data_master_ruang.replace(
                                  ":id",
                                  String(ruang.id_ruangan),
                                ),
                              )
                            }
                          >
                            <Edit3 className="h-4 w-4" />
                          </button>
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                            title="Hapus"
                            onClick={() => handleOpenDeleteConfirm(ruang.id_ruangan)}
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
                        <Building2 className="h-10 w-10 text-slate-300" />
                        <p className="text-base font-medium text-slate-900">
                          Tidak ada ruang ujian ditemukan
                        </p>
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
        title="Konfirmasi Hapus Ruang Ujian"
        message="Data ruang ujian yang dihapus tidak bisa dikembalikan. Lanjutkan?"
        onClose={() => {
          if (sedangHapus) return;
          setModalKonfirmasiTerbuka(false);
          setIdRuangAkanDihapus(null);
        }}
        onConfirm={handleConfirmDelete}
        isLoading={sedangHapus}
        confirmLabel="Ya, Hapus"
        loadingLabel="Menghapus..."
      />
    </div>
  );
};

export default RuangUjianTables;
