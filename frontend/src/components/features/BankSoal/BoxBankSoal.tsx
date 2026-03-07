import React from "react";

import SvgIcons from "@/assets/SvgIcons/svgIcons";
type BoxBankSoalProps = {
  idBankSoal: number;
  namaBankSoal?: string;
  guruLabel?: string;
  mapelLabel?: string;
  materi?: string;
  kelasLabel?: string;
  tglBuat?: string;
  soalUploaded?: boolean;
  onPreview?: (idBankSoal: number) => void;
  onUpload?: () => void;
  onKelola?: () => void;
  onHapus?: () => void;
  className?: string;
};

const kelasLabelClass = (kelasLabel?: string) => {
  const kk = String(kelasLabel ?? "");
  const base = "rounded-md px-2 py-1 text-sm font-bold text-white shadow-sm";

  const map: Record<string, string> = {
    "Kelas 10": "bg-green-700",
    "Kelas 11": "bg-green-800",
    "Kelas 12": "bg-green-900",
  };

  return [base, map[kk] ?? "bg-gray-400"].join(" ");
};

const BoxBankSoal: React.FC<BoxBankSoalProps> = ({
  idBankSoal,
  namaBankSoal = "Bank Soal",
  guruLabel = "-",
  mapelLabel = "-",
  materi = "-",
  kelasLabel = "-",
  tglBuat,
  soalUploaded = false,
  onUpload,
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
            {SvgIcons.calendar("h-4 w-4 text-gray-400")}
            {tglBuat || "-"}
          </span>

          <span className={kelasLabelClass(kelasLabel)}>{kelasLabel}</span>
        </div>

        {/* Title (Lebih Besar) */}
        <h3 className="mb-2 line-clamp-2 text-lg font-bold leading-snug text-[#37513d] transition-colors group-hover:text-[#397e50]">
          {namaBankSoal}
        </h3>

        {/* Guru pembuat soal */}
        <h4 className="mb-7 truncate  text-gray-700 text-sm flex items-center gap-2">
          <span className="font-semibold">Dibuat Oleh :</span> {guruLabel}
        </h4>

        {/* Info Content (Text SM agar terbaca jelas) */}
        <div className="mb-4 space-y-2 text-sm">
          <div className="flex gap-3">
            <span className="min-w-20 font-medium text-gray-500">Mapel:</span>
            <span className="truncate font-semibold text-gray-700">
              {mapelLabel}
            </span>
          </div>

          <div className="flex gap-3">
            <span className="min-w-20 font-medium text-gray-500">Materi:</span>
            <span className="truncate font-semibold text-gray-700">
              {materi}
            </span>
          </div>
        </div>

        {/* Divider & Actions */}
        <div className="mt-auto border-t border-gray-100 pt-4">
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2">
              {/* Primary Button: Preview */}
              <button
                type="button"
                onClick={() => onPreview?.(idBankSoal)}
                disabled={!onPreview || !soalUploaded}
                className={[
                  "inline-flex cursor-pointer items-center gap-2 rounded-full px-5 py-2 text-sm font-bold text-white shadow-md transition-all",
                  "bg-[#397e50] hover:bg-[#2f5c3f] hover:shadow-lg active:scale-95",
                  "disabled:cursor-not-allowed disabled:opacity-50",
                ].join(" ")}
              >
                {SvgIcons.eye("h-4 w-4")}
                Preview
              </button>

              <button
                type="button"
                onClick={onUpload}
                disabled={!onUpload}
                className="inline-flex cursor-pointer items-center rounded-full border border-[#397e50]/30 px-4 py-2 text-sm font-semibold text-[#397e50] transition-colors hover:bg-[#397e50]/10"
              >
                {soalUploaded ? "Perbarui soal" : "Upload bank soal"}
              </button>
            </div>

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
                {SvgIcons.edit("h-5 w-5")}
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
                {SvgIcons.trash("h-5 w-5")}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BoxBankSoal;
