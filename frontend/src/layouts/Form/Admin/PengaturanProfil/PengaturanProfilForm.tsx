import { useEffect, useRef, useState } from "react";

import type { ProfilSekolahFormValues } from "@/types/ProfilSekolah/ProfilSekolah";

const initialValues: ProfilSekolahFormValues = {
  nama_sekolah: "",
  alamat_sekolah: "",
  no_telp_sekolah: "",
  email_sekolah: "",
  kepala_sekolah: "",
  wakil_kepala_sekolah: "",
  logo_sekolah: null,
};

const inputBase =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";
const labelBase = "text-xs font-medium text-slate-600";
const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const PengaturanProfilForm = () => {
  const [values, setValues] = useState<ProfilSekolahFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [logoUrl, setLogoUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (!values.logo_sekolah) {
      setLogoUrl("");
      return;
    }

    const url = URL.createObjectURL(values.logo_sekolah);
    setLogoUrl(url);

    return () => {
      URL.revokeObjectURL(url);
    };
  }, [values.logo_sekolah]);

  const setField = <K extends keyof ProfilSekolahFormValues>(
    key: K,
    value: ProfilSekolahFormValues[K]
  ) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  const onBlur = (name: keyof ProfilSekolahFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = (v: ProfilSekolahFormValues) => {
    const errors: Partial<Record<keyof ProfilSekolahFormValues, string>> = {};

    if (!v.nama_sekolah.trim()) errors.nama_sekolah = "Nama sekolah wajib diisi.";
    if (!v.alamat_sekolah.trim())
      errors.alamat_sekolah = "Alamat sekolah wajib diisi.";
    if (!v.no_telp_sekolah.trim())
      errors.no_telp_sekolah = "Nomor telepon sekolah wajib diisi.";
    if (!v.email_sekolah.trim())
      errors.email_sekolah = "Email sekolah wajib diisi.";
    if (v.email_sekolah && !/^\S+@\S+\.\S+$/.test(v.email_sekolah))
      errors.email_sekolah = "Format email tidak valid.";
    if (!v.kepala_sekolah.trim())
      errors.kepala_sekolah = "Nama kepala sekolah wajib diisi.";
    if (!v.wakil_kepala_sekolah.trim())
      errors.wakil_kepala_sekolah = "Nama wakil kepala sekolah wajib diisi.";

    if (v.logo_sekolah) {
      const maxBytes = 2 * 1024 * 1024;
      if (v.logo_sekolah.size > maxBytes)
        errors.logo_sekolah = "Ukuran logo maksimal 2MB.";
      if (!v.logo_sekolah.type.startsWith("image/"))
        errors.logo_sekolah = "File harus berupa gambar.";
    }

    return errors;
  };

  const errors = validate(values);
  const hasError = (name: keyof ProfilSekolahFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    setTouched({
      nama_sekolah: true,
      alamat_sekolah: true,
      no_telp_sekolah: true,
      email_sekolah: true,
      kepala_sekolah: true,
      wakil_kepala_sekolah: true,
      logo_sekolah: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      return;
    }

    const formData = new FormData();
    formData.append("nama_sekolah", values.nama_sekolah);
    formData.append("alamat_sekolah", values.alamat_sekolah);
    formData.append("no_telp_sekolah", values.no_telp_sekolah);
    formData.append("email_sekolah", values.email_sekolah);
    formData.append("kepala_sekolah", values.kepala_sekolah);
    formData.append("wakil_kepala_sekolah", values.wakil_kepala_sekolah);
    if (values.logo_sekolah)
      formData.append("logo_sekolah", values.logo_sekolah);

    console.log("READY_TO_SUBMIT", Object.fromEntries(formData.entries()));
  };

  const clearLogo = () => {
    setField("logo_sekolah", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Pengaturan Profil Sekolah
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Perbarui informasi sekolah dan unggah logo resmi.
            </p>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Sekolah</h2>
              <p className={helperText}>
                Data utama sekolah yang akan tampil di profil.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className={labelBase} htmlFor="nama_sekolah">
                  Nama Sekolah
                </label>
                <input
                  id="nama_sekolah"
                  className={`${inputBase} ${
                    hasError("nama_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Contoh: SMA Negeri 1"
                  value={values.nama_sekolah}
                  onChange={(e) => setField("nama_sekolah", e.target.value)}
                  onBlur={() => onBlur("nama_sekolah")}
                />
                {hasError("nama_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.nama_sekolah}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="no_telp_sekolah">
                  Nomor Telepon
                </label>
                <input
                  id="no_telp_sekolah"
                  className={`${inputBase} ${
                    hasError("no_telp_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Contoh: 021-123456"
                  value={values.no_telp_sekolah}
                  onChange={(e) => setField("no_telp_sekolah", e.target.value)}
                  onBlur={() => onBlur("no_telp_sekolah")}
                />
                {hasError("no_telp_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.no_telp_sekolah}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="email_sekolah">
                  Email Sekolah
                </label>
                <input
                  id="email_sekolah"
                  type="email"
                  className={`${inputBase} ${
                    hasError("email_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Contoh: info@sekolah.sch.id"
                  value={values.email_sekolah}
                  onChange={(e) => setField("email_sekolah", e.target.value)}
                  onBlur={() => onBlur("email_sekolah")}
                />
                {hasError("email_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.email_sekolah}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="kepala_sekolah">
                  Kepala Sekolah
                </label>
                <input
                  id="kepala_sekolah"
                  className={`${inputBase} ${
                    hasError("kepala_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Nama kepala sekolah"
                  value={values.kepala_sekolah}
                  onChange={(e) => setField("kepala_sekolah", e.target.value)}
                  onBlur={() => onBlur("kepala_sekolah")}
                />
                {hasError("kepala_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.kepala_sekolah}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="wakil_kepala_sekolah">
                  Wakil Kepala Sekolah
                </label>
                <input
                  id="wakil_kepala_sekolah"
                  className={`${inputBase} ${
                    hasError("wakil_kepala_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Nama wakil kepala sekolah"
                  value={values.wakil_kepala_sekolah}
                  onChange={(e) =>
                    setField("wakil_kepala_sekolah", e.target.value)
                  }
                  onBlur={() => onBlur("wakil_kepala_sekolah")}
                />
                {hasError("wakil_kepala_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.wakil_kepala_sekolah}
                  </p>
                )}
              </div>

              <div className="md:col-span-2">
                <label className={labelBase} htmlFor="alamat_sekolah">
                  Alamat Sekolah
                </label>
                <textarea
                  id="alamat_sekolah"
                  className={`${inputBase} min-h-24 ${
                    hasError("alamat_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Alamat lengkap sekolah"
                  value={values.alamat_sekolah}
                  onChange={(e) => setField("alamat_sekolah", e.target.value)}
                  onBlur={() => onBlur("alamat_sekolah")}
                />
                {hasError("alamat_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.alamat_sekolah}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Logo Sekolah</h2>
              <p className={helperText}>
                Unggah logo berformat PNG/JPG dengan ukuran maksimal 2MB.
              </p>
            </div>

            <div className="flex flex-col gap-4 md:flex-row md:items-center">
              <div className="flex h-28 w-28 items-center justify-center overflow-hidden rounded-xl border border-dashed border-slate-200 bg-slate-50">
                {logoUrl ? (
                  <img
                    src={logoUrl}
                    alt="Preview Logo Sekolah"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <span className="text-xs text-slate-400">
                    Preview logo
                  </span>
                )}
              </div>

              <div className="flex-1 space-y-2">
                <input
                  ref={fileInputRef}
                  id="logo_sekolah"
                  type="file"
                  accept="image/*"
                  className="hidden"
                  onChange={(e) =>
                    setField(
                      "logo_sekolah",
                      e.target.files ? e.target.files[0] : null
                    )
                  }
                  onBlur={() => onBlur("logo_sekolah")}
                />

                <div className="flex flex-wrap gap-2">
                  <label
                    htmlFor="logo_sekolah"
                    className="inline-flex cursor-pointer items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 text-xs font-medium text-slate-600 shadow-sm transition hover:border-[#397e50] hover:text-[#397e50]"
                  >
                    Pilih File
                  </label>
                  {values.logo_sekolah && (
                    <button
                      type="button"
                      onClick={clearLogo}
                      className="rounded-lg border border-rose-200  px-3 py-2 text-xs font-medium text-rose-500 cursor-pointer hover:bg-rose-50 hover:border-rose-300 transition duration-100"
                    >
                      Hapus Logo
                    </button>
                  )}
                </div>

                <p className="text-xs text-slate-500">
                  {values.logo_sekolah
                    ? values.logo_sekolah.name
                    : "Belum ada file yang dipilih."}
                </p>
                {hasError("logo_sekolah") && (
                  <p className="text-xs text-rose-500">
                    {errors.logo_sekolah}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3">
            <button
              type="button"
              className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600"
            >
              Batal
            </button>
            <button
              type="submit"
              className="rounded-lg bg-[#397e50] px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-[#2f6541]"
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
