import React, { useEffect, useRef, useState } from "react";

import ImageUpload from "@/components/features/ImageUpload/ImageUpload";
import InputField from "@/components/common/Input/InputField";

import type { ProfilSekolahFormValues } from "@/types/ProfilSekolah/ProfilSekolah";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  emailFormat,
  fileMaxSize,
  fileTypeStartsWith,
  requiredString,
} from "@/helper/validate/validateForm";

const initialValues: ProfilSekolahFormValues = {
  nama_sekolah: "",
  alamat_sekolah: "",
  no_telp_sekolah: "",
  email_sekolah: "",
  kepala_sekolah: "",
  wakil_kepala_sekolah: "",
  logo_sekolah: null,
  semester: "Ganjil",
  tahun_ajaran: "2025/2026",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const PengaturanProfilForm = () => {
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

  const setField = createSetField(setValues);

  const onBlur = (name: keyof ProfilSekolahFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<ProfilSekolahFormValues>({
    nama_sekolah: [requiredString("Nama sekolah wajib diisi.")],
    alamat_sekolah: [requiredString("Alamat sekolah wajib diisi.")],
    no_telp_sekolah: [requiredString("Nomor telepon sekolah wajib diisi.")],
    email_sekolah: [
      requiredString("Email sekolah wajib diisi."),
      emailFormat("Format email tidak valid."),
    ],
    kepala_sekolah: [requiredString("Nama kepala sekolah wajib diisi.")],
    wakil_kepala_sekolah: [
      requiredString("Nama wakil kepala sekolah wajib diisi."),
    ],
    logo_sekolah: [
      fileMaxSize(2 * 1024 * 1024, "Ukuran logo maksimal 2MB."),
      fileTypeStartsWith("image/", "File harus berupa gambar."),
    ],
  });

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
    if (Object.keys(currentErrors).length > 0) return;

    const formData = new FormData();
    formData.append("nama_sekolah", values.nama_sekolah.trim());
    formData.append("alamat_sekolah", values.alamat_sekolah.trim());
    formData.append("no_telp_sekolah", values.no_telp_sekolah.trim());
    formData.append("email_sekolah", values.email_sekolah.trim());
    formData.append("kepala_sekolah", values.kepala_sekolah.trim());
    formData.append("wakil_kepala_sekolah", values.wakil_kepala_sekolah.trim());
    if (values.logo_sekolah)
      formData.append("logo_sekolah", values.logo_sekolah);

    console.log("READY_TO_SUBMIT", Object.fromEntries(formData.entries()));
  };

  const clearLogo = () => {
    setField("logo_sekolah", null);
    onBlur("logo_sekolah");
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
          {/* INFORMASI SEKOLAH */}
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Sekolah</h2>
              <p className={helperText}>
                Data utama sekolah yang akan tampil di profil.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="nama_sekolah"
                  label="Nama Sekolah"
                  value={values.nama_sekolah}
                  onChange={(v) => setField("nama_sekolah", v)}
                  onBlur={() => onBlur("nama_sekolah")}
                  placeholder="Contoh: SMA Negeri 1"
                  inputClassName={
                    hasError("nama_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("nama_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.nama_sekolah}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="no_telp_sekolah"
                  label="Nomor Telepon"
                  value={values.no_telp_sekolah}
                  onChange={(v) => setField("no_telp_sekolah", v)}
                  onBlur={() => onBlur("no_telp_sekolah")}
                  placeholder="Contoh: 021-123456"
                  inputClassName={
                    hasError("no_telp_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("no_telp_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.no_telp_sekolah}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="email_sekolah"
                  type="email"
                  label="Email Sekolah"
                  value={values.email_sekolah}
                  onChange={(v) => setField("email_sekolah", v)}
                  onBlur={() => onBlur("email_sekolah")}
                  placeholder="Contoh: info@sekolah.sch.id"
                  inputClassName={
                    hasError("email_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("email_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.email_sekolah}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="kepala_sekolah"
                  label="Kepala Sekolah"
                  value={values.kepala_sekolah}
                  onChange={(v) => setField("kepala_sekolah", v)}
                  onBlur={() => onBlur("kepala_sekolah")}
                  placeholder="Nama kepala sekolah"
                  inputClassName={
                    hasError("kepala_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("kepala_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.kepala_sekolah}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="wakil_kepala_sekolah"
                  label="Wakil Kepala Sekolah"
                  value={values.wakil_kepala_sekolah}
                  onChange={(v) => setField("wakil_kepala_sekolah", v)}
                  onBlur={() => onBlur("wakil_kepala_sekolah")}
                  placeholder="Nama wakil kepala sekolah"
                  inputClassName={
                    hasError("wakil_kepala_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("wakil_kepala_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.wakil_kepala_sekolah}
                  </p>
                )}
              </div>

              <div className="md:col-span-2">
                <label
                  className="text-xs font-medium text-slate-600"
                  htmlFor="alamat_sekolah"
                >
                  Alamat Sekolah
                </label>
                <textarea
                  id="alamat_sekolah"
                  className={`min-h-24 w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("alamat_sekolah")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Alamat lengkap sekolah"
                  value={values.alamat_sekolah}
                  onChange={(e) => setField("alamat_sekolah", e.target.value)}
                  onBlur={() => onBlur("alamat_sekolah")}
                  required
                />
                {hasError("alamat_sekolah") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.alamat_sekolah}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* LOGO SEKOLAH (pakai komponen ImageUpload) */}
          <div>
            <ImageUpload
              ref={fileInputRef}
              sectionTitle="Logo Sekolah"
              helperText="Unggah logo berformat PNG/JPG dengan ukuran maksimal 2MB."
              formatText="Format: PNG/JPG"
              optionalText="Logo akan ditampilkan pada profil sekolah."
              imgSrc={logoUrl || undefined}
              imgAlt="Preview Logo Sekolah"
              type="file"
              accept="image/*"
              imageFileCheck={!!values.logo_sekolah}
              fileName={values.logo_sekolah?.name}
              size={
                values.logo_sekolah
                  ? Number(
                      (values.logo_sekolah.size / (1024 * 1024)).toFixed(2)
                    )
                  : undefined
              }
              onChange={(e) => {
                const file = e.target.files?.[0] ?? null;
                setField("logo_sekolah", file);
                onBlur("logo_sekolah");
              }}
              onClick={clearLogo}
            />

            {hasError("logo_sekolah") && (
              <p className="mt-2 text-xs text-rose-500">
                {errors.logo_sekolah}
              </p>
            )}
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

export default PengaturanProfilForm;
