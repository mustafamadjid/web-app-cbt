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

const kelasLabelClass = (k?: number | string) => {
  const kk = String(k ?? "");
  const base = "rounded-md px-2 py-1 text-sm font-bold text-white shadow-sm";

  const map: Record<string, string> = {
    "10": "bg-linear-to-r from-emerald-600 to-emerald-800",
    "11": "bg-linear-to-r from-green-700 to-green-900",
    "12": "bg-linear-to-r from-teal-600 to-teal-800",
  };

  return [base, map[kk] ?? "bg-gray-400"].join(" ");
};


export const BoxBankSoal: React.FC<BoxBankSoalProps> = ({
  nama_banksoal = "Bank Soal",
  guru = "-",
  mata_pelajaran = "-",
  materi = "-",
  kelas,
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
        "group flex w-full flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:-translate-y-1 hover:border-[#397e50]/50 hover:shadow-lg hover:shadow-[#397e50]/10",
        className,
      ].join(" ")}
    >
      {/* Top Accent Line */}
      <div className="h-1.5 w-full bg-linear-to-r from-[#397e50] to-[#37513d]" />

      <div className="flex flex-1 flex-col p-5">
        {/* Header: Date & Class */}
        <div className="mb-3 flex items-center justify-between text-xs font-medium text-gray-500">
          <span className="flex items-center gap-1.5">
            <svg
              className="h-4 w-4 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
            {formatTanggal(tgl_buat)}
          </span>

          <span className={kelasLabelClass(kelas)}>Kelas {kelas ?? "-"}</span>
        </div>

        {/* Title (Lebih Besar) */}
        <h3 className="mb-2 line-clamp-2 text-lg font-bold leading-snug text-[#37513d] transition-colors group-hover:text-[#397e50]">
          {nama_banksoal}
        </h3>

        {/* Guru pembuat soal */}
        <h4 className="mb-7 truncate  text-gray-700 text-sm flex items-center gap-2">
          <span className="font-semibold">Dibuat Oleh :</span> {guru}
        </h4>

        {/* Info Content (Text SM agar terbaca jelas) */}
        <div className="mb-4 space-y-2 text-sm">
          <div className="flex gap-3">
            <span className="min-w-20 font-medium text-gray-500">Mapel:</span>
            <span className="truncate font-semibold text-gray-700">
              {mata_pelajaran}
            </span>
          </div>

          <div className="flex gap-3">
            <span className="min-w-20 font-medium text-gray-500">Materi:</span>
            <span className="truncate font-semibold text-gray-700">
              {materi}
            </span>
          </div>
        </div>

        {/* Stats Pills */}
        <div className="mb-5 flex items-center gap-2">
          <span className="inline-flex items-center gap-1.5 rounded-full border border-emerald-100 bg-emerald-50 px-2.5 py-1 text-xs font-bold text-emerald-800">
            <span className="h-2 w-2 rounded-full bg-emerald-600" />
            PG: {jumlah_soal_pg}
          </span>

          <span className="inline-flex items-center gap-1.5 rounded-full border border-amber-100 bg-amber-50 px-2.5 py-1 text-xs font-bold text-amber-800">
            <span className="h-2 w-2 rounded-full bg-amber-500" />
            Essay: {jumlah_soal_essay}
          </span>
        </div>

        {/* Divider & Actions */}
        <div className="mt-auto border-t border-gray-100 pt-4">
          <div className="flex items-center justify-between gap-3">
            {/* Primary Button: Preview */}
            <button
              type="button"
              onClick={onPreview}
              disabled={!onPreview}
              className={[
                "inline-flex cursor-pointer items-center gap-2 rounded-full px-5 py-2 text-sm font-bold text-white shadow-md transition-all",
                "bg-[#397e50] hover:bg-[#2f5c3f] hover:shadow-lg active:scale-95",
                "disabled:cursor-not-allowed disabled:opacity-50",
              ].join(" ")}
            >
              <svg
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth="2.5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                />
              </svg>
              Preview
            </button>

            {/* Icon Buttons Group */}
            <div className="flex items-center gap-1">
              {/* Kelola */}
              <button
                type="button"
                onClick={onKelola}
                disabled={!onKelola}
                className={[
                  "cursor-pointer rounded-full p-2 text-gray-500 transition-colors",
                  "hover:bg-emerald-50 hover:text-[#397e50]",
                  "disabled:cursor-not-allowed disabled:opacity-50",
                ].join(" ")}
                title="Kelola"
              >
                <svg
                  className="h-5 w-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  />
                </svg>
              </button>

              {/* Hapus */}
              <button
                type="button"
                onClick={onHapus}
                disabled={!onHapus}
                className={[
                  "cursor-pointer rounded-full p-2 text-gray-400 transition-colors",
                  "hover:bg-rose-50 hover:text-rose-600",
                  "disabled:cursor-not-allowed disabled:opacity-50",
                ].join(" ")}
                title="Hapus"
              >
                <svg
                  className="h-5 w-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
