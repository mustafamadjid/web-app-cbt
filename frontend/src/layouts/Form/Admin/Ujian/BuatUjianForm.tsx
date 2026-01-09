import React, { useEffect, useMemo, useState } from "react";

import { InputField } from "@/components/common/Input/InputField";

import type {
  BankSoalOption,
  BuatUjianFormValues,
  GuruPengawasOption,
  SesiUjianOption,
  SiswaPreviewItem,
  TipeUjian,
} from "@/types/Ujian/BuatUjian";
import type { TingkatKelasOption } from "@/types/DataMaster/Kelas";
import {
  getUjianBankSoalOptions,
  getUjianGuruPengawasOptions,
  getUjianSesiOptions,
  getUjianSiswaPreview,
  submitBuatUjian,
} from "@/services/Api/features-api/Ujian/ujian.service";

import { getRuangUjianOptions } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";

import { ApiError } from "@/services/Api/api";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";

// helper
import { createSetField } from "@/helper/setField/setField";
import { calculateDuration } from "@/helper/CalculateDuration/calculateDuration";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";

const initialValues: BuatUjianFormValues = {
  nama_ujian: "",
  deskripsi_ujian: "",
  tipe_ujian: "PILIHAN_GANDA",
  kelas_id: "",
  bank_soal_id: "",
  jumlah_soal: 0,
  tanggal_ujian: "",
  waktu_mulai: "",
  waktu_selesai: "",
  durasi_menit: 0,
  ruang_ujian_id: "",
  acak_soal: true,
  guru_pengawas_id: "",
  sesi_id: "",
  token_ujian: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const selectBaseClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";


export const BuatUjianForm = () => {
  const [values, setValues] = useState<BuatUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [loadingBankSoal, setLoadingBankSoal] = useState(false);
  const [loadingSiswa, setLoadingSiswa] = useState(false);

  const [kelasOptions, setKelasOptions] = useState<TingkatKelasOption[]>([]);
  const [bankSoalOptions, setBankSoalOptions] = useState<BankSoalOption[]>([]);
  const [ruangOptions, setRuangOptions] = useState<RuangUjianRow[]>([]);
  const [guruOptions, setGuruOptions] = useState<GuruPengawasOption[]>([]);
  const [sesiOptions, setSesiOptions] = useState<SesiUjianOption[]>([]);
  const [siswaPreview, setSiswaPreview] = useState<SiswaPreviewItem[]>([]);

  const setField = createSetField(setValues);

  const onBlur = (name: keyof BuatUjianFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const selectedKelasId = useMemo(() => {
    const tingkatKelas =
      values.kelas_id === "" ? undefined : getTingkatKelasById(values.kelas_id);
    return tingkatKelas ?? undefined;
  }, [values.kelas_id]);

  const bankSoalById = useMemo(() => {
    const map = new Map(bankSoalOptions.map((x) => [x.id, x]));
    return map;
  }, [bankSoalOptions]);

  const selectedBankSoal =
    values.bank_soal_id === "" ? undefined : bankSoalById.get(values.bank_soal_id);

  useEffect(() => {
    let active = true;

    const loadOptions = async () => {
      try {
        const [kelas, ruang, guru, sesi] = await Promise.all([
          getTingkatKelasOptions(),
          getRuangUjianOptions(),
          getUjianGuruPengawasOptions(),
          getUjianSesiOptions(),
        ]);

        if (!active) return;
        setKelasOptions(kelas);
        setRuangOptions(ruang);
        setGuruOptions(guru);
        setSesiOptions(sesi);
      } catch (error) {
        if (!active) return;
        setSubmitError("Gagal memuat data pendukung ujian.");
      }
    };

    loadOptions();
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    let active = true;
    const loadBankSoal = async () => {
      setLoadingBankSoal(true);
      try {
        const data = await getUjianBankSoalOptions({
          tingkatKelasId: values.kelas_id === "" ? undefined : values.kelas_id,
        });
        if (!active) return;
        setBankSoalOptions(data);
      } finally {
        if (active) setLoadingBankSoal(false);
      }
    };

    if (selectedKelasId != null) {
      loadBankSoal();
    } else {
      setBankSoalOptions([]);
      setField("bank_soal_id", "");
      setField("jumlah_soal", 0);
    }

    return () => {
      active = false;
    };
  }, [selectedKelasId]);

  useEffect(() => {
    let active = true;
    const loadSiswa = async () => {
      setLoadingSiswa(true);
      try {
        const data = await getUjianSiswaPreview({
          tingkatKelasId: values.kelas_id === "" ? undefined : values.kelas_id,
        });
        if (!active) return;
        // harusnya data dari server sudah di sorting
        setSiswaPreview(data);
      } finally {
        if (active) setLoadingSiswa(false);
      }
    };

    if (selectedKelasId != null) {
      loadSiswa();
    } else {
      setSiswaPreview([]);
    }

    return () => {
      active = false;
    };
  }, [selectedKelasId]);

  useEffect(() => {
    if (selectedBankSoal) {
      setField("jumlah_soal", selectedBankSoal.total_soal);
    } else {
      setField("jumlah_soal", 0);
    }
  }, [selectedBankSoal]);

  useEffect(() => {
    const duration = calculateDuration(
      values.waktu_mulai,
      values.waktu_selesai
    );
    setField("durasi_menit", duration);
  }, [values.waktu_mulai, values.waktu_selesai]);

  const validate = (v: BuatUjianFormValues) => {
    const errors: Partial<Record<keyof BuatUjianFormValues, string>> = {};

    if (!v.nama_ujian.trim()) errors.nama_ujian = "Nama ujian wajib diisi.";
    if (!v.deskripsi_ujian.trim())
      errors.deskripsi_ujian = "Deskripsi ujian wajib diisi.";
    if (!v.tipe_ujian) errors.tipe_ujian = "Tipe ujian wajib dipilih.";
    if (v.kelas_id === "") errors.kelas_id = "Tingkat kelas wajib dipilih.";
    if (!v.bank_soal_id) errors.bank_soal_id = "Bank soal wajib dipilih.";
    if (!v.tanggal_ujian) errors.tanggal_ujian = "Tanggal ujian wajib diisi.";
    if (!v.waktu_mulai) errors.waktu_mulai = "Waktu mulai wajib diisi.";
    if (!v.waktu_selesai) errors.waktu_selesai = "Waktu selesai wajib diisi.";
    if (v.waktu_mulai && v.waktu_selesai && v.durasi_menit <= 0) {
      errors.waktu_selesai = "Waktu selesai harus setelah waktu mulai.";
    }
    if (!v.ruang_ujian_id) errors.ruang_ujian_id = "Ruang ujian wajib dipilih.";
    if (!v.guru_pengawas_id)
      errors.guru_pengawas_id = "Guru pengawas wajib dipilih.";
    if (!v.sesi_id) errors.sesi_id = "Sesi ujian wajib dipilih.";
    if (!v.token_ujian.trim()) errors.token_ujian = "Token ujian wajib diisi.";

    return errors;
  };

  const errors = validate(values);
  const hasError = (name: keyof BuatUjianFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    const nextTouched: Record<string, boolean> = {};
    Object.keys(values).forEach((key) => {
      nextTouched[key] = true;
    });
    setTouched(nextTouched);

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError(
        "Periksa kembali input yang masih kosong atau belum valid."
      );
      return;
    }

    try {
      await submitBuatUjian(values);
      alert("Ujian berhasil dibuat.");
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      } else {
        setSubmitError("Terjadi kesalahan saat menyimpan ujian.");
      }
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-6xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">Buat Ujian</h1>
          <p className="mt-1 text-sm text-slate-500">
            Lengkapi detail ujian, jadwal, serta bank soal yang digunakan.
          </p>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Ujian</h2>
              <p className={helperText}>
                Isi nama ujian, deskripsi, dan tipe soal yang digunakan.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="md:col-span-2">
                <InputField
                  id="nama_ujian"
                  label="Nama Ujian"
                  value={values.nama_ujian}
                  onChange={(v) => setField("nama_ujian", v)}
                  onBlur={() => onBlur("nama_ujian")}
                  placeholder="Contoh: Ujian Tengah Semester Matematika"
                  inputClassName={
                    hasError("nama_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("nama_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.nama_ujian}
                  </p>
                )}
              </div>

              <div className="md:col-span-2">
                <label
                  htmlFor="deskripsi_ujian"
                  className="text-xs font-medium text-slate-600"
                >
                  Deskripsi Ujian
                </label>
                <textarea
                  id="deskripsi_ujian"
                  className={`min-h-[100px] w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] ${
                    hasError("deskripsi_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Jelaskan ujian ini untuk apa dan cakupan materinya."
                  value={values.deskripsi_ujian}
                  onChange={(e) => setField("deskripsi_ujian", e.target.value)}
                  onBlur={() => onBlur("deskripsi_ujian")}
                  required
                />
                {hasError("deskripsi_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.deskripsi_ujian}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="tipe_ujian"
                  className="text-xs font-medium text-slate-600"
                >
                  Tipe Ujian
                </label>
                <select
                  id="tipe_ujian"
                  className={`${selectBaseClass} ${
                    hasError("tipe_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.tipe_ujian}
                  onChange={(e) =>
                    setField("tipe_ujian", e.target.value as TipeUjian)
                  }
                  onBlur={() => onBlur("tipe_ujian")}
                >
                  <option value="PILIHAN_GANDA">Pilihan Ganda</option>
                  <option value="ESSAY">Essay</option>
                  <option value="CAMPURAN">Pilihan Ganda & Essay</option>
                </select>
                {hasError("tipe_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tipe_ujian}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Tingkat Kelas & Bank Soal</h2>
              <p className={helperText}>
                Pilih tingkat kelas, bank soal, dan lihat jumlah soal otomatis.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label
                  htmlFor="kelas_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="kelas_id"
                  className={`${selectBaseClass} ${
                    hasError("kelas_id") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelas_id}
                  onChange={(e) =>
                    setField(
                      "kelas_id",
                      e.target.value === "" ? "" : Number(e.target.value)
                    )
                  }
                  onBlur={() => onBlur("kelas_id")}
                  required
                >
                  <option value="">Pilih tingkat kelas</option>
                  {kelasOptions.map((tingkat) => (
                    <option
                      key={tingkat.id_tingkat_kelas}
                      value={tingkat.id_tingkat_kelas}
                    >
                      Kelas {tingkat.tingkat_kelas}
                    </option>
                  ))}
                </select>
                {hasError("kelas_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.kelas_id}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="bank_soal_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Bank Soal
                </label>
                <select
                  id="bank_soal_id"
                  className={`${selectBaseClass} ${
                    hasError("bank_soal_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.bank_soal_id}
                  onChange={(e) =>
                    setField(
                      "bank_soal_id",
                      e.target.value === "" ? "" : Number(e.target.value)
                    )
                  }
                  onBlur={() => onBlur("bank_soal_id")}
                  disabled={values.kelas_id === "" || loadingBankSoal}
                >
                  <option value="">
                    {loadingBankSoal
                      ? "Memuat bank soal..."
                      : "Pilih bank soal"}
                  </option>
                  {bankSoalOptions.map((bank) => (
                    <option key={bank.id} value={bank.id}>
                      {bank.nama}
                    </option>
                  ))}
                </select>
                {hasError("bank_soal_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.bank_soal_id}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="jumlah_soal"
                  type="number"
                  label="Jumlah Soal"
                  value={String(values.jumlah_soal)}
                  onChange={(v) => setField("jumlah_soal", Number(v))}
                  onBlur={() => onBlur("jumlah_soal")}
                  inputClassName="bg-slate-50 text-slate-500"
                  disabled
                  required
                />
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Jadwal Ujian</h2>
              <p className={helperText}>
                Tentukan tanggal, waktu mulai, dan waktu berakhir ujian.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
              <div>
                <InputField
                  id="tanggal_ujian"
                  type="date"
                  label="Tanggal Ujian"
                  value={values.tanggal_ujian}
                  onChange={(v) => setField("tanggal_ujian", v)}
                  onBlur={() => onBlur("tanggal_ujian")}
                  inputClassName={
                    hasError("tanggal_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("tanggal_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.tanggal_ujian}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="waktu_mulai"
                  type="time/datetime-local"
                  label="Waktu Mulai"
                  value={values.waktu_mulai}
                  onChange={(v) => setField("waktu_mulai", v)}
                  onBlur={() => onBlur("waktu_mulai")}
                  inputClassName={
                    hasError("waktu_mulai")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("waktu_mulai") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.waktu_mulai}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="waktu_selesai"
                  type="time/datetime-local"
                  label="Waktu Selesai"
                  value={values.waktu_selesai}
                  onChange={(v) => setField("waktu_selesai", v)}
                  onBlur={() => onBlur("waktu_selesai")}
                  inputClassName={
                    hasError("waktu_selesai")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("waktu_selesai") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.waktu_selesai}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="durasi_menit"
                  type="number"
                  label="Durasi (menit)"
                  value={String(values.durasi_menit)}
                  onChange={(v) => setField("durasi_menit", Number(v))}
                  onBlur={() => onBlur("durasi_menit")}
                  inputClassName="bg-slate-50 text-slate-500"
                  disabled
                  required
                />
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Ruang, Sesi, Pengawas</h2>
              <p className={helperText}>
                Atur ruang ujian, sesi, dan guru pengawas.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div>
                <label
                  htmlFor="ruang_ujian_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Ruang Ujian
                </label>
                <select
                  id="ruang_ujian_id"
                  className={`${selectBaseClass} ${
                    hasError("ruang_ujian_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.ruang_ujian_id}
                  onChange={(e) =>
                    setField(
                      "ruang_ujian_id",
                      e.target.value === "" ? "" : Number(e.target.value)
                    )
                  }
                  onBlur={() => onBlur("ruang_ujian_id")}
                  required
                >
                  <option value="">Pilih ruang ujian</option>
                  {ruangOptions.map((ruang) => (
                    <option key={ruang.id} value={ruang.id}>
                      {ruang.namaRuangan}
                    </option>
                  ))}
                </select>
                {hasError("ruang_ujian_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.ruang_ujian_id}
                  </p>
                )}
              </div>

              <div>
                <label
                  htmlFor="sesi_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Sesi Ujian
                </label>
                <select
                  id="sesi_id"
                  className={`${selectBaseClass} ${
                    hasError("sesi_id") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.sesi_id}
                  onChange={(e) =>
                    setField(
                      "sesi_id",
                      e.target.value === "" ? "" : Number(e.target.value)
                    )
                  }
                  onBlur={() => onBlur("sesi_id")}
                  required
                >
                  <option value="">Pilih sesi ujian</option>
                  {sesiOptions.map((sesi) => (
                    <option key={sesi.id} value={sesi.id}>
                      {sesi.kode} - {sesi.nama}
                    </option>
                  ))}
                </select>
                {hasError("sesi_id") && (
                  <p className="mt-1 text-xs text-rose-600">{errors.sesi_id}</p>
                )}
              </div>

              <div>
                <label
                  htmlFor="guru_pengawas_id"
                  className="text-xs font-medium text-slate-600"
                >
                  Guru Pengawas
                </label>
                <select
                  id="guru_pengawas_id"
                  className={`${selectBaseClass} ${
                    hasError("guru_pengawas_id")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  value={values.guru_pengawas_id}
                  onChange={(e) =>
                    setField(
                      "guru_pengawas_id",
                      e.target.value === "" ? "" : Number(e.target.value)
                    )
                  }
                  onBlur={() => onBlur("guru_pengawas_id")}
                  required
                >
                  <option value="">Pilih guru pengawas</option>
                  {guruOptions.map((guru) => (
                    <option key={guru.id} value={guru.id}>
                      {guru.nama} {guru.mapel ? `- ${guru.mapel}` : ""}
                    </option>
                  ))}
                </select>
                {hasError("guru_pengawas_id") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.guru_pengawas_id}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Keamanan & Token</h2>
              <p className={helperText}>
                Atur pengacakan soal dan token ujian.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label
                  htmlFor="acak_soal"
                  className="text-xs font-medium text-slate-600"
                >
                  Acak Soal
                </label>
                <select
                  id="acak_soal"
                  className={selectBaseClass}
                  value={values.acak_soal ? "ya" : "tidak"}
                  onChange={(e) =>
                    setField("acak_soal", e.target.value === "ya")
                  }
                  onBlur={() => onBlur("acak_soal")}
                >
                  <option value="ya">Ya, acak soal</option>
                  <option value="tidak">Tidak, urutan tetap</option>
                </select>
              </div>

              <div>
                <InputField
                  id="token_ujian"
                  label="Token Ujian"
                  value={values.token_ujian}
                  onChange={(v) => setField("token_ujian", v)}
                  onBlur={() => onBlur("token_ujian")}
                  placeholder="Masukkan token ujian"
                  inputClassName={
                    hasError("token_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("token_ujian") && (
                  <p className="mt-1 text-xs text-rose-600">
                    {errors.token_ujian}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="mb-4 flex items-center justify-between">
                <div>
                  <h2 className={sectionTitle}>Preview Daftar Siswa</h2>
                  <p className={helperText}>
                    Siswa terdaftar berdasarkan tingkat kelas yang dipilih.
                  </p>
                </div>
                <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
                  {siswaPreview.length} siswa
                </span>
              </div>

              {values.kelas_id === "" && (
                <p className="text-sm text-slate-500">
                  Pilih tingkat kelas untuk melihat daftar siswa.
                </p>
              )}

              {values.kelas_id !== "" && loadingSiswa && (
                <p className="text-sm text-slate-500">Memuat data siswa...</p>
              )}

              {values.kelas_id !== "" &&
                !loadingSiswa &&
                siswaPreview.length === 0 && (
                  <p className="text-sm text-slate-500">
                    Belum ada siswa pada tingkat kelas ini.
                  </p>
                )}

              {siswaPreview.length > 0 && (
                <div className="max-h-64 overflow-y-auto rounded-lg border border-slate-200">
                  <table className="w-full text-left text-xs text-slate-600">
                    <thead className="bg-slate-50 text-[11px] uppercase text-slate-400">
                      <tr>
                        <th className="px-3 py-2 font-medium">Nama</th>
                        <th className="px-3 py-2 font-medium">Absen</th>
                        <th className="px-3 py-2 font-medium">Kelas</th>
                        <th className="px-3 py-2 font-medium">Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      {siswaPreview.map((siswa) => (
                        <tr
                          key={siswa.id}
                          className="border-t border-slate-100"
                        >
                          <td className="px-3 py-2 font-medium text-slate-700">
                            {siswa.nama}
                          </td>
                          <td className="px-3 py-2">{siswa.no_absen}</td>
                          <td className="px-3 py-2">{siswa.kelas}</td>
                          <td className="px-3 py-2 capitalize">
                            {siswa.status_akun}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="mb-4">
                <h2 className={sectionTitle}>Preview Bank Soal</h2>
                <p className={helperText}>
                  Detail bank soal yang dipilih akan tampil di sini.
                </p>
              </div>

              {!selectedBankSoal && (
                <p className="text-sm text-slate-500">
                  Pilih bank soal untuk melihat detailnya.
                </p>
              )}

              {selectedBankSoal && (
                <div className="space-y-3 text-sm text-slate-600">
                  <div>
                    <p className="text-xs text-slate-400">Nama Bank Soal</p>
                    <p className="font-semibold text-slate-800">
                      {selectedBankSoal.nama}
                    </p>
                  </div>
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <p className="text-xs text-slate-400">Mata Pelajaran</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.mata_pelajaran ?? "-"}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400">Materi</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.materi ?? "-"}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400">Jumlah PG</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.jumlah_soal_pg ?? 0}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-slate-400">Jumlah Essay</p>
                      <p className="font-medium text-slate-700">
                        {selectedBankSoal.jumlah_soal_essay ?? 0}
                      </p>
                    </div>
                  </div>
                  <div className="rounded-lg bg-slate-50 px-3 py-2 text-sm font-semibold text-slate-700">
                    Total Soal: {selectedBankSoal.total_soal}
                  </div>
                </div>
              )}
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
              className="inline-flex items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({});
                setSubmitError(null);
                setBankSoalOptions([]);
                setSiswaPreview([]);
              }}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43]"
            >
              Simpan Ujian
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
