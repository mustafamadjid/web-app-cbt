import React, { useState } from "react";
import { useNavigate } from "react-router";

import InputField from "@/components/common/Input/InputField";
import { createSetField } from "@/helper/setField/setField";
import { createValidator, requiredString } from "@/helper/validate/validateForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { createRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import type { RuangUjianFormValues } from "@/types/DataMaster/RuangUjian";

const initialValues: RuangUjianFormValues = {
  kode_ruang: "",
  nama_ruangan: "",
};

const DataRuangForm = () => {
  const navigate = useNavigate();
  const [values, setValues] = useState<RuangUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<keyof RuangUjianFormValues, boolean>>({
    kode_ruang: false,
    nama_ruangan: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const setField = createSetField(setValues);

  const normalize = (s: string) => s.trim().replace(/\s+/g, " ");

  const validate = createValidator<RuangUjianFormValues>({
    kode_ruang: [requiredString("Kode ruang wajib diisi.")],
    nama_ruangan: [requiredString("Nama ruangan wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof RuangUjianFormValues) => !!errors[name] && touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setTouched({ kode_ruang: true, nama_ruangan: true });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong atau tidak valid.");
      return;
    }

    const payload: RuangUjianFormValues = {
      kode_ruang: normalize(values.kode_ruang),
      nama_ruangan: normalize(values.nama_ruangan),
    };

    setSubmitting(true);
    try {
      await createRuangUjian(payload);
      navigate(paths.dashboard.data_master_ruang);
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.message === "bad request: kode ruang ujian already exist"
            ? "Kode ruang ujian sudah ada."
            : "Ruang ujian gagal ditambahkan."
          : "Ruang ujian gagal ditambahkan.";
      setSubmitError(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={onSubmit} className="space-y-6">
      <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <InputField
            id="kode_ruang"
            label="Kode Ruang"
            value={values.kode_ruang}
            onChange={(v) => setField("kode_ruang", v)}
            onBlur={() => setTouched((p) => ({ ...p, kode_ruang: true }))}
            placeholder="Contoh: RU-01"
            inputClassName={hasError("kode_ruang") ? "border-rose-300 ring-rose-100" : ""}
            required
          />
          <InputField
            id="nama_ruangan"
            label="Nama Ruangan"
            value={values.nama_ruangan}
            onChange={(v) => setField("nama_ruangan", v)}
            onBlur={() => setTouched((p) => ({ ...p, nama_ruangan: true }))}
            placeholder="Contoh: Lab Komputer"
            inputClassName={hasError("nama_ruangan") ? "border-rose-300 ring-rose-100" : ""}
            required
          />
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
            setTouched({ kode_ruang: false, nama_ruangan: false });
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
          Simpan Ruangan Ujian
        </button>
      </div>
    </form>
  );
};

export default DataRuangForm;
