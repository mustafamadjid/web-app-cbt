import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router";

import InputField from "@/components/common/Input/InputField";
import ImageUpload from "@/components/features/Upload/ImageUpload";

import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import type { StudentRegisterFormValues } from "@/types/KelolaAkun/AkunSiswa";
import type { KelasRow, TingkatKelas } from "@/types/DataMaster/Kelas";

import { submitStudentRegister } from "@/services/Api/features-api/KelolaAkun/akunsiswa.service";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { getKelas } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getTingkatKelass } from "@/services/Api/features-api/GetOptions/options.service";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  fileMaxSize,
  fileTypeStartsWith,
  minLength,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";

/**
 * Catatan penting:
 * - Kamu minta tetap pakai number untuk field numeric: noAbsen, angkatan, id_tingkat_kelas.
 * - Karena InputField biasanya mengembalikan string, kita konversi dengan aman di onChange.
 */

const initialValues: StudentRegisterFormValues = {
  role: "SISWA",
  namaLengkap: "",
  username: "",
  password: "",
  jenisKelamin: "LAKI_LAKI",
  email: "",
  noHp: "",
  noAbsen: 0,
  angkatan: 0,
  tempatLahir: "",
  tanggalLahir: "",
  id_tingkat_kelas: 0, // 0 = belum pilih
  id_nama_kelas: "", // string id kelas
  fotoProfil: null,
  statusAkun: "aktif",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

/** Konversi input string menjadi integer. Kalau invalid, pertahankan nilai sebelumnya. */
function toIntOrPrev(prev: number, raw: string): number {
  const s = String(raw ?? "").trim();
  if (s === "") return 0; // reset
  if (!/^\d+$/.test(s)) return prev; // tolak selain digit
  const n = Number(s);
  return Number.isSafeInteger(n) ? n : prev;
}

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

  const [TingkatKelass, setTingkatKelass] = useState<TingkatKelas[]>([]);
  const [daftarKelas, setDaftarKelas] = useState<KelasRow[]>([]);
  const [namaKelasOptions, setNamaKelasOptions] = useState<KelasRow[]>([]);

  useEffect(() => {
    let active = true;
    const loadKelas = async () => {
      const [tingkat, kelas] = await Promise.all([
        getTingkatKelass(),
        getKelas(),
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
      (kelas) => String(kelas.id) === values.id_nama_kelas,
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
    password: [
      requiredString("Password wajib diisi."),
      minLength(8, "Password minimal 8 karakter."),
    ],
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
        const cleaned = value.trim().replace(/\s+/g, "");
        return /^[0-9+]{8,16}$/.test(cleaned)
          ? null
          : "No HP tidak valid (gunakan angka, 8-16 digit).";
      },
    ],

    // ===== numeric: noAbsen =====
    noAbsen: [
      (value) => {
        if (!Number.isInteger(value) || value <= 0) {
          return "No absen wajib diisi dan harus angka > 0.";
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

    tempatLahir: [requiredString("Tempat lahir wajib diisi.")],
    tanggalLahir: [requiredString("Tanggal lahir wajib diisi.")],

    // ===== select numeric: id_tingkat_kelas =====
    id_tingkat_kelas: [
      (value) => (value && value > 0 ? null : "Tingkat kelas wajib dipilih."),
    ],

    // ===== select string: id_nama_kelas =====
    id_nama_kelas: [requiredValue("Nama kelas wajib dipilih.")],

    fotoProfil: [
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
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih salah / kosong.");
      return;
    }

    try {
      setSubmitting(true);

      // FIX: cari kelas berdasarkan id_tingkat_kelas (number) dan id kelas (string)
      const kelasTerpilih = daftarKelas.find(
        (kelas) =>
          kelas.id_tingkat_kelas === values.id_tingkat_kelas &&
          String(kelas.id) === values.id_nama_kelas,
      );

      if (!kelasTerpilih) {
        setSubmitError("Kelas yang dipilih tidak ditemukan.");
        return;
      }

      const payload: StudentRegisterFormValues = {
        ...values,
        // Pastikan konsisten dengan yang ditemukan
        id_tingkat_kelas: kelasTerpilih.id_tingkat_kelas,
        id_nama_kelas: String(kelasTerpilih.id),
      };

      await submitStudentRegister(payload);
      alert("Akun siswa berhasil dibuat.");
      setTimeout(
        () =>
          navigate(
            `dashboard/administrator/${paths.dashboard.kelola_akun_siswa}`,
          ),
        1500,
      );
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      } else {
        setSubmitError("Terjadi kesalahan saat menyimpan data.");
      }
    } finally {
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
                  value={values.noAbsen ? String(values.noAbsen) : ""}
                  onChange={(v) => {
                    setValues((prev) => ({
                      ...prev,
                      noAbsen: toIntOrPrev(prev.noAbsen, v),
                    }));
                  }}
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

              {/* Tingkat Kelas */}
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
                  onChange={(e) => {
                    const n = Number(e.target.value);
                    setField("id_tingkat_kelas", Number.isFinite(n) ? n : 0);
                  }}
                  onBlur={() => onBlur("id_tingkat_kelas")}
                  required
                >
                  <option value={0} disabled>
                    Pilih tingkat kelas...
                  </option>

                  {/* FIX: value harus id_tingkat_kelas (bukan tingkat_kelas) */}
                  {TingkatKelass.map((tingkat) => (
                    <option
                      key={tingkat.id_tingkat_kelas}
                      value={tingkat.id_tingkat_kelas}
                    >
                      Kelas {tingkat.tingkat_kelas}
                    </option>
                  ))}
                </select>

                <p className="mt-1 text-xs text-slate-500">
                  Tingkat kelas akan diambil dari Data Master.
                </p>

                {hasError("id_tingkat_kelas") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.id_tingkat_kelas}
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
                  disabled={!values.id_tingkat_kelas}
                >
                  <option value="" disabled>
                    Pilih nama kelas...
                  </option>
                  {namaKelasOptions.map((kelas) => (
                    <option key={kelas.id} value={String(kelas.id)}>
                      {kelas.nama_kelas}
                    </option>
                  ))}
                </select>
                <p className="mt-1 text-xs text-slate-500">
                  Nama kelas mengikuti tingkat yang dipilih.
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
