import { useEffect, useMemo, useRef, useState } from "react";
import BoxHasilUjian from "@/components/features/Ujian/BoxHasilUjian";
import { getHasilUjianList } from "@/services/Api/features-api/Ujian/hasilUjian.service";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { paths } from "@/routes/paths";

const HasilUjian = () => {
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [daftarUjian, setDaftarUjian] = useState<JadwalUjianItem[]>([]);
  const requestSeq = useRef(0);

  useEffect(() => {
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getHasilUjianList();
        if (seq !== requestSeq.current) return;
        setDaftarUjian(data);
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data hasil ujian.");
        setDaftarUjian([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, []);

  const daftarSelesai = useMemo(
    () => daftarUjian.filter((ujian) => ujian.status_ujian === "selesai"),
    [daftarUjian]
  );

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="flex flex-col gap-6">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold text-slate-800">Hasil Ujian</h1>
              <p className="text-sm text-slate-500">
                Daftar ujian yang sudah selesai beserta akses hasilnya.
              </p>
            </div>
            <span className="inline-flex items-center rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">
              Total selesai: {daftarSelesai.length}
            </span>
          </div>
        </div>

        {errorMsg ? (
          <div className="rounded-2xl border border-rose-200 bg-rose-50 p-6 text-sm font-semibold text-rose-700">
            {errorMsg}
          </div>
        ) : null}

        {loading ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500">
            Memuat daftar hasil ujian...
          </div>
        ) : null}

        {!loading && daftarSelesai.length === 0 ? (
          <div className="rounded-2xl border border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
            Belum ada ujian selesai yang memiliki hasil.
          </div>
        ) : null}

        <div className="grid gap-6 lg:grid-cols-2">
          {daftarSelesai.map((ujian) => (
            <BoxHasilUjian
              key={ujian.id}
              {...ujian}
              linkHasil={paths.dashboard.hasil_ujian_detail.replace(
                ":id",
                String(ujian.id)
              )}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

export default HasilUjian;
