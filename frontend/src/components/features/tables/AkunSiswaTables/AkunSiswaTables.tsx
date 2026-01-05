import React, { useEffect, useRef, useState } from "react";
import {
  Search,
  Eye,
  EyeOff,
  ChevronDown,
  Trash2,
  Archive,
  Edit3,
  User,
  ShieldAlert,
  CheckCircle2,
  Trash,
  XCircle,
} from "lucide-react";

import { AddButton } from "@/components/common/Button/AddButton";
import type { StatusAkun, JenisKelamin } from "@/types/OpsiTypes/Option";
import { useNavigate } from "react-router";

import {
  getAngkatanOptions,
  getJenisKelaminOptions,
  getKelasOptions,
  getSiswa,
  type BarisSiswa,
} from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";

import type { KelasOption } from "@/types/DataMaster/MataPelajaran";

/** ===== Helpers ===== */
const getStatusBadge = (status: StatusAkun) => {
  switch (status) {
    case "aktif":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700 ring-1 ring-inset ring-emerald-600/20">
          <CheckCircle2 className="h-3 w-3" /> Aktif
        </span>
      );
    case "nonaktif":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-rose-50 px-2.5 py-0.5 text-xs font-medium text-rose-700 ring-1 ring-inset ring-rose-600/20">
          <XCircle className="h-3 w-3" /> Nonaktif
        </span>
      );
    case "dibekukan":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-slate-50 px-2.5 py-0.5 text-xs font-medium text-slate-600 ring-1 ring-inset ring-slate-500/20">
          <ShieldAlert className="h-3 w-3" /> Dibekukan
        </span>
      );
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
export const AkunSiswaTables: React.FC = () => {
  const navigate = useNavigate();

  const [dropdownAksiTerbuka, setDropdownAksiTerbuka] = useState(false);
  const [kataKunci, setKataKunci] = useState("");

  // FILTERS (backend-driven)
  const [angkatan, setAngkatan] = useState<string>("");
  const [kelasId, setKelasId] = useState<string>("");
  const [jenisKelamin, setJenisKelamin] = useState<string>("");

  const debouncedKataKunci = useDebouncedValue(kataKunci, 300);

  // OPTIONS from server
  const [opsiAngkatan, setOpsiAngkatan] = useState<number[]>([]);
  const [opsiKelas, setOpsiKelas] = useState<KelasOption[]>([]);
  const [opsiGender, setOpsiGender] = useState<
    Array<{ value: JenisKelamin; label: string }>
  >([]);

  // DATA from server
  const [daftarSiswa, setDaftarSiswa] = useState<BarisSiswa[]>([]);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  // Selection / privacy tetap
  const [idTerpilih, setIdTerpilih] = useState<Set<string>>(new Set());
  const [samarkanDataSensitif, setSamarkanDataSensitif] = useState(true);

  // Anti race condition
  const requestSeq = useRef(0);

  // ===== load filter options (sekali)
  useEffect(() => {
    let mounted = true;

    (async () => {
      try {
        setErrorMsg("");
        const [a, k, g] = await Promise.all([
          getAngkatanOptions(),
          getKelasOptions(),
          getJenisKelaminOptions(),
        ]);
        if (!mounted) return;

        setOpsiAngkatan(a);
        setOpsiKelas(k);
        setOpsiGender(g);
      } catch {
        if (!mounted) return;
        setErrorMsg("Gagal memuat opsi filter dari server.");
      }
    })();

    return () => {
      mounted = false;
    };
  }, []);

  // ===== fetch siswa saat filter berubah
  useEffect(() => {
    const seq = ++requestSeq.current;

    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");

        const data = await getSiswa({
          q: debouncedKataKunci.trim() || undefined,
          angkatan: angkatan ? Number(angkatan) : undefined,
          kelasId: kelasId || undefined,
          jenisKelamin: (jenisKelamin as JenisKelamin) || undefined,
        });

        if (seq !== requestSeq.current) return;

        setDaftarSiswa(data);

        // bersihkan selection yang tidak ada lagi di hasil
        setIdTerpilih((prev) => {
          if (prev.size === 0) return prev;
          const ids = new Set(data.map((x) => x.id));
          const next = new Set<string>();
          prev.forEach((id) => {
            if (ids.has(id)) next.add(id);
          });
          return next;
        });
      } catch (e) {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data siswa.");
        setDaftarSiswa([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, [debouncedKataKunci, angkatan, kelasId, jenisKelamin]);

  const siswaTersaring = daftarSiswa;

  const semuaTerlihatTerpilih =
    siswaTersaring.length > 0 &&
    siswaTersaring.every((s) => idTerpilih.has(s.id));

  const togglePilihSemuaTerlihat = () => {
    setIdTerpilih((prev) => {
      const next = new Set(prev);
      if (semuaTerlihatTerpilih)
        siswaTersaring.forEach((s) => next.delete(s.id));
      else siswaTersaring.forEach((s) => next.add(s.id));
      return next;
    });
  };

  const togglePilihBaris = (id: string) => {
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
    setAngkatan("");
    setKelasId("");
    setJenisKelamin("");
  };

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
            onClick={() =>
              navigate(`/dashboard/administrator/kelola-akun/tambah-siswa`)
            }
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
                {opsiAngkatan.map((a) => (
                  <option key={a} value={String(a)}>
                    {a}
                  </option>
                ))}
              </select>
            </div>

            {/* Kelas */}
            <div>
              <label className="text-xs font-medium text-slate-600">
                Kelas
              </label>
              <select
                value={kelasId}
                onChange={(e) => setKelasId(e.target.value)}
                className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 focus:border-[#397e50] focus:outline-none focus:ring-1 focus:ring-[#397e50]"
              >
                <option value="">Semua</option>
                {opsiKelas.map((k) => (
                  <option key={k.id} value={k.id}>
                    {k.label}
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
                {opsiGender.map((g) => (
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
                          onClick={() => setDropdownAksiTerbuka(false)}
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
                  Kelas & Absen
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
              {siswaTersaring.length > 0 ? (
                siswaTersaring.map((s) => {
                  const hpRaw = (s.nomorHp ?? "").trim();
                  const emailRaw = (s.email ?? "").trim();

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

                  const ttlTampil = `${s.tempatLahir}, ${formatTanggalIndo(
                    s.tanggalLahir
                  )}`;

                  return (
                    <tr
                      key={s.id}
                      className={`transition-colors hover:bg-slate-50 ${
                        idTerpilih.has(s.id) ? "bg-indigo-50/30" : ""
                      }`}
                    >
                      <td className="p-4">
                        <div className="flex items-center">
                          <input
                            type="checkbox"
                            checked={idTerpilih.has(s.id)}
                            onChange={() => togglePilihBaris(s.id)}
                            className="h-4 w-4 cursor-pointer rounded border-slate-300 text-[#397e50] focus:ring-[#397e50]"
                          />
                        </div>
                      </td>

                      {/* Siswa */}
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          <img
                            className="h-10 w-10 rounded-full object-cover ring-2 ring-white"
                            src={s.urlGambarProfil}
                            alt=""
                          />
                          <div className="flex flex-col">
                            <span className="font-semibold text-slate-900">
                              {s.namaLengkap}
                            </span>
                            <span className="text-xs text-slate-500">
                              @{s.username} • {labelGender[s.jenisKelamin]}
                            </span>
                          </div>
                        </div>
                      </td>

                      {/* Kelas & Absen */}
                      <td className="px-6 py-4">
                        <div className="flex flex-col gap-1">
                          <span className="text-slate-900">{s.kelas}</span>
                          <span className="text-xs text-slate-500">
                            No Absen:{" "}
                            <span className="font-medium text-slate-700">
                              {s.noAbsen}
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
                        {getStatusBadge(s.statusAkun)}
                      </td>

                      {/* Aksi */}
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600"
                            title="Edit"
                            onClick={() => console.log("Edit", s.id)}
                          >
                            <Edit3 className="h-4 w-4" />
                          </button>
                          <button
                            className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                            title="Detail"
                            onClick={() => console.log("Detail", s.id)}
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
                  <td colSpan={7} className="px-6 py-12 text-center">
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

        {/* Footer */}
        <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
          <div className="flex flex-1 items-center justify-between">
            <p className="text-sm text-slate-700">
              Menampilkan <span className="font-medium">1</span> sampai{" "}
              <span className="font-medium">{siswaTersaring.length}</span> dari{" "}
              <span className="font-medium">{siswaTersaring.length}</span> hasil
            </p>

            <p className="hidden sm:block text-xs text-slate-500">
              Geser tabel ke kanan/kiri untuk melihat kolom lainnya.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
