import React, { useState } from "react";
import { useNavigate } from "react-router";

import InputField from "@/components/common/Input/InputField";
import { createSetField } from "@/helper/setField/setField";
import { createValidator, requiredString } from "@/helper/validate/validateForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { createSesi } from "@/services/Api/features-api/DataMaster/sesi.service";
import type { SesiFormValues } from "@/types/DataMaster/Sesi";

const initialValues: SesiFormValues = {
  kode_sesi: "",
  nama_sesi: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const DataSesiForm = () => {
  const navigate = useNavigate();
  const [values, setValues] = useState<SesiFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<keyof SesiFormValues, boolean>>(
    {
      kode_sesi: false,
      nama_sesi: false,
    },
  );
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof SesiFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const normalizeText = (text: string) => text.trim().replace(/\s+/g, " ");

  const validate = createValidator<SesiFormValues>({
    kode_sesi: [requiredString("Kode sesi wajib diisi.")],
    nama_sesi: [requiredString("Nama sesi wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof SesiFormValues) => !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setTouched({
      kode_sesi: true,
      nama_sesi: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong atau tidak valid.");
      return;
    }

    const payload: SesiFormValues = {
      kode_sesi: normalizeText(values.kode_sesi),
      nama_sesi: normalizeText(values.nama_sesi),
    };

    setSubmitting(true);
    try {
      await createSesi(payload);
      navigate(paths.dashboard.data_master_sesi);
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.message === "bad request: kode sesi already exist"
            ? "Kode sesi sudah ada."
            : "Data sesi gagal ditambahkan."
          : "Data sesi gagal ditambahkan.";
      setSubmitError(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">Data Sesi Ujian</h1>
            <p className="mt-1 text-sm text-slate-500">
              Lengkapi informasi sesi ujian yang akan ditambahkan.
            </p>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Sesi</h2>
              <p className={helperText}>Isi kode dan nama sesi ujian.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="kode_sesi"
                  label="Kode Sesi"
                  value={values.kode_sesi}
                  onChange={(v) => setField("kode_sesi", v)}
                  onBlur={() => onBlur("kode_sesi")}
                  placeholder="Contoh: SESI-01"
                  inputClassName={hasError("kode_sesi") ? "border-rose-300 ring-rose-100" : ""}
                  disabled={submitting}
                  required
                />
                {hasError("kode_sesi") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.kode_sesi}</p>
                )}
              </div>

              <div>
                <InputField
                  id="nama_sesi"
                  label="Nama Sesi"
                  value={values.nama_sesi}
                  onChange={(v) => setField("nama_sesi", v)}
                  onBlur={() => onBlur("nama_sesi")}
                  placeholder="Contoh: Sesi Pagi"
                  inputClassName={hasError("nama_sesi") ? "border-rose-300 ring-rose-100" : ""}
                  disabled={submitting}
                  required
                />
                {hasError("nama_sesi") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.nama_sesi}</p>
                )}
              </div>
            </div>
          </div>

          {submitError && (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
              {submitError}
            </div>
          )}

          <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({
                  kode_sesi: false,
                  nama_sesi: false,
                });
                setSubmitError(null);
              }}
              disabled={submitting}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
              disabled={submitting}
            >
              Simpan Sesi
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default DataSesiForm;
