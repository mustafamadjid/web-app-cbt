import { useMemo, useState } from "react";
import { useNavigate } from "react-router";
import toast from "react-hot-toast";

import InputField from "@/components/common/Input/InputField";
import { paths } from "@/routes/paths";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import { createBankSoal } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import { useGetMataPelajaranOptions } from "@/services/Api/features-api/GetOptions/options.service";
import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import { useAuth } from "@/contexts/AuthContext";

import type { BankSoalFormValues } from "@/types/BankSoal/BankSoal";
import type { TingkatKelas } from "@/types/DataMaster/Kelas";
import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";

const initialValues: BankSoalFormValues = {
  namaBankSoal: "",
  kelasId: "",
  mapelId: "",
  materi: "",
  deskripsi: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";
const selectClassName =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";

const BankSoalForm = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [values, setValues] = useState<BankSoalFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const {
    data: kelasData,
    loading: loadingKelasOptions,
    error: kelasOptionsError,
  } = useGetDataKelasFull();
  const {
    data: mapelData,
    loading: loadingMapelOptions,
    error: mapelOptionsError,
  } = useGetMataPelajaranOptions({ source: "dataMaster" });

  const kelasOptions: TingkatKelas[] = useMemo(
    () => kelasData?.item_tingkat_kelas ?? [],
    [kelasData],
  );
  const mapelOptions: MataPelajaranOption[] = useMemo(
    () => mapelData ?? [],
    [mapelData],
  );
  const loadingOptions = loadingKelasOptions || loadingMapelOptions;
  const optionsErrorMsg = kelasOptionsError || mapelOptionsError;

  const setField = createSetField(setValues);

  const selectedKelas = useMemo(
    () =>
      kelasOptions.find((item) => item.id_tingkat_kelas === values.kelasId),
    [kelasOptions, values.kelasId],
  );

  const selectedMapel = useMemo(
    () => mapelOptions.find((item) => item.id === values.mapelId),
    [mapelOptions, values.mapelId],
  );

  const onBlur = (name: keyof BankSoalFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<BankSoalFormValues>({
    namaBankSoal: [requiredString("Nama bank soal wajib diisi.")],
    kelasId: [requiredValue("Kelas wajib dipilih.")],
    mapelId: [requiredValue("Mata pelajaran wajib dipilih.")],
    materi: [requiredString("Materi wajib diisi.")],
    deskripsi: [requiredString("Deskripsi wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof BankSoalFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setTouched({
      namaBankSoal: true,
      kelasId: true,
      mapelId: true,
      materi: true,
      deskripsi: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali field yang belum valid.");
      return;
    }

    if (!selectedKelas || !selectedMapel) {
      setSubmitError("Pilihan kelas atau mata pelajaran tidak ditemukan.");
      return;
    }

    setSubmitting(true);
    try {
      await createBankSoal({
        id_pengguna: user?.id_pengguna,
        nama_bank_soal: values.namaBankSoal.trim(),
        id_kelas: selectedKelas.id_tingkat_kelas,
        id_mapel: selectedMapel.id,
        deskripsi: values.deskripsi.trim(),
        materi: values.materi.trim(),
      });

      toast.success("Bank soal berhasil dibuat.");
      if (user?.role === "GURU") {
        navigate(paths.dashboard.bank_soal_guru);
      } else {
        navigate(paths.dashboard.bank_soal);
      }
    } catch (error) {
      setSubmitError(
        getUserFriendlyErrorMessage(error, {
          action: "create",
          entity: "bank soal",
          fallbackMessage: "Terjadi kesalahan saat menyimpan bank soal.",
        }),
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">
            Buat Bank Soal
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Lengkapi informasi bank soal sebelum mengunggah dokumen soal.
          </p>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Bank Soal</h2>
              <p className={helperText}>
                Pilih tingkat kelas, mata pelajaran, lalu isi deskripsi singkat bank
                soal beserta materi yang dicakup.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="md:col-span-2">
                <InputField
                  id="namaBankSoal"
                  label="Nama Bank Soal"
                  value={values.namaBankSoal}
                  onChange={(v) => setField("namaBankSoal", v)}
                  onBlur={() => onBlur("namaBankSoal")}
                  placeholder="Contoh: Bank Soal Ulangan Harian Matematika"
                  inputClassName={
                    hasError("namaBankSoal")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("namaBankSoal") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.namaBankSoal}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="kelasId"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="kelasId"
                  className={`${selectClassName} ${
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
                  disabled={loadingOptions || submitting}
                  required
                >
                  <option value="">
                    {loadingOptions
                      ? "Memuat tingkat kelas..."
                      : "Pilih tingkat kelas"}
                  </option>
                  {kelasOptions.map((kelas) => (
                    <option
                      key={kelas.id_tingkat_kelas}
                      value={kelas.id_tingkat_kelas}
                    >
                      Kelas {kelas.tingkat_kelas}
                    </option>
                  ))}
                </select>
                {hasError("kelasId") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.kelasId}</p>
                )}
              </div>

              <div>
                <label
                  htmlFor="mapelId"
                  className="text-xs font-medium text-slate-600"
                >
                  Mapel Bank Soal
                </label>
                <select
                  id="mapelId"
                  className={`${selectClassName} ${
                    hasError("mapelId") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.mapelId}
                  onChange={(e) =>
                    setField(
                      "mapelId",
                      e.target.value === "" ? "" : Number(e.target.value),
                    )
                  }
                  onBlur={() => onBlur("mapelId")}
                  disabled={loadingOptions || submitting}
                  required
                >
                  <option value="">
                    {loadingOptions
                      ? "Memuat mata pelajaran..."
                      : "Pilih mata pelajaran"}
                  </option>
                  {mapelOptions.map((mapel) => (
                    <option key={mapel.id} value={mapel.id}>
                      {mapel.label}
                    </option>
                  ))}
                </select>
                {hasError("mapelId") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.mapelId}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <InputField
                  id="materi"
                  label="Materi"
                  value={values.materi}
                  onChange={(v) => setField("materi", v)}
                  onBlur={() => onBlur("materi")}
                  placeholder="Contoh: Persamaan Linear Dua Variabel"
                  inputClassName={
                    hasError("materi") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={submitting}
                  required
                />
                {hasError("materi") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.materi}</p>
                )}
              </div>
            </div>

            <div className="mt-4">
              <label
                htmlFor="deskripsi"
                className="text-xs font-medium text-slate-600"
              >
                Deskripsi
              </label>
              <textarea
                id="deskripsi"
                rows={4}
                className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                  hasError("deskripsi") ? "border-rose-300 ring-rose-100" : ""
                }`}
                placeholder="Tuliskan ringkasan isi bank soal."
                value={values.deskripsi}
                onChange={(e) => setField("deskripsi", e.target.value)}
                onBlur={() => onBlur("deskripsi")}
                disabled={submitting}
                required
              />
              {hasError("deskripsi") && (
                <p className="mt-1 text-xs text-rose-500">
                  {errors.deskripsi}
                </p>
              )}
            </div>
          </div>

          {(submitError || optionsErrorMsg) && (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
              {submitError || optionsErrorMsg}
            </div>
          )}

          <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              className="inline-flex items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({});
                setSubmitError(null);
              }}
              disabled={submitting}
            >
              Reset
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="inline-flex items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:bg-[#397e50]/70"
            >
              {submitting ? "Menyimpan..." : "Simpan Bank Soal"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default BankSoalForm;
