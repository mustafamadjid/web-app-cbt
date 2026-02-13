import React, { useEffect, useRef, useState } from "react";
import {
  Archive,
  Building2,
  ChevronDown,
  Edit3,
  Search,
  Trash2,
} from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import { getRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { paths } from "@/routes/paths";

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

  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");

  const [daftarRuangUjian, setDaftarRuangUjian] = useState<RuangUjianRow[]>([]);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const [idTerpilih, setIdTerpilih] = useState<Set<number>>(new Set());
  const [batasData, setBatasData] = useState(12);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);
  const requestSeq = useRef(0);

  useEffect(() => {
    const seq = ++requestSeq.current;

    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");

        const data = await getRuangUjian({
          q: debouncedKataKunci.trim() || undefined,
          limit: batasData,
          offset: (halamanSaatIni - 1) * batasData,
        });

        if (seq !== requestSeq.current) return;

        setDaftarRuangUjian(data);
        setIdTerpilih((prev) => {
          if (prev.size === 0) return prev;
          const ids = new Set(data.map((ruang) => ruang.id));
          const next = new Set<number>();
          prev.forEach((id) => {
            if (ids.has(id)) next.add(id);
          });
          return next;
        });
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data ruang ujian.");
        setDaftarRuangUjian([]);
      } finally {
        if (seq === requestSeq.current) {
          setLoading(false);
        }
      }
    })();
  }, [debouncedKataKunci, batasData, halamanSaatIni]);

  const totalData = daftarRuangUjian.length;
  const dataTerlihat = daftarRuangUjian;

  const semuaTerlihatTerpilih =
    dataTerlihat.length > 0 && dataTerlihat.every((r) => idTerpilih.has(r.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih) {
        dataTerlihat.forEach((r) => next.delete(r.id));
      } else {
        dataTerlihat.forEach((r) => next.add(r.id));
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
    setHalamanSaatIni(1);
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
            Data Ruang Ujian
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola daftar ruangan yang digunakan untuk ujian.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <AddButton
            label="Tambah Ruang"
            onClick={() =>
              navigate(`${paths.dashboard.tambah_data_master_ruang}`)
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
                placeholder="Cari nama ruang..."
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
              Menampilkan{" "}
              <span className="font-medium">{daftarRuangUjian.length}</span>{" "}
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
                  Nama Ruang
                </th>
                <th scope="col" className="px-6 py-3 text-right font-semibold">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {dataTerlihat.length > 0 ? (
                dataTerlihat.map((ruang) => (
                  <tr
                    key={ruang.id}
                    className={`transition-colors hover:bg-slate-50 ${
                      idTerpilih.has(ruang.id) ? "bg-indigo-50/30" : ""
                    }`}
                  >
                    <td className="p-4">
                      <input
                        type="checkbox"
                        checked={idTerpilih.has(ruang.id)}
                        onChange={() => togglePilihBaris(ruang.id)}
                        className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                      />
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col">
                        <span className="font-semibold text-slate-900">
                          {ruang.namaRuangan}
                        </span>
                        <span className="text-xs text-slate-500">
                          ID: {ruang.id}
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
                              paths.dashboard.edit_data_master_ruang.replace(
                                ":id",
                                String(ruang.id),
                              ),
                            )
                          }
                        >
                          <Edit3 className="h-4 w-4" />
                        </button>
                        <button
                          className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                          title="Hapus"
                          onClick={() => setIdTerpilih(new Set([ruang.id]))}
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
    </div>
  );
};

export default RuangUjianTables;
