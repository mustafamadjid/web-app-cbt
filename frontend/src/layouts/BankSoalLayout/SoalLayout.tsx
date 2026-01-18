import React from "react";

import type { SoalUjianItem } from "@/types/BankSoal/BankSoal";
import type { TipeUjian } from "@/types/Ujian/BuatUjian";

type SoalLayoutUser = {
  nama: string;
  kelas?: string;
  avatarUrl?: string;
};

type SoalLayoutProps = {
  title: string;
  currentNumber: number;
  totalSoal: number;
  sisaWaktu: string;
  soal: SoalUjianItem;
  user?: SoalLayoutUser;
  selectedOption?: string;
  onSelectOption?: (value: string) => void;
  onPrev?: () => void;
  onNext?: () => void;
  onBack?: () => void;
  className?: string;
};

const SoalLayout: React.FC<SoalLayoutProps> = ({
  title,
  currentNumber,
  totalSoal,
  sisaWaktu,
  soal,
  user,
  selectedOption,
  onSelectOption,
  onPrev,
  onNext,
  onBack,
  className = "",
}) => {
  const options = [
    { key: "opsi_a", label: "A", value: soal.opsi_a },
    { key: "opsi_b", label: "B", value: soal.opsi_b },
    { key: "opsi_c", label: "C", value: soal.opsi_c },
    { key: "opsi_d", label: "D", value: soal.opsi_d },
    { key: "opsi_e", label: "E", value: soal.opsi_e },
  ].filter((option) => Boolean(option.value));

  const tipeSoal = (label: TipeUjian) => {
    if (label === "PILIHAN_GANDA") {
      return "Pilihan Ganda";
    }
    if (label === "ESSAY") {
      return "Essay";
    }
    return "-";
  };

  const canBack = Boolean(onBack);
  const canPrev = Boolean(onPrev);
  const canNext = Boolean(onNext);

  return (
    <div className={["w-full min-h-screen bg-slate-50", className].join(" ")}>
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-4 py-6">
        <header className="flex flex-wrap items-center justify-between gap-4 rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-sm">
          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={onBack}
              disabled={!canBack}
              className={[
                "flex h-9 w-9 items-center justify-center rounded-full border border-slate-200 text-lg text-[#397e50] transition",
                canBack
                  ? "cursor-pointer hover:border-[#397e50]"
                  : "cursor-not-allowed opacity-50",
              ].join(" ")}
              aria-label="Kembali"
            >
              ←
            </button>
            <div>
              <p className="text-sm font-semibold text-slate-800">{title}</p>
              <p className="text-xs text-slate-400">Bank Soal</p>
            </div>
          </div>

          {user && (
            <div className="flex items-center gap-3 rounded-full border border-slate-200 bg-white px-3 py-2">
              <div className="text-right">
                <p className="text-xs font-semibold text-slate-700">
                  {user.nama}
                </p>
                {user.kelas && (
                  <p className="text-[11px] text-slate-400">{user.kelas}</p>
                )}
              </div>

              <div className="h-9 w-9 overflow-hidden rounded-full border border-slate-200">
                {user.avatarUrl ? (
                  <img
                    src={user.avatarUrl}
                    alt={user.nama}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center bg-slate-100 text-xs font-semibold text-slate-500">
                    {user.nama.slice(0, 2).toUpperCase()}
                  </div>
                )}
              </div>
            </div>
          )}
        </header>

        <section className="flex min-h-[640px] flex-col justify-between gap-6 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="space-y-6">
            <div className="flex flex-wrap items-center justify-between gap-4 border-b border-slate-100 pb-4">
              <p className="text-sm font-semibold text-slate-700">
                Soal {currentNumber} dari {totalSoal} Soal
              </p>
              <p className="text-sm font-semibold text-[#397e50]">
                Sisa Waktu : {sisaWaktu}
              </p>
            </div>

            <div className="space-y-4">
              <div>
                <p className="text-xs uppercase tracking-wide text-slate-400">
                  {tipeSoal(soal.tipe_soal)}
                </p>
                <p className="mt-2 text-sm leading-relaxed text-slate-600">
                  {soal.pertanyaan}
                </p>
              </div>

              {soal.urlGambar && (
                <div className="overflow-hidden rounded-xl border border-slate-200">
                  <img
                    src={soal.urlGambar}
                    alt="Ilustrasi soal"
                    className="h-52 w-full object-cover"
                  />
                </div>
              )}

              <div className="space-y-3">
                <p className="text-xs font-semibold text-slate-400">
                  Pilihan Jawaban
                </p>

                <div className="space-y-2">
                  {options.map((option) => {
                    const isSelected = selectedOption === option.label;
                    const isClickable = Boolean(option.value && onSelectOption);

                    return (
                      <button
                        key={option.key}
                        type="button"
                        onClick={() =>
                          option.value && onSelectOption?.(option.label)
                        }
                        disabled={!isClickable}
                        className={[
                          "flex w-full items-start gap-3 rounded-xl border px-4 py-3 text-left text-sm transition",
                          isSelected
                            ? "border-[#397e50] bg-[#397e50]/10 text-[#397e50]"
                            : "border-slate-200 text-slate-600 hover:border-[#397e50]",
                          isClickable
                            ? "cursor-pointer"
                            : "cursor-not-allowed opacity-50",
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
                        <span className="flex-1">{option.value}</span>
                      </button>
                    );
                  })}
                </div>
              </div>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 border-t border-slate-100 pt-4">
            <button
              type="button"
              onClick={onPrev}
              disabled={!canPrev}
              className={[
                "rounded-lg border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-500 transition",
                canPrev
                  ? "cursor-pointer hover:border-[#397e50] hover:text-[#397e50]"
                  : "cursor-not-allowed opacity-50",
              ].join(" ")}
            >
              Sebelumnya
            </button>

            <button
              type="button"
              onClick={onNext}
              disabled={!canNext}
              className={[
                "rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white transition hover:bg-[#326f45]",
                canNext ? "cursor-pointer" : "cursor-not-allowed opacity-50",
              ].join(" ")}
            >
              Selanjutnya
            </button>
          </div>
        </section>
      </div>
    </div>
  );
};

export default SoalLayout;
