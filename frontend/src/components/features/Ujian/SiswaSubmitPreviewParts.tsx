import React from "react";

import RichContentRenderer from "@/components/common/RichContentRenderer";
import {
  formatSoalTypeLabel,
  isPilihanGandaSoal,
} from "@/helper/Ujian/soalType";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type { JawabanUjianSiswaItem } from "@/types/Ujian/ujianSiswa";

export type GetJawabanBySoalId = (
  soalId: number,
) => JawabanUjianSiswaItem | undefined;
export type IsAnswered = (soalId: number) => boolean;

type SubmitPreviewHeaderProps = {
  title: string;
  sisaWaktu: string;
  answeredCount: number;
  totalSoal: number;
};

export const SubmitPreviewHeader: React.FC<SubmitPreviewHeaderProps> = ({
  title,
  sisaWaktu,
  answeredCount,
  totalSoal,
}) => (
  <header className="rounded-xl border border-slate-200 bg-white px-6 py-5 shadow-sm">
    <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
      <div>
        <p className="text-sm font-semibold text-slate-800">{title}</p>
        <h1 className="mt-1 text-2xl font-bold text-slate-900">
          Preview Submit Ujian
        </h1>
        <p className="mt-2 text-sm text-slate-500">
          Tinjau seluruh jawaban tersimpan sebelum submit final.
        </p>
      </div>

      <div className="grid gap-3 sm:grid-cols-2">
        <SummaryCard label="Sisa Waktu" value={sisaWaktu} highlight />
        <SummaryCard
          label="Status Jawaban"
          value={`${answeredCount} / ${totalSoal} terjawab`}
        />
      </div>
    </div>
  </header>
);

type SummaryCardProps = {
  label: string;
  value: string;
  highlight?: boolean;
};

const SummaryCard: React.FC<SummaryCardProps> = ({
  label,
  value,
  highlight = false,
}) => (
  <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
    <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
      {label}
    </p>
    <p
      className={[
        "mt-1 text-lg font-semibold",
        highlight ? "text-[#397e50]" : "text-slate-900",
      ].join(" ")}
    >
      {value}
    </p>
  </div>
);

type QuestionListProps = {
  soalPreview: SoalPreviewItem[];
  getJawabanBySoalId: GetJawabanBySoalId;
  isAnswered: IsAnswered;
};

export const QuestionList: React.FC<QuestionListProps> = ({
  soalPreview,
  getJawabanBySoalId,
  isAnswered,
}) => (
  <section className="space-y-4">
    {soalPreview.map((soal) => (
      <QuestionCard
        key={soal.id}
        soal={soal}
        jawaban={getJawabanBySoalId(soal.id)}
        isAnswered={isAnswered(soal.id)}
      />
    ))}
  </section>
);

type QuestionCardProps = {
  soal: SoalPreviewItem;
  jawaban?: JawabanUjianSiswaItem;
  isAnswered: boolean;
};

const QuestionCard: React.FC<QuestionCardProps> = ({
  soal,
  jawaban,
  isAnswered,
}) => {
  const selectedOptionId = jawaban?.id_pilihan ?? null;
  const essayAnswer = jawaban?.jawaban_essay?.trim() ?? "";
  const isPilihanGanda = isPilihanGandaSoal(soal.tipe);

  return (
    <article
      id={`submit-preview-soal-${soal.id}`}
      className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm md:p-6"
    >
      <QuestionHeader soal={soal} isAnswered={isAnswered} />

      <RichContentRenderer
        content={soal.pertanyaan_content}
        fallbackText={soal.pertanyaan}
        className="mt-4"
        paragraphClassName="whitespace-pre-wrap text-sm leading-relaxed text-slate-600"
      />

      {soal.gambar_url && (
        <QuestionImage gambarUrl={soal.gambar_url} nomor={soal.nomor} />
      )}

      {isPilihanGanda ? (
        <OptionList
          options={soal.opsi}
          selectedOptionId={selectedOptionId}
          isAnswered={isAnswered}
        />
      ) : (
        <EssayPreview answer={essayAnswer} />
      )}
    </article>
  );
};

type QuestionHeaderProps = {
  soal: SoalPreviewItem;
  isAnswered: boolean;
};

const QuestionHeader: React.FC<QuestionHeaderProps> = ({
  soal,
  isAnswered,
}) => (
  <div className="flex flex-wrap items-start justify-between gap-3">
    <div>
      <p className="text-xs uppercase tracking-wide text-slate-400">
        {formatSoalTypeLabel(soal.tipe)}
      </p>
      <h2 className="mt-1 text-lg font-semibold text-slate-900">
        Soal {soal.nomor}
      </h2>
    </div>

    <span
      className={[
        "rounded-full px-3 py-1 text-xs font-semibold",
        isAnswered
          ? "bg-[#397e50]/10 text-[#397e50]"
          : "bg-amber-100 text-amber-700",
      ].join(" ")}
    >
      {isAnswered ? "Terjawab" : "Belum Dijawab"}
    </span>
  </div>
);

type QuestionImageProps = {
  gambarUrl: string;
  nomor: number;
};

const QuestionImage: React.FC<QuestionImageProps> = ({ gambarUrl, nomor }) => (
  <div className="mt-4 overflow-hidden rounded-xl border border-slate-200 bg-slate-50">
    <img
      src={gambarUrl}
      alt={`Ilustrasi soal ${nomor}`}
      className="max-h-[420px] w-full object-contain"
    />
  </div>
);

