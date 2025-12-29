import React, { useEffect, useRef, useState } from "react";

type JenisKelamin = "LAKI_LAKI" | "PEREMPUAN";

type KelasOption = {
  id: string;
  nama: string; // contoh: "X IPA 1"
};

type StudentRegisterFormValues = {
  namaLengkap: string;
  username: string;
  password: string;

  jenisKelamin: JenisKelamin;

  email?: string; // optional
  noHp?: string; // optional

  noAbsen: string; // string supaya input bebas, divalidasi angka
  angkatan: string; // contoh: 2025

  tempatLahir: string;
  tanggalLahir: string; // YYYY-MM-DD dari input type="date"

  kelasId: string; // dropdown dari data master nanti

  fotoProfil: File | null; // upload + preview
};

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

const inputBase =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";
const labelBase = "text-xs font-medium text-slate-600";
const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const AkunSiswaForm = () => {
  const [values, setValues] =
    useState<StudentRegisterFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  // Foto preview (aman: create+revoke di effect, bukan di memo)
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

  // Dummy options - nanti ganti dari API Data Master
  const kelasOptions: KelasOption[] = [
    { id: "kelas-1", nama: "X IPA 1" },
    { id: "kelas-2", nama: "X IPA 2" },
    { id: "kelas-3", nama: "X IPS 1" },
  ];

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

    // Optional email: kalau diisi harus valid
    if (v.email && v.email.trim() && !/^\S+@\S+\.\S+$/.test(v.email.trim())) {
      errors.email = "Format email tidak valid.";
    }

    // Optional noHp: kalau diisi minimal angka (simple)
    if (v.noHp && v.noHp.trim()) {
      const cleaned = v.noHp.trim().replace(/\s+/g, "");
      if (!/^[0-9+]{8,16}$/.test(cleaned)) {
        errors.noHp = "No HP tidak valid (gunakan angka, 8-16 digit).";
      }
    }

    // No absen wajib, harus angka > 0
    if (!v.noAbsen.trim()) {
      errors.noAbsen = "No absen wajib diisi.";
    } else if (!/^\d+$/.test(v.noAbsen.trim())) {
      errors.noAbsen = "No absen harus berupa angka.";
    } else if (parseInt(v.noAbsen.trim(), 10) <= 0) {
      errors.noAbsen = "No absen minimal 1.";
    }

    // Angkatan wajib, contoh 4 digit
    if (!v.angkatan.trim()) {
      errors.angkatan = "Angkatan wajib diisi.";
    } else if (!/^\d{4}$/.test(v.angkatan.trim())) {
      errors.angkatan = "Angkatan harus 4 digit (contoh: 2025).";
    }

    if (!v.tempatLahir.trim()) errors.tempatLahir = "Tempat lahir wajib diisi.";
    if (!v.tanggalLahir.trim())
      errors.tanggalLahir = "Tanggal lahir wajib diisi.";

    // Kelas wajib
    if (!v.kelasId.trim()) errors.kelasId = "Kelas wajib dipilih.";

    // Foto profil (opsional, tapi validasi kalau diisi)
    if (v.fotoProfil) {
      const maxBytes = 2 * 1024 * 1024; // 2MB
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

  const onSubmit = (e: React.FormEvent) => {
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

    const formData = new FormData();
    formData.append("namaLengkap", values.namaLengkap.trim());
    formData.append("username", values.username.trim());
    formData.append("password", values.password);
    formData.append("jenisKelamin", values.jenisKelamin);

    if (values.email?.trim()) formData.append("email", values.email.trim());
    if (values.noHp?.trim()) formData.append("noHp", values.noHp.trim());

    formData.append("noAbsen", values.noAbsen.trim());
    formData.append("angkatan", values.angkatan.trim());
    formData.append("tempatLahir", values.tempatLahir.trim());
    formData.append("tanggalLahir", values.tanggalLahir);
    formData.append("kelasId", values.kelasId);

    if (values.fotoProfil) formData.append("fotoProfil", values.fotoProfil);

    // TODO: panggil API
    // await api.post("/siswa/register", formData)

    console.log("READY_TO_SUBMIT", Object.fromEntries(formData.entries()));
    alert("Form valid. Lanjutkan submit ke API.");
  };

  const clearFoto = () => {
    setField("fotoProfil", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <div className="min-h-screen w-full bg-slate-50 py-8">
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
                <div className="h-2 w-[100%] rounded-full bg-[#397e50]" />
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
                <label className={labelBase} htmlFor="namaLengkap">
                  Nama Lengkap
                </label>
                <input
                  id="namaLengkap"
                  className={`${inputBase} ${
                    hasError("namaLengkap")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.namaLengkap}
                  onChange={(e) => setField("namaLengkap", e.target.value)}
                  onBlur={() => onBlur("namaLengkap")}
                  placeholder="Contoh: Siti Aminah"
                />
                {hasError("namaLengkap") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.namaLengkap}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="username">
                  Username
                </label>
                <input
                  id="username"
                  className={`${inputBase} ${
                    hasError("username") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.username}
                  onChange={(e) => setField("username", e.target.value)}
                  onBlur={() => onBlur("username")}
                  placeholder="contoh: siti.aminah"
                />
                {hasError("username") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.username}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="password">
                  Password
                </label>
                <input
                  id="password"
                  type="password"
                  className={`${inputBase} ${
                    hasError("password") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.password}
                  onChange={(e) => setField("password", e.target.value)}
                  onBlur={() => onBlur("password")}
                  placeholder="Minimal 8 karakter"
                />
                {hasError("password") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.password}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="jenisKelamin">
                  Jenis Kelamin
                </label>
                <select
                  id="jenisKelamin"
                  className={`${inputBase} ${
                    hasError("jenisKelamin")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.jenisKelamin}
                  onChange={(e) =>
                    setField("jenisKelamin", e.target.value as JenisKelamin)
                  }
                  onBlur={() => onBlur("jenisKelamin")}
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
                <label className={labelBase} htmlFor="email">
                  Email <span className="text-slate-400">(opsional)</span>
                </label>
                <input
                  id="email"
                  type="email"
                  className={`${inputBase} ${
                    hasError("email") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.email ?? ""}
                  onChange={(e) => setField("email", e.target.value)}
                  onBlur={() => onBlur("email")}
                  placeholder="nama@gmail.com"
                />
                {hasError("email") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.email}</p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="noHp">
                  Nomor HP <span className="text-slate-400">(opsional)</span>
                </label>
                <input
                  id="noHp"
                  inputMode="tel"
                  className={`${inputBase} ${
                    hasError("noHp") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.noHp ?? ""}
                  onChange={(e) => setField("noHp", e.target.value)}
                  onBlur={() => onBlur("noHp")}
                  placeholder="Contoh: 081234567890"
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
                <label className={labelBase} htmlFor="noAbsen">
                  No Absen
                </label>
                <input
                  id="noAbsen"
                  inputMode="numeric"
                  className={`${inputBase} ${
                    hasError("noAbsen") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.noAbsen}
                  onChange={(e) => setField("noAbsen", e.target.value)}
                  onBlur={() => onBlur("noAbsen")}
                  placeholder="Contoh: 12"
                />
                {hasError("noAbsen") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.noAbsen}</p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="angkatan">
                  Angkatan
                </label>
                <input
                  id="angkatan"
                  inputMode="numeric"
                  className={`${inputBase} ${
                    hasError("angkatan") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.angkatan}
                  onChange={(e) => setField("angkatan", e.target.value)}
                  onBlur={() => onBlur("angkatan")}
                  placeholder="Contoh: 2025"
                />
                {hasError("angkatan") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.angkatan}
                  </p>
                )}
              </div>

              <div className="md:col-span-2">
                <label className={labelBase} htmlFor="kelasId">
                  Kelas
                </label>
                <select
                  id="kelasId"
                  className={`${inputBase} ${
                    hasError("kelasId") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelasId}
                  onChange={(e) => setField("kelasId", e.target.value)}
                  onBlur={() => onBlur("kelasId")}
                >
                  <option value="" disabled>
                    Pilih kelas...
                  </option>
                  {kelasOptions.map((k) => (
                    <option key={k.id} value={k.id}>
                      {k.nama}
                    </option>
                  ))}
                </select>
                <p className="mt-1 text-xs text-slate-500">
                  Kelas akan diambil dari Data Master.
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
                <label className={labelBase} htmlFor="tempatLahir">
                  Tempat Lahir
                </label>
                <input
                  id="tempatLahir"
                  className={`${inputBase} ${
                    hasError("tempatLahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.tempatLahir}
                  onChange={(e) => setField("tempatLahir", e.target.value)}
                  onBlur={() => onBlur("tempatLahir")}
                  placeholder="Contoh: Bandung"
                />
                {hasError("tempatLahir") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tempatLahir}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="tanggalLahir">
                  Tanggal Lahir
                </label>
                <input
                  id="tanggalLahir"
                  type="date"
                  className={`${inputBase} ${
                    hasError("tanggalLahir")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.tanggalLahir}
                  onChange={(e) => setField("tanggalLahir", e.target.value)}
                  onBlur={() => onBlur("tanggalLahir")}
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
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4 flex items-start justify-between gap-4">
              <div>
                <h2 className={sectionTitle}>Foto Profil</h2>
                <p className={helperText}>Unggah foto profil (maks. 2MB).</p>
              </div>

              <div className="text-right">
                <p className="text-xs text-slate-500">Format: JPG/PNG</p>
              </div>
            </div>

            <div className="flex flex-col gap-4 md:flex-row md:items-center">
              {/* Preview */}
              <div className="flex items-center gap-3">
                <div className="h-24 w-24 overflow-hidden rounded-md border border-slate-200 bg-slate-100">
                  {fotoUrl ? (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img
                      src={fotoUrl}
                      alt="Preview foto profil"
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div className="flex h-full w-full items-center justify-center text-xs text-slate-400">
                      No Photo
                    </div>
                  )}
                </div>

                <div>
                  <p className="text-sm font-medium text-slate-700">
                    Foto profil
                  </p>
                  <p className="text-xs text-slate-500">
                    Akan ditampilkan pada akun siswa.
                  </p>
                </div>
              </div>

              {/* Upload + Chip */}
              <div className="flex flex-1 flex-col gap-2 md:flex-row md:items-center md:justify-end">
                <label className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm transition hover:bg-slate-50">
                  Pilih Foto
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={(e) => {
                      const file = e.target.files?.[0] ?? null;
                      setField("fotoProfil", file);
                      onBlur("fotoProfil");
                    }}
                  />
                </label>

                {values.fotoProfil ? (
                  <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-slate-50 px-3 py-2">
                    <div className="min-w-0">
                      <p className="truncate text-xs font-medium text-slate-700">
                        {values.fotoProfil.name}
                      </p>
                      <p className="text-[11px] text-slate-500">
                        {(values.fotoProfil.size / (1024 * 1024)).toFixed(2)} MB
                      </p>
                    </div>
                    <button
                      type="button"
                      onClick={clearFoto}
                      className="rounded-md px-2 py-1 text-xs font-medium text-slate-600 hover:bg-white"
                      aria-label="Hapus foto"
                      title="Hapus foto"
                    >
                      X
                    </button>
                  </div>
                ) : (
                  <p className="text-xs text-slate-500 md:text-right">
                    Belum ada file dipilih.
                  </p>
                )}
              </div>
            </div>

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
                    clearFoto(); // sekaligus clear input file
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
