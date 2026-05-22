import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";

import InformasiUjianSection from "./BuatUjianForm/InformasiUjianSection";
import JadwalUjianSection from "./BuatUjianForm/JadwalUjianSection";
import KeamananTokenSection from "./BuatUjianForm/KeamananTokenSection";
import KelasBankSoalSection from "./BuatUjianForm/KelasBankSoalSection";
import PreviewSiswaSection from "./BuatUjianForm/PreviewSiswaSection";
import RuangSesiPengawasSection from "./BuatUjianForm/RuangSesiPengawasSection";
import { useBuatUjianForm } from "./BuatUjianForm/useBuatUjianForm";

type BuatUjianFormProps = {
  title?: string;
  description?: string;
  submitLabel?: string;
  onSubmit: (values: BuatUjianFormValues) => Promise<void>;
};

const BuatUjianForm = ({
  title,
  description,
  submitLabel,
  onSubmit,
}: BuatUjianFormProps) => {
  const form = useBuatUjianForm({ onSubmit });
  const pageTitle = title ?? "Buat Ujian";
  const pageDescription = description ?? "Silakan lengkapi data ujian";
  const actionLabel = submitLabel ?? "Simpan Ujian";

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-6xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">
            {pageTitle}
          </h1>
          <p className="mt-1 text-sm text-slate-500">{pageDescription}</p>
        </div>

        <form onSubmit={form.handleSubmit} className="space-y-6">
          <InformasiUjianSection
            errors={form.errors}
            hasError={form.hasError}
            onBlur={form.onBlur}
            setField={form.setField}
            values={form.values}
          />

          <KelasBankSoalSection
            bankSoalOptions={form.bankSoalOptions}
            bankSoalPlaceholder={form.bankSoalPlaceholder}
            errors={form.errors}
            hasError={form.hasError}
            loadingBankSoal={form.loadingBankSoal}
            loadingKelas={form.loadingKelas}
            loadingMapel={form.loadingMapel}
            mapelOptions={form.mapelOptions}
            namaKelasOptions={form.namaKelasOptions}
            onBlur={form.onBlur}
            selectedMapelId={form.selectedMapelId}
            setField={form.setField}
            setSelectedMapelId={form.setSelectedMapelId}
            submitting={form.submitting}
            tingkatKelasOptions={form.tingkatKelasOptions}
            values={form.values}
          />

          <JadwalUjianSection
            durasiMenit={form.durasiMenit}
            errors={form.errors}
            hasError={form.hasError}
            onBlur={form.onBlur}
            setField={form.setField}
            values={form.values}
          />

          <RuangSesiPengawasSection
            errors={form.errors}
            guruOptions={form.guruOptions}
            hasError={form.hasError}
            loadingGuru={form.loadingGuru}
            loadingRuang={form.loadingRuang}
            loadingSesi={form.loadingSesi}
            onBlur={form.onBlur}
            ruangOptions={form.ruangOptions}
            sesiOptions={form.sesiOptions}
            setField={form.setField}
            submitting={form.submitting}
            values={form.values}
          />

          <KeamananTokenSection
            errors={form.errors}
            hasError={form.hasError}
            onBlur={form.onBlur}
            setField={form.setField}
            submitting={form.submitting}
            values={form.values}
          />

          <PreviewSiswaSection
            loadingSiswa={form.loadingSiswa}
            siswaPreview={form.siswaPreview}
            siswaPreviewEnabled={form.siswaPreviewEnabled}
            values={form.values}
          />

          {(form.submitError || form.loadErrorMessage) && (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
              {form.submitError || form.loadErrorMessage}
            </div>
          )}

          <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-70"
              onClick={form.handleReset}
              disabled={form.submitting}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
              disabled={form.submitting}
            >
              {form.submitting ? "Menyimpan..." : actionLabel}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default BuatUjianForm;
