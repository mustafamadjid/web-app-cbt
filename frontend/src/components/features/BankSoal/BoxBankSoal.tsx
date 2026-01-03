import React from "react";

import type { BankSoalItem } from "@/types/DataMaster/BankSoal";

type BoxBankSoalProps = Omit<BankSoalItem, "id"> & {
  onPreview?: () => void;
  onKelola?: () => void;
  onHapus?: () => void;
  className?: string;
};

const formatTanggal = (iso?: string) => {
  if (!iso) return "-";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
};

export const BoxBankSoal: React.FC<BoxBankSoalProps> = ({
  nama_banksoal = "Bank Soal",
  mata_pelajaran = "-",
  materi = "-",
  kelas,
  deskripsi = "-",
  tgl_buat,
  jumlah_soal_pg = 0,
  jumlah_soal_essay = 0,
  onKelola,
  onPreview,
  onHapus,
  className = "",
}) => {
  return (
    <div
      className={[
        "group relative w-full overflow-hidden rounded-2xl border",
        "border-emerald-200/70 bg-white shadow-sm transition hover:shadow-md",
        "p-4 sm:p-5",
        className,
      ].join(" ")}
    >
      {/* accent bar */}
      <div className="pointer-events-none absolute inset-x-0 top-0 h-1.5 bg-[#397e50]" />

      {/* subtle tint */}
      <div className="pointer-events-none absolute -right-24 -top-24 h-56 w-56 rounded-full bg-emerald-500/10 blur-3xl" />

      <div className="flex gap-4">
        {/* Left content */}
        <div className="min-w-0 flex-1">
          {/* Title + date */}
          <div className="flex items-start justify-between gap-3">
            <h3 className="truncate text-base font-semibold text-emerald-950 sm:text-lg">
              {nama_banksoal}
            </h3>

            <span className="shrink-0 rounded-lg border border-emerald-200/70 bg-emerald-50 px-3 py-1.5 text-sm text-emerald-900">
              <span className="font-semibold">Tanggal Upload:</span>{" "}
              {formatTanggal(tgl_buat)}
            </span>
          </div>

          {/* Meta */}
          <div className="mt-3 space-y-2 text-sm text-emerald-950/80">
            <div className="flex flex-wrap gap-x-6 gap-y-1">
              <div className="flex items-center gap-2">
                <span className="text-emerald-900/60">Kelas</span>
                <span className="font-semibold text-emerald-950">
                  {kelas ?? "-"}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-emerald-900/60">Mata Pelajaran</span>
                <span className="font-semibold text-emerald-950">
                  {mata_pelajaran}
                </span>
              </div>
            </div>

            <div className="flex flex-wrap gap-x-6 gap-y-1">
              <div className="flex items-center gap-2">
                <span className="text-emerald-900/60">Materi</span>
                <span className="font-semibold text-emerald-950">{materi}</span>
              </div>
            </div>

            {/* Jumlah Soal */}
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <span className="text-emerald-900/60">Jumlah Soal</span>

              <span className="inline-flex items-center gap-2 rounded-full border border-emerald-200/70 bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-900">
                <span className="h-1.5 w-1.5 rounded-full bg-emerald-600" />
                PG: <span className="text-emerald-950">{jumlah_soal_pg}</span>
              </span>

              <span className="inline-flex items-center gap-2 rounded-full border border-emerald-200/70 bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-900">
                <span className="h-1.5 w-1.5 rounded-full bg-emerald-600" />
                Essay:{" "}
                <span className="text-emerald-950">{jumlah_soal_essay}</span>
              </span>
            </div>
          </div>

          {/* Description */}
          <div className="mt-4">
            <div className="text-sm font-semibold text-emerald-950">
              Deskripsi
            </div>
            <p className="mt-1 line-clamp-3 text-sm leading-relaxed text-emerald-950/70">
              {deskripsi}
            </p>
          </div>

          {/* Actions */}
          <div className="mt-5 flex flex-wrap gap-2">
            {/* Preview (primary) */}
            <button
              type="button"
              onClick={onPreview}
              disabled={!onPreview}
              className={[
                "inline-flex cursor-pointer items-center gap-2 rounded-xl px-3 py-2 text-sm font-semibold",
                "bg-emerald-800 text-white shadow-sm transition",
                "hover:bg-emerald-900 active:scale-[0.99]",
                "disabled:cursor-not-allowed disabled:opacity-60",
              ].join(" ")}
            >
              {/* eye icon */}
              <svg
                className="h-5 w-5"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
              >
                <path
                  d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinejoin="round"
                />
                <path
                  d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z"
                  stroke="currentColor"
                  strokeWidth="2"
                />
              </svg>
              Preview
            </button>

            {/* Kelola */}
            <button
              type="button"
              onClick={onKelola}
              disabled={!onKelola}
              className={[
                "inline-flex cursor-pointer items-center gap-2 rounded-xl border px-3 py-2 text-sm font-semibold shadow-sm transition",
                "border-emerald-200/70 bg-white text-emerald-900",
                "hover:bg-emerald-50 active:scale-[0.99]",
                "disabled:cursor-not-allowed disabled:opacity-60",
              ].join(" ")}
            >
              {/* gear icon (react-safe attrs) */}
              <svg
                className="h-5 w-5 text-emerald-700"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M21 13v-2a1 1 0 0 0-1-1h-.757l-.707-1.707.535-.536a1 1 0 0 0 0-1.414l-1.414-1.414a1 1 0 0 0-1.414 0l-.536.535L14 4.757V4a1 1 0 0 0-1-1h-2a1 1 0 0 0-1 1v.757l-1.707.707-.536-.535a1 1 0 0 0-1.414 0L4.929 6.343a1 1 0 0 0 0 1.414l.536.536L4.757 10H4a1 1 0 0 0-1 1v2a1 1 0 0 0 1 1h.757l.707 1.707-.535.536a1 1 0 0 0 0 1.414l1.414 1.414a1 1 0 0 0 1.414 0l.536-.535 1.707.707V20a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1v-.757l1.707-.708.536.536a1 1 0 0 0 1.414 0l1.414-1.414a1 1 0 0 0 0-1.414l-.535-.536.707-1.707H20a1 1 0 0 0 1-1Z"
                />
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z"
                />
              </svg>
              Kelola
            </button>

            {/* Hapus */}
            <button
              type="button"
              onClick={onHapus}
              disabled={!onHapus}
              className={[
                "inline-flex cursor-pointer items-center gap-2 rounded-xl border px-3 py-2 text-sm font-semibold shadow-sm transition",
                "border-rose-200 bg-rose-50 text-rose-700",
                "hover:bg-rose-100 active:scale-[0.99]",
                "disabled:cursor-not-allowed disabled:opacity-60",
              ].join(" ")}
            >
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                className="text-rose-600"
                fill="none"
                aria-hidden="true"
              >
                <path
                  d="M4 7h16M10 11v7M14 11v7M6 7l1 14h10l1-14M9 7V5h6v2"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
              Hapus
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
