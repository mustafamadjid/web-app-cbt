import React from "react";

import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import {
  formatSoalTypeLabel,
  isEssaySoal,
  isPilihanGandaSoal,
} from "@/helper/Ujian/soalType";
import type { HasilJawabanUjianItem } from "@/types/Ujian/HasilJawabanUjian";
import type { SubmitKoreksiEssayRequest } from "@/types/Ujian/SubmitKoreksiEssay";

type HasilJawabanUjianContentProps = {
  title: string;
  subtitle?: string;
  nilaiAkhir: number | null;
  hasilJawabanUjian: HasilJawabanUjianItem[];
  canGradeEssay?: boolean;
  submitEssayCorrection?: (
    payload: SubmitKoreksiEssayRequest,
  ) => void | Promise<void>;
  submitDisabled?: boolean;
  submitLoading?: boolean;
};

type EssayCorrectionMap = Record<number, boolean | null>;
type QuestionTone = "success" | "danger" | "warning";

type QuestionStatus = {
  label: string;
  tone: QuestionTone;
};

const STATUS_THEME: Record<
  QuestionTone,
  {
    articleClassName: string;
    badgeClassName: string;
    navClassName: string;
  }
> = {
  success: {
    articleClassName: "border-emerald-200 bg-emerald-50/50",
    badgeClassName: "bg-emerald-100 text-emerald-700",
    navClassName: "border-emerald-600 bg-emerald-600 text-white",
  },
  danger: {
    articleClassName: "border-rose-200 bg-rose-50/50",
    badgeClassName: "bg-rose-100 text-rose-700",
    navClassName: "border-rose-600 bg-rose-600 text-white",
  },
  warning: {
    articleClassName: "border-amber-200 bg-amber-50/60",
    badgeClassName: "bg-amber-100 text-amber-700",
    navClassName: "border-amber-500 bg-amber-400 text-slate-900",
  },
};

const buildInitialEssayCorrections = (
  items: HasilJawabanUjianItem[],
): EssayCorrectionMap => {
  const next: EssayCorrectionMap = {};

  items.forEach((item) => {
    if (!isEssaySoal(item.tipe_soal) || item.jawaban_siswa.id_jawaban === null) {
      return;
    }

    next[item.jawaban_siswa.id_jawaban] = item.jawaban_siswa.essay_is_benar;
  });

  return next;
};