type OptionListProps = {
  options: SoalPreviewItem["opsi"];
  selectedOptionId: number | null;
  isAnswered: boolean;
};

const OptionList: React.FC<OptionListProps> = ({
  options,
  selectedOptionId,
  isAnswered,
}) => (
  <div className="mt-5 space-y-2">
    {options.map((option) => (
      <OptionRow
        key={option.id}
        option={option}
        isSelected={selectedOptionId === option.id}
      />
    ))}

    {!isAnswered && (
      <p className="pt-1 text-sm font-medium text-amber-700">
        Belum dijawab.
      </p>
    )}
  </div>
);

type OptionRowProps = {
  option: SoalPreviewItem["opsi"][number];
  isSelected: boolean;
};

const OptionRow: React.FC<OptionRowProps> = ({ option, isSelected }) => (
  <div
    className={[
      "flex items-start gap-3 rounded-xl border px-4 py-3 text-sm",
      isSelected
        ? "border-[#397e50] bg-[#397e50]/10 text-[#397e50]"
        : "border-slate-200 text-slate-500",
    ].join(" ")}
  >
    <span
      className={[
        "mt-0.5 flex h-6 w-6 items-center justify-center rounded-full border text-xs font-semibold",
        isSelected
          ? "border-[#397e50] bg-[#397e50] text-white"
          : "border-slate-200 text-slate-500",
      ].join(" ")}
    >
      {option.label}
    </span>
    <span className="flex-1">
      <RichContentRenderer
        content={option.content}
        fallbackText={option.text}
        inline
        paragraphClassName="whitespace-pre-wrap text-sm leading-relaxed text-inherit"
      />
    </span>
  </div>
);

type EssayPreviewProps = {
  answer: string;
};

const EssayPreview: React.FC<EssayPreviewProps> = ({ answer }) => (
  <div className="mt-5 whitespace-pre-wrap rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm leading-relaxed text-slate-600">
    {answer !== "" ? answer : "Belum dijawab."}
  </div>
);

type QuestionNavigatorProps = {
  soalPreview: SoalPreviewItem[];
  isAnswered: IsAnswered;
  onNavigate: (soalId: number) => void;
};

export const QuestionNavigator: React.FC<QuestionNavigatorProps> = ({
  soalPreview,
  isAnswered,
  onNavigate,
}) => (
  <aside className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm md:p-5 lg:sticky lg:top-6">
    <div className="space-y-1 border-b border-slate-100 pb-3">
      <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
        Navigasi
      </p>
      <p className="text-sm font-semibold text-slate-700">Nomor Soal</p>
      <p className="text-xs text-slate-500">
        Klik nomor untuk lompat ke soal terkait.
      </p>
    </div>

    <div className="mt-4 grid grid-cols-4 gap-2">
      {soalPreview.map((soal) => (
        <QuestionNavigatorButton
          key={soal.id}
          soal={soal}
          isAnswered={isAnswered(soal.id)}
          onNavigate={onNavigate}
        />
      ))}
    </div>
  </aside>
);

type QuestionNavigatorButtonProps = {
  soal: SoalPreviewItem;
  isAnswered: boolean;
  onNavigate: (soalId: number) => void;
};

const QuestionNavigatorButton: React.FC<QuestionNavigatorButtonProps> = ({
  soal,
  isAnswered,
  onNavigate,
}) => (
  <button
    type="button"
    onClick={() => onNavigate(soal.id)}
    className={[
      "flex h-9 w-9 items-center justify-center rounded-lg border text-xs font-semibold transition",
      isAnswered
        ? "border-[#397e50] bg-[#397e50] text-white"
        : "border-slate-200 text-slate-500 hover:border-[#397e50] hover:text-[#397e50]",
    ].join(" ")}
    aria-label={`Lompat ke soal nomor ${soal.nomor}`}
  >
    {soal.nomor}
  </button>
);

type SubmitActionsProps = {
  onBackToQuestionMode: () => void;
  onSubmitFinal: () => void;
  backDisabled: boolean;
  submitDisabled: boolean;
  submitLoading: boolean;
};

export const SubmitActions: React.FC<SubmitActionsProps> = ({
  onBackToQuestionMode,
  onSubmitFinal,
  backDisabled,
  submitDisabled,
  submitLoading,
}) => (
  <div className="sticky bottom-4 flex flex-col gap-3 rounded-2xl border border-slate-200 bg-white p-4 shadow-lg md:flex-row md:items-center md:justify-end">
    <button
      type="button"
      onClick={onBackToQuestionMode}
      disabled={backDisabled}
      className={[
        "rounded-lg border px-4 py-2 text-sm font-semibold transition",
        backDisabled
          ? "cursor-not-allowed border-slate-200 text-slate-300"
          : "cursor-pointer border-slate-200 text-slate-600 hover:border-[#397e50] hover:text-[#397e50]",
      ].join(" ")}
    >
      Kembali ke Pengerjaan
    </button>

    <button
      type="button"
      onClick={onSubmitFinal}
      disabled={submitDisabled}
      className={[
        "rounded-lg px-4 py-2 text-sm font-semibold text-white transition",
        submitDisabled
          ? "cursor-not-allowed bg-slate-300"
          : "cursor-pointer bg-[#397e50] hover:bg-[#326f45]",
      ].join(" ")}
    >
      {submitLoading ? "Submit..." : "Submit Final"}
    </button>
  </div>
);
