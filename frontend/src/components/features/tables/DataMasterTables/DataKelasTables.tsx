import React, { useEffect, useMemo, useState } from "react";
import {
  Archive,
  ChevronDown,
  Edit3,
  GraduationCap,
  Search,
  Trash2,
} from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import {
  deleteNamaKelas,
  useGetDataKelasFull,
} from "@/services/Api/features-api/DataMaster/kelas.service";
import { paths } from "@/routes/paths";
import toast from "react-hot-toast";

import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";

function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

const DataKelasTables: React.FC = () => {
  const navigate = useNavigate();

  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");
  const [tingkatKelas, setTingkatKelas] = useState<number | null>(null);

  const [idTerpilih, setIdTerpilih] = useState<Set<number>>(new Set());
  const [batasData, setBatasData] = useState(12);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);
  const [modalKonfirmasiTerbuka, setModalKonfirmasiTerbuka] = useState(false);
  const [sedangMemprosesKonfirmasi, setSedangMemprosesKonfirmasi] =
    useState(false);
  const [aksiKonfirmasi, setAksiKonfirmasi] = useState<
    null | (() => Promise<void>)
  >(null);

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);

  // Hook: fetch kelas options on mount (no filters)
  const { data: kelasOptionsData } = useGetDataKelasFull();
  const opsiTingkat = kelasOptionsData?.item_tingkat_kelas ?? [];

  // Hook: fetch filtered kelas data
  const {
    data: kelasFilteredData,
    loading,
    error: errorMsg,
    refetch: refetchKelas,
  } = useGetDataKelasFull({
    search: debouncedKataKunci.trim() || undefined,
    tingkatKelas: tingkatKelas || undefined,
    limit: batasData,
    offset: (halamanSaatIni - 1) * batasData,
  });

  const daftarKelas = kelasFilteredData?.item_nama_kelas ?? [];

  const tingkatById = useMemo(
    () =>
      new Map(
        opsiTingkat.map((tingkat) => [
          tingkat.id_tingkat_kelas,
          tingkat.tingkat_kelas,
        ]),
      ),
    [opsiTingkat],
  );

  const totalData = daftarKelas.length;
  const dataTerlihat = daftarKelas;

  const semuaTerlihatTerpilih =
    dataTerlihat.length > 0 &&
    dataTerlihat.every((k) => idTerpilih.has(k.id_nama_kelas));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        dataTerlihat.forEach((kelas) => next.delete(kelas.id_nama_kelas));
      } else {
        dataTerlihat.forEach((kelas) => next.add(kelas.id_nama_kelas));
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

  const handleDeleteKelas = async (idNamaKelas: number) => {
    try {
      await deleteNamaKelas(idNamaKelas);
      await refetchKelas();
      setIdTerpilih((prev) => {
        const next = new Set(prev);
        next.delete(idNamaKelas);
        return next;
      });
      toast.success("Berhasil menghapus data kelas");
    } catch (e) {
      toast.error(
        getUserFriendlyErrorMessage(e, {
          action: "delete",
          entity: "kelas",
        }),
      );
    }
  };

  const handleBulkDelete = async () => {
    if (idTerpilih.size === 0) return;
    const ids = Array.from(idTerpilih);
    try {
      await Promise.all(ids.map((id) => deleteNamaKelas(id)));
      await refetchKelas();
      setIdTerpilih(new Set());
      setDropdownAksiTerbuka(false);
      toast.success("Berhasil menghapus data kelas terpilih");
    } catch {
      toast.error("Gagal menghapus data kelas terpilih");
    }
  };

  const pesanKonfirmasiHapusKelas =
    "Apakah anda yakin ingin menghapus data kelas ini?";

  const bukaModalKonfirmasiHapus = (action: () => Promise<void>) => {
    setAksiKonfirmasi(() => action);
    setModalKonfirmasiTerbuka(true);
  };

  const tutupModalKonfirmasi = () => {
    if (sedangMemprosesKonfirmasi) return;
    setModalKonfirmasiTerbuka(false);
    setAksiKonfirmasi(null);
  };

  const jalankanAksiKonfirmasi = async () => {
    if (!aksiKonfirmasi) return;

    setSedangMemprosesKonfirmasi(true);
    try {
      await aksiKonfirmasi();
      setModalKonfirmasiTerbuka(false);
      setAksiKonfirmasi(null);
    } finally {
      setSedangMemprosesKonfirmasi(false);
    }
  };

  const jumlahTerpilih = idTerpilih.size;

  const resetFilter = () => {
    setKataKunci("");
    setHalamanSaatIni(1);
    setTingkatKelas(null);
  };

  useEffect(() => {
    setHalamanSaatIni(1);
  }, [debouncedKataKunci, batasData, tingkatKelas]);

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
            Data Kelas
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola daftar kelas, tingkat, serta informasi terkait.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <AddButton
            label="Tambah Kelas"
            onClick={() =>
              navigate(`${paths.dashboard.tambah_data_master_kelas}`)
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
                placeholder="Cari nama kelas atau tingkat..."
              />
            </div>
          </div>

          <div className="grid w-full grid-cols-1 gap-3 sm:grid-cols-2">
            <div>
              <label className="text-xs font-medium text-slate-600">
                Tingkat Kelas
              </label>
              <select
                value={tingkatKelas ?? ""}
                onChange={(e) =>
                  setTingkatKelas(
                    e.target.value === "" ? null : Number(e.target.value),
                  )
                }
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {opsiTingkat.map((tingkat) => (
                  <option
                    key={tingkat.tingkat_kelas}
                    value={tingkat.tingkat_kelas}
                  >
                    {tingkat.tingkat_kelas}
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
                        onClick={() => {
                          setDropdownAksiTerbuka(false);
                          bukaModalKonfirmasiHapus(handleBulkDelete);
                        }}
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
              Menampilkan <span className="font-medium">{totalData}</span>{" "}
              hasil.
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
                  <div className="flex items-center">
                    <input
                      type="checkbox"
                      checked={semuaTerlihatTerpilih}
                      onChange={togglePilihSemuaTerlihat}
                      className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                    />
                  </div>
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Tingkat
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Nama Kelas
                </th>
                <th scope="col" className="px-6 py-3 text-right font-semibold">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {dataTerlihat.length > 0 ? (
                dataTerlihat.map((kelas) => (
                  <tr
                    key={kelas.id_nama_kelas}
                    className={`transition-colors hover:bg-slate-50 ${
                      idTerpilih.has(kelas.id_nama_kelas)
                        ? "bg-indigo-50/30"
                        : ""
                    }`}
                  >
                    <td className="p-4">
                      <input
                        type="checkbox"
                        checked={idTerpilih.has(kelas.id_nama_kelas)}
                        onChange={() => togglePilihBaris(kelas.id_nama_kelas)}
                        className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                      />
                    </td>
                    <td className="px-6 py-4 font-medium text-slate-900">
                      {tingkatById.get(kelas.id_tingkat_kelas) ?? "-"}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col">
                        <span className="font-semibold text-slate-900">
                          {kelas.nama_kelas}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600"
                          title="Edit"
                          onClick={() =>
                            navigate(
                              paths.dashboard.edit_data_master_kelas
                                .replace(
                                  ":idTingkatKelas",
                                  String(kelas.id_tingkat_kelas),
                                )
                                .replace(
                                  ":idNamaKelas",
                                  String(kelas.id_nama_kelas),
                                ),
                            )
                          }
                        >
                          <Edit3 className="h-4 w-4" />
                        </button>
                        <button
                          className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                          title="Hapus"
                          onClick={() =>
                            bukaModalKonfirmasiHapus(() =>
                              handleDeleteKelas(kelas.id_nama_kelas),
                            )
                          }
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4} className="px-6 py-12 text-center">
                    <div className="flex flex-col items-center justify-center gap-2">
                      <GraduationCap className="h-10 w-10 text-slate-300" />
                      <p className="text-base font-medium text-slate-900">
                        Tidak ada kelas ditemukan
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
        isOpen={modalKonfirmasiTerbuka}
        message={pesanKonfirmasiHapusKelas}
        onClose={tutupModalKonfirmasi}
        onConfirm={() => void jalankanAksiKonfirmasi()}
        isLoading={sedangMemprosesKonfirmasi}
      />
    </div>
  );
};

export default DataKelasTables;
