import React, { useEffect, useRef, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { StudentRegisterFormValues } from "@/types/KelolaAkun/AkunSiswa";
import type { JenisKelamin, StatusAkun } from "@/types/OpsiTypes/Option";
import type { NamaKelas, TingkatKelas } from "@/types/DataMaster/Kelas";

import { getNamaKelas, getTingkatKelas } from "@/services/Api/features-api/DataMaster/kelas.service";
import { ApiError } from "@/services/Api/api";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  fileMaxSize,
  fileTypeStartsWith,
  minLength,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";

type EditAkunSiswaFormProps = {
  initialValues: StudentRegisterFormValues;
  initialFotoUrl?: string;
  onSubmit: (values: StudentRegisterFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

function toIntOrPrev(prev: number, raw: string): number {
  const s = String(raw ?? "").trim();
  if (s === "") return 0;
  if (!/^\d+$/.test(s)) return prev;
  const n = Number(s);
  return Number.isSafeInteger(n) ? n : prev;
}

const EditAkunSiswaForm = ({
  initialValues,
  initialFotoUrl,
  onSubmit,
  loading = false,
  submitting = false,
}: EditAkunSiswaFormProps) => {
  const [values, setValues] = useState<StudentRegisterFormValues>(initialValues);
  const [touched, setTouched] = useState<
    Partial<Record<keyof StudentRegisterFormValues, boolean>>
  >({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const [fotoUrl, setFotoUrl] = useState<string>(initialFotoUrl ?? "");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [tingkatKelass, setTingkatKelass] = useState<TingkatKelas[]>([]);
  const [daftarKelas, setDaftarKelas] = useState<NamaKelas[]>([]);
  const [namaKelasOptions, setNamaKelasOptions] = useState<NamaKelas[]>([]);

  useEffect(() => {
    setValues(initialValues);
    setTouched({});
  }, [initialValues]);

  useEffect(() => {
    if (values.fotoProfil) {
      const url = URL.createObjectURL(values.fotoProfil);
      setFotoUrl(url);
      return () => URL.revokeObjectURL(url);
    }

    setFotoUrl(initialFotoUrl ?? "");
    return undefined;
  }, [initialFotoUrl, values.fotoProfil]);

  useEffect(() => {
    let active = true;
    const loadKelas = async () => {
      const [tingkat, kelas] = await Promise.all([
        getTingkatKelas(),
        getNamaKelas(),
      ]);
      if (!active) return;
      setTingkatKelass(tingkat);
      setDaftarKelas(kelas);
    };
    loadKelas();
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    if (!values.id_tingkat_kelas) {
      setNamaKelasOptions([]);
      if (values.id_nama_kelas !== "") {
        setValues((prev) => ({ ...prev, id_nama_kelas: "" }));
      }
      return;
    }

    const filtered = daftarKelas.filter(
      (kelas) => kelas.id_tingkat_kelas === values.id_tingkat_kelas,
    );
    setNamaKelasOptions(filtered);

    const isValid = filtered.some(
      (kelas) => String(kelas.id_nama_kelas) === values.id_nama_kelas,
    );
    if (!isValid && values.id_nama_kelas !== "") {
      setValues((prev) => ({ ...prev, id_nama_kelas: "" }));
    }
  }, [daftarKelas, values.id_nama_kelas, values.id_tingkat_kelas]);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof StudentRegisterFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<StudentRegisterFormValues>({
    namaLengkap: [requiredString("Nama lengkap wajib diisi.")],
    username: [requiredString("Username wajib diisi.")],
    password: [minLength(8, "Password minimal 8 karakter.")],
    jenisKelamin: [requiredValue("Jenis kelamin wajib dipilih.")],
    email: [
      (value) => {
        if (!value || !value.trim()) return null;
        return /^\S+@\S+\.\S+$/.test(value.trim())
          ? null
          : "Format email tidak valid.";
      },
    ],
    noHp: [
      (value) => {
        if (!value || !value.trim()) return null;
        return /^\d{10,15}$/.test(value.trim())
          ? null
          : "No HP harus berupa 10-15 digit angka.";
      },
    ],
    noAbsen: [
      (value) => (value <= 0 ? "No absen harus lebih dari 0." : null),
    ],
    angkatan: [
      (value) => (value <= 0 ? "Angkatan wajib diisi." : null),
    ],
    tempatLahir: [requiredString("Tempat lahir wajib diisi.")],
    tanggalLahir: [requiredString("Tanggal lahir wajib diisi.")],
    id_tingkat_kelas: [requiredValue("Tingkat kelas wajib dipilih.")],
    id_nama_kelas: [requiredValue("Nama kelas wajib dipilih.")],
    fotoProfil: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran foto maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
    statusAkun: [requiredValue("Status akun wajib dipilih.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof StudentRegisterFormValues) =>
    !!errors[name] && !!touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      namaLengkap: true,
      username: true,
      password: true,
      jenisKelamin: true,
      email: true,
      noHp: true,
      noAbsen: true,
      angkatan: true,
      tempatLahir: true,
      tanggalLahir: true,
      id_tingkat_kelas: true,
      id_nama_kelas: true,
      fotoProfil: true,
      statusAkun: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih salah / kosong.");
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

  const clearFoto = () => {
    setField("fotoProfil", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const isDisabled = loading || submitting;

  return (
    <div className="min-h-screen w-full bg-slate-50 py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Edit Akun Siswa
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Perbarui data akun siswa sesuai kebutuhan.
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Dasar</h2>
              <p className={helperText}>Data akun dan identitas siswa.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="namaLengkap"
                  label="Nama Lengkap"
                  value={values.namaLengkap}
                  onChange={(v) => setField("namaLengkap", v)}
                  onBlur={() => onBlur("namaLengkap")}
                  placeholder="Contoh: Siti Aminah"
                  inputClassName={
                    hasError("namaLengkap")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("namaLengkap") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.namaLengkap}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="username"
                  label="Username"
                  value={values.username}
                  onChange={(v) => setField("username", v)}
                  onBlur={() => onBlur("username")}
                  placeholder="Contoh: sita.aminah"
                  inputClassName={
                    hasError("username")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("username") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.username}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="password"
                  type="password"
                  label="Password Baru"
                  value={values.password}
                  onChange={(v) => setField("password", v)}
                  onBlur={() => onBlur("password")}
                  placeholder="Biarkan kosong jika tidak diganti"
                  inputClassName={
                    hasError("password")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required={false}
                />
                {hasError("password") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.password}
                  </p>
                )}
              </div>

              <div>
                <label className="text-xs font-medium text-slate-600" htmlFor="jenisKelamin">
                  Jenis Kelamin
                </label>
                <select
                  id="jenisKelamin"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("jenisKelamin")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.jenisKelamin}
                  onChange={(e) =>
                    setField("jenisKelamin", e.target.value as JenisKelamin)
                  }
                  onBlur={() => onBlur("jenisKelamin")}
                  disabled={isDisabled}
                >
                  <option value="LAKI_LAKI">Laki-laki</option>
                  <option value="PEREMPUAN">Perempuan</option>
                </select>
                {hasError("jenisKelamin") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.jenisKelamin}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="email"
                  type="email"
                  label="Email"
                  value={values.email ?? ""}
                  onChange={(v) => setField("email", v)}
                  onBlur={() => onBlur("email")}
                  placeholder="nama@email.com"
                  inputClassName={
                    hasError("email") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required={false}
                />
                {hasError("email") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.email}</p>
                )}
              </div>

              <div>
                <InputField
                  id="noHp"
                  type="tel"
                  label="Nomor HP"
                  value={values.noHp ?? ""}
                  onChange={(v) => setField("noHp", v)}
                  onBlur={() => onBlur("noHp")}
                  placeholder="Contoh: 081234567890"
                  inputClassName={
                    hasError("noHp") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required={false}
                />
                {hasError("noHp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noHp}</p>
                )}
              </div>

              <div>
                <label className="text-xs font-medium text-slate-600" htmlFor="statusAkun">
                  Status Akun
                </label>
                <select
                  id="statusAkun"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("statusAkun")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.statusAkun}
                  onChange={(e) =>
                    setField("statusAkun", e.target.value as StatusAkun)
                  }
                  onBlur={() => onBlur("statusAkun")}
                  disabled={isDisabled}
                >
                  <option value="aktif">Aktif</option>
                  <option value="nonaktif">Nonaktif</option>
                  <option value="dibekukan">Dibekukan</option>
                </select>
                {hasError("statusAkun") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.statusAkun}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Data Akademik</h2>
              <p className={helperText}>Identitas akademik siswa.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="noAbsen"
                  type="number"
                  label="Nomor Absen"
                  value={String(values.noAbsen)}
                  onChange={(v) =>
                    setField("noAbsen", toIntOrPrev(values.noAbsen, v))
                  }
                  onBlur={() => onBlur("noAbsen")}
                  inputClassName={
                    hasError("noAbsen") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("noAbsen") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noAbsen}</p>
                )}
              </div>

              <div>
                <InputField
                  id="angkatan"
                  type="number"
                  label="Angkatan"
                  value={String(values.angkatan)}
                  onChange={(v) =>
                    setField("angkatan", toIntOrPrev(values.angkatan, v))
                  }
                  onBlur={() => onBlur("angkatan")}
                  inputClassName={
                    hasError("angkatan")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("angkatan") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.angkatan}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="tempatLahir"
                  label="Tempat Lahir"
                  value={values.tempatLahir}
                  onChange={(v) => setField("tempatLahir", v)}
                  onBlur={() => onBlur("tempatLahir")}
                  inputClassName={
                    hasError("tempatLahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("tempatLahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tempatLahir}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="tanggalLahir"
                  type="date"
                  label="Tanggal Lahir"
                  value={values.tanggalLahir}
                  onChange={(v) => setField("tanggalLahir", v)}
                  onBlur={() => onBlur("tanggalLahir")}
                  inputClassName={
                    hasError("tanggalLahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("tanggalLahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tanggalLahir}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="id_tingkat_kelas"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="id_tingkat_kelas"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("id_tingkat_kelas")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.id_tingkat_kelas}
                  onChange={(e) =>
                    setField(
                      "id_tingkat_kelas",
                      e.target.value === "" ? "" : Number(e.target.value)
                    )
                  }
                  onBlur={() => onBlur("id_tingkat_kelas")}
                  disabled={isDisabled}
                >
                  <option value="">Pilih tingkat kelas</option>
                  {tingkatKelass.map((tingkat) => (
                    <option
                      key={tingkat.id_tingkat_kelas}
                      value={tingkat.id_tingkat_kelas}
                    >
                      Kelas {tingkat.tingkat_kelas}
                    </option>
                  ))}
                </select>
                {hasError("id_tingkat_kelas") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.id_tingkat_kelas}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="id_nama_kelas"
                  className="text-xs font-medium text-slate-600"
                >
                  Nama Kelas
                </label>
                <select
                  id="id_nama_kelas"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("id_nama_kelas")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.id_nama_kelas}
                  onChange={(e) => setField("id_nama_kelas", e.target.value)}
                  onBlur={() => onBlur("id_nama_kelas")}
                  disabled={isDisabled || namaKelasOptions.length === 0}
                >
                  <option value="">Pilih nama kelas</option>
                  {namaKelasOptions.map((kelas) => (
                    <option
                      key={kelas.id_nama_kelas}
                      value={kelas.id_nama_kelas}
                    >
                      {kelas.nama_kelas}
                    </option>
                  ))}
                </select>
                {hasError("id_nama_kelas") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.id_nama_kelas}
                  </p>
                )}
              </div>
            </div>
          </div>

          <ImageUpload
            sectionTitle="Foto Profil"
            helperText="Unggah foto profil (maks. 2MB)."
            formatText="Format: JPG/PNG"
            imgSrc={fotoUrl}
            fileName={values.fotoProfil?.name}
            size={values.fotoProfil?.size ? Number((values.fotoProfil.size / (1024 * 1024)).toFixed(2)) : undefined}
            imageFileCheck={!!values.fotoProfil}
            onChange={(e) => {
              const file = e.target.files?.[0] ?? null;
              setField("fotoProfil", file);
              onBlur("fotoProfil");
            }}
            onClick={clearFoto}
          />
          {hasError("fotoProfil") && (
            <p className="-mt-4 text-xs text-rose-600">{errors.fotoProfil}</p>
          )}

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
                clearFoto();
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

export default EditAkunSiswaForm;
