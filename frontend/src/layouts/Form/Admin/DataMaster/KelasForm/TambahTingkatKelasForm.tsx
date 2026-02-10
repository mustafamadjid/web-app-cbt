import React, { useState } from "react";

import { ApiError } from "@/services/Api/api";
import { createTingkatKelas } from "@/services/Api/features-api/DataMaster/kelas.service";
import { createValidator, integerNumber, minNumber, requiredValue } from "@/helper/validate/validateForm";
import toast from "react-hot-toast";
import { useNavigate } from "react-router";
import { paths } from "@/routes/paths";

const validate = createValidator<{ tingkat_kelas: number | "" }>({
  tingkat_kelas: [
    requiredValue("Tingkat kelas wajib diisi."),
    integerNumber("Tingkat kelas harus bilangan bulat."),
    minNumber(1, "Tingkat kelas tidak valid."),
  ],
});

const TambahTingkatKelasForm = () => {
  const [tingkatKelas, setTingkatKelas] = useState<number | "">("");
  const [touched, setTouched] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const errors = validate({ tingkat_kelas: tingkatKelas });
  const hasError = touched && !!errors.tingkat_kelas;

  const navigate = useNavigate();

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setTouched(true);
    setSubmitError(null);

    const currentErrors = validate({ tingkat_kelas: tingkatKelas });
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali tingkat kelas yang diinput.");
      return;
    }

    try {
      await createTingkatKelas({
        tingkat_kelas: Number(tingkatKelas),
      });
      toast.success("Tingkat kelas berhasil ditambahkan.");
      setTingkatKelas("");
      setTouched(false);
      setTimeout(() => {
        navigate(`${paths.dashboard.data_master_kelas}`);
      })
    } catch (error) {
      const message =
        e instanceof ApiError
          ? e.message === "bad request: tingkat kelas already exist"
            ? "Tingkat kelas sudah ada."
            : "Tingkat kelas gagal ditambahkan"
          : "Tingkat kelas gagal ditambahkan";
      setSubmitError(message);
    }
  };

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div>
        <label
          htmlFor="tingkat_kelas"
          className="text-xs font-medium text-slate-600"
        >
          Tingkat Kelas
        </label>
        <input
          id="tingkat_kelas"
          type="number"
          inputMode="numeric"
          className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] ${
            hasError ? "border-rose-300 ring-rose-100" : ""
          }`}
          placeholder="Contoh: 10"
          min={1}
          step={1}
          value={tingkatKelas}
          onChange={(e) => {
            const raw = e.target.value;
            setTingkatKelas(raw === "" ? "" : Number(raw));
          }}
          onBlur={() => setTouched(true)}
          required
        />

        {hasError && (
          <p className="mt-1 text-xs text-rose-500">{errors.tingkat_kelas}</p>
        )}
      </div>

      {submitError && (
        <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
          {submitError}
        </div>
      )}

      <div className="flex justify-end">
        <button
          type="submit"
          className="inline-flex items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43]"
        >
          Simpan Tingkat Kelas
        </button>
      </div>
    </form>
  );
};

export default TambahTingkatKelasForm;
