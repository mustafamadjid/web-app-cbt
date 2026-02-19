import React, { useEffect, useRef, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { StudentUpdateFormValues } from "@/types/KelolaAkun/AkunSiswa";
import type { ResetPasswordFormValues } from "@/types/KelolaAkun/ResetPassword";
import type { JenisKelamin, StatusAkun } from "@/types/OpsiTypes/Option";
import type { NamaKelas } from "@/types/DataMaster/Kelas";

import { GetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import { ApiError } from "@/services/Api/api";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  fileMaxSize,
  fileTypeStartsWith,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";

type EditAkunSiswaFormProps = {
  initialValues: StudentUpdateFormValues;
  initialFotoUrl?: string;
  onSubmit: (values: StudentUpdateFormValues) => Promise<void>;
  onSubmitResetPassword: (values: ResetPasswordFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";
const NISN_LENGTH = 10;

function toIntOrPrev(prev: number, raw: string): number {
  const s = String(raw ?? "").trim();
  if (s === "") return 0;
  if (!/^\d+$/.test(s)) return prev;
  const n = Number(s);
  return Number.isSafeInteger(n) ? n : prev;
}

const normalizeNisnInput = (value: string) => {
  const trimmed = value.trim();
  if (trimmed === "-") return "-";
  const digitsOnly = trimmed.replace(/\D/g, "");
  return digitsOnly.slice(0, NISN_LENGTH);
};

const EditAkunSiswaForm = ({
  initialValues,
  initialFotoUrl,
  onSubmit,
  onSubmitResetPassword,
  loading = false,
  submitting = false,
}: EditAkunSiswaFormProps) => {
  const [values, setValues] =
    useState<StudentUpdateFormValues>(initialValues);
  const [touched, setTouched] = useState<
    Partial<Record<keyof StudentUpdateFormValues, boolean>>
  >({});
  const [submitError, setSubmitError] = useState<string | null>(null);


  const [resetPasswordValues, setResetPasswordValues] = useState<ResetPasswordFormValues>({
    password: "",
    konfirmasi_password: "",
  });
  const [resetPasswordTouched, setResetPasswordTouched] = useState<Record<string, boolean>>({});
  const [resetPasswordError, setResetPasswordError] = useState<string | null>(null);
  const [fotoUrl, setFotoUrl] = useState<string>(initialFotoUrl ?? "");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [daftarKelas, setDaftarKelas] = useState<NamaKelas[]>([]);

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

  useEffect(() => {
    let active = true;
    const loadKelas = async () => {
      const data = await GetDataKelasFull();
      if (!active) return;
      setDaftarKelas(data.item_nama_kelas ?? []);
    };
    loadKelas();
    return () => {
      active = false;
    };
  }, []);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof StudentUpdateFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<StudentUpdateFormValues>({
    nama_lengkap: [requiredString("Nama lengkap wajib diisi.")],
    username: [requiredString("Username wajib diisi.")],
    jenis_kelamin: [requiredValue("Jenis kelamin wajib dipilih.")],
    email: [
      (value) => {
        if (!value || !value.trim()) return null;
        return /^\S+@\S+\.\S+$/.test(value.trim())
          ? null
          : "Format email tidak valid.";
      },
    ],
    no_hp: [
      (value) => {
        if (!value || !value.trim()) return null;
        return /^\d{10,15}$/.test(value.trim())
          ? null
          : "No HP harus berupa 10-15 digit angka.";
      },
    ],
    nisn: [
      requiredString("NISN wajib diisi."),
      (value) => {
        const trimmed = value.trim();
        if (trimmed === "-") return null;
        if (trimmed.length !== NISN_LENGTH) {
          return "NISN belum valid (kurang dari 10 digit).";
        }
        return null;
      },
    ],
    no_absen: [
      (value) => (value <= 0 ? "No absen harus lebih dari 0." : null),
    ],
    angkatan: [
      (value) => (value <= 0 ? "Angkatan wajib diisi." : null),
    ],
    tempat_lahir: [requiredString("Tempat lahir wajib diisi.")],
    tanggal_lahir: [requiredString("Tanggal lahir wajib diisi.")],
    id_nama_kelas: [requiredValue("Nama kelas wajib dipilih.")],
    foto_profil: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran foto maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
    status_akun: [requiredValue("Status akun wajib dipilih.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof StudentUpdateFormValues) =>
    !!errors[name] && !!touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      nama_lengkap: true,
      nisn: true,
      username: true,
      jenis_kelamin: true,
      email: true,
      no_hp: true,
      no_absen: true,
      angkatan: true,
      tempat_lahir: true,
      tanggal_lahir: true,
      id_nama_kelas: true,
      foto_profil: true,
      status_akun: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih salah / kosong.");
      return;
    }

    const kelasTerpilih = daftarKelas.find(
      (kelas) => String(kelas.id_nama_kelas) === values.id_nama_kelas,
    );

    if (!kelasTerpilih) {
      setSubmitError("Kelas yang dipilih tidak ditemukan.");
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
                  id="nama_lengkap"
                  label="Nama Lengkap"
                  value={values.nama_lengkap}
                  onChange={(v) => setField("nama_lengkap", v)}
                  onBlur={() => onBlur("nama_lengkap")}
                  placeholder="Contoh: Siti Aminah"
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
                    setField("jenis_kelamin", e.target.value as JenisKelamin)
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
                <InputField
                  id="nisn"
                  type="text"
                  label="NISN (Jika tidak ada, isi dengan - )"
                  value={values.nisn ?? ""}
                  onChange={(v) => setField("nisn", normalizeNisnInput(v))}
                  onBlur={() => onBlur("nisn")}
                  placeholder="Contoh: 1234567890"
                  inputClassName={
                    hasError("nisn") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required={false}
                />
                {hasError("nisn") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.nisn}</p>
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
                  id="no_hp"
                  type="tel"
                  label="Nomor HP"
                  value={values.no_hp ?? ""}
                  onChange={(v) => setField("no_hp", v)}
                  onBlur={() => onBlur("no_hp")}
                  placeholder="Contoh: 081234567890"
                  inputClassName={
                    hasError("no_hp") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required={false}
                />
                {hasError("no_hp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.no_hp}</p>
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
                  id="no_absen"
                  type="number"
                  label="Nomor Absen"
                  value={String(values.no_absen)}
                  onChange={(v) =>
                    setField("no_absen", toIntOrPrev(values.no_absen, v))
                  }
                  onBlur={() => onBlur("no_absen")}
                  inputClassName={
                    hasError("no_absen") ? "border-rose-300 ring-rose-100" : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("no_absen") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.no_absen}</p>
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
                  id="tempat_lahir"
                  label="Tempat Lahir"
                  value={values.tempat_lahir}
                  onChange={(v) => setField("tempat_lahir", v)}
                  onBlur={() => onBlur("tempat_lahir")}
                  inputClassName={
                    hasError("tempat_lahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("tempat_lahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tempat_lahir}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="tanggal_lahir"
                  type="date"
                  label="Tanggal Lahir"
                  value={values.tanggal_lahir}
                  onChange={(v) => setField("tanggal_lahir", v)}
                  onBlur={() => onBlur("tanggal_lahir")}
                  inputClassName={
                    hasError("tanggal_lahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  disabled={isDisabled}
                  required
                />
                {hasError("tanggal_lahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tanggal_lahir}
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
                  disabled={isDisabled || daftarKelas.length === 0}
                >
                  <option value="">Pilih nama kelas</option>
                  {daftarKelas.map((kelas) => (
                    <option
                      key={kelas.id_nama_kelas}
                      value={String(kelas.id_nama_kelas)}
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
              <h2 className={sectionTitle}>Update Password</h2>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="password_baru_siswa"
                  type="password"
                  label="Password Baru"
                  value={resetPasswordValues.password}
                  onChange={(v) => setResetPasswordValues((prev) => ({ ...prev, password: v }))}
                  onBlur={() => setResetPasswordTouched((prev) => ({ ...prev, password: true }))}
                  disabled={isDisabled}
                  required
                />
                {!!resetPasswordTouched.password && !resetPasswordValues.password.trim() && (
                  <p className="mt-1 text-xs text-rose-600">Password wajib diisi.</p>
                )}
              </div>

              <div>
                <InputField
                  id="konfirmasi_password_baru_siswa"
                  type="password"
                  label="Konfirmasi Password"
                  value={resetPasswordValues.konfirmasi_password}
                  onChange={(v) => setResetPasswordValues((prev) => ({ ...prev, konfirmasi_password: v }))}
                  onBlur={() => setResetPasswordTouched((prev) => ({ ...prev, konfirmasi_password: true }))}
                  disabled={isDisabled}
                  required
                />
                {!!resetPasswordTouched.konfirmasi_password && resetPasswordValues.password !== resetPasswordValues.konfirmasi_password && (
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

export default EditAkunSiswaForm;
