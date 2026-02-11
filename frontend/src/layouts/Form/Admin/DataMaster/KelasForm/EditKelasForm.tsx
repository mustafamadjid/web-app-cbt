import React, { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import type { KelasFormValues } from "@/types/DataMaster/Kelas";
import { ApiError } from "@/services/Api/api";
import {
  GetDataKelasFull,
  getKelasByIdsRequest,
  updateNamaKelasPartial,
} from "@/services/Api/features-api/DataMaster/kelas.service";
import { paths } from "@/routes/paths";
import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  integerNumber,
  minNumber,
  requiredValue,
  requiredString,
} from "@/helper/validate/validateForm";

import EditNamaKelasForm from "./EditNamaKelasForm";
import EditTingkatKelasForm from "./EditTingkatKelasForm";

const buildInitialValues = (): KelasFormValues => ({
  tingkat_kelas: "",
  nama_kelas: "",
});

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const validateTingkatKelas = createValidator<Pick<KelasFormValues, "tingkat_kelas">>({
  tingkat_kelas: [
    requiredValue("Tingkat kelas wajib diisi."),
    integerNumber("Tingkat kelas harus bilangan bulat."),
    minNumber(1, "Tingkat kelas tidak valid."),
  ],
});

const validateNamaKelas = createValidator<Pick<KelasFormValues, "nama_kelas">>({
  nama_kelas: [requiredString("Nama kelas wajib diisi.")],
});

const EditKelasForm = () => {
  const { idTingkatKelas, idNamaKelas } = useParams();
  const navigate = useNavigate();

  const [initialValues, setInitialValues] =
    useState<KelasFormValues>(buildInitialValues());
  const [values, setValues] = useState<KelasFormValues>(buildInitialValues());
  const [currentIdTingkatKelas, setCurrentIdTingkatKelas] = useState<number | null>(null);
  const [touchedTingkatKelas, setTouchedTingkatKelas] = useState<boolean>(false);
  const [touchedNamaKelas, setTouchedNamaKelas] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [pageError, setPageError] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const tingkatKelasId = useMemo(
    () => Number(idTingkatKelas),
    [idTingkatKelas],
  );
  const namaKelasId = useMemo(() => Number(idNamaKelas), [idNamaKelas]);

  const setField = createSetField(setValues);

  useEffect(() => {
    let active = true;

    const loadKelasByIds = async () => {
      if (
        !idTingkatKelas ||
        !idNamaKelas ||
        Number.isNaN(tingkatKelasId) ||
        Number.isNaN(namaKelasId)
      ) {
        setPageError("Parameter edit kelas tidak valid.");
        setLoading(false);
        return;
      }

      try {
        const dataKelas = await getKelasByIdsRequest(tingkatKelasId, namaKelasId);

        if (!active) return;

        if (!dataKelas) {
          setPageError("Data kelas tidak ditemukan.");
          return;
        }

        const nextValues = {
          tingkat_kelas: dataKelas.item_tingkat_kelas.tingkat_kelas,
          nama_kelas: dataKelas.item_nama_kelas.nama_kelas,
        } satisfies KelasFormValues;

        setCurrentIdTingkatKelas(dataKelas.item_tingkat_kelas.id_tingkat_kelas);
        setInitialValues(nextValues);
        setValues(nextValues);
      } catch (error) {
        if (!active) return;
        if (error instanceof ApiError) {
          setPageError(error.message);
          return;
        }
        setPageError("Gagal memuat data kelas.");
      } finally {
        if (active) setLoading(false);
      }
    };

    loadKelasByIds();

    return () => {
      active = false;
    };
  }, [idTingkatKelas, idNamaKelas, tingkatKelasId, namaKelasId]);

  const tingkatKelasErrors = validateTingkatKelas({
    tingkat_kelas: values.tingkat_kelas,
  });
  const namaKelasErrors = validateNamaKelas({
    nama_kelas: values.nama_kelas,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setTouchedTingkatKelas(true);
    setTouchedNamaKelas(true);

    if (tingkatKelasErrors.tingkat_kelas || namaKelasErrors.nama_kelas) {
      setSubmitError("Periksa kembali data kelas yang diisi.");
      return;
    }

    if (!idNamaKelas || Number.isNaN(namaKelasId)) {
      setSubmitError("ID nama kelas tidak ditemukan.");
      return;
    }

    setSubmitting(true);
    try {
      const payload: { id_tingkat_kelas?: number; nama_kelas?: string } = {};

      if (values.nama_kelas !== initialValues.nama_kelas) {
        payload.nama_kelas = values.nama_kelas;
      }

      if (values.tingkat_kelas !== initialValues.tingkat_kelas) {
        const dataKelas = await GetDataKelasFull();
        const matchedTingkatKelas = dataKelas.item_tingkat_kelas.find(
          (item) => item.tingkat_kelas === values.tingkat_kelas,
        );

        if (!matchedTingkatKelas) {
          setSubmitError("Tingkat kelas tidak ditemukan.");
          return;
        }

        if (matchedTingkatKelas.id_tingkat_kelas !== currentIdTingkatKelas) {
          payload.id_tingkat_kelas = matchedTingkatKelas.id_tingkat_kelas;
        }
      }

      if (Object.keys(payload).length === 0) {
        setSubmitError("Tidak ada perubahan data untuk disimpan.");
        return;
      }

      await updateNamaKelasPartial(namaKelasId, payload);
      setInitialValues(values);
      if (payload.id_tingkat_kelas) {
        setCurrentIdTingkatKelas(payload.id_tingkat_kelas);
      }
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
          <h1 className="text-base font-semibold text-slate-900">
            Edit Data Kelas
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Ubah tingkat kelas dan nama kelas, lalu simpan keduanya sekaligus.
          </p>
        </div>

        {pageError && (
          <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
            {pageError}
          </div>
        )}

        <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <form onSubmit={handleSubmit} className="space-y-8">
            <div className="space-y-4">
              <div className="mb-4">
                <h2 className={sectionTitle}>Edit Tingkat Kelas</h2>
                <p className={helperText}>
                  Perbarui tingkat kelas dalam format angka.
                </p>
              </div>
              <EditTingkatKelasForm
                value={values.tingkat_kelas}
                error={
                  touchedTingkatKelas
                    ? tingkatKelasErrors.tingkat_kelas
                    : undefined
                }
                onChange={(value) => setField("tingkat_kelas", value)}
                onBlur={() => setTouchedTingkatKelas(true)}
                disabled={isDisabled}
              />
            </div>

            <div className="space-y-4">
              <div className="mb-4">
                <h2 className={sectionTitle}>Edit Nama Kelas</h2>
                <p className={helperText}>
                  Perbarui nama kelas sesuai tingkat kelas yang dipilih.
                </p>
              </div>
              <EditNamaKelasForm
                value={values.nama_kelas}
                error={touchedNamaKelas ? namaKelasErrors.nama_kelas : undefined}
                onChange={(value) => setField("nama_kelas", value)}
                onBlur={() => setTouchedNamaKelas(true)}
                disabled={isDisabled}
              />
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
                  setTouchedTingkatKelas(false);
                  setTouchedNamaKelas(false);
                  setSubmitError(null);
                }}
                disabled={isDisabled}
              >
                Reset Form
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
        </section>

        <div className="flex justify-end">
          <button
            type="button"
            className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
            onClick={() => navigate(paths.dashboard.data_master_kelas)}
            disabled={isDisabled}
          >
            Kembali ke daftar kelas
          </button>
        </div>
      </div>
    </div>
  );
};

export default EditKelasForm;
