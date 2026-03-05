import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import type { StudentRegisterFormValues } from "@/types/KelolaAkun/AkunSiswa";

import { submitStudentRegister } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  fileMaxSize,
  fileTypeStartsWith,
  minLength,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";
import toast from "react-hot-toast";

const initialValues: StudentRegisterFormValues = {
  nama_lengkap: "",
  nisn: "",
  username: "",
  password: "",
  jenis_kelamin: "LAKI_LAKI",
  email: "",
  no_hp: "",
  no_absen: 0,
  angkatan: 0,
  tempat_lahir: "",
  tanggal_lahir: "",
  id_nama_kelas: "",
  foto_profil: null,
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";
const NISN_LENGTH = 10;
const DUPLICATE_ACCOUNT_MESSAGES: Record<string, string> = {
  USERNAME_TAKEN: "Username sudah terdaftar. Gunakan username lain yang unik.",
  EMAIL_TAKEN: "Email sudah terdaftar. Gunakan email lain yang unik.",
  NO_HP_TAKEN: "Nomor HP sudah terdaftar. Gunakan nomor HP lain yang unik.",
  NISN_TAKEN: "NISN sudah terdaftar. Gunakan NISN lain yang unik.",
  NIP_TAKEN: "NIP sudah terdaftar. Gunakan NIP lain yang unik.",
  CONFLICT: "Data yang diinputkan sudah ada sebelumnya. Pastikan data unik.",
};

const uniqueConstraintMessage = (error: ApiError) => {
  if (!error.code) return error.message;
  return DUPLICATE_ACCOUNT_MESSAGES[error.code] ?? error.message;
};

/** Konversi input string menjadi integer. Kalau invalid, pertahankan nilai sebelumnya. */
function toIntOrPrev(prev: number, raw: string): number {
  const s = String(raw ?? "").trim();
  if (s === "") return 0; // reset
  if (!/^\d+$/.test(s)) return prev; // tolak selain digit
  const n = Number(s);
  return Number.isSafeInteger(n) ? n : prev;
}

const normalizeNisnInput = (value: string) => {
  const trimmed = value.trim();
  if (trimmed === "-") return "-";
  const digitsOnly = trimmed.replace(/\D/g, "");
  return digitsOnly.slice(0, NISN_LENGTH);
};

const AkunSiswaForm = () => {
  const navigate = useNavigate();

  const [submitting, setSubmitting] = useState<boolean>(false);
  const [values, setValues] =
    useState<StudentRegisterFormValues>(initialValues);

  // Lebih aman kalau touched mengikuti key form (bukan string bebas)
  const [touched, setTouched] = useState<
    Partial<Record<keyof StudentRegisterFormValues, boolean>>
  >({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const [fotoUrl, setFotoUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (!values.foto_profil) {
      setFotoUrl("");
      return;
    }

    const url = URL.createObjectURL(values.foto_profil);
    setFotoUrl(url);

    return () => {
      URL.revokeObjectURL(url);
    };
  }, [values.foto_profil]);

  const { data: daftarKelasData } = useGetDataKelasFull();

  const setField = createSetField(setValues);

  const onBlur = (name: keyof StudentRegisterFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<StudentRegisterFormValues>({
    nama_lengkap: [requiredString("Nama lengkap wajib diisi.")],
    username: [requiredString("Username wajib diisi.")],
    password: [
      requiredString("Password wajib diisi."),
      minLength(8, "Password minimal 8 karakter."),
    ],
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
        const cleaned = value.trim().replace(/\s+/g, "");
        return /^[0-9+]{8,16}$/.test(cleaned)
          ? null
          : "No HP tidak valid (gunakan angka, 8-16 digit).";
      },
    ],

    // ===== numeric: no_absen =====
    no_absen: [
      (value) => {
        if (!Number.isInteger(value) || value <= 0) {
          return "No absen wajib diisi dan harus angka > 0.";
        }
        return null;
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

    // ===== numeric: angkatan (4 digit + range) =====
    angkatan: [
      (value) => {
        if (!Number.isInteger(value))
          return "Angkatan harus berupa angka bulat.";
        const currentYear = new Date().getFullYear();
        if (value < 2000 || value > currentYear) return "Angkatan tidak valid.";
        if (value < 1000 || value > 9999)
          return "Angkatan harus 4 digit angka.";
        return null;
      },
    ],

    tempat_lahir: [requiredString("Tempat lahir wajib diisi.")],
    tanggal_lahir: [requiredString("Tanggal lahir wajib diisi.")],

    // ===== select string: id_nama_kelas =====
    id_nama_kelas: [requiredValue("Nama kelas wajib dipilih.")],

    foto_profil: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran foto maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof StudentRegisterFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      nisn: true,
      nama_lengkap: true,
      username: true,
      password: true,
      jenis_kelamin: true,
      email: true,
      no_hp: true,
      no_absen: true,
      angkatan: true,
      tempat_lahir: true,
      tanggal_lahir: true,
      id_nama_kelas: true,
      foto_profil: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih salah / kosong.");
      return;
    }

    try {
      setSubmitting(true);

      const kelasTerpilih = daftarKelasData?.item_nama_kelas.find(
        (kelas) => String(kelas.id_nama_kelas) === values.id_nama_kelas,
      );

      if (!kelasTerpilih) {
        setSubmitError("Kelas yang dipilih tidak ditemukan.");
        return;
      }

      const payload: StudentRegisterFormValues = {
        ...values,
        id_nama_kelas: String(kelasTerpilih.id_nama_kelas),
      };

      await submitStudentRegister(payload);
      toast.success("Berhasil membuat akun siswa.");
      setTimeout(
        () =>
          navigate(
            `${paths.dashboard.kelola_akun_siswa}`,
          ),
        1500,
      );
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(uniqueConstraintMessage(error));
      } else {
        setSubmitError("Terjadi kesalahan saat menyimpan data.");
      }
    } finally {
      setSubmitting(false);
    }
  };

  const clearFoto = () => {
    setField("foto_profil", null);
    onBlur("foto_profil");
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-base font-semibold text-slate-900">
                Tambah Akun Siswa
              </h1>
              <p className="mt-1 text-sm text-slate-500">
                Lengkapi data siswa untuk membuat akun.
              </p>
            </div>

            <div className="hidden w-56 md:block">
              <div className="h-2 w-full rounded-full bg-slate-100">
                <div className="h-2 w-full rounded-full bg-[#397e50]" />
              </div>
              <p className="mt-2 text-right text-xs text-slate-500">Step 1/1</p>
            </div>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          {/* INFORMASI DASAR */}
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
                  placeholder="contoh: siti.aminah"
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

              <div>
                <label
                  htmlFor="jenis_kelamin"
                  className="text-xs font-medium text-slate-600"
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
                  label="Email (opsional)"
                  value={values.email ?? ""}
                  onChange={(v) => setField("email", v)}
                  onBlur={() => onBlur("email")}
                  placeholder="nama@gmail.com"
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
                  id="no_hp"
                  type="tel"
                  label="Nomor HP (opsional)"
                  value={values.no_hp ?? ""}
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

          {/* DATA AKADEMIK */}
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Data Akademik</h2>
              <p className={helperText}>Data kelas dan identitas sekolah.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="no_absen"
                  type="text"
                  label="No Absen"
                  value={values.no_absen ? String(values.no_absen) : ""}
                  onChange={(v) => {
                    setValues((prev) => ({
                      ...prev,
                      no_absen: toIntOrPrev(prev.no_absen, v),
                    }));
                  }}
                  onBlur={() => onBlur("no_absen")}
                  placeholder="Contoh: 12"
                  inputClassName={
                    hasError("no_absen") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                {hasError("no_absen") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.no_absen}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="angkatan"
                  type="text"
                  label="Angkatan"
                  value={values.angkatan ? String(values.angkatan) : ""}
                  onChange={(v) => {
                    setValues((prev) => ({
                      ...prev,
                      angkatan: toIntOrPrev(prev.angkatan, v),
                    }));
                  }}
                  onBlur={() => onBlur("angkatan")}
                  placeholder="Contoh: 2025"
                  inputClassName={
                    hasError("angkatan") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                {hasError("angkatan") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.angkatan}
                  </p>
                )}
              </div>

              {/* Nama Kelas */}
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
                  required
                >
                  <option value="" disabled>
                    Pilih nama kelas...
                  </option>
                  {(daftarKelasData?.item_nama_kelas ?? []).map((kelas) => (
                    <option
                      key={kelas.id_nama_kelas}
                      value={String(kelas.id_nama_kelas)}
                    >
                      {kelas.nama_kelas}
                    </option>
                  ))}
                </select>
                <p className="mt-1 text-xs text-slate-500">
                  Nama kelas diambil dari Data Master.
                </p>
                {hasError("id_nama_kelas") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.id_nama_kelas}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* DATA KELAHIRAN */}
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Data Kelahiran</h2>
              <p className={helperText}>Tempat dan tanggal lahir siswa.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="tempat_lahir"
                  label="Tempat Lahir"
                  value={values.tempat_lahir}
                  onChange={(v) => setField("tempat_lahir", v)}
                  onBlur={() => onBlur("tempat_lahir")}
                  placeholder="Contoh: Bandung"
                  inputClassName={
                    hasError("tempat_lahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
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
                  id="tanggalLahir"
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
                  required
                />
                {hasError("tanggal_lahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tanggal_lahir}
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
              optionalText="Akan ditampilkan pada akun siswa."
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
                    setFotoUrl("");
                  }}
                >
                  Reset
                </button>

                <button
                  type="submit"
                  disabled={submitting}
                  className="rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-emerald-600 cursor-pointer disabled:opacity-70"
                >
                  {submitting ? "Menyimpan..." : "Tambah Siswa"}
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AkunSiswaForm;
