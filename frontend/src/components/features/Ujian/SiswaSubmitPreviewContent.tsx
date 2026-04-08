import React from "react";

import {
  formatSoalTypeLabel,
  isPilihanGandaSoal,
} from "@/helper/Ujian/soalType";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type { JawabanUjianSiswaItem } from "@/types/Ujian/ujianSiswa";

type SiswaSubmitPreviewContentProps = {
  title: string;
  sisaWaktu: string;
  soalPreview: SoalPreviewItem[];
  jawabanItems: JawabanUjianSiswaItem[];
  onBackToQuestionMode: () => void;
  onSubmitFinal: () => void;
  backDisabled?: boolean;
  submitDisabled?: boolean;
  submitLoading?: boolean;
};

const SiswaSubmitPreviewContent: React.FC<SiswaSubmitPreviewContentProps> = ({
  title,
  sisaWaktu,
  soalPreview,
  jawabanItems,
  onBackToQuestionMode,
  onSubmitFinal,
  backDisabled = false,
  submitDisabled = false,
  submitLoading = false,
}) => {
  const getJawabanBySoalId = React.useCallback(
    (soalId: number) =>
      jawabanItems.find((item) => String(item.id_soal) === String(soalId)),
    [jawabanItems],
  );

  const isAnswered = React.useCallback(
    (soalId: number) => {
      const jawaban = getJawabanBySoalId(soalId);
      return Boolean(
        jawaban &&
        (jawaban.id_pilihan !== null ||
          (jawaban.jawaban_essay?.trim() ?? "") !== ""),
      );
    },
    [getJawabanBySoalId],
  );

  const answeredCount = soalPreview.filter((soal) =>
    isAnswered(soal.id),
  ).length;

  const scrollToQuestion = React.useCallback((soalId: number) => {
    const element = document.getElementById(`submit-preview-soal-${soalId}`);
    element?.scrollIntoView({ behavior: "smooth", block: "start" });
  }, []);

  return (
    <div className="min-h-screen w-full bg-slate-50">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-4 py-6 md:px-6">
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
              <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                  Sisa Waktu
                </p>
                <p className="mt-1 text-lg font-semibold text-[#397e50]">
                  {sisaWaktu}
                </p>
              </div>

              <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                  Status Jawaban
                </p>
                <p className="mt-1 text-lg font-semibold text-slate-900">
                  {answeredCount} / {soalPreview.length} terjawab
                </p>
              </div>
            </div>
          </div>
        </header>

        <div className="grid gap-6 lg:grid-cols-[minmax(0,820px)_240px] lg:justify-center lg:items-start">
          <section className="space-y-4">
            {soalPreview.map((soal) => {
              const jawaban = getJawabanBySoalId(soal.id);
              const selectedOptionId = jawaban?.id_pilihan ?? null;
              const essayAnswer = jawaban?.jawaban_essay?.trim() ?? "";
              const soalAnswered = isAnswered(soal.id);
              const isPilihanGanda = isPilihanGandaSoal(soal.tipe);

              return (
                <article
                  id={`submit-preview-soal-${soal.id}`}
                  key={soal.id}
                  className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm md:p-6"
                >
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
                        soalAnswered
                          ? "bg-[#397e50]/10 text-[#397e50]"
                          : "bg-amber-100 text-amber-700",
                      ].join(" ")}
                    >
                      {soalAnswered ? "Terjawab" : "Belum Dijawab"}
                    </span>
                  </div>

                  <p className="mt-4 whitespace-pre-wrap text-sm leading-relaxed text-slate-600">
                    {soal.pertanyaan}
                  </p>

                  {soal.gambar_url && (
                    <div className="mt-4 overflow-hidden rounded-xl border border-slate-200 bg-slate-50">
                      <img
                        src={soal.gambar_url}
                        alt={`Ilustrasi soal ${soal.nomor}`}
                        className="max-h-[420px] w-full object-contain"
                      />
                    </div>
                  )}

                  {isPilihanGanda ? (
                    <div className="mt-5 space-y-2">
                      {soal.opsi.map((option) => {
                        const isSelected = selectedOptionId === option.id;
                        return (
                          <div
                            key={option.id}
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
                            <span className="flex-1">{option.text}</span>
                          </div>
                        );
                      })}

                      {!soalAnswered && (
                        <p className="pt-1 text-sm font-medium text-amber-700">
                          Belum dijawab.
                        </p>
                      )}
                    </div>
                  ) : (
                    <div className="mt-5 whitespace-pre-wrap rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm leading-relaxed text-slate-600">
                      {essayAnswer !== "" ? essayAnswer : "Belum dijawab."}
                    </div>
                  )}
                </article>
              );
            })}
          </section>

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
                <button
                  key={soal.id}
                  type="button"
                  onClick={() => scrollToQuestion(soal.id)}
                  className={[
                    "flex h-9 w-9 items-center justify-center rounded-lg border text-xs font-semibold transition",
                    isAnswered(soal.id)
                      ? "border-[#397e50] bg-[#397e50] text-white"
                      : "border-slate-200 text-slate-500 hover:border-[#397e50] hover:text-[#397e50]",
                  ].join(" ")}
                  aria-label={`Lompat ke soal nomor ${soal.nomor}`}
                >
                  {soal.nomor}
                </button>
              ))}
            </div>
          </aside>
        </div>

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
      </div>
    </div>
  );
};

export default SiswaSubmitPreviewContent;
