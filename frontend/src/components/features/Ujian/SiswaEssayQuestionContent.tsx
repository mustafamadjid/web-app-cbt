import React from "react";

import { formatSoalTypeLabel } from "@/helper/Ujian/soalType";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";

type SiswaEssayQuestionContentProps = {
  soal: SoalPreviewItem;
  value: string;
  onChange: (value: string) => void;
};

const SiswaEssayQuestionContent: React.FC<SiswaEssayQuestionContentProps> = ({
  soal,
  value,
  onChange,
}) => {
  return (
    <div className="space-y-5">
      <div>
        <p className="text-xs uppercase tracking-wide text-slate-400">
          {formatSoalTypeLabel(soal.tipe)}
        </p>
        <p className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-slate-600">
          {soal.pertanyaan}
        </p>
      </div>

      {soal.gambar_url && (
        <div className="overflow-hidden rounded-xl border border-slate-200 bg-slate-50">
          <img
            src={soal.gambar_url}
            alt={`Ilustrasi soal ${soal.nomor}`}
            className="max-h-[420px] w-full object-contain"
          />
        </div>
      )}

      <div className="space-y-3">
        <label
          htmlFor={`essay-answer-${soal.id}`}
          className="block text-xs font-semibold text-slate-400"
        >
          Jawaban Essay
        </label>

        <textarea
          id={`essay-answer-${soal.id}`}
          value={value}
          onChange={(event) => onChange(event.target.value)}
          rows={10}
          className="w-full resize-y rounded-xl border border-slate-200 px-4 py-3 text-sm leading-relaxed text-slate-600 outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50]/20"
          placeholder="Tulis jawaban essay di sini..."
        />
      </div>
    </div>
  );
};

export default SiswaEssayQuestionContent;
