import { useEffect, useState } from "react";

import InputField from "@/components/common/Input/InputField";

import type { MataPelajaranFormValues } from "@/types/DataMaster/MataPelajaran";
import type { TingkatKelas } from "@/types/DataMaster/Kelas";
import { getTingkatKelas } from "@/services/Api/features-api/DataMaster/kelas.service";
import { ApiError } from "@/services/Api/api";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";

type EditMapelFormProps = {
  initialValues: MataPelajaranFormValues;
  onSubmit: (values: MataPelajaranFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const EditMapelForm = ({
  initialValues,
  onSubmit,
  loading = false,
  submitting = false,
}: EditMapelFormProps) => {
  const [values, setValues] = useState<MataPelajaranFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [kelasOptions, setKelasOptions] = useState<TingkatKelas[]>([]);

  useEffect(() => {
    setValues(initialValues);
    setTouched({});
  }, [initialValues]);

  useEffect(() => {
    let active = true;
    const loadKelas = async () => {
      const data = await getTingkatKelas();
      if (!active) return;
      setKelasOptions(data);
    };
    loadKelas();
    return () => {
      active = false;
    };
  }, []);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof MataPelajaranFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<MataPelajaranFormValues>({
    kelasId: [requiredValue("Tingkat kelas wajib dipilih.")],
    kodeMapel: [requiredString("Kode mapel wajib diisi.")],
    namaMapel: [requiredString("Nama mata pelajaran wajib diisi.")],
    deskripsiMapel: [requiredString("Deskripsi mata pelajaran wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof MataPelajaranFormValues) =>
    !!errors[name] && !!touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      kelasId: true,
      kodeMapel: true,
      namaMapel: true,
      deskripsiMapel: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong.");
      return;
    }

    try {
      await onSubmit(values);
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      }
    }
  };

  const isDisabled = loading || submitting;

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Edit Data Mata Pelajaran
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Perbarui informasi mata pelajaran.
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Mata Pelajaran</h2>
              <p className={helperText}>
                Pilih tingkat kelas dan isi detail mata pelajaran.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label
                  htmlFor="kelasId"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="kelasId"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("kelasId") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelasId}
                  onChange={(e) =>
                    setField(
                      "kelasId",
                      e.target.value === "" ? "" : Number(e.target.value),
                    )
                  }
                  onBlur={() => onBlur("kelasId")}
                  disabled={isDisabled}
                  required
                >
                  <option value="">Pilih tingkat kelas</option>
                  {kelasOptions.map((tingkat) => (
                    <option
                      key={tingkat.id_tingkat_kelas}
                      value={tingkat.id_tingkat_kelas}
                    >
                      Kelas {tingkat.tingkat_kelas}
                    </option>
                  ))}
                </select>

                {hasError("kelasId") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.kelasId}</p>
                )}
              </div>

              <div>
                <InputField
                  id="kodeMapel"
                  label="Kode Mapel"
                  value={values.kodeMapel}
                  onChange={(v) => setField("kodeMapel", v)}
                  onBlur={() => onBlur("kodeMapel")}
                  placeholder="Contoh: MAT-10-01"
                  inputClassName={
                    hasError("kodeMapel") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("kodeMapel") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.kodeMapel}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="namaMapel"
                  label="Nama Mapel"
                  value={values.namaMapel}
                  onChange={(v) => setField("namaMapel", v)}
                  onBlur={() => onBlur("namaMapel")}
                  placeholder="Contoh: Matematika"
                  inputClassName={
                    hasError("namaMapel") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("namaMapel") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.namaMapel}
                  </p>
                )}
              </div>

              <div className="md:col-span-3">
                <InputField
                  id="deskripsiMapel"
                  label="Deskripsi"
                  value={values.deskripsiMapel}
                  onChange={(v) => setField("deskripsiMapel", v)}
                  onBlur={() => onBlur("deskripsiMapel")}
                  placeholder="Contoh: Aljabar dasar, geometri, dan statistika."
                  inputClassName={
                    hasError("deskripsiMapel")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("deskripsiMapel") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.deskripsiMapel}
                  </p>
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
              className="inline-flex items-center justify-center cursor-pointer rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({});
                setSubmitError(null);
              }}
              disabled={isDisabled}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
              disabled={isDisabled}
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditMapelForm;
