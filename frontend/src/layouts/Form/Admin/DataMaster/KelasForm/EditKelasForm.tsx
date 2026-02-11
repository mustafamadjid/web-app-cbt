import React, { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import type { KelasFormValues } from "@/types/DataMaster/Kelas";
import { ApiError } from "@/services/Api/api";
import {
  getNamaKelasById,
  getTingkatKelasByIdRequest,
  updateKelas,
} from "@/services/Api/features-api/DataMaster/kelas.service";
import { paths } from "@/routes/paths";
import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  integerNumber,
  minNumber,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";

import EditNamaKelasForm from "./EditNamaKelasForm";
import EditTingkatKelasForm from "./EditTingkatKelasForm";

const buildInitialValues = (): KelasFormValues => ({
  tingkat_kelas: "",
  nama_kelas: "",
});

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const validate = createValidator<KelasFormValues>({
  tingkat_kelas: [
    requiredValue("Tingkat kelas wajib diisi."),
    integerNumber("Tingkat kelas harus bilangan bulat."),
    minNumber(1, "Tingkat kelas tidak valid."),
  ],
  nama_kelas: [requiredString("Nama kelas wajib diisi.")],
});

const EditKelasForm = () => {
  const { idTingkatKelas, idNamaKelas } = useParams();
  const navigate = useNavigate();

  const [initialValues, setInitialValues] = useState<KelasFormValues>(buildInitialValues());
  const [values, setValues] = useState<KelasFormValues>(buildInitialValues());
  const [touched, setTouched] = useState<Record<keyof KelasFormValues, boolean>>({
    tingkat_kelas: false,
    nama_kelas: false,
  });
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const tingkatKelasId = useMemo(() => Number(idTingkatKelas), [idTingkatKelas]);
  const namaKelasId = useMemo(() => Number(idNamaKelas), [idNamaKelas]);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof KelasFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  useEffect(() => {
    let active = true;

    const loadKelasByIds = async () => {
      if (
        !idTingkatKelas ||
        !idNamaKelas ||
        Number.isNaN(tingkatKelasId) ||
        Number.isNaN(namaKelasId)
      ) {
        setSubmitError("Parameter edit kelas tidak valid.");
        setLoading(false);
        return;
      }

      try {
        const [tingkatKelas, namaKelas] = await Promise.all([
          getTingkatKelasByIdRequest(tingkatKelasId),
          getNamaKelasById(namaKelasId),
        ]);

        if (!active) return;

        if (!tingkatKelas || !namaKelas) {
          setSubmitError("Data kelas tidak ditemukan.");
          return;
        }

        const nextValues = {
          tingkat_kelas: tingkatKelas.tingkat_kelas,
          nama_kelas: namaKelas.nama_kelas,
        } satisfies KelasFormValues;

        setInitialValues(nextValues);
        setValues(nextValues);
      } catch (error) {
        if (!active) return;
        if (error instanceof ApiError) {
          setSubmitError(error.message);
          return;
        }
        setSubmitError("Gagal memuat data kelas.");
      } finally {
        if (active) setLoading(false);
      }
    };

    loadKelasByIds();

    return () => {
      active = false;
    };
  }, [idTingkatKelas, idNamaKelas, tingkatKelasId, namaKelasId]);

  const errors = validate(values);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      tingkat_kelas: true,
      nama_kelas: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong atau tidak valid.");
      return;
    }

    if (!idNamaKelas || Number.isNaN(namaKelasId)) {
      setSubmitError("ID nama kelas tidak ditemukan.");
      return;
    }

    setSubmitting(true);
    try {
      await updateKelas(namaKelasId, values);
      navigate(paths.dashboard.data_master_kelas);
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const isDisabled = loading || submitting;

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4 space-y-6">
        <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">Edit Data Kelas</h1>
          <p className="mt-1 text-sm text-slate-500">
            Perbarui tingkat kelas dan nama kelas pada section yang tersedia.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Edit Tingkat Kelas</h2>
              <p className={helperText}>Perbarui tingkat kelas dalam format angka.</p>
            </div>
            <EditTingkatKelasForm
              value={values.tingkat_kelas}
              error={touched.tingkat_kelas ? errors.tingkat_kelas : undefined}
              onChange={(value) => setField("tingkat_kelas", value)}
              onBlur={() => onBlur("tingkat_kelas")}
              disabled={isDisabled}
            />
          </section>

          <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Edit Nama Kelas</h2>
              <p className={helperText}>Perbarui nama kelas sesuai tingkat kelas yang dipilih.</p>
            </div>
            <EditNamaKelasForm
              value={values.nama_kelas}
              error={touched.nama_kelas ? errors.nama_kelas : undefined}
              onChange={(value) => setField("nama_kelas", value)}
              onBlur={() => onBlur("nama_kelas")}
              disabled={isDisabled}
            />
          </section>

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
                setTouched({ tingkat_kelas: false, nama_kelas: false });
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

export default EditKelasForm;