const formatNilaiAkhir = (value: number | null) => {
  if (typeof value !== "number") {
    return "-";
  }

  return value.toLocaleString("id-ID", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
};

const getEssayAnswer = (item: HasilJawabanUjianItem) =>
  item.jawaban_siswa.jawaban_essay?.trim() ?? "";

const isItemAnswered = (item: HasilJawabanUjianItem) =>
  item.jawaban_siswa.id_pilihan !== null || getEssayAnswer(item) !== "";

const getSelectedOption = (item: HasilJawabanUjianItem) =>
  item.opsi_jawaban.find(
    (option) => option.id_pilihan_ganda === item.jawaban_siswa.id_pilihan,
  );

const getCorrectOptions = (item: HasilJawabanUjianItem) =>
  item.opsi_jawaban
    .map((option, index) => ({
      ...option,
      optionLabel: String.fromCharCode(65 + index),
    }))
    .filter((option) => option.is_benar);

const getQuestionStatus = (
  item: HasilJawabanUjianItem,
  essayCorrections: EssayCorrectionMap,
): QuestionStatus => {
  if (isPilihanGandaSoal(item.tipe_soal)) {
    if (item.jawaban_siswa.id_pilihan === null) {
      return {
        label: "Belum dijawab",
        tone: "danger",
      };
    }

    return getSelectedOption(item)?.is_benar
      ? {
          label: "Benar",
          tone: "success",
        }
      : {
          label: "Salah",
          tone: "danger",
        };
  }

  const essayAnswer = getEssayAnswer(item);
  if (essayAnswer === "" || item.jawaban_siswa.id_jawaban === null) {
    return {
      label: "Belum dijawab",
      tone: "danger",
    };
  }

  const essayResult = essayCorrections[item.jawaban_siswa.id_jawaban] ?? null;

  if (essayResult === true) {
    return {
      label: "Benar",
      tone: "success",
    };
  }

  if (essayResult === false) {
    return {
      label: "Salah",
      tone: "danger",
    };
  }

  return {
    label: "Belum dikoreksi",
    tone: "warning",
  };
};

const HasilJawabanUjianContent: React.FC<HasilJawabanUjianContentProps> = ({
  title,
  subtitle = "Tinjau jawaban siswa dan lakukan koreksi untuk soal essay.",
  nilaiAkhir,
  hasilJawabanUjian,
  canGradeEssay = false,
  submitEssayCorrection,
  submitDisabled = false,
  submitLoading = false,
}) => {
  const [essayCorrections, setEssayCorrections] = React.useState<EssayCorrectionMap>(
    () => buildInitialEssayCorrections(hasilJawabanUjian),
  );

  React.useEffect(() => {
    setEssayCorrections(buildInitialEssayCorrections(hasilJawabanUjian));
  }, [hasilJawabanUjian]);

  const essayItems = hasilJawabanUjian.filter((item) => isEssaySoal(item.tipe_soal));
  const answeredCount = hasilJawabanUjian.filter(isItemAnswered).length;
  const reviewedEssayCount = essayItems.filter((item) => {
    const essayId = item.jawaban_siswa.id_jawaban;
    if (essayId === null || getEssayAnswer(item) === "") {
      return false;
    }

    return essayCorrections[essayId] !== null;
  }).length;
  const pendingEssayItems = essayItems.filter((item) => {
    const essayId = item.jawaban_siswa.id_jawaban;
    return (
      essayId !== null &&
      getEssayAnswer(item) !== "" &&
      item.jawaban_siswa.essay_is_benar === null
    );
  });
  const pendingEssaySelections = pendingEssayItems
    .map((item) => {
      const essayId = item.jawaban_siswa.id_jawaban;
      if (essayId === null) {
        return null;
      }

      const result = essayCorrections[essayId];
      if (result === null || typeof result !== "boolean") {
        return null;
      }

      return {
        id_jawaban: essayId,
        essay_is_benar: result,
      };
    })
    .filter(
      (item): item is SubmitKoreksiEssayRequest["jawaban"][number] => item !== null,
    );
  const isPendingEssayComplete =
    pendingEssayItems.length > 0 &&
    pendingEssaySelections.length === pendingEssayItems.length;
  const shouldShowSubmit =
    canGradeEssay &&
    typeof submitEssayCorrection === "function" &&
    pendingEssayItems.length > 0;

  const scrollToQuestion = (soalId: number) => {
    const element = document.getElementById(`hasil-jawaban-soal-${soalId}`);
    element?.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  const handleSubmitEssayCorrection = () => {
    if (!submitEssayCorrection || pendingEssaySelections.length === 0) {
      return;
    }

    void submitEssayCorrection({
      jawaban: pendingEssaySelections,
    });
  };

  if (hasilJawabanUjian.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
        Hasil jawaban ujian belum tersedia.
      </div>
    );
  }



  return (
    <div className="flex w-full flex-col gap-6">
      <header className="rounded-2xl border border-slate-200 bg-white px-6 py-5 shadow-sm">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
              Hasil Jawaban
            </p>
            <h1 className="mt-1 text-2xl font-bold text-slate-900">{title}</h1>
            <p className="mt-2 text-sm text-slate-500">{subtitle}</p>
          </div>

          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
            <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                Nilai Akhir
              </p>
              <p className="mt-1 text-lg font-semibold text-slate-900">
                {formatNilaiAkhir(nilaiAkhir)}
              </p>
            </div>

            <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                Total Soal
              </p>
              <p className="mt-1 text-lg font-semibold text-slate-900">
                {hasilJawabanUjian.length}
              </p>
            </div>

            <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                Terjawab
              </p>
              <p className="mt-1 text-lg font-semibold text-[#397e50]">
                {answeredCount} / {hasilJawabanUjian.length}
              </p>
            </div>

            <div className="rounded-xl border border-slate-200 bg-slate-50 px-4 py-3">
              <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                Essay Dikoreksi
              </p>
              <p className="mt-1 text-lg font-semibold text-slate-900">
                {reviewedEssayCount} / {essayItems.length}
              </p>
            </div>
          </div>
        </div>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_240px] lg:items-start">
        <section className="space-y-4">
          {hasilJawabanUjian.map((item, index) => {
            const status = getQuestionStatus(item, essayCorrections);
            const theme = STATUS_THEME[status.tone];
            const selectedOptionId = item.jawaban_siswa.id_pilihan;
            const essayAnswer = getEssayAnswer(item);
            const gambarUrl = resolveImageUrl(item.gambar) || "";
            const essayId = item.jawaban_siswa.id_jawaban;
            const essayCorrectionValue =
              essayId === null ? null : (essayCorrections[essayId] ?? null);
            const correctOptions = getCorrectOptions(item);
            const canEditEssay =
              canGradeEssay &&
              essayId !== null &&
              essayAnswer !== "" &&
              item.jawaban_siswa.essay_is_benar === null;

            return (
              <article
                id={`hasil-jawaban-soal-${item.id_soal}`}
                key={item.id_soal}
                className={[
                  "rounded-2xl border p-5 shadow-sm md:p-6",
                  theme.articleClassName,
                ].join(" ")}
              >
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="text-xs uppercase tracking-wide text-slate-500">
                      {formatSoalTypeLabel(item.tipe_soal)}
                    </p>
                    <h2 className="mt-1 text-lg font-semibold text-slate-900">
                      Soal {item.no_urut_soal || index + 1}
                    </h2>
                    <p className="mt-1 text-xs text-slate-500">
                      Bobot soal: {item.bobot_soal}
                    </p>
                  </div>

                  <span
                    className={[
                      "rounded-full px-3 py-1 text-xs font-semibold",
                      theme.badgeClassName,
                    ].join(" ")}
                  >
                    {status.label}
                  </span>
                </div>

                <p className="mt-4 whitespace-pre-wrap text-sm leading-relaxed text-slate-700">
                  {item.pertanyaan}
                </p>

                {gambarUrl && (
                  <div className="mt-4 overflow-hidden rounded-xl border border-slate-200 bg-white">
                    <img
                      src={gambarUrl}
                      alt={`Ilustrasi soal ${item.no_urut_soal || index + 1}`}
                      className="max-h-[420px] w-full object-contain"
                    />
                  </div>
                )}

                {isPilihanGandaSoal(item.tipe_soal) ? (
                  <div className="mt-5 space-y-3">
                    <div className="space-y-2">
                      {item.opsi_jawaban.map((option, optionIndex) => {
                        const isSelected = selectedOptionId === option.id_pilihan_ganda;
                        const optionLabel = String.fromCharCode(65 + optionIndex);
                        const optionClass = option.is_benar
                          ? "border-emerald-200 bg-emerald-50 text-emerald-700"
                          : isSelected
                            ? "border-rose-300 bg-rose-50 text-rose-700"
                            : "border-slate-200 bg-white text-slate-600";
                        const bulletClass = option.is_benar
                          ? "border-emerald-500 bg-emerald-500 text-white"
                          : isSelected
                            ? "border-rose-500 bg-rose-500 text-white"
                            : "border-slate-300 text-slate-500";

                        return (
                          <div
                            key={option.id_pilihan_ganda}
                            className={["rounded-xl border px-4 py-3 text-sm", optionClass].join(
                              " ",
                            )}
                          >
                            <div className="flex items-start gap-3">
                              <span
                                className={[
                                  "mt-0.5 flex h-6 w-6 items-center justify-center rounded-full border text-xs font-semibold",
                                  bulletClass,
                                ].join(" ")}
                              >
                                {optionLabel}
                              </span>

                              <div className="flex-1 space-y-2">
                                <p>{option.isi_pilihan}</p>

                                <div className="flex flex-wrap gap-2">
                                  {isSelected && (
                                    <span className="rounded-full bg-white/90 px-2.5 py-1 text-xs font-semibold">
                                      Dipilih siswa
                                    </span>
                                  )}
                                  {option.is_benar && (
                                    <span className="rounded-full bg-white/90 px-2.5 py-1 text-xs font-semibold">
                                      Opsi benar
                                    </span>
                                  )}
                                </div>
                              </div>
                            </div>
                          </div>
                        );
                      })}
                    </div>

                    <div className="rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800">
                      <p className="text-xs font-semibold uppercase tracking-wide text-emerald-700">
                        Kunci Jawaban
                      </p>
                      <div className="mt-2 space-y-1">
                        {correctOptions.length > 0 ? (
                          correctOptions.map((option) => (
                            <p key={option.id_pilihan_ganda}>
                              {option.optionLabel}. {option.isi_pilihan}
                            </p>
                          ))
                        ) : (
                          <p>Kunci jawaban belum tersedia.</p>
                        )}
                      </div>
                    </div>

                    {selectedOptionId === null && (
                      <p className="text-sm font-medium text-rose-700">
                        Siswa belum memilih jawaban untuk soal ini.
                      </p>
                    )}
                  </div>
                ) : (
                  <div className="mt-5 space-y-4">
                    <div className="rounded-xl border border-slate-200 bg-white px-4 py-3">
                      <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                        Jawaban Essay Siswa
                      </p>
                      <p className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-slate-700">
                        {essayAnswer !== "" ? essayAnswer : "Belum dijawab."}
                      </p>
                    </div>

                    {essayId !== null && essayAnswer !== "" && (
                      <div className="rounded-xl border border-slate-200 bg-white px-4 py-4">
                        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                          <div>
                            <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                              Koreksi Essay
                            </p>
                            <p className="mt-1 text-sm text-slate-600">
                              {canEditEssay
                                ? "Tentukan apakah jawaban essay siswa benar atau salah."
                                : essayCorrectionValue === true
                                  ? "Jawaban essay sudah dinilai benar."
                                  : essayCorrectionValue === false
                                    ? "Jawaban essay sudah dinilai salah."
                                    : "Jawaban essay ini masih menunggu koreksi."}
                            </p>
                          </div>

                          <span className="text-xs font-semibold text-slate-500">
                            ID Jawaban: {essayId}
                          </span>
                        </div>

                        {canEditEssay ? (
                          <div className="mt-4 flex flex-wrap gap-2">
                            <button
                              type="button"
                              onClick={() =>
                                setEssayCorrections((prev) => ({
                                  ...prev,
                                  [essayId]: true,
                                }))
                              }
                              className={[
                                "cursor-pointer rounded-lg border px-3 py-2 text-sm font-semibold transition",
                                essayCorrectionValue === true
                                  ? "border-emerald-600 bg-emerald-600 text-white"
                                  : "border-emerald-200 bg-white text-emerald-700 hover:border-emerald-400",
                              ].join(" ")}
                            >
                              Benar
                            </button>

                            <button
                              type="button"
                              onClick={() =>
                                setEssayCorrections((prev) => ({
                                  ...prev,
                                  [essayId]: false,
                                }))
                              }
                              className={[
                                "cursor-pointer rounded-lg border px-3 py-2 text-sm font-semibold transition",
                                essayCorrectionValue === false
                                  ? "border-rose-600 bg-rose-600 text-white"
                                  : "border-rose-200 bg-white text-rose-700 hover:border-rose-400",
                              ].join(" ")}
                            >
                              Salah
                            </button>
                          </div>
                        ) : (
                          <div className="mt-4">
                            <span
                              className={[
                                "inline-flex rounded-full px-3 py-1 text-xs font-semibold",
                                STATUS_THEME[status.tone].badgeClassName,
                              ].join(" ")}
                            >
                              {status.label}
                            </span>
                          </div>
                        )}
                      </div>
                    )}
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
              Klik nomor untuk lompat ke hasil jawaban soal.
            </p>
          </div>

          <div className="mt-4 grid grid-cols-4 gap-2">
            {hasilJawabanUjian.map((item, index) => {
              const status = getQuestionStatus(item, essayCorrections);
              const theme = STATUS_THEME[status.tone];

              return (
                <button
                  key={item.id_soal}
                  type="button"
                  onClick={() => scrollToQuestion(item.id_soal)}
                  className={[
                    "flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg border text-xs font-semibold transition hover:border-[#397e50] hover:bg-[#397e50] hover:text-white",
                    theme.navClassName,
                  ].join(" ")}
                  aria-label={`Lompat ke soal nomor ${item.no_urut_soal || index + 1}`}
                >
                  {item.no_urut_soal || index + 1}
                </button>
              );
            })}
          </div>
        </aside>
      </div>

      {shouldShowSubmit && (
        <div className="sticky bottom-4 flex flex-col gap-3 rounded-2xl border border-slate-200 bg-white p-4 shadow-lg md:flex-row md:items-center md:justify-between">
          <p className="text-sm text-slate-600">
            {isPendingEssayComplete
              ? `Siap mengirim ${pendingEssaySelections.length} koreksi essay.`
              : `Lengkapi ${pendingEssayItems.length - pendingEssaySelections.length} essay yang belum dipilih.`}
          </p>

          <button
            type="button"
            onClick={handleSubmitEssayCorrection}
            disabled={submitDisabled || submitLoading || !isPendingEssayComplete}
            className={[
              "rounded-lg px-4 py-2 text-sm font-semibold text-white transition",
              submitDisabled || submitLoading || !isPendingEssayComplete
                ? "cursor-not-allowed bg-slate-300"
                : "cursor-pointer bg-[#397e50] hover:bg-[#326f45]",
            ].join(" ")}
          >
            {submitLoading ? "Menyimpan..." : "Simpan Koreksi Essay"}
          </button>
        </div>
      )}
    </div>
  );
};

export default HasilJawabanUjianContent;
