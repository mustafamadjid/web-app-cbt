import React, { useEffect, useRef, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { TeacherUpdateFormValues } from "@/types/KelolaAkun/AkunGuru";
import type { ResetPasswordFormValues } from "@/types/KelolaAkun/ResetPassword";
import type { StatusAkun } from "@/types/OpsiTypes/Option";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  emailFormat,
  fileMaxSize,
  fileTypeStartsWith,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import { ApiError } from "@/services/Api/api";

type EditAkunGuruFormProps = {
  initialValues: TeacherUpdateFormValues;
  initialFotoUrl?: string;
  onSubmit: (values: TeacherUpdateFormValues) => Promise<void>;
  onSubmitResetPassword: (values: ResetPasswordFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";
const NIP_LENGTH = 18;

const normalizeNipInput = (value: string) => {
  const trimmed = value.trim();
  if (trimmed === "-") return "-";
  const digitsOnly = trimmed.replace(/\D/g, "");
  return digitsOnly.slice(0, NIP_LENGTH);
};

const EditAkunGuruForm = ({
  initialValues,
  initialFotoUrl,
  onSubmit,
  onSubmitResetPassword,
  loading = false,
  submitting = false,
}: EditAkunGuruFormProps) => {
  const [values, setValues] = useState<TeacherUpdateFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const [resetPasswordValues, setResetPasswordValues] = useState<ResetPasswordFormValues>({
    password: "",
    konfirmasi_password: "",
  });
  const [resetPasswordTouched, setResetPasswordTouched] = useState<Record<string, boolean>>({});
  const [resetPasswordError, setResetPasswordError] = useState<string | null>(null);
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

  const onBlur = (name: keyof TeacherUpdateFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<TeacherUpdateFormValues>({
    role: [requiredString("Role wajib diisi.")],
    nama_lengkap: [requiredString("Nama lengkap wajib diisi.")],
    email: [
      requiredString("Email wajib diisi."),
      emailFormat("Format email tidak valid."),
    ],
    username: [requiredString("Username wajib diisi.")],
    no_hp: [requiredString("Nomor HP wajib diisi.")],
    nip: [
      requiredString("NIP wajib diisi."),
      (value) => {
        const trimmed = value.trim();
        if (trimmed === "-") return null;
        if (trimmed.length !== NIP_LENGTH) {
          return "NIP belum valid (kurang dari 18 digit).";
        }
        return null;
      },
    ],
    jabatan: [requiredString("Jabatan wajib diisi.")],
    bidang_studi: [requiredString("Bidang studi wajib diisi.")],
    jenis_kelamin: [requiredValue("Jenis kelamin wajib dipilih.")],
    status_akun: [requiredValue("Status akun wajib dipilih.")],
    foto_profil: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran foto maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof TeacherUpdateFormValues) =>
    !!errors[name] && !!touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      role: true,
      nama_lengkap: true,
      email: true,
      username: true,
      no_hp: true,
      jenis_kelamin: true,
      status_akun: true,
      nip: true,
      jabatan: true,
      bidang_studi: true,
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


  const handleSubmitResetPassword = async () => {
    setResetPasswordError(null);
    setResetPasswordTouched({ password: true, konfirmasi_password: true });

    if (!resetPasswordValues.password.trim()) {
      setResetPasswordError("Password wajib diisi.");
      return;
    }

    if (resetPasswordValues.password !== resetPasswordValues.konfirmasi_password) {
      setResetPasswordError("Konfirmasi password tidak sama.");
      return;
    }

    try {
      await onSubmitResetPassword(resetPasswordValues);
      setResetPasswordValues({ password: "", konfirmasi_password: "" });
      setResetPasswordTouched({});
    } catch (error) {
      if (error instanceof ApiError) {
        setResetPasswordError(error.message);
      }
    }
  };

  const clearFoto = () => {
    setField("foto_profil", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const isDisabled = loading || submitting;
  const isResetPasswordStarted =
    !!resetPasswordValues.password.trim() ||
    !!resetPasswordValues.konfirmasi_password.trim();

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
                  id="nama_lengkap"
                  label="Nama Lengkap"
                  value={values.nama_lengkap}
                  onChange={(v) => setField("nama_lengkap", v)}
                  onBlur={() => onBlur("nama_lengkap")}
                  placeholder="Contoh: Budi Santoso"
                  inputClassName={
                    hasError("nama_lengkap")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("nama_lengkap") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.nama_lengkap}
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
                <label className="text-xs font-medium text-slate-600" htmlFor="jenis_kelamin">
                  Jenis Kelamin
                </label>
                <select
                  id="jenis_kelamin"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("jenis_kelamin")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.jenis_kelamin}
                  onChange={(e) =>
                    setField(
                      "jenis_kelamin",
                      e.target.value as "LAKI_LAKI" | "PEREMPUAN"
                    )
                  }
                  onBlur={() => onBlur("jenis_kelamin")}
                  disabled={isDisabled}
                >
                  <option value="LAKI_LAKI">Laki-laki</option>
                  <option value="PEREMPUAN">Perempuan</option>
                </select>
                {hasError("jenis_kelamin") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.jenis_kelamin}
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
                <label className="text-xs font-medium text-slate-600" htmlFor="status_akun">
                  Status Akun
                </label>
                <select
                  id="status_akun"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("status_akun")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.status_akun}
                  onChange={(e) =>
                    setField("status_akun", e.target.value as StatusAkun)
                  }
                  onBlur={() => onBlur("status_akun")}
                  disabled={isDisabled}
                >
                  <option value="AKTIF">Aktif</option>
                  <option value="NONAKTIF">Nonaktif</option>
                </select>
                {hasError("status_akun") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.status_akun}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="no_hp"
                  type="tel"
                  label="Nomor HP"
                  value={values.no_hp}
                  onChange={(v) => setField("no_hp", v)}
                  onBlur={() => onBlur("no_hp")}
                  placeholder="Contoh: 081234567890"
                  inputClassName={
                    hasError("no_hp") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("no_hp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.no_hp}</p>
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
                  onChange={(v) => setField("nip", normalizeNipInput(v))}
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
                  id="bidang_studi"
                  label="Bidang Studi"
                  value={values.bidang_studi}
                  onChange={(v) => setField("bidang_studi", v)}
                  onBlur={() => onBlur("bidang_studi")}
                  placeholder="Contoh: Matematika"
                  inputClassName={
                    hasError("bidang_studi")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("bidang_studi") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.bidang_studi}
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



          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="password_baru_guru"
                  type="password"
                  label="Password Baru"
                  value={resetPasswordValues.password}
                  onChange={(v) => setResetPasswordValues((prev) => ({ ...prev, password: v }))}
                  onBlur={() => setResetPasswordTouched((prev) => ({ ...prev, password: true }))}
                  disabled={isDisabled}
                />
                {!!resetPasswordTouched.password && isResetPasswordStarted && !resetPasswordValues.password.trim() && (
                  <p className="mt-1 text-xs text-rose-600">Password wajib diisi.</p>
                )}
              </div>

              <div>
                <InputField
                  id="konfirmasi_password_baru_guru"
                  type="password"
                  label="Konfirmasi Password"
                  value={resetPasswordValues.konfirmasi_password}
                  onChange={(v) => setResetPasswordValues((prev) => ({ ...prev, konfirmasi_password: v }))}
                  onBlur={() => setResetPasswordTouched((prev) => ({ ...prev, konfirmasi_password: true }))}
                  disabled={isDisabled}
                />
                {!!resetPasswordTouched.konfirmasi_password && isResetPasswordStarted && resetPasswordValues.password !== resetPasswordValues.konfirmasi_password && (
                  <p className="mt-1 text-xs text-rose-600">Konfirmasi password tidak sama.</p>
                )}
              </div>
            </div>

            {resetPasswordError && (
              <div className="mt-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
                {resetPasswordError}
              </div>
            )}

            <div className="mt-4 flex justify-end">
              <button
                type="button"
                onClick={handleSubmitResetPassword}
                className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
                disabled={isDisabled}
              >
                Update Password
              </button>
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
