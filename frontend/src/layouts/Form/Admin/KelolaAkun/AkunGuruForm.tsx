import React, { useEffect, useRef, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { TeacherRegisterFormValues } from "@/types/KelolaAkun/AkunGuru";
import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import { submitTeacherRegister } from "@/services/Api/features-api/KelolaAkun/akunguru.service";
import { ApiError } from "@/services/Api/api";
import { useNavigate } from "react-router";
import { paths } from "@/routes/paths";

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

const initialValues: TeacherRegisterFormValues = {
  role: "guru",
  namaLengkap: "",
  email: "",
  username: "",
  password: "",
  noHp: "",
  jenisKelamin: "LAKI_LAKI",
  statusAkun: "aktif",
  nip: "",
  jabatan: "",
  bidangStudi: "",
  fotoProfil: null,
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

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
    if (!values.fotoProfil) {
      setFotoUrl("");
      return;
    }

    const url = URL.createObjectURL(values.fotoProfil);
    setFotoUrl(url);

    return () => URL.revokeObjectURL(url);
  }, [values.fotoProfil]);

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
    password: [
      requiredString("Password wajib diisi."),
      minLength(8, "Password minimal 8 karakter."),
    ],
    noHp: [requiredString("Nomor HP wajib diisi.")],
    nip: [requiredString("NIP wajib diisi.")],
    jabatan: [requiredString("Jabatan wajib diisi.")],
    bidangStudi: [requiredString("Bidang studi wajib diisi.")],
    jenisKelamin: [requiredValue("Jenis kelamin wajib dipilih.")],
    fotoProfil: [
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
      fotoProfil: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih salah / kosong.");
      return;
    }

    try {
      setSubmitting(true);
      await submitTeacherRegister(values);

      alert("Akun guru berhasil dibuat.");
      setTimeout(
        () =>
          navigate(
            `dashboard/administrator/${paths.dashboard.kelola_akun_guru}`
          ),
        1500
      );
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const clearFoto = () => {
    setField("fotoProfil", null);
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
                    hasError("username") ? "border-rose-300 ring-rose-100" : ""
                  }
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
                  htmlFor="jenisKelamin"
                >
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
                  required
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

              {/* Role */}
              <div>
                <label
                  className="text-xs font-medium text-slate-600"
                  htmlFor="jenisKelamin"
                >
                  Role Akun
                </label>
                <select
                  id="role"
                  className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500"
                  value={values.role}
                  onChange={(e) => setField("role", e.target.value)}
                  onBlur={() => onBlur("role")}
                  required
                >
                  <option value="GURU">Guru</option>
                  <option value="ADMIN">Admin</option>
                </select>
                {hasError("role") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.role}</p>
                )}
              </div>

              <div>
                <InputField
                  id="statusAkun"
                  label="Status Akun"
                  value={values.statusAkun}
                  onChange={() => {}}
                  disabled
                  required={false}
                />
                <p className="mt-1 text-xs text-slate-500">
                  Default aktif saat registrasi.
                </p>
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
                  required
                />
                {hasError("noHp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noHp}</p>
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
                  onChange={(v) => setField("nip", v)}
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

          {/* FOTO PROFIL */}
          <div>
            <ImageUpload
              ref={fileInputRef}
              sectionTitle="Foto Profil"
              helperText="Unggah foto profil (maks. 2MB)."
              formatText="Format: JPG/PNG"
              optionalText="Akan ditampilkan pada akun guru."
              imgSrc={fotoUrl || undefined}
              imgAlt="Preview foto profil"
              type="file"
              accept="image/*"
              imageFileCheck={!!values.fotoProfil}
              fileName={values.fotoProfil?.name}
              size={
                values.fotoProfil
                  ? Number((values.fotoProfil.size / (1024 * 1024)).toFixed(2))
                  : undefined
              }
              onChange={(e) => {
                const file = e.target.files?.[0] ?? null;
                setField("fotoProfil", file);
                onBlur("fotoProfil");
              }}
              onClick={clearFoto}
            />

            {hasError("fotoProfil") && (
              <p className="mt-2 text-xs text-rose-600">{errors.fotoProfil}</p>
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
                >
                  Reset
                </button>

                <button
                  type="submit"
                  className="rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-emerald-600 cursor-pointer"
                >
                  Daftarkan Guru
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
