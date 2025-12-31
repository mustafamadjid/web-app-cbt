import {  useState,useEffect,useRef } from "react";

type StatusAkun = "AKTIF" | "NONAKTIF";
type JenisKelamin = "LAKI_LAKI" | "PEREMPUAN";

type TeacherEditFormValues = {
  namaLengkap: string;
  email: string;
  username: string;
  password: string;
  noHp: string;
  jenisKelamin: JenisKelamin;
  statusAkun: StatusAkun;
  nip: string;
  jabatan: string;
  bidangStudi: string;
  fotoProfil: File | null;
};

const initialValues: TeacherEditFormValues = {
  namaLengkap: "",
  email: "",
  username: "",
  password: "",
  noHp: "",
  jenisKelamin: "LAKI_LAKI",
  statusAkun: "AKTIF",
  nip: "",
  jabatan: "",
  bidangStudi: "",
  fotoProfil: null,
};

// Type untuk props initialValues dan onSubmit

const inputBase =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";
const labelBase = "text-xs font-medium text-slate-600";
const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const EditAkunGuruForm= () => {
  const [values, setValues] =
    useState<TeacherEditFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [fotoUrl, setFotoUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  // TODO : Panggil service GetGuruById() di halaman editGuruPage
  // Return nya akan jadi isi initialValues
  //  Buat props untuk onSubmit dan initialValues


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

  const setField = <K extends keyof TeacherEditFormValues>(
    key: K,
    value: TeacherEditFormValues[K]
  ) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  const onBlur = (name: keyof TeacherEditFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = (v: TeacherEditFormValues) => {
    const errors: Partial<Record<keyof TeacherEditFormValues, string>> = {};

    if (!v.namaLengkap.trim()) errors.namaLengkap = "Nama lengkap wajib diisi.";
    if (!v.email.trim()) errors.email = "Email wajib diisi.";
    if (v.email && !/^\S+@\S+\.\S+$/.test(v.email))
      errors.email = "Format email tidak valid.";
    if (!v.username.trim()) errors.username = "Username wajib diisi.";
    if (!v.password.trim()) errors.password = "Password wajib diisi.";
    if (v.password && v.password.length < 8)
      errors.password = "Password minimal 8 karakter.";
    if (!v.noHp.trim()) errors.noHp = "Nomor HP wajib diisi.";
    if (!v.nip.trim()) errors.nip = "NIP wajib diisi.";
    if (!v.jabatan.trim()) errors.jabatan = "Jabatan wajib diisi.";
    if (!v.bidangStudi.trim()) errors.bidangStudi = "Bidang studi wajib diisi.";

    // Jenis kelamin wajib (default sudah ada, tapi tetap jaga)
    if (!v.jenisKelamin) errors.jenisKelamin = "Jenis kelamin wajib dipilih.";

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
  const hasError = (name: keyof TeacherEditFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    // mark all touched to show errors
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

    // Contoh payload (biasanya pakai FormData untuk upload file)
    const formData = new FormData();
    formData.append("namaLengkap", values.namaLengkap);
    formData.append("email", values.email);
    formData.append("username", values.username);
    formData.append("password", values.password);
    formData.append("noHp", values.noHp);
    formData.append("jenisKelamin", values.jenisKelamin);
    formData.append("statusAkun", values.statusAkun);
    formData.append("nip", values.nip);
    formData.append("jabatan", values.jabatan);
    formData.append("bidangStudi", values.bidangStudi);
    if (values.fotoProfil) formData.append("fotoProfil", values.fotoProfil);

    // TODO: panggil API
    // await api.put("/guru/register", formData)

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
                  placeholder="Contoh: Budi Santoso"
                />
                {hasError("namaLengkap") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.namaLengkap}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="email">
                  Email
                </label>
                <input
                  id="email"
                  type="email"
                  className={`${inputBase} ${
                    hasError("email") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.email}
                  onChange={(e) => setField("email", e.target.value)}
                  onBlur={() => onBlur("email")}
                  placeholder="nama@sekolah.sch.id"
                />
                {hasError("email") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.email}</p>
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
                  placeholder="contoh: budi.santoso"
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
                <label className={labelBase} htmlFor="statusAkun">
                  Status Akun
                </label>
                <input
                  id="statusAkun"
                  className={inputBase}
                  value={values.statusAkun}
                  disabled
                  readOnly
                />
                <p className="mt-1 text-xs text-slate-500">
                  Default aktif saat registrasi.
                </p>
              </div>

              <div>
                <label className={labelBase} htmlFor="noHp">
                  Nomor HP
                </label>
                <input
                  id="noHp"
                  inputMode="tel"
                  className={`${inputBase} ${
                    hasError("noHp") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.noHp}
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

          {/* DATA KEPEGAWAIAN */}
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Data Kepegawaian</h2>
              <p className={helperText}>Identitas resmi dan peran mengajar.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className={labelBase} htmlFor="nip">
                  NIP
                </label>
                <input
                  id="nip"
                  className={`${inputBase} ${
                    hasError("nip") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.nip}
                  onChange={(e) => setField("nip", e.target.value)}
                  onBlur={() => onBlur("nip")}
                  placeholder="Nomor Induk Pegawai"
                />
                {hasError("nip") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.nip}</p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="jabatan">
                  Jabatan
                </label>
                <input
                  id="jabatan"
                  className={`${inputBase} ${
                    hasError("jabatan") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.jabatan}
                  onChange={(e) => setField("jabatan", e.target.value)}
                  onBlur={() => onBlur("jabatan")}
                  placeholder="Contoh: Guru Tetap"
                />
                {hasError("jabatan") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.jabatan}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label className={labelBase} htmlFor="bidangStudi">
                  Bidang Studi
                </label>
                <input
                  id="bidangStudi"
                  className={`${inputBase} ${
                    hasError("bidangStudi")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.bidangStudi}
                  onChange={(e) => setField("bidangStudi", e.target.value)}
                  onBlur={() => onBlur("bidangStudi")}
                  placeholder="Contoh: Matematika"
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
                    Akan ditampilkan pada akun guru.
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
}

