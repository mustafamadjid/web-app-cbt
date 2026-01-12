import { useEffect, useMemo, useState, type ReactNode } from "react";
import { Link, useParams } from "react-router";
import {
  Calendar,
  Clock,
  GraduationCap,
  Hash,
  MapPin,
  ShieldCheck,
  User,
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
import type { DetailUjianItem } from "@/types/Ujian/DetailUjian";

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
  belum_dimulai: "bg-amber-100 text-amber-700 border-amber-200",
  berlangsung: "bg-emerald-100 text-emerald-700 border-emerald-200",
  selesai: "bg-slate-100 text-slate-600 border-slate-200",
};

type InfoItem = {
  label: string;
  value: string;
  icon?: ReactNode;
};

const InfoCard = ({ title, items }: { title: string; items: InfoItem[] }) => (
  <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
    <h2 className="text-sm font-semibold text-slate-800">{title}</h2>
    <div className="mt-4 grid gap-4 sm:grid-cols-2">
      {items.map((item) => (
        <div key={item.label} className="flex gap-3">
          {item.icon && (
            <div className="mt-0.5 text-[#397e50]">{item.icon}</div>
          )}
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
              {item.label}
            </p>
            <p className="text-sm font-medium text-slate-700">{item.value}</p>
          </div>
        </div>
      ))}
    </div>
  </div>
);

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
    if (!detail || detail.guru_pengawas_id === "") {
      return detail?.pengawas_ujian ?? "-";
    }
    const guru = guruOptions.find((item) => item.id === detail.guru_pengawas_id);
    return guru ? guru.nama : detail.pengawas_ujian ?? "-";
  }, [detail, guruOptions]);

  const sesiLabel = useMemo(() => {
    if (!detail || detail.sesi_id === "") {
      return detail?.sesi_ujian ? `Sesi ${detail.sesi_ujian}` : "-";
    }
    const sesi = sesiOptions.find((item) => item.id === detail.sesi_id);
    return sesi ? `${sesi.kode} - ${sesi.nama}` : `Sesi #${detail.sesi_id}`;
  }, [detail, sesiOptions]);

  const ruangLabel = useMemo(() => {
    if (!detail || detail.ruang_ujian_id === "") {
      return detail?.ruang_ujian ?? "-";
    }
    const ruang = ruangOptions.find((item) => item.id === detail.ruang_ujian_id);
    return ruang ? ruang.namaRuangan : detail.ruang_ujian ?? "-";
  }, [detail, ruangOptions]);

  const bankSoalLabel = useMemo(() => {
    if (!detail || detail.bank_soal_id === "") return "-";
    const bank = bankSoalOptions.find((item) => item.id === detail.bank_soal_id);
    return bank ? bank.nama : `Bank Soal #${detail.bank_soal_id}`;
  }, [detail, bankSoalOptions]);

  const detailStatusLabel = detail?.status_ujian
    ? statusLabelMap[detail.status_ujian] ?? detail.status_ujian
    : "Status belum tersedia";

  const infoUjian: InfoItem[] = detail
    ? [
        {
          label: "Nama Ujian",
          value: detail.nama_ujian,
          icon: <GraduationCap size={16} />,
        },
        {
          label: "Tipe Ujian",
          value: tipeUjianLabel[detail.tipe_ujian],
        },
        {
          label: "Deskripsi",
          value: detail.deskripsi_ujian || "-",
        },
        {
          label: "Status",
          value: detailStatusLabel,
        },
      ]
    : [];

  const infoKelasBank: InfoItem[] = detail
    ? [
        { label: "Tingkat Kelas", value: kelasLabel, icon: <Hash size={16} /> },
        { label: "Bank Soal", value: bankSoalLabel },
        {
          label: "Jumlah Soal",
          value: `${detail.jumlah_soal} soal`,
        },
      ]
    : [];

  const infoJadwal: InfoItem[] = detail
    ? [
        {
          label: "Tanggal",
          value: detail.tanggal_ujian || detail.tgl_ujian || "-",
          icon: <Calendar size={16} />,
        },
        {
          label: "Waktu",
          value: `${detail.waktu_mulai} - ${detail.waktu_selesai}`,
          icon: <Clock size={16} />,
        },
        {
          label: "Durasi",
          value: `${detail.durasi_menit} menit`,
        },
      ]
    : [];

  const infoRuang: InfoItem[] = detail
    ? [
        { label: "Ruang Ujian", value: ruangLabel, icon: <MapPin size={16} /> },
        { label: "Sesi", value: sesiLabel },
        { label: "Guru Pengawas", value: guruLabel, icon: <User size={16} /> },
      ]
    : [];

  const infoKeamanan: InfoItem[] = detail
    ? [
        {
          label: "Acak Soal",
          value: detail.acak_soal ? "Ya, acak soal" : "Tidak, urutan tetap",
          icon: <ShieldCheck size={16} />,
        },
        {
          label: "Token Ujian",
          value: detail.token_ujian || "-",
        },
      ]
    : [];

  return (
    <div className="mx-auto max-w-5xl px-4 py-10 sm:px-8">
      <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
            Detail Ujian
          </p>
          <h1 className="text-2xl font-bold text-slate-800">
            {detail?.nama_ujian ?? "Detail Ujian"}
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Informasi ujian berdasarkan data yang telah dibuat.
          </p>
        </div>
        <Link
          to={paths.dashboard.jadwal_ujian}
          className="rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
        >
          Kembali ke Jadwal
        </Link>
      </div>

      {loading && (
        <div className="flex min-h-[200px] items-center justify-center rounded-2xl border border-dashed border-slate-300 text-slate-500">
          <div className="flex flex-col items-center gap-2">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-[#397e50] border-t-transparent" />
            <span>Memuat detail ujian...</span>
          </div>
        </div>
      )}

      {!loading && errorMsg && (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {errorMsg}
        </div>
      )}

      {!loading && detail && (
        <div className="flex flex-col gap-5">
          <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex flex-wrap items-start justify-between gap-4">
              <div className="space-y-2">
                <span className="inline-flex items-center gap-2 rounded-md border border-slate-200 bg-slate-50 px-2 py-1 text-xs font-semibold text-slate-600">
                  Kelas {detail.tingkat_kelas ?? "-"} -{" "}
                  {detail.nama_kelas ?? "-"}
                </span>
                <p className="text-sm text-slate-500">{detail.deskripsi_ujian}</p>
              </div>
              <span
                className={`rounded-full border px-3 py-1 text-xs font-bold uppercase tracking-wider ${
                  detail.status_ujian
                    ? statusColorMap[detail.status_ujian] ??
                      "bg-slate-100 text-slate-600 border-slate-200"
                    : "bg-slate-100 text-slate-600 border-slate-200"
                }`}
              >
                {detailStatusLabel}
              </span>
            </div>
          </div>

          <InfoCard title="Informasi Ujian" items={infoUjian} />
          <InfoCard title="Tingkat Kelas & Bank Soal" items={infoKelasBank} />
          <InfoCard title="Jadwal Ujian" items={infoJadwal} />
          <InfoCard title="Ruang, Sesi, Pengawas" items={infoRuang} />
          <InfoCard title="Keamanan & Token" items={infoKeamanan} />
        </div>
      )}
    </div>
  );
};

export default DetailUjian;
