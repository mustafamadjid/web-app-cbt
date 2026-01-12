import React, { useEffect, useRef, useState } from "react";
import {
  Archive,
  ChevronDown,
  Clock3,
  Edit3,
  Search,
  Trash2,
} from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import type { SesiRow } from "@/types/DataMaster/Sesi";
import { getSesi } from "@/services/Api/features-api/DataMaster/sesi.service";
import { paths } from "@/routes/paths";

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

  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");

  const [daftarSesi, setDaftarSesi] = useState<SesiRow[]>([]);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const [idTerpilih, setIdTerpilih] = useState<Set<number>>(new Set());

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);
  const requestSeq = useRef(0);

  useEffect(() => {
    const seq = ++requestSeq.current;

    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");

        const data = await getSesi({
          q: debouncedKataKunci.trim() || undefined,
        });

        if (seq !== requestSeq.current) return;

        setDaftarSesi(data);
        setIdTerpilih((prev) => {
          if (prev.size === 0) return prev;
          const ids = new Set(data.map((sesi) => sesi.id));
          const next = new Set<number>();
          prev.forEach((id) => {
            if (ids.has(id)) next.add(id);
          });
          return next;
        });
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data sesi.");
        setDaftarSesi([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, [debouncedKataKunci]);

  const semuaTerlihatTerpilih =
    daftarSesi.length > 0 && daftarSesi.every((s) => idTerpilih.has(s.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        daftarSesi.forEach((s) => next.delete(s.id));
      } else {
        daftarSesi.forEach((s) => next.add(s.id));
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

  const jumlahTerpilih = idTerpilih.size;

  const resetFilter = () => {
    setKataKunci("");
  };

  return (
    <div className="w-full space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Data Sesi
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola daftar sesi ujian sesuai kebutuhan jadwal.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <AddButton
            label="Tambah Sesi"
            onClick={() =>
              navigate(`${paths.dashboard.tambah_data_master_sesi}`)
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
                placeholder="Cari kode atau nama sesi..."
              />
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
                        onClick={() => setDropdownAksiTerbuka(false)}
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
              Menampilkan <span className="font-medium">{daftarSesi.length}</span>{" "}
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
                  <input
                    type="checkbox"
                    checked={semuaTerlihatTerpilih}
                    onChange={togglePilihSemuaTerlihat}
                    className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                  />
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Kode Sesi
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Nama Sesi
                </th>
                <th scope="col" className="px-6 py-3 text-right font-semibold">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {daftarSesi.length > 0 ? (
                daftarSesi.map((sesi) => (
                  <tr
                    key={sesi.id}
                    className={`transition-colors hover:bg-slate-50 ${
                      idTerpilih.has(sesi.id) ? "bg-indigo-50/30" : ""
                    }`}
                  >
                    <td className="p-4">
                      <input
                        type="checkbox"
                        checked={idTerpilih.has(sesi.id)}
                        onChange={() => togglePilihBaris(sesi.id)}
                        className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                      />
                    </td>
                    <td className="px-6 py-4 font-semibold text-slate-900">
                      {sesi.kodeSesi}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col">
                        <span className="font-semibold text-slate-900">
                          {sesi.namaSesi}
                        </span>
                        <span className="text-xs text-slate-500">
                          ID: {sesi.id}
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
                              "data-master/tambah-sesi"
                            )
                          }
                        >
                          <Edit3 className="h-4 w-4" />
                        </button>
                        <button
                          className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                          title="Hapus"
                          onClick={() => setIdTerpilih(new Set([sesi.id]))}
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
                      <Clock3 className="h-10 w-10 text-slate-300" />
                      <p className="text-base font-medium text-slate-900">
                        Tidak ada sesi ditemukan
                      </p>
                      <p className="text-sm text-slate-500">
                        Coba sesuaikan kata kunci pencarian Anda.
                      </p>
                    </div>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
          <div className="flex flex-1 items-center justify-between">
            <p className="text-sm text-slate-700">
              Menampilkan <span className="font-medium">1</span> sampai{" "}
              <span className="font-medium">{daftarSesi.length}</span> dari{" "}
              <span className="font-medium">{daftarSesi.length}</span> hasil
            </p>
            <p className="hidden text-xs text-slate-500 sm:block">
              Geser tabel ke kanan/kiri untuk melihat kolom lainnya.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DataSesiTables;
