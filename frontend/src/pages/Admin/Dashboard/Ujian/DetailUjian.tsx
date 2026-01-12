import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router";
import {
  Calendar,
  GraduationCap,
  MapPin,
  ShieldCheck,
  BookOpen,
  Layers,
  ArrowLeft,
  BookAIcon
} from "lucide-react";

import { paths } from "@/routes/paths";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getRuangUjianOptions } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import {
  getUjianBankSoalOptions,
  getUjianGuruPengawasOptions,
  getUjianSesiOptions,
} from "@/services/Api/features-api/Ujian/ujian.service";
import { getJadwalUjianDetail } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { TingkatKelasOption } from "@/types/DataMaster/Kelas";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type {
  BankSoalOption,
  GuruPengawasOption,
  SesiUjianOption,
  TipeUjian,
} from "@/types/Ujian/BuatUjian";
import type { DetailUjianItem} from "@/types/Ujian/DetailUjian";



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
  const params = useParams();
  const ujianId = Number(params.id);

  const [detail, setDetail] = useState<DetailUjianItem | null>(null);
  const [loading, setLoading] = useState(true);
  const [errorMsg, setErrorMsg] = useState("");

  const [kelasOptions, setKelasOptions] = useState<TingkatKelasOption[]>([]);
  const [ruangOptions, setRuangOptions] = useState<RuangUjianRow[]>([]);
  const [guruOptions, setGuruOptions] = useState<GuruPengawasOption[]>([]);
  const [sesiOptions, setSesiOptions] = useState<SesiUjianOption[]>([]);
  const [bankSoalOptions, setBankSoalOptions] = useState<BankSoalOption[]>([]);

  // ... (Keep existing useEffect & logic for options and detail fetching)
  useEffect(() => {
    let active = true;
    const loadOptions = async () => {
      try {
        const [kelas, ruang, guru, sesi] = await Promise.all([
          getTingkatKelasOptions(),
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
      if (!detail || detail.kelas_id === "") {
        setBankSoalOptions([]);
        return;
      }
      const data = await getUjianBankSoalOptions({
        tingkatKelasId: detail.kelas_id,
      });
      if (active) setBankSoalOptions(data);
    };
    loadBankSoal();
    return () => {
      active = false;
    };
  }, [detail]);

  const kelasLabel = useMemo(() => {
    if (!detail || detail.kelas_id === "") return "-";
    const kelas = kelasOptions.find(
      (item) => item.id_tingkat_kelas === detail.kelas_id
    );
    return kelas ? `Kelas ${kelas.tingkat_kelas}` : `Kelas #${detail.kelas_id}`;
  }, [detail, kelasOptions]);

  const guruLabel = useMemo(() => {
    if (!detail || detail.guru_pengawas_id === "")
      return detail?.pengawas_ujian ?? "-";
    const guru = guruOptions.find(
      (item) => item.id === detail.guru_pengawas_id
    );
    return guru ? guru.nama : detail.pengawas_ujian ?? "-";
  }, [detail, guruOptions]);

  const sesiLabel = useMemo(() => {
    if (!detail || detail.sesi_id === "")
      return detail?.sesi_ujian ? `Sesi ${detail.sesi_ujian}` : "-";
    const sesi = sesiOptions.find((item) => item.id === detail.sesi_id);
    return sesi ? `${sesi.kode} - ${sesi.nama}` : `Sesi #${detail.sesi_id}`;
  }, [detail, sesiOptions]);

  const ruangLabel = useMemo(() => {
    if (!detail || detail.ruang_ujian_id === "")
      return detail?.ruang_ujian ?? "-";
    const ruang = ruangOptions.find(
      (item) => item.id === detail.ruang_ujian_id
    );
    return ruang ? ruang.namaRuangan : detail.ruang_ujian ?? "-";
  }, [detail, ruangOptions]);

  const bankSoalLabel = useMemo(() => {
    if (!detail || detail.bank_soal_id === "") return "-";
    const bank = bankSoalOptions.find(
      (item) => item.id === detail.bank_soal_id
    );
    return bank ? bank.nama : `Bank Soal #${detail.bank_soal_id}`;
  }, [detail, bankSoalOptions]);

  const detailStatusLabel = detail?.status_ujian
    ? statusLabelMap[detail.status_ujian] ?? detail.status_ujian
    : "Status belum tersedia";

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      {/* Header Section */}
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
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
            {detail?.nama_ujian ?? "Detail Ujian"}
          </h1>
          <p className="text-sm text-slate-500">ID Ujian: #{ujianId}</p>
        </div>

        {!loading && detail && (
          <div
            className={`inline-flex items-center gap-2 rounded-full border px-4 py-1.5 text-sm font-bold uppercase tracking-wider ${
              statusColorMap[detail.status_ujian] || statusColorMap.selesai
            }`}
          >
            <span
              className={`h-2 w-2 rounded-full ${
                detail.status_ujian === "berlangsung"
                  ? "animate-pulse bg-emerald-500"
                  : "bg-current"
              }`}
            />
            {detailStatusLabel}
          </div>
        )}
      </div>

      {loading ? (
        <div className="flex min-h-[400px] flex-col items-center justify-center rounded-3xl border-2 border-dashed border-slate-200 bg-white p-12 text-slate-400">
          <div className="h-10 w-10 animate-spin rounded-full border-4 border-[#397e50] border-t-transparent" />
          <p className="mt-4 font-medium">Menyiapkan detail informasi...</p>
        </div>
      ) : errorMsg ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center">
          <p className="font-semibold text-red-700">{errorMsg}</p>
        </div>
      ) : (
        detail && (
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
            {/* Main Info Card (Left Column) */}
            <div className="lg:col-span-2 space-y-6">
              <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
                <div className="border-b border-slate-100 bg-slate-50/50 px-6 py-4">
                  <h2 className="flex items-center gap-2 font-bold text-slate-800">
                    <BookOpen size={18} className="text-[#397e50]" />
                    Informasi Umum
                  </h2>
                </div>
                <div className="p-6">
                  <div className="grid gap-6 sm:grid-cols-2">
                    <div className="space-y-1">
                      <p className="text-2xs font-bold uppercase tracking-widest text-slate-400">
                        Nama Ujian
                      </p>
                      <p className="text-lg font-bold text-slate-800">
                        {detail.nama_ujian}
                      </p>
                    </div>
                    <div className="space-y-1">
                      <p className="text-2xs font-bold uppercase tracking-widest text-slate-400">
                        Tipe Ujian
                      </p>
                      <p className="text-lg font-medium text-slate-700">
                        {tipeUjianLabel[detail.tipe_ujian]}
                      </p>
                    </div>
                    <div className="sm:col-span-2 space-y-1">
                      <p className="text-2xs font-bold uppercase tracking-widest text-slate-400">
                        Deskripsi
                      </p>
                      <p className="text-sm leading-relaxed text-slate-600">
                        {detail.deskripsi_ujian || "Tidak ada deskripsi."}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Layout Grid for Sub-Cards */}
              <div className="grid gap-6 sm:grid-cols-2">
                {/* Jadwal Card */}
                <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                  <div className="mb-4 flex items-center gap-2 border-b border-slate-50 pb-3 font-bold text-slate-800">
                    <Calendar size={18} className="text-[#397e50]" />
                    Jadwal Pelaksanaan
                  </div>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-slate-500">Tanggal</span>
                      <span className="text-sm font-bold text-slate-700">
                        {detail.tanggal_ujian || detail.tgl_ujian || "-"}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-slate-500">Waktu</span>
                      <span className="text-sm font-bold text-slate-700">
                        {detail.waktu_mulai} - {detail.waktu_selesai}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-slate-500">Durasi</span>
                      <span className="rounded bg-slate-100 px-2 py-0.5 text-xs font-bold text-slate-600">
                        {detail.durasi_menit} Menit
                      </span>
                    </div>
                  </div>
                </div>

                {/* Ruangan Card */}
                <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                  <div className="mb-4 flex items-center gap-2 border-b border-slate-50 pb-3 font-bold text-slate-800">
                    <MapPin size={18} className="text-[#397e50]" />
                    Lokasi & Pengawas
                  </div>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-slate-500">Ruangan</span>
                      <span className="text-sm font-bold text-[#397e50]">
                        {ruangLabel}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-slate-500">Sesi</span>
                      <span className="text-sm font-bold text-slate-700">
                        {sesiLabel}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-slate-500">Pengawas</span>
                      <span className="text-sm font-bold text-slate-700">
                        {guruLabel}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Sidebar Info (Right Column) */}
            <div className="space-y-6">
              {/* Kelas & Bank Soal */}
              <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                <h3 className="mb-4 text-xs font-bold uppercase tracking-widest text-slate-400">
                  Kurikulum & Soal
                </h3>
                <div className="space-y-5">
                  <div className="flex items-start gap-3">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-[#397e50]">
                      <GraduationCap size={18} />
                    </div>
                    <div>
                      <p className="text-2xs font-bold uppercase text-slate-400">
                        Kelas
                      </p>
                      <p className="text-sm font-bold text-slate-700">
                        {kelasLabel} - {detail.nama_kelas}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-green-800">
                      <BookAIcon size={18} />
                    </div>
                    <div>
                      <p className="text-2xs font-bold uppercase text-slate-400">
                        Bank Soal
                      </p>
                      <p className="text-sm font-bold text-slate-700">
                        {bankSoalLabel}
                      </p>
                      <p className="text-xs text-slate-500">
                        {detail.jumlah_soal} Butir Soal
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Keamanan */}
              <div className="rounded-2xl border border-[#397e50]/20 bg-emerald-50/30 p-6 shadow-sm">
                <h3 className="mb-4 text-xs font-bold uppercase tracking-widest text-emerald-700">
                  Konfigurasi & Keamanan
                </h3>
                <div className="space-y-4">
                  <div className="flex items-center gap-3">
                    <ShieldCheck size={18} className="text-[#397e50]" />
                    <span className="text-sm font-medium text-slate-700">
                      {detail.acak_soal ? "Soal Diacak" : "Urutan Tetap"}
                    </span>
                  </div>
                  <div className="rounded-xl border border-[#397e50]/20 bg-white p-3 text-center">
                    <p className="text-2xs font-bold uppercase text-slate-400">
                      Token Ujian
                    </p>
                    <p className="text-xl font-mono font-black tracking-widest text-[#397e50]">
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
