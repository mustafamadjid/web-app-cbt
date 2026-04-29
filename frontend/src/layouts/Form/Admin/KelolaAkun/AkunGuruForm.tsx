import React, { useEffect, useRef, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { TeacherRegisterFormValues } from "@/types/KelolaAkun/AkunGuru";
import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import { submitTeacherRegister } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { ApiError } from "@/services/Api/api";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import { useNavigate } from "react-router";
import { paths } from "@/routes/paths";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  emailFormat,
  fileMaxSize,
  fileTypeStartsWith,
  maxLength,
  minLength,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import toast from "react-hot-toast";
import {
  USERNAME_HELPER_TEXT,
  USERNAME_LENGTH_INVALID_MESSAGE,
  USERNAME_MAX_LENGTH,
  USERNAME_MIN_LENGTH,
} from "@/constants/username";

const initialValues: TeacherRegisterFormValues = {
  nama_lengkap: "",
  email: "",
  username: "",
  password: "",
  no_hp: "",
  jenis_kelamin: "LAKI_LAKI",
  nip: "",
  jabatan: "",
  bidang_studi: "",
  foto_profil: null,
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";
const NIP_LENGTH = 18;
const DUPLICATE_ACCOUNT_MESSAGES: Record<string, string> = {
  USERNAME_TAKEN: "Username sudah terdaftar. Gunakan username lain yang unik.",
  USERNAME_LENGTH_INVALID: USERNAME_LENGTH_INVALID_MESSAGE,
  EMAIL_TAKEN: "Email sudah terdaftar. Gunakan email lain yang unik.",
  NO_HP_TAKEN: "Nomor HP sudah terdaftar. Gunakan nomor HP lain yang unik.",
  NISN_TAKEN: "NISN sudah terdaftar. Gunakan NISN lain yang unik.",
  NIP_TAKEN: "NIP sudah terdaftar. Gunakan NIP lain yang unik.",
  CONFLICT: "Data yang diinputkan sudah ada sebelumnya. Pastikan data unik.",
};

const uniqueConstraintMessage = (error: ApiError) => {
  return (
    (error.code ? DUPLICATE_ACCOUNT_MESSAGES[error.code] : undefined) ??
    getUserFriendlyErrorMessage(error, {
      action: "create",
      entity: "akun guru",
    })
  );
};

const normalizeNipInput = (value: string) => {
  const trimmed = value.trim();
  if (trimmed === "-") return "-";
  const digitsOnly = trimmed.replace(/\D/g, "");
  return digitsOnly.slice(0, NIP_LENGTH);
};

const AkunGuruForm = () => {
  const navigate = useNavigate();

  const [values, setValues] =
    useState<TeacherRegisterFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [fotoUrl, setFotoUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (!values.foto_profil) {
      setFotoUrl("");
      return;
    }

    const url = URL.createObjectURL(values.foto_profil);
    setFotoUrl(url);

    return () => URL.revokeObjectURL(url);
  }, [values.foto_profil]);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof TeacherRegisterFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<TeacherRegisterFormValues>({
    nama_lengkap: [requiredString("Nama lengkap wajib diisi.")],
    email: [emailFormat("Format email tidak valid.")],
    username: [
      requiredString("Username wajib diisi."),
      minLength(USERNAME_MIN_LENGTH, USERNAME_LENGTH_INVALID_MESSAGE),
      maxLength(USERNAME_MAX_LENGTH, USERNAME_LENGTH_INVALID_MESSAGE),
    ],
    password: [
      requiredString("Password wajib diisi."),
      minLength(8, "Password minimal 8 karakter."),
    ],
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
    foto_profil: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran foto maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof TeacherRegisterFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      nama_lengkap: true,
      email: true,
      username: true,
      password: true,
      no_hp: true,
      jenis_kelamin: true,
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
      setSubmitting(true);
      await submitTeacherRegister(values);

      toast.success("Registrasi akun guru berhasil.");
      setTimeout(
        () =>
          navigate(
            `${paths.dashboard.kelola_akun_guru}`
          ),
        1500
      );
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(uniqueConstraintMessage(error));
      } else {
        setSubmitError(
          getUserFriendlyErrorMessage(error, {
            action: "create",
            entity: "akun guru",
          }),
        );
      }
    } finally {
      setSubmitting(false);
    }
  };

  const clearFoto = () => {
    setField("foto_profil", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        {/* Header + Progress */}
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-base font-semibold text-slate-900">
                Registrasi Akun Guru
              </h1>
              <p className="mt-1 text-sm text-slate-500">
                Lengkapi data guru untuk membuat akun.
              </p>
            </div>

            <div className="hidden w-56 md:block">
              <div className="h-2 w-full rounded-full bg-slate-100">
                <div className="h-2 w-[45%] rounded-full bg-[#397e50]" />
              </div>
              <p className="mt-2 text-right text-xs text-slate-500">Step 1/2</p>
            </div>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          {/* INFORMASI DASAR */}
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
                  required={false}
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
                  maxLength={USERNAME_MAX_LENGTH}
                  inputClassName={
                    hasError("username") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                <p className="mt-1 text-xs text-slate-500">
                  {USERNAME_HELPER_TEXT}
                </p>
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
                  label="Password"
                  value={values.password}
                  onChange={(v) => setField("password", v)}
                  onBlur={() => onBlur("password")}
                  placeholder="Minimal 8 karakter"
                  inputClassName={
                    hasError("password") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                {hasError("password") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.password}
                  </p>
                )}
              </div>

              {/* Jenis kelamin (select native) */}
              <div>
                <label
                  className="text-xs font-medium text-slate-600"
                  htmlFor="jenis_kelamin"
                >
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
                    setField("jenis_kelamin", e.target.value as JenisKelamin)
                  }
                  onBlur={() => onBlur("jenis_kelamin")}
                  required
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
                  required={false}
                />
                {hasError("no_hp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.no_hp}</p>
                )}
              </div>
            </div>
          </div>

          {/* DATA KEPEGAWAIAN */}
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

          {/* FOTO PROFIL */}
          <div>
            <ImageUpload
              ref={fileInputRef}
              sectionTitle="Foto Profil"
              helperText="Unggah foto profil (maks. 2MB)."
              formatText="Format: JPG/PNG"
              optionalText="Opsional, akan ditampilkan pada akun guru."
              imgSrc={fotoUrl || undefined}
              imgAlt="Preview foto profil"
              type="file"
              accept="image/*"
              imageFileCheck={!!values.foto_profil}
              fileName={values.foto_profil?.name}
              size={
                values.foto_profil
                  ? Number((values.foto_profil.size / (1024 * 1024)).toFixed(2))
                  : undefined
              }
              onChange={(e) => {
                const file = e.target.files?.[0] ?? null;
                setField("foto_profil", file);
                onBlur("foto_profil");
              }}
              onClick={clearFoto}
            />

            {hasError("foto_profil") && (
              <p className="mt-2 text-xs text-rose-600">{errors.foto_profil}</p>
            )}
          </div>

          {/* SUBMIT */}
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            {submitError && (
              <div className="mb-4 rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700">
                {submitError}
              </div>
            )}

            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <p className="text-xs text-slate-500">
                Dengan mendaftar, data akan disimpan sesuai kebijakan sekolah.
              </p>

              <div className="flex gap-2">
                <button
                  type="button"
                  className="rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50 cursor-pointer"
                  onClick={() => {
                    setValues(initialValues);
                    setTouched({});
                    setSubmitError(null);
                    if (fileInputRef.current) fileInputRef.current.value = "";
                  }}
                  disabled={submitting}
                >
                  Reset
                </button>

                <button
                  type="submit"
                  className="rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-emerald-600 cursor-pointer"
                  disabled={submitting}
                >
                  {submitting ? "Mendaftarkan..." : "Daftarkan Guru"}
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AkunGuruForm;
