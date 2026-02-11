import React, { useEffect, useState } from "react";
import {
  Search,
  Eye,
  EyeOff,
  ChevronDown,
  Trash2,
  Archive,
  Edit3,
  User,
  CheckCircle2,
  XCircle,
  Trash,
} from "lucide-react";
import AddButton from "@/components/common/Button/AddButton";
import { useNavigate } from "react-router";

import type { StatusAkun } from "@/types/OpsiTypes/Option";
import type { DataGuru } from "@/types/KelolaAkun/AkunGuru";
import { paths } from "@/routes/paths";
import { GetAllGuru } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import {
  DeletePengguna,
  DeletePenggunaBulk,
} from "@/services/Api/features-api/KelolaAkun/akun.service";
import toast from "react-hot-toast";

// --- Helper Functions ---
const getStatusBadge = (status: StatusAkun) => {
  switch (status) {
    case "AKTIF":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700 ring-1 ring-inset ring-emerald-600/20">
          <CheckCircle2 className="h-3 w-3" /> Aktif
        </span>
      );
    case "NONAKTIF":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-rose-50 px-2.5 py-0.5 text-xs font-medium text-rose-700 ring-1 ring-inset ring-rose-600/20">
          <XCircle className="h-3 w-3" /> Nonaktif
        </span>
      );
  }
};

function samarkanNomorHp(nomorHp?: string | null) {
  const digit = String(nomorHp ?? "").replace(/\s+/g, "");
  if (digit.length <= 4) return digit;
  const terlihat = digit.slice(-3);
  return `${digit.slice(0, 3)}****${terlihat}`;
}

function samarkanNip(nip?: string | null) {
  const digit = String(nip ?? "").replace(/\s+/g, "");
  if (digit.length <= 6) return digit;
  const terlihat = digit.slice(-4);
  return `${digit.slice(0, 4)}****${terlihat}`;
}

