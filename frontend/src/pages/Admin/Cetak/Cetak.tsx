import { Link } from "react-router";

import { paths } from "@/routes/paths";

const Cetak = () => {
  return (
    <div className="mx-auto max-w-3xl px-4 py-10 sm:px-8">
      <div className="rounded-2xl border border-slate-200 bg-white p-6 text-center shadow-sm">
        <h1 className="text-2xl font-bold text-slate-800">
          Cetak Dokumen Ujian
        </h1>
        <p className="mt-2 text-sm text-slate-500">
          Cetak daftar hadir, berita acara, dan kartu peserta sekarang tersedia
          langsung pada detail ujian di halaman jadwal.
        </p>
        <Link
          to={paths.dashboard.jadwal_ujian}
          className="mt-6 inline-flex items-center justify-center rounded-xl bg-[#397e50] px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-[#2f6842] focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#397e50]/30"
        >
          Buka Jadwal Ujian
        </Link>
      </div>
    </div>
  );
};

export default Cetak;
