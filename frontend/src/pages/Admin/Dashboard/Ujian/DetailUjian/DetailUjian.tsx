import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router";
import {
  Calendar,
  GraduationCap,
  MapPin,
  ShieldCheck,
  BookOpen,
  ArrowLeft,
  BookAIcon,
  Hash,
  Layout,
  AlertCircle,
} from "lucide-react";

import PrintButton from "@/components/common/Input/PrintButton";
import { paths } from "@/routes/paths";
import {
  getRuangUjianOptions,
  getTingkatKelass,
  getUjianBankSoalOptions,
  getUjianGuruPengawasOptions,
  getUjianSesiOptions,
} from "@/services/Api/features-api/GetOptions/options.service";
import { getJadwalUjianDetail } from "@/services/Api/features-api/Ujian/jadwalujian.service";

import type { TingkatKelas } from "@/types/DataMaster/Kelas";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type {
  BankSoalOption,
  GuruPengawasOption,
  SesiUjianOption,
  TipeUjian,
} from "@/types/Ujian/BuatUjian";
import type { DetailUjianItem, PrintJenis } from "@/types/Ujian/DetailUjian";

const tipeUjianLabel: Record<TipeUjian, string> = {
  PILIHAN_GANDA: "Pilihan Ganda",
  ESSAY: "Essay",
  CAMPURAN: "Pilihan Ganda & Essay",
};

const statusLabelMap: Record<string, string> = {
  belum_dimulai: "Belum Dimulai",
  berlangsung: "Berlangsung",
  selesai: "Selesai",
};

const statusColorMap: Record<string, string> = {
  belum_dimulai: "bg-amber-50 text-amber-700 border-amber-200",
  berlangsung: "bg-emerald-50 text-emerald-700 border-emerald-200",
  selesai: "bg-slate-50 text-slate-600 border-slate-200",
};

