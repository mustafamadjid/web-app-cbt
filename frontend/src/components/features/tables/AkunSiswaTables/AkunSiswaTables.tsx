import React, { useEffect, useMemo, useState } from "react";
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
  Trash,
  XCircle,
} from "lucide-react";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import type { StatusAkun, JenisKelamin } from "@/types/OpsiTypes/Option";
import { useNavigate } from "react-router";

import { useGetListSiswa } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import {
  DeletePengguna,
  DeletePenggunaBulk,
} from "@/services/Api/features-api/KelolaAkun/akun.service";
import { useGetJenisKelaminOptions } from "@/services/Api/features-api/GetOptions/options.service";

import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";

import { paths } from "@/routes/paths";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { tahunOption } from "@/helper/TahunOption/TahunOption";

/** ===== Helpers ===== */
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
    default:
      return null;
  }
};

function samarkanNomorHp(nomorHp: string) {
  const digit = nomorHp.replace(/\s+/g, "");
  if (digit.length <= 4) return digit;
  const terlihat = digit.slice(-3);
  return `${digit.slice(0, 3)}****${terlihat}`;
}

function samarkanEmail(email: string) {
  const e = email.trim();
  const [user, domain] = e.split("@");
  if (!domain) return e;
  if (user.length <= 2) return `**@${domain}`;
  return `${user.slice(0, 2)}****@${domain}`;
}

function formatTanggalIndo(yyyyMmDd: string) {
  if (!yyyyMmDd) return "-";
  const [y, m, d] = yyyyMmDd.split("-");
  if (!y || !m || !d) return yyyyMmDd;
  return `${d}/${m}/${y}`;
}

const labelGender: Record<JenisKelamin, string> = {
  LAKI_LAKI: "Laki-laki",
  PEREMPUAN: "Perempuan",
};

function useDebouncedValue<T>(value: T, delayMs = 300) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(t);
  }, [value, delayMs]);
  return debounced;
}

