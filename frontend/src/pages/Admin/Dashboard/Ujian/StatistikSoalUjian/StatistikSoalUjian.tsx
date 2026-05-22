import { useMemo, useState } from "react";
import { Link, useParams } from "react-router";
import {
  ArrowLeft,
  CheckCircle2,
  FileBarChart,
  RefreshCw,
  SortDesc,
  XCircle,
} from "lucide-react";

import RichContentRenderer from "@/components/common/RichContentRenderer";
import { useAuth } from "@/contexts/AuthContext";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { formatSoalTypeLabel } from "@/helper/Ujian/soalType";
import { paths } from "@/routes/paths";
import { useGetAnalisisSoal } from "@/services/Api/features-api/Ujian/analisisSoal.service";
import type { AnalisisSoalItem } from "@/types/Ujian/AnalisisSoal";

type SortMode = "salah" | "benar";

const SORT_OPTIONS: { value: SortMode; label: string }[] = [
  { value: "salah", label: "Paling banyak salah" },
  { value: "benar", label: "Paling banyak benar" },
];

const getStableQuestionOrder = (item: AnalisisSoalItem) =>
  item.no_urut_soal || item.id_soal;

const sortAnalisisSoal = (
  items: AnalisisSoalItem[],
  sortMode: SortMode,
) =>
  [...items].sort((left, right) => {
    const leftValue =
      sortMode === "salah"
        ? left.jumlah_jawaban_salah
        : left.jumlah_jawaban_benar;
    const rightValue =
      sortMode === "salah"
        ? right.jumlah_jawaban_salah
        : right.jumlah_jawaban_benar;

    if (rightValue !== leftValue) {
      return rightValue - leftValue;
    }

    return getStableQuestionOrder(left) - getStableQuestionOrder(right);
  });

const StatistikSoalUjian = () => {
  const { user } = useAuth();
  const { id } = useParams();
  const [sortMode, setSortMode] = useState<SortMode>("salah");

  const jadwalUjianId = useMemo(() => Number(id), [id]);
  const isJadwalUjianIdValid =
    Number.isInteger(jadwalUjianId) && jadwalUjianId > 0;

  const {
    data,
    loading,
    error,
    refetch,
  } = useGetAnalisisSoal(jadwalUjianId, isJadwalUjianIdValid);

  const analisisSoal = useMemo(
    () => sortAnalisisSoal(data?.analisis_soal ?? [], sortMode),
    [data?.analisis_soal, sortMode],
  );

  const backPath =
    user?.role === "GURU"
      ? paths.dashboard.hasil_ujian_detail_guru.replace(
          ":id",
          String(jadwalUjianId),
        )
      : paths.dashboard.hasil_ujian_detail.replace(
          ":id",
          String(jadwalUjianId),
        );

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4">
        <Link
          to={backPath}
          className="group inline-flex w-fit cursor-pointer items-center gap-2 text-sm font-medium text-slate-500 transition-colors hover:text-[#397e50]"
        >
          <ArrowLeft
            size={16}
            className="transition-transform group-hover:-translate-x-1"
          />
          Kembali ke Detail Hasil
        </Link>

        <div className="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:flex-row lg:items-center lg:justify-between">
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-[#397e50]/10 text-[#397e50]">
              <FileBarChart size={20} />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-slate-800">
                Statistik Soal
              </h1>
              <p className="mt-1 text-sm text-slate-500">
                Ringkasan jumlah jawaban benar dan salah untuk setiap soal.
              </p>
            </div>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
              <SortDesc size={16} />
              Urutkan
            </label>
            <select
              value={sortMode}
              onChange={(event) => setSortMode(event.target.value as SortMode)}
              className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm font-semibold text-slate-700 transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
            >
              {SORT_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        </div>
      </header>

      {!isJadwalUjianIdValid && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm font-medium text-rose-700">
          Data jadwal ujian tidak valid.
        </div>
      )}

      {error && (
        <div className="rounded-2xl border border-rose-200 bg-rose-50 p-5 text-sm font-semibold text-rose-700">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <span>Gagal memuat statistik soal: {error}</span>
            <button
              type="button"
              onClick={() => void refetch()}
              className="inline-flex w-fit cursor-pointer items-center gap-2 rounded-xl border border-rose-200 bg-white px-4 py-2 text-xs font-semibold text-rose-700 transition hover:border-rose-300 hover:bg-rose-100"
            >
              <RefreshCw size={14} />
              Coba Lagi
            </button>
          </div>
        </div>
      )}

      {loading && (
        <div className="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500">
          <div className="flex items-center gap-3">
            <div className="h-7 w-7 animate-spin rounded-full border-2 border-slate-200 border-t-[#397e50]" />
            Memuat statistik soal...
          </div>
        </div>
      )}

      {!loading && !error && isJadwalUjianIdValid && analisisSoal.length === 0 && (
        <div className="rounded-2xl border border-slate-200 bg-white p-8 text-center text-sm text-slate-500">
          Statistik soal belum tersedia untuk jadwal ujian ini.
        </div>
      )}

      {!loading && !error && analisisSoal.length > 0 && (
        <section className="grid gap-4">
          {analisisSoal.map((item, index) => {
            const gambarUrl = resolveImageUrl(item.gambar);
            const nomorSoal = item.no_urut_soal || index + 1;

            return (
              <article
                key={item.id_soal}
                className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm transition hover:border-[#397e50]/40"
              >
                <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="rounded-full bg-[#397e50]/10 px-3 py-1 text-xs font-bold text-[#397e50]">
                        Soal {nomorSoal}
                      </span>
                      <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-600">
                        {formatSoalTypeLabel(item.tipe_soal)}
                      </span>
                    </div>

                    <RichContentRenderer
                      content={item.pertanyaan_content}
                      fallbackText={item.pertanyaan}
                      className="mt-4"
                      paragraphClassName="whitespace-pre-wrap text-sm leading-relaxed text-slate-700"
                    />

                    {gambarUrl && (
                      <div className="mt-4 overflow-hidden rounded-xl border border-slate-200 bg-slate-50">
                        <img
                          src={gambarUrl}
                          alt={`Ilustrasi soal ${nomorSoal}`}
                          className="max-h-[360px] w-full object-contain"
                          loading="lazy"
                        />
                      </div>
                    )}
                  </div>

                  <div className="grid min-w-full grid-cols-2 gap-3 sm:min-w-72 lg:w-72">
                    <div className="rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3">
                      <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-emerald-700">
                        <CheckCircle2 size={15} />
                        Dijawab Benar
                      </div>
                      <p className="mt-2 text-2xl font-bold text-emerald-700">
                        {item.jumlah_jawaban_benar}
                      </p>
                    </div>

                    <div className="rounded-xl border border-rose-200 bg-rose-50 px-4 py-3">
                      <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-rose-700">
                        <XCircle size={15} />
                        Dijawab Salah
                      </div>
                      <p className="mt-2 text-2xl font-bold text-rose-700">
                        {item.jumlah_jawaban_salah}
                      </p>
                    </div>
                  </div>
                </div>
              </article>
            );
          })}
        </section>
      )}
    </div>
  );
};

export default StatistikSoalUjian;
