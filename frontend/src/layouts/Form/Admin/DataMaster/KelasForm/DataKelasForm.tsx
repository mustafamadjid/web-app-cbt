import React from "react";

import TambahNamaKelasForm from "./TambahNamaKelasForm";
import TambahTingkatKelasForm from "./TambahTingkatKelasForm";

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const DataKelasForm = () => {
  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4 space-y-6">
        <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">Data Kelas</h1>
          <p className="mt-1 text-sm text-slate-500">
            Tambahkan tingkat kelas dan nama kelas pada section yang tersedia.
          </p>
        </div>

        <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="mb-4">
            <h2 className={sectionTitle}>Tambah Tingkat Kelas</h2>
            <p className={helperText}>
              Isi satu field tingkat kelas dalam format angka.
            </p>
          </div>
          <TambahTingkatKelasForm />
        </section>

        <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="mb-4">
            <h2 className={sectionTitle}>Tambah Nama Kelas</h2>
            <p className={helperText}>
              Pilih tingkat kelas lalu isi nama kelas.
            </p>
          </div>
          <TambahNamaKelasForm />
        </section>
      </div>
    </div>
  );
};

export default DataKelasForm;
