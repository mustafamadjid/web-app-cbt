import React, { useEffect, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import type { TingkatKelas } from "@/types/DataMaster/Kelas";
import { ApiError } from "@/services/Api/api";
import {
  GetDataKelasFull,
  submitKelasResponse,
} from "@/services/Api/features-api/DataMaster/kelas.service";
import { createValidator, requiredString, requiredValue } from "@/helper/validate/validateForm";

const validate = createValidator<{ tingkat_kelas: number | ""; nama_kelas: string }>({
  tingkat_kelas: [requiredValue("Tingkat kelas wajib dipilih.")],
  nama_kelas: [requiredString("Nama kelas wajib diisi.")],
});

const TambahNamaKelasForm = () => {
  const [opsiTingkat, setOpsiTingkat] = useState<TingkatKelas[]>([]);
  const [tingkatKelas, setTingkatKelas] = useState<number | "">("");
  const [namaKelas, setNamaKelas] = useState("");
  const [touched, setTouched] = useState({
    tingkat_kelas: false,
    nama_kelas: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  useEffect(() => {
    let akitf = true;
    const loadKelas = async () => {
      const data = await GetDataKelasFull();
      if (akitf) {
        setOpsiTingkat(data.item_tingkat_kelas ?? []);
      }
    };
    loadKelas();
    return () => {
      akitf = false;
    };
  }, []);

  const errors = validate({ tingkat_kelas: tingkatKelas, nama_kelas: namaKelas });
  const hasError = (name: "tingkat_kelas" | "nama_kelas") =>
    !!errors[name] && touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setTouched({ tingkat_kelas: true, nama_kelas: true });
    setSubmitError(null);

    const currentErrors = validate({ tingkat_kelas: tingkatKelas, nama_kelas: namaKelas });
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali data nama kelas yang diinput.");
      return;
    }

    try {
      await submitKelasResponse({
        tingkat_kelas: Number(tingkatKelas),
        nama_kelas: namaKelas,
      });
      alert("Nama kelas berhasil ditambahkan.");
      setTingkatKelas("");
      setNamaKelas("");
      setTouched({ tingkat_kelas: false, nama_kelas: false });
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      }
    }
  };

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div>
        <label
          htmlFor="opsi_tingkat_kelas"
          className="text-xs font-medium text-slate-600"
        >
          Tingkat Kelas
        </label>
        <select
          id="opsi_tingkat_kelas"
          className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] ${
            hasError("tingkat_kelas") ? "border-rose-300 ring-rose-100" : ""
          }`}
          value={tingkatKelas}
          onChange={(e) => {
            const value = e.target.value;
            setTingkatKelas(value === "" ? "" : Number(value));
          }}
          onBlur={() => setTouched((prev) => ({ ...prev, tingkat_kelas: true }))}
          required
        >
          <option value="">Pilih tingkat kelas</option>
          {opsiTingkat.map((item) => (
            <option key={item.id_tingkat_kelas} value={item.tingkat_kelas}>
              {item.tingkat_kelas}
            </option>
          ))}
        </select>
        {hasError("tingkat_kelas") && (
          <p className="mt-1 text-xs text-rose-500">{errors.tingkat_kelas}</p>
        )}
      </div>

      <div>
        <InputField
          id="nama_kelas"
          label="Nama Kelas"
          value={namaKelas}
          onChange={setNamaKelas}
          onBlur={() => setTouched((prev) => ({ ...prev, nama_kelas: true }))}
          placeholder="Contoh: X IPA 1"
          inputClassName={
            hasError("nama_kelas") ? "border-rose-300 ring-rose-100" : ""
          }
          required
        />
        {hasError("nama_kelas") && (
          <p className="mt-1 text-xs text-rose-500">{errors.nama_kelas}</p>
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
          Simpan Nama Kelas
        </button>
      </div>
    </form>
  );
};

export default TambahNamaKelasForm;
