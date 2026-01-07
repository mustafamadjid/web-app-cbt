import React, { useEffect, useRef, useState } from "react";

import { InputField } from "@/components/common/Input/InputField";
import { ImageUpload } from "@/components/features/ImageUpload/ImageUpload";

import type { JenisKelamin} from "@/types/OpsiTypes/Option";
import type { StudentRegisterFormValues } from "@/types/KelolaAkun/AkunSiswa";
import { submitStudentRegister } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { paths } from "@/routes/paths";
import { useNavigate } from "react-router";
import { ApiError } from "@/services/Api/api";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";


const initialValues: StudentRegisterFormValues = {
  namaLengkap: "",
  username: "",
  password: "",
  jenisKelamin: "LAKI_LAKI",
  email: "",
  noHp: "",
  noAbsen: "",
  angkatan: "",
  tempatLahir: "",
  tanggalLahir: "",
  kelasId: "",
  fotoProfil: null,
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const AkunSiswaForm = () => {
  const navigate = useNavigate();

  const [submitting, setSubmitting] = useState<boolean>(false);
  const [values, setValues] =
    useState<StudentRegisterFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const [fotoUrl, setFotoUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (!values.fotoProfil) {
      setFotoUrl("");
      return;
    }

    const url = URL.createObjectURL(values.fotoProfil);
    setFotoUrl(url);

    return () => {
      URL.revokeObjectURL(url);
    };
  }, [values.fotoProfil]);

  const [kelasOptions, setKelasOptions] = useState<number[]>([]);

  useEffect(() => {
    let active = true;
    const loadKelas = async () => {
      const data = await getTingkatKelasOptions();
      if (!active) return;
      setKelasOptions(data);
    };
    loadKelas();
    return () => {
      active = false;
    };
  }, []);

  const setField = <K extends keyof StudentRegisterFormValues>(
    key: K,
    value: StudentRegisterFormValues[K]
  ) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  const onBlur = (name: keyof StudentRegisterFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = (v: StudentRegisterFormValues) => {
    const errors: Partial<Record<keyof StudentRegisterFormValues, string>> = {};

    if (!v.namaLengkap.trim()) errors.namaLengkap = "Nama lengkap wajib diisi.";
    if (!v.username.trim()) errors.username = "Username wajib diisi.";
    if (!v.password.trim()) errors.password = "Password wajib diisi.";
    if (v.password && v.password.length < 8)
      errors.password = "Password minimal 8 karakter.";

    if (!v.jenisKelamin) errors.jenisKelamin = "Jenis kelamin wajib dipilih.";

    if (v.email && v.email.trim() && !/^\S+@\S+\.\S+$/.test(v.email.trim())) {
      errors.email = "Format email tidak valid.";
    }

    if (v.noHp && v.noHp.trim()) {
      const cleaned = v.noHp.trim().replace(/\s+/g, "");
      if (!/^[0-9+]{8,16}$/.test(cleaned)) {
        errors.noHp = "No HP tidak valid (gunakan angka, 8-16 digit).";
      }
    }

    if (!v.noAbsen.trim()) {
      errors.noAbsen = "No absen wajib diisi.";
    } else if (!/^\d+$/.test(v.noAbsen.trim())) {
      errors.noAbsen = "No absen harus berupa angka.";
    } else if (parseInt(v.noAbsen.trim(), 10) <= 0) {
      errors.noAbsen = "No absen minimal 1.";
    }

    if (!v.angkatan.trim()) {
      errors.angkatan = "Angkatan wajib diisi.";
    } else if (!/^\d{4}$/.test(v.angkatan.trim())) {
      errors.angkatan = "Angkatan harus 4 digit (contoh: 2025).";
    }

    if (!v.tempatLahir.trim()) errors.tempatLahir = "Tempat lahir wajib diisi.";
    if (!v.tanggalLahir.trim())
      errors.tanggalLahir = "Tanggal lahir wajib diisi.";

    if (!v.kelasId.trim()) errors.kelasId = "Tingkat kelas wajib dipilih.";

    if (v.fotoProfil) {
      const maxBytes = 2 * 1024 * 1024;
      if (v.fotoProfil.size > maxBytes)
        errors.fotoProfil = "Ukuran foto maksimal 2MB.";
      if (!v.fotoProfil.type.startsWith("image/"))
        errors.fotoProfil = "File harus berupa gambar.";
    }

    return errors;
  };

  const errors = validate(values);
  const hasError = (name: keyof StudentRegisterFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
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
      kelasId: true,
      fotoProfil: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih salah / kosong.");
      return;
    }

    try {
      setSubmitting(true);
      await submitStudentRegister(values);
      alert("Akun siswa berhasil dibuat.");
      setTimeout(() => navigate(`dashboard/administrator/${paths.dashboard.kelola_akun_siswa}`), 1500);
    } catch (error) {
      if (error instanceof ApiError){
              setSubmitError(error.message);
        } 
    }finally{
      setSubmitting(false);
    }
  };

  const clearFoto = () => {
    setField("fotoProfil", null);
    onBlur("fotoProfil");
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

              {/* Jenis Kelamin (select native) */}
              <div>
                <label
                  htmlFor="jenisKelamin"
                  className="text-xs font-medium text-slate-600"
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
                  id="noHp"
                  type="tel"
                  label="Nomor HP (opsional)"
                  value={values.noHp ?? ""}
                  onChange={(v) => setField("noHp", v)}
                  onBlur={() => onBlur("noHp")}
                  placeholder="Contoh: 081234567890"
                  inputClassName={
                    hasError("noHp") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required={false}
                />
                {hasError("noHp") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noHp}</p>
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
                  id="noAbsen"
                  type="text"
                  label="No Absen"
                  value={values.noAbsen}
                  onChange={(v) => setField("noAbsen", v)}
                  onBlur={() => onBlur("noAbsen")}
                  placeholder="Contoh: 12"
                  inputClassName={
                    hasError("noAbsen") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                {hasError("noAbsen") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noAbsen}</p>
                )}
              </div>

              <div>
                <InputField
                  id="angkatan"
                  type="text"
                  label="Angkatan"
                  value={values.angkatan}
                  onChange={(v) => setField("angkatan", v)}
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

              {/* Kelas (select native) */}
              <div className="md:col-span-2">
                <label
                  htmlFor="kelasId"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="kelasId"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("kelasId") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelasId}
                  onChange={(e) => setField("kelasId", e.target.value)}
                  onBlur={() => onBlur("kelasId")}
                  required
                >
                  <option value="" disabled>
                    Pilih tingkat kelas...
                  </option>
                  {kelasOptions.map((tingkat) => (
                    <option key={tingkat} value={String(tingkat)}>
                      Kelas {tingkat}
                    </option>
                  ))}
                </select>

                <p className="mt-1 text-xs text-slate-500">
                  Tingkat kelas akan diambil dari Data Master.
                </p>

                {hasError("kelasId") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.kelasId}</p>
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
                  id="tempatLahir"
                  label="Tempat Lahir"
                  value={values.tempatLahir}
                  onChange={(v) => setField("tempatLahir", v)}
                  onBlur={() => onBlur("tempatLahir")}
                  placeholder="Contoh: Bandung"
                  inputClassName={
                    hasError("tempatLahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
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
                  required
                />
                {hasError("tanggalLahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tanggalLahir}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* FOTO PROFIL (pakai komponen ImageUpload) */}
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
                    setFotoUrl("");
                  }}
                >
                  Reset
                </button>

                <button
                  type="submit"
                  className="rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-emerald-600 cursor-pointer"
                >
                  Tambah Siswa
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};