// --- Komponen Utama ---
const AkunGuruTables: React.FC = () => {
  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");
  const [kataKunciDebounce, setKataKunciDebounce] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusAkun>("AKTIF");
  const [idTerpilih, setIdTerpilih] = useState<Set<number>>(new Set());
  const [samarkanDataSensitif, setSamarkanDataSensitif] = useState(true);
  const [batasData, setBatasData] = useState(10);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);
  const [modalKonfirmasiTerbuka, setModalKonfirmasiTerbuka] = useState(false);
  const [sedangMemprosesKonfirmasi, setSedangMemprosesKonfirmasi] =
    useState(false);
  const [aksiKonfirmasi, setAksiKonfirmasi] = useState<null | (() => Promise<void>)>(
    null,
  );

  const [daftarPengguna, setDaftarPengguna] = useState<DataGuru[]>([]);

  useEffect(() => {
    const handle = window.setTimeout(() => {
      setKataKunciDebounce(kataKunci.trim());
    }, 500);

    return () => {
      window.clearTimeout(handle);
    };
  }, [kataKunci]);

  useEffect(() => {
    setHalamanSaatIni(1);
  }, [kataKunciDebounce, statusFilter, batasData]);

  useEffect(() => {
    let aktif = true;

    const fetchGuru = async () => {
      const data = await GetAllGuru({
        q: kataKunciDebounce,
        status: statusFilter,
        limit: batasData,
        offset: (halamanSaatIni - 1) * batasData,
      });
      if (aktif) {
        setDaftarPengguna(data);
      }
    };

    fetchGuru();

    return () => {
      aktif = false;
    };
  }, [kataKunciDebounce, statusFilter, batasData, halamanSaatIni]);

  const penggunaTersaring = daftarPengguna;

  const semuaTerlihatTerpilih =
    penggunaTersaring.length > 0 &&
    penggunaTersaring.every((p) => idTerpilih.has(p.id_pengguna));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((sebelumnya) => {
      const berikutnya = new Set(sebelumnya);
      if (semuaTerlihatTerpilih) {
        penggunaTersaring.forEach((p) => berikutnya.delete(p.id_pengguna));
      } else {
        penggunaTersaring.forEach((p) => berikutnya.add(p.id_pengguna));
      }
      return berikutnya;
    });
  };

  const togglePilihBaris = (id: number) => {
    setIdTerpilih((sebelumnya) => {
      const berikutnya = new Set(sebelumnya);
      if (berikutnya.has(id)) berikutnya.delete(id);
      else berikutnya.add(id);
      return berikutnya;
    });
  };

  const navigate = useNavigate();

  const fetchGuruList = async () => {
    return GetAllGuru({
      q: kataKunciDebounce,
      status: statusFilter,
      limit: batasData,
      offset: (halamanSaatIni - 1) * batasData,
    });
  };

  const DeleteUser = async (id: number) => {
    try {
      await DeletePengguna(id);

      const data = await fetchGuruList();
      setDaftarPengguna(data);
      setIdTerpilih((sebelumnya) => {
        const berikutnya = new Set(sebelumnya);
        berikutnya.delete(id);
        return berikutnya;
      });

      toast.success("Berhasil menghapus akun guru");
    } catch {
      toast.error("Gagal menghapus akun guru");
    }
  };

  const handleBulkDelete = async () => {
    if (idTerpilih.size === 0) return;
    const ids = Array.from(idTerpilih);
    try {
      await DeletePenggunaBulk(ids);

      const data = await fetchGuruList();
      setDaftarPengguna(data);
      setIdTerpilih(new Set());
      setDropdownAksiTerbuka(false);

      toast.success("Berhasil menghapus akun guru terpilih");
    } catch {
      toast.error("Gagal menghapus akun guru terpilih");
    }
  };

  const pesanKonfirmasiHapusGuru =
    "Apakah anda yakin ingin menghapus akun ini? semua data yang berkaitan dengan akun ini akan dihapus (data bank soal, data aktivitas, data ujian yang dibuat, dan data hasil ujian yang dibuat  )";

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
  const totalTerlihat = daftarPengguna.length;
  const awalData = totalTerlihat === 0 ? 0 : (halamanSaatIni - 1) * batasData + 1;
  const akhirData =
    totalTerlihat === 0 ? 0 : (halamanSaatIni - 1) * batasData + totalTerlihat;
  const bisaSebelumnya = halamanSaatIni > 1;
  const bisaSelanjutnya = totalTerlihat === batasData;

  return (
    <div className="w-full space-y-6">
      {/* --- Header Section --- */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Manajemen Akun Guru
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola data guru, status akun, dan informasi kepegawaian.
          </p>
        </div>
        <div className="flex items-center gap-3">
          {/* Privacy Toggle */}
          <button
            onClick={() => setSamarkanDataSensitif(!samarkanDataSensitif)}
            className={`group inline-flex cursor-pointer items-center gap-2 rounded-lg border px-3 py-2 text-sm font-medium transition-all ${
              samarkanDataSensitif
                ? "border-slate-200 bg-white text-slate-700 hover:bg-slate-50"
                : "border-slate-200 bg-slate-100 text-slate-900 hover:bg-slate-50"
            }`}
          >
            {samarkanDataSensitif ? (
              <>
                <EyeOff className="h-4 w-4 text-slate-400 group-hover:text-slate-600" />
                <span className="hidden sm:inline">Tampilkan Data</span>
              </>
            ) : (
              <>
                <Eye className="h-4 w-4 text-slate-400 group-hover:text-slate-600" />
                <span className="hidden sm:inline">Samarkan Data</span>
              </>
            )}
          </button>
          <AddButton
            label="Tambah Akun Guru"
            onClick={() => navigate(`${paths.dashboard.tambah_guru}`)}
          />
        </div>
      </div>

      {/* --- Filter & Action Bar --- */}
      <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between">
        <div className="flex w-full flex-col gap-3 sm:w-auto sm:flex-row sm:items-center">
          <div className="relative w-full sm:w-80">
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
              <Search className="h-4 w-4 text-slate-400" />
            </div>
            <input
              type="text"
              value={kataKunci}
              onChange={(e) => setKataKunci(e.target.value)}
              className="block w-full rounded-lg border border-slate-200 bg-slate-50 py-2 pl-10 pr-3 text-sm text-slate-900 placeholder:text-slate-400 focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              placeholder="Cari nama, NIP, atau email..."
            />
          </div>
          <div className="flex w-full flex-wrap items-center gap-2 sm:w-auto">
            {(["AKTIF", "NONAKTIF"] as const).map((status) => (
              <button
                key={status}
                type="button"
                onClick={() => setStatusFilter(status)}
                className={`rounded-full border px-3 py-1.5 text-xs font-semibold transition ${
                  statusFilter === status
                    ? "border-[#397e50] bg-[#397e50]/10 text-[#2f6b44]"
                    : "border-slate-200 text-slate-500 hover:border-slate-300 hover:text-slate-700"
                }`}
              >
                {status === "AKTIF" ? "Aktif" : "Nonaktif"}
              </button>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* Bulk Actions Dropdown (Muncul jika ada yang dipilih) */}
          {jumlahTerpilih > 0 && (
            <div className="relative">
              <button
                onClick={() => setDropdownAksiTerbuka(!dropdownAksiTerbuka)}
                className="inline-flex items-center gap-2 rounded-lg bg-slate-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-slate-800"
              >
                {jumlahTerpilih} Terpilih
                <ChevronDown className="h-4 w-4" />
              </button>

              {dropdownAksiTerbuka && (
                <>
                  <div
                    className="fixed inset-0 z-10"
                    onClick={() => setDropdownAksiTerbuka(false)}
                  />
                  <div className="absolute right-0 z-20 mt-2 w-48 origin-top-right rounded-lg border border-slate-100 bg-white shadow-xl ring-1 ring-black ring-opacity-5 focus:outline-none">
                    <div className="p-1">
                      <button className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-slate-700 hover:bg-slate-50">
                        <Archive className="h-4 w-4 text-slate-400" /> Arsipkan
                      </button>
                      <button
                        className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-rose-600 hover:bg-rose-50"
                        onClick={() => {
                          setDropdownAksiTerbuka(false);
                          bukaModalKonfirmasiHapus(handleBulkDelete);
                        }}
                      >
                        <Trash2 className="h-4 w-4" /> Hapus Data
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          )}
        </div>
      </div>

      {/* --- Table Section --- */}
      <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="border-b border-slate-200 bg-slate-50 text-xs uppercase text-slate-500">
              <tr>
                <th scope="col" className="p-4 w-4">
                  <div className="flex items-center">
                    <input
                      type="checkbox"
                      checked={semuaTerlihatTerpilih}
                      onChange={togglePilihSemuaTerlihat}
                      className="h-4 w-4 rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                    />
                  </div>
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Guru
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  NIP & Kontak
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Jabatan
                </th>

                <th scope="col" className="px-6 py-3 font-semibold">
                  Status
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Role Akun
                </th>
                <th scope="col" className="px-6 py-3 text-right font-semibold">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {penggunaTersaring.length > 0 ? (
                penggunaTersaring.map((p) => (
                  <tr
                    key={p.id_pengguna}
                    className={`transition-colors hover:bg-slate-50 ${
                      idTerpilih.has(p.id_pengguna) ? "bg-indigo-50/30" : ""
                    }`}
                  >
                    <td className="p-4">
                      <div className="flex items-center">
                        <input
                          type="checkbox"
                          checked={idTerpilih.has(p.id_pengguna)}
                          onChange={() => togglePilihBaris(p.id_pengguna)}
                          className="h-4 w-4 rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                        />
                      </div>
                    </td>
                    <td className="min-w-[260px] px-6 py-4">
                      <div className="flex items-center gap-3">
                        <button
                          className="cursor-pointer shrink-0"
                          onClick={() =>
                            window.open(
                              `${import.meta.env.VITE_API_URL}${p.foto_profil}`,
                            )
                          }
                        >
                          <img
                            className="h-10 w-10 shrink-0 rounded-full object-cover ring-2 ring-white"
                            src={`${import.meta.env.VITE_API_URL}${p.foto_profil}`}
                            alt=""
                          />
                        </button>
                        <div className="min-w-0 flex flex-col">
                          <span className="font-semibold text-slate-900">
                            {p.nama_lengkap}
                          </span>
                          <span className="text-xs text-slate-500">
                            @{p.username}
                          </span>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col gap-1">
                        <span className="font-mono text-xs font-medium text-slate-700">
                          {samarkanDataSensitif ? samarkanNip(p.nip) : p.nip}
                        </span>
                        <span className="text-xs text-slate-500">
                          {p.email}
                        </span>
                        <span className="text-xs text-slate-500">
                          {samarkanDataSensitif
                            ? samarkanNomorHp(p.no_hp)
                            : p.no_hp}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col">
                        <span className="text-slate-900">{p.jabatan}</span>
                        <span className="text-xs text-slate-500">
                          {p.bidang_studi}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      {getStatusBadge(p.status_akun)}
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2">
                        <span className="text-xs font-semibold text-slate-500 text-center">
                          {p.role}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          className="rounded-lg cursor-pointer p-2 text-slate-400 hover:bg-slate-100 hover:text-green-600 transition-colors"
                          title="Edit"
                          onClick={() =>
                            navigate(
                              paths.dashboard.edit_guru.replace(
                                ":id",
                                String(p.id_pengguna),
                              ),
                            )
                          }
                        >
                          <Edit3 className="h-4 w-4" />
                        </button>
                        <button
                          className="rounded-lg p-2 text-slate-400 cursor-pointer hover:bg-slate-100 hover:text-red-600 transition-colors"
                          title="Detail"
                          onClick={() =>
                            bukaModalKonfirmasiHapus(() =>
                              DeleteUser(Number(p.id_pengguna)),
                            )
                          }
                        >
                          <Trash className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center">
                    <div className="flex flex-col items-center justify-center gap-2">
                      <User className="h-10 w-10 text-slate-300" />
                      <p className="text-base font-medium text-slate-900">
                        Tidak ada pengguna ditemukan
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
        {/* Pagination Dummy */}
        <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
          <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm text-slate-700">
                Menampilkan <span className="font-medium">{awalData}</span>{" "}
                sampai <span className="font-medium">{akhirData}</span> dari{" "}
                <span className="font-medium">{akhirData}</span> hasil
              </p>
            </div>
            <div className="flex items-center gap-3">
              <div className=" cursor-pointer flex items-center gap-2 text-sm text-slate-600">
                <span>Tampilkan</span>
                <select
                  value={batasData}
                  onChange={(event) => setBatasData(Number(event.target.value))}
                  className="cursor-pointer appearance-none rounded-lg border border-slate-200 bg-white px-7 py-1 text-sm text-slate-700 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
                >
                  {[10, 20, 30, 40, 50].map((opsi) => (
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
        </div>
      </div>

      {modalKonfirmasiTerbuka && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/55 px-4">
          <div className="w-full max-w-xl rounded-xl border border-slate-200 bg-white p-6 shadow-2xl">
            <h3 className="text-lg font-semibold text-slate-900">
              Konfirmasi Hapus Akun
            </h3>
            <p className="mt-2 text-sm leading-relaxed text-slate-600">
              {pesanKonfirmasiHapusGuru}
            </p>

            <div className="mt-6 flex justify-end gap-3">
              <button
                type="button"
                onClick={tutupModalKonfirmasi}
                disabled={sedangMemprosesKonfirmasi}
                className="cursor-pointer rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Batal
              </button>
              <button
                type="button"
                onClick={() => void jalankanAksiKonfirmasi()}
                disabled={sedangMemprosesKonfirmasi}
                className="cursor-pointer rounded-lg bg-rose-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-rose-700 disabled:cursor-not-allowed disabled:bg-rose-400"
              >
                {sedangMemprosesKonfirmasi ? "Menghapus..." : "Ya, Hapus"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AkunGuruTables;