/** ===== Component ===== */
const AkunSiswaTables: React.FC = () => {
  const navigate = useNavigate();

  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");

  // FILTERS (backend-driven)
  const [angkatan, setAngkatan] = useState<string>("");
  const [tingkatKelasId, setTingkatKelasId] = useState<number | null>(null);
  const [jenisKelamin, setJenisKelamin] = useState<string>("");

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);

  // OPTIONS from server via hooks
  const opsiAngkatan = useMemo(() => tahunOption(), []);
  const { data: kelasData } = useGetDataKelasFull();
  const opsiTingkatKelas = kelasData?.item_tingkat_kelas ?? [];
  const { data: opsiGender } = useGetJenisKelaminOptions();

  // Selection / privacy
  const [idTerpilih, setIdTerpilih] = useState<Set<number>>(new Set());
  const [samarkanDataSensitif, setSamarkanDataSensitif] = useState(true);
  const [batasData, setBatasData] = useState(12);
  const [halamanSaatIni, setHalamanSaatIni] = useState(1);
  const [modalKonfirmasiTerbuka, setModalKonfirmasiTerbuka] = useState(false);
  const [sedangMemprosesKonfirmasi, setSedangMemprosesKonfirmasi] =
    useState(false);
  const [aksiKonfirmasi, setAksiKonfirmasi] = useState<
    null | (() => Promise<void>)
  >(null);

  // Hook: fetch siswa list with filters
  const {
    data: rawSiswa,
    loading,
    error: errorMsg,
    refetch: refetchSiswa,
  } = useGetListSiswa({
    q: debouncedKataKunci.trim() || undefined,
    angkatan: angkatan ? Number(angkatan) : undefined,
    tingkatKelas: tingkatKelasId ?? undefined,
    jenisKelamin: (jenisKelamin as JenisKelamin) || undefined,
    limit: batasData,
    offset: (halamanSaatIni - 1) * batasData,
  });

  const daftarSiswa = rawSiswa ?? [];
  const siswaTersaring = daftarSiswa;
  const totalData = siswaTersaring.length;
  const siswaTerlihat = siswaTersaring;

  const semuaTerlihatTerpilih =
    siswaTerlihat.length > 0 &&
    siswaTerlihat.every((s) => idTerpilih.has(s.id_pengguna));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih)
        siswaTerlihat.forEach((s) => next.delete(s.id_pengguna));
      else siswaTerlihat.forEach((s) => next.add(s.id_pengguna));
      return next;
    });
  };

  const togglePilihBaris = (id_pengguna: number) => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (next.has(id_pengguna)) next.delete(id_pengguna);
      else next.add(id_pengguna);
      return next;
    });
  };

  const jumlahTerpilih = idTerpilih.size;

  const resetFilter = () => {
    setKataKunci("");
    setAngkatan("");
    setTingkatKelasId(null);
    setJenisKelamin("");
    setHalamanSaatIni(1);
  };

  const handleDeleteSiswa = async (id_pengguna: number) => {
    await DeletePengguna(id_pengguna);
    await refetchSiswa();
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      next.delete(id_pengguna);
      return next;
    });
  };

  const handleBulkDelete = async () => {
    if (idTerpilih.size === 0) return;
    const ids = Array.from(idTerpilih);
    await DeletePenggunaBulk(ids);
    await refetchSiswa();
    setIdTerpilih(new Set());
  };

  const pesanKonfirmasiHapusSiswa =
    "Apakah anda yakin ingin menghapus akun? semua data yang berkaitan dengan akun ini akan ikut terhapus (data ujian, data hasil ujian, dan data kartu ujian )";

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

  useEffect(() => {
    setHalamanSaatIni(1);
  }, [debouncedKataKunci, angkatan, tingkatKelasId, jenisKelamin, batasData]);

  const awalData =
    totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + 1;
  const akhirData =
    totalData === 0 ? 0 : (halamanSaatIni - 1) * batasData + totalData;
  const bisaSebelumnya = halamanSaatIni > 1;
  const bisaSelanjutnya = totalData === batasData;

  return (
    <div className="w-full space-y-6">
      {/* ===== Header Section ===== */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Manajemen Akun Siswa
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola data siswa, status akun, dan informasi kelas.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => setSamarkanDataSensitif((v) => !v)}
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
            label="Tambah Akun Siswa"
            onClick={() => navigate(`${paths.dashboard.tambah_siswa}`)}
          />
        </div>
      </div>

      {/* ===== Filter & Action Bar ===== */}
      <div className="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <div className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
          {/* Search */}
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
                placeholder="Cari nama, kelas, email, atau nomor absen..."
              />
            </div>
          </div>

          {/* Filters */}
          <div className="grid w-full grid-cols-1 gap-3 sm:grid-cols-3">
            {/* Angkatan */}
            <div>
              <label className="text-xs font-medium text-slate-600">
                Angkatan
              </label>
              <select
                value={angkatan}
                onChange={(e) => setAngkatan(e.target.value)}
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {(opsiAngkatan ?? []).map((a) => (
                  <option key={a} value={String(a)}>
                    {a}
                  </option>
                ))}
              </select>
            </div>

            {/* Tingkat Kelas */}
            <div>
              <label className="text-xs font-medium text-slate-600">
                Tingkat Kelas
              </label>
              <select
                value={tingkatKelasId ?? ""}
                onChange={(e) =>
                  setTingkatKelasId(
                    e.target.value === "" ? null : Number(e.target.value),
                  )
                }
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {opsiTingkatKelas.map((tingkat) => (
                  <option
                    key={tingkat.id_tingkat_kelas}
                    value={tingkat.tingkat_kelas}
                  >
                    Kelas {tingkat.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>

            {/* Jenis Kelamin */}
            <div>
              <label className="text-xs font-medium text-slate-600">
                Jenis Kelamin
              </label>
              <select
                value={jenisKelamin}
                onChange={(e) => setJenisKelamin(e.target.value)}
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {(opsiGender ?? []).map((g) => (
                  <option key={g.value} value={g.value}>
                    {g.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Right side actions */}
          <div className="flex items-center justify-between gap-2 lg:justify-end">
            <button
              type="button"
              onClick={resetFilter}
              className="cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
            >
              Reset Filter
            </button>

            {/* Bulk Actions Dropdown (muncul jika ada yang dipilih) */}
            {jumlahTerpilih > 0 && (
              <div className="relative">
                <button
                  onClick={() => setDropdownAksiTerbuka((v) => !v)}
                  className="inline-flex cursor-pointer items-center gap-2 rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-slate-800"
                >
                  {jumlahTerpilih} Terpilih
                  <ChevronDown className="h-4 w-4" />
                </button>

                {dropdownAksiTerbuka && (
                  <>
                    <div
                      className="fixed inset-0 z-10 cursor-pointer"
                      onClick={() => setDropdownAksiTerbuka(false)}
                    />
                    <div className="absolute right-0 z-20 mt-2 w-48 origin-top-right rounded-lg border border-slate-100 bg-white shadow-xl ring-1 ring-black/5 focus:outline-none">
                      <div className="p-1">
                        <button
                          className="flex w-full cursor-pointer items-center gap-2 rounded-md px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
                          onClick={() => setDropdownAksiTerbuka(false)}
                        >
                          <Archive className="h-4 w-4 text-slate-400" />{" "}
                          Arsipkan
                        </button>
                        <button
                          className="flex w-full cursor-pointer items-center gap-2 rounded-md px-3 py-2 text-sm text-rose-600 hover:bg-rose-50"
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

        {/* Status line */}
        <div className="text-sm text-slate-600">
          {loading ? (
            <span>Memuat data...</span>
          ) : errorMsg ? (
            <span className="text-rose-600">{errorMsg}</span>
          ) : (
            <span>
              Menampilkan{" "}
              <span className="font-medium">{siswaTersaring.length}</span>{" "}
              hasil.
            </span>
          )}
        </div>
      </div>

      {/* ===== Table Section ===== */}
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
                  Siswa
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  NISN
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Jenis Kelamin
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Tingkat Kelas
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Nama Kelas
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  No. Absen & Angkatan
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  TTL
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Kontak
                </th>
                <th scope="col" className="px-6 py-3 font-semibold">
                  Status
                </th>
                <th scope="col" className="px-6 py-3 text-right font-semibold">
                  Aksi
                </th>
              </tr>
            </thead>

            <tbody className="divide-y divide-slate-200">
              {siswaTerlihat.length > 0 ? (
                siswaTerlihat.map((s) => {
                  const hpRaw = (s.no_hp ?? "").trim();
                  const emailRaw = (s.email ?? "").trim();
                  const fotoUrl = resolveImageUrl(s.foto_profil);

                  const hpTampil = !hpRaw
                    ? "-"
                    : samarkanDataSensitif
                      ? samarkanNomorHp(hpRaw)
                      : hpRaw;

                  const emailTampil = !emailRaw
                    ? "-"
                    : samarkanDataSensitif
                      ? samarkanEmail(emailRaw)
                      : emailRaw;

                  const ttlTampil = `${s.tempat_lahir}, ${formatTanggalIndo(
                    s.tanggal_lahir,
                  )}`;

                  return (
                    <tr
                      key={s.id_pengguna}
                      className={`transition-colors hover:bg-slate-50 ${
                        idTerpilih.has(s.id_pengguna) ? "bg-indigo-50/30" : ""
                      }`}
                    >
                      <td className="p-4">
                        <div className="flex items-center">
                          <input
                            type="checkbox"
                            checked={idTerpilih.has(s.id_pengguna)}
                            onChange={() => togglePilihBaris(s.id_pengguna)}
                            className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                          />
                        </div>
                      </td>

                      {/* Siswa */}
                      <td className="min-w-[260px] px-6 py-4">
                        <div className="flex items-center gap-3">
                          <button
                            className="cursor-pointer shrink-0"
                            onClick={() => window.open(fotoUrl)}
                          >
                            <img
                              className="h-10 w-10 shrink-0 rounded-full object-cover ring-2 ring-white"
                              src={
                                fotoUrl ||
                                `https://ui-avatars.com/api/?name=${encodeURIComponent(
                                  s.nama_lengkap,
                                )}&background=random`
                              }
                              alt=""
                            />
                          </button>

                          <div className="min-w-0 flex flex-col">
                            <span className="font-semibold text-slate-900">
                              {s.nama_lengkap}
                            </span>
                            <span className="text-xs text-slate-500">
                              {s.username}
                            </span>
                          </div>
                        </div>
                      </td>
                      {/* NISN */}
                      <td className="px-6 py-4">
                        <div className="flex flex-col gap-1">
                          <span className="text-slate-900">{s.nisn}</span>
                        </div>
                      </td>

                      {/* Jenis Kelamin */}
                      <td className="px-6 py-4">
                        <span className="text-slate-900">
                          {labelGender[s.jenis_kelamin]}
                        </span>
                      </td>

                      {/* Tingkat Kelas */}
                      <td className="px-6 py-4">
                        <span className="text-slate-900">
                          {s.tingkat_kelas ? `Kelas ${s.tingkat_kelas}` : "-"}
                        </span>
                      </td>

                      {/* Nama Kelas */}
                      <td className="px-6 py-4">
                        <span className="text-slate-900">
                          {s.nama_kelas || s.kelas || "-"}
                        </span>
                      </td>

                      {/* No. Absen & Angkatan */}
                      <td className="px-6 py-4">
                        <div className="flex flex-col gap-1">
                          <span className="text-xs text-slate-500">
                            No Absen:{" "}
                            <span className="font-medium text-slate-700">
                              {s.no_absen}
                            </span>{" "}
                            • Angkatan:{" "}
                            <span className="font-medium text-slate-700">
                              {s.angkatan}
                            </span>
                          </span>
                        </div>
                      </td>

                      {/* TTL */}
                      <td className="px-6 py-4">
                        <span className="text-slate-700">{ttlTampil}</span>
                      </td>

                      {/* Kontak */}
                      <td className="px-6 py-4">
                        <div className="flex flex-col gap-1">
                          <span className="text-xs text-slate-500">
                            {emailTampil}
                          </span>
                          <span className="text-xs text-slate-500">
                            {hpTampil}
                          </span>
                        </div>
                      </td>

                      {/* Status */}
                      <td className="px-6 py-4">
                        {getStatusBadge(s.status_akun)}
                      </td>

                      {/* Aksi */}
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600"
                            title="Edit"
                            onClick={() =>
                              navigate(
                                paths.dashboard.edit_siswa.replace(
                                  ":id",
                                  String(s.id_pengguna),
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
                                handleDeleteSiswa(s.id_pengguna),
                              )
                            }
                          >
                            <Trash className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td colSpan={11} className="px-6 py-12 text-center">
                    <div className="flex flex-col items-center justify-center gap-2">
                      <User className="h-10 w-10 text-slate-300" />
                      <p className="text-base font-medium text-slate-900">
                        Tidak ada siswa ditemukan
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
        </div>
      </div>
      <ConfirmAlert
        isOpen={modalKonfirmasiTerbuka}
        message={pesanKonfirmasiHapusSiswa}
        onClose={tutupModalKonfirmasi}
        onConfirm={() => void jalankanAksiKonfirmasi()}
        isLoading={sedangMemprosesKonfirmasi}
      />
    </div>
  );
};

export default AkunSiswaTables;