const DetailUjian = () => {
  // TODO : Saat sudah integrasi dengan API, maka cukup panggil service detail ujian dengan param id
  // TODO : Hanya cukup detail const [detail, setDetail] = useState<DetailUjianItem | null>(null);

  const params = useParams();
  const ujianId = Number(params.id);

  const [detail, setDetail] = useState<DetailUjianItem | null>(null);
  const [loading, setLoading] = useState(true);
  const [errorMsg, setErrorMsg] = useState("");

  const [kelasOptions, setKelasOptions] = useState<TingkatKelas[]>([]);
  const [ruangOptions, setRuangOptions] = useState<RuangUjianRow[]>([]);
  const [guruOptions, setGuruOptions] = useState<GuruPengawasOption[]>([]);
  const [sesiOptions, setSesiOptions] = useState<SesiUjianOption[]>([]);
  const [bankSoalOptions, setBankSoalOptions] = useState<BankSoalOption[]>([]);

  useEffect(() => {
    let active = true;
    const loadOptions = async () => {
      try {
        const [kelas, ruang, guru, sesi] = await Promise.all([
          getTingkatKelass(),
          getRuangUjianOptions(),
          getUjianGuruPengawasOptions(),
          getUjianSesiOptions(),
        ]);
        if (!active) return;
        setKelasOptions(kelas);
        setRuangOptions(ruang);
        setGuruOptions(guru);
        setSesiOptions(sesi);
      } catch {
        if (!active) return;
        setKelasOptions([]);
        setRuangOptions([]);
        setGuruOptions([]);
        setSesiOptions([]);
      }
    };
    loadOptions();
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    let active = true;
    const loadDetail = async () => {
      if (!Number.isFinite(ujianId)) {
        setErrorMsg("ID ujian tidak valid.");
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getJadwalUjianDetail(ujianId);
        if (!active) return;
        setDetail(data);
      } catch {
        if (!active) return;
        setErrorMsg("Detail ujian tidak ditemukan.");
        setDetail(null);
      } finally {
        if (active) setLoading(false);
      }
    };
    loadDetail();
    return () => {
      active = false;
    };
  }, [ujianId]);

  useEffect(() => {
    let active = true;
    const loadBankSoal = async () => {
      if (!detail || !detail.kelas_id) {
        setBankSoalOptions([]);
        return;
      }
      try {
        const data = await getUjianBankSoalOptions({
          tingkatKelasId: detail.kelas_id,
        });
        if (active) setBankSoalOptions(data);
      } catch {
        if (active) setBankSoalOptions([]);
      }
    };
    loadBankSoal();
    return () => {
      active = false;
    };
  }, [detail]);

  const kelasLabel = useMemo(() => {
    if (!detail?.kelas_id) return "-";
    const kelas = kelasOptions.find(
      (item) => item.id_tingkat_kelas === detail.kelas_id,
    );
    return kelas ? `Kelas ${kelas.tingkat_kelas}` : `Kelas #${detail.kelas_id}`;
  }, [detail, kelasOptions]);

  const guruLabel = useMemo(() => {
    if (!detail?.guru_pengawas_id) return detail?.pengawas_ujian ?? "-";
    const guru = guruOptions.find(
      (item) => item.id === detail.guru_pengawas_id,
    );
    return guru ? guru.nama : (detail.pengawas_ujian ?? "-");
  }, [detail, guruOptions]);

  const sesiLabel = useMemo(() => {
    if (!detail?.sesi_id)
      return detail?.sesi_ujian ? `Sesi ${detail.sesi_ujian}` : "-";
    const sesi = sesiOptions.find((item) => item.id === detail.sesi_id);
    return sesi ? `${sesi.kode} - ${sesi.nama}` : `Sesi #${detail.sesi_id}`;
  }, [detail, sesiOptions]);

  const ruangLabel = useMemo(() => {
    if (!detail?.ruang_ujian_id) return detail?.ruang_ujian ?? "-";
    const ruang = ruangOptions.find(
      (item) => item.id === detail.ruang_ujian_id,
    );
    return ruang ? ruang.namaRuangan : (detail.ruang_ujian ?? "-");
  }, [detail, ruangOptions]);

  const bankSoalLabel = useMemo(() => {
    if (!detail?.bank_soal_id) return "-";
    const bank = bankSoalOptions.find(
      (item) => item.id === detail.bank_soal_id,
    );
    return bank ? bank.nama : `Bank Soal #${detail.bank_soal_id}`;
  }, [detail, bankSoalOptions]);

  const handlePrint = (jenis: PrintJenis) => {
    if (typeof window !== "undefined") window.print();
  };

  const statusLabel = detail?.status_ujian
    ? statusLabelMap[detail.status_ujian]
    : "Tidak Diketahui";

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      {/* 1. Header Section */}
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b border-slate-100 pb-6">
        <div>
          <Link
            to={paths.dashboard.jadwal_ujian}
            className="group mb-2 inline-flex items-center gap-2 text-sm font-medium text-slate-500 transition hover:text-[#397e50]"
          >
            <ArrowLeft
              size={16}
              className="transition-transform group-hover:-translate-x-1"
            />
            Kembali ke Jadwal
          </Link>
          <h1 className="text-3xl font-extrabold tracking-tight text-slate-900">
            {detail?.nama_ujian ?? "Memuat Detail..."}
          </h1>
          <div className="mt-1 flex items-center gap-3 text-sm text-slate-500">
            <span className="flex items-center gap-1.5 font-medium">
              <Hash size={14} className="text-slate-400" /> ID: {ujianId}
            </span>
            <span className="text-slate-300">|</span>
            <span className="font-semibold text-[#397e50]">
              {detail ? tipeUjianLabel[detail.tipe_ujian] : "-"}
            </span>
          </div>
        </div>

        {!loading && detail && (
          <div
            className={`inline-flex items-center gap-2 rounded-lg border px-4 py-2 text-sm font-bold uppercase tracking-wider ${
              statusColorMap[detail.status_ujian] || statusColorMap.selesai
            }`}
          >
            <span
              className={`h-2.5 w-2.5 rounded-full ${
                detail.status_ujian === "berlangsung"
                  ? "animate-pulse bg-emerald-500"
                  : "bg-current"
              }`}
            />
            {statusLabel}
          </div>
        )}
      </div>

      {loading ? (
        <div className="flex min-h-[400px] flex-col items-center justify-center rounded-3xl border border-slate-200 bg-white p-12 text-slate-400">
          <div className="h-10 w-10 animate-spin rounded-full border-4 border-[#397e50] border-t-transparent" />
          <p className="mt-4 font-medium">Sinkronisasi data ujian...</p>
        </div>
      ) : errorMsg ? (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-red-100 bg-red-50 p-10 text-center">
          <AlertCircle size={40} className="mb-3 text-red-500" />
          <p className="text-lg font-semibold text-red-700">{errorMsg}</p>
          <Link
            to={paths.dashboard.jadwal_ujian}
            className="mt-4 text-sm font-bold text-red-600 underline"
          >
            Kembali
          </Link>
        </div>
      ) : (
        detail && (
          <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
            {/* Main Column (Left) */}
            <div className="lg:col-span-2 space-y-6">
              {/* Deskripsi */}
              <div className="rounded-2xl border border-slate-200 bg-white shadow-sm overflow-hidden">
                <div className="bg-slate-50/50 px-6 py-4 border-b border-slate-100 flex items-center gap-2 font-bold text-slate-800 text-sm">
                  <BookOpen size={16} className="text-[#397e50]" />
                  Deskripsi Ujian
                </div>
                <div className="p-6">
                  <p className="text-sm leading-relaxed text-slate-600">
                    {detail.deskripsi_ujian ||
                      "Tidak ada deskripsi tambahan untuk ujian ini."}
                  </p>
                </div>
              </div>

              {/* Grid Sub-Info */}
              <div className="grid gap-6 sm:grid-cols-2">
                <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                  <div className="mb-4 flex items-center gap-2 border-b border-slate-50 pb-3 font-bold text-slate-800 text-sm uppercase tracking-wide">
                    <Calendar size={16} className="text-[#397e50]" />
                    Waktu
                  </div>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Tanggal</span>
                      <span className="font-bold text-slate-700">
                        {detail.tanggal_ujian || detail.tgl_ujian}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Jam</span>
                      <span className="font-bold text-slate-700">
                        {detail.waktu_mulai} - {detail.waktu_selesai}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Durasi</span>
                      <span className="font-bold text-[#397e50]">
                        {detail.durasi_menit} Menit
                      </span>
                    </div>
                  </div>
                </div>

                <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                  <div className="mb-4 flex items-center gap-2 border-b border-slate-50 pb-3 font-bold text-slate-800 text-sm uppercase tracking-wide">
                    <MapPin size={16} className="text-[#397e50]" />
                    Lokasi
                  </div>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Ruangan</span>
                      <span className="font-bold text-[#397e50]">
                        {ruangLabel}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Sesi</span>
                      <span className="font-bold text-slate-700">
                        {sesiLabel}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Pengawas</span>
                      <span className="font-bold text-slate-700">
                        {guruLabel}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {/* 4. Action Section (PRINT BUTTONS REVISED POSITION) */}
              <div className="rounded-2xl p-5">
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                  <PrintButton
                    label="Daftar Hadir"
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-lg bg-[#397e50] text-xs font-bold text-white transition hover:bg-[#2d633f]"
                    onClick={() => handlePrint("daftar-hadir")}
                  />
                  <PrintButton
                    label="Berita Acara"
                    variant="outline"
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-lg border border-slate-300 bg-white text-xs font-bold text-slate-600 transition hover:bg-slate-50"
                    onClick={() => handlePrint("berita-acara")}
                  />
                  <PrintButton
                    label="Kartu Peserta"
                    variant="outline"
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-lg border border-slate-300 bg-white text-xs font-bold text-slate-600 transition hover:bg-slate-50"
                    onClick={() => handlePrint("kartu-peserta")}
                  />
                </div>
              </div>
            </div>

            {/* Sidebar (Right) */}
            <div className="space-y-6">
              {/* Target Akademik */}
              <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                <h3 className="mb-5 text-2xs font-bold uppercase tracking-[0.15em] text-slate-400">
                  Target Akademik
                </h3>
                <div className="space-y-6">
                  <div className="flex items-start gap-3">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-[#397e50]">
                      <GraduationCap size={18} />
                    </div>
                    <div>
                      <p className="text-2xs font-bold uppercase text-slate-400">
                        Kelas
                      </p>
                      <p className="text-xs font-bold text-slate-800 leading-tight">
                        {kelasLabel}
                      </p>
                      <p className="text-2xs text-slate-500">
                        {detail.nama_kelas}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                      <BookAIcon size={18} />
                    </div>
                    <div>
                      <p className="text-2xs font-bold uppercase text-slate-400">
                        Bank Soal
                      </p>
                      <p className="text-xs font-bold text-slate-800 leading-tight">
                        {bankSoalLabel}
                      </p>
                      <div className="mt-1 inline-flex items-center gap-1 rounded bg-blue-100/50 px-1.5 py-0.5 text-[9px] font-bold text-blue-700">
                        <Layout size={10} /> {detail.jumlah_soal} SOAL
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Keamanan & Token */}
              <div className="rounded-2xl border border-emerald-100 bg-emerald-50/20 p-6">
                <h3 className="mb-5 text-2xs font-bold uppercase tracking-[0.15em] text-emerald-800/60">
                  Sistem & Keamanan
                </h3>
                <div className="space-y-5">
                  <div className="flex items-center gap-2 text-xs font-semibold text-slate-600">
                    <ShieldCheck size={16} className="text-[#397e50]" />
                    {detail.acak_soal ? "Soal Diacak" : "Urutan Statis"}
                  </div>
                  <div className="rounded-xl border border-white bg-white p-5 text-center shadow-sm">
                    <p className="mb-2 text-[9px] font-bold uppercase tracking-widest text-slate-400">
                      Token Akses
                    </p>
                    <p className="text-2xl font-mono font-black tracking-[0.2em] text-[#397e50]">
                      {detail.token_ujian || "------"}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )
      )}
    </div>
  );
};

export default DetailUjian;
