import React, { useEffect, useRef, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { TeacherRegisterFormValues } from "@/types/KelolaAkun/AkunGuru";
import type { StatusAkun } from "@/types/OpsiTypes/Option";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  emailFormat,
  fileMaxSize,
  fileTypeStartsWith,
  minLength,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import { ApiError } from "@/services/Api/api";

type EditAkunGuruFormProps = {
  initialValues: TeacherRegisterFormValues;
  initialFotoUrl?: string;
  onSubmit: (values: TeacherRegisterFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const EditAkunGuruForm = ({
  initialValues,
  initialFotoUrl,
  onSubmit,
  loading = false,
  submitting = false,
}: EditAkunGuruFormProps) => {
  const [values, setValues] = useState<TeacherRegisterFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [fotoUrl, setFotoUrl] = useState<string>(initialFotoUrl ?? "");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    setValues(initialValues);
    setTouched({});
  }, [initialValues]);

  useEffect(() => {
    if (values.foto_profil) {
      const url = URL.createObjectURL(values.foto_profil);
      setFotoUrl(url);
      return () => URL.revokeObjectURL(url);
    }

    setFotoUrl(initialFotoUrl ?? "");
    return undefined;
  }, [initialFotoUrl, values.foto_profil]);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof TeacherRegisterFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<TeacherRegisterFormValues>({
    role: [requiredString("Role wajib diisi.")],
    namaLengkap: [requiredString("Nama lengkap wajib diisi.")],
    email: [
      requiredString("Email wajib diisi."),
      emailFormat("Format email tidak valid."),
    ],
    username: [requiredString("Username wajib diisi.")],
    password: [minLength(8, "Password minimal 8 karakter.")],
    noHp: [requiredString("Nomor HP wajib diisi.")],
    nip: [requiredString("NIP wajib diisi.")],
    jabatan: [requiredString("Jabatan wajib diisi.")],
    bidangStudi: [requiredString("Bidang studi wajib diisi.")],
    jenisKelamin: [requiredValue("Jenis kelamin wajib dipilih.")],
    statusAkun: [requiredValue("Status akun wajib dipilih.")],
    foto_profil: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran foto maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof TeacherRegisterFormValues) =>
    !!errors[name] && !!touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      role: true,
      namaLengkap: true,
      email: true,
      username: true,
      password: true,
      noHp: true,
      jenisKelamin: true,
      statusAkun: true,
      nip: true,
      jabatan: true,
      bidangStudi: true,
      foto_profil: true,
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
    setField("foto_profil", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const isDisabled = loading || submitting;

  return (
    <div className="min-h-screen w-full bg-slate-50 py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Edit Akun Guru
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Perbarui data akun guru sesuai kebutuhan.
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Dasar</h2>
              <p className={helperText}>Data akun dan identitas guru.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="namaLengkap"
                  label="Nama Lengkap"
                  value={values.namaLengkap}
                  onChange={(v) => setField("namaLengkap", v)}
                  onBlur={() => onBlur("namaLengkap")}
                  placeholder="Contoh: Budi Santoso"
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
                  id="email"
                  type="email"
                  label="Email"
                  value={values.email}
                  onChange={(v) => setField("email", v)}
                  onBlur={() => onBlur("email")}
                  placeholder="nama@sekolah.sch.id"
                  inputClassName={
                    hasError("email") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("email") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.email}</p>
                )}
              </div>

              <div>
                <InputField
                  id="username"
                  label="Username"
                  value={values.username}
                  onChange={(v) => setField("username", v)}
                  onBlur={() => onBlur("username")}
                  placeholder="contoh: budi.santoso"
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
                    setField("jenisKelamin", e.target.value as "LAKI_LAKI" | "PEREMPUAN")
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
                <label className="text-xs font-medium text-slate-600" htmlFor="role">
                  Role Akun
                </label>
                <select
                  id="role"
                  className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500"
                  value={values.role}
                  onChange={(e) => setField("role", e.target.value)}
                  onBlur={() => onBlur("role")}
                  disabled={isDisabled}
                >
                  <option value="GURU">Guru</option>
                  <option value="ADMIN">Admin</option>
                </select>
                {hasError("role") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.role}</p>
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

              <div>
                <InputField
                  id="noHp"
                  type="tel"
                  label="Nomor HP"
                  value={values.noHp}
                  onChange={(v) => setField("noHp", v)}
                  onBlur={() => onBlur("noHp")}
                  placeholder="Contoh: 081234567890"
                  inputClassName={
                    hasError("noHp") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("noHp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noHp}</p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Data Kepegawaian</h2>
              <p className={helperText}>Identitas resmi dan peran mengajar.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="nip"
                  label="NIP"
                  value={values.nip}
                  onChange={(v) => setField("nip", v)}
                  onBlur={() => onBlur("nip")}
                  placeholder="Nomor Induk Pegawai"
                  inputClassName={
                    hasError("nip") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("nip") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.nip}</p>
                )}
              </div>

              <div>
                <InputField
                  id="jabatan"
                  label="Jabatan"
                  value={values.jabatan}
                  onChange={(v) => setField("jabatan", v)}
                  onBlur={() => onBlur("jabatan")}
                  placeholder="Contoh: Guru Tetap"
                  inputClassName={
                    hasError("jabatan") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("jabatan") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.jabatan}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <InputField
                  id="bidangStudi"
                  label="Bidang Studi"
                  value={values.bidangStudi}
                  onChange={(v) => setField("bidangStudi", v)}
                  onBlur={() => onBlur("bidangStudi")}
                  placeholder="Contoh: Matematika"
                  inputClassName={
                    hasError("bidangStudi")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("bidangStudi") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.bidangStudi}
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
            fileName={values.foto_profil?.name}
            size={values.foto_profil?.size ? Number((values.foto_profil.size / (1024 * 1024)).toFixed(2)) : undefined}
            imageFileCheck={!!values.foto_profil}
            onChange={(e) => {
              const file = e.target.files?.[0] ?? null;
              setField("foto_profil", file);
              onBlur("foto_profil");
            }}
            onClick={clearFoto}
          />
          {hasError("foto_profil") && (
            <p className="-mt-4 text-xs text-rose-600">{errors.foto_profil}</p>
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

export default EditAkunGuruForm;
