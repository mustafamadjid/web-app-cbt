import { useEffect, useMemo, useRef, useState } from "react";
import BoxHasilUjian from "@/components/features/Ujian/BoxHasilUjian";
import {
  getHasilUjianList,
} from "@/services/Api/features-api/Ujian/hasilUjian.service";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { paths } from "@/routes/paths";
import type { TingkatKelasOption } from "@/types/DataMaster/Kelas";
import { tahunOption } from "@/helper/TahunOption/TahunOption";
import { Calendar, Layers } from "lucide-react";

const HasilUjian = () => {
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [daftarUjian, setDaftarUjian] = useState<JadwalUjianItem[]>([]);
  const [tingkatKelasOptions, setTingkatKelasOptions] = useState<
    TingkatKelasOption[]
  >([]);
  const [tahunOptions, setTahunOptions] = useState<string[]>([]);
  const [selectedTingkatId, setSelectedTingkatId] = useState<number | null>(
    null
  );
  const [selectedTahun, setSelectedTahun] = useState<string | null>(null);
  const requestSeq = useRef(0);

  useEffect(() => {
    (async () => {
      try {
        const [options, tahunOptionsData] = await Promise.all([
          getTingkatKelasOptions(),
          tahunOption(),
        ]);
        setTingkatKelasOptions(options);
        setTahunOptions(tahunOptionsData);
      } catch {
        setTingkatKelasOptions([]);
        setTahunOptions([]);
      }
    })();
  }, []);

  useEffect(() => {
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getHasilUjianList({
          tingkatKelasId: selectedTingkatId ?? undefined,
          tahun: selectedTahun ?? undefined,
        });
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
  }, [selectedTingkatId, selectedTahun]);

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
          <div className="mt-6 grid gap-4 sm:max-w-md sm:grid-cols-2">
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Layers size={16} /> Kelas
              </label>
              <select
                value={selectedTingkatId ?? ""}
                onChange={(e) =>
                  setSelectedTingkatId(
                    e.target.value === "" ? null : Number(e.target.value)
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Kelas</option>
                {tingkatKelasOptions.map((tingkat) => (
                  <option
                    key={tingkat.id_tingkat_kelas}
                    value={tingkat.id_tingkat_kelas}
                  >
                    Kelas {tingkat.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Calendar size={16} /> Tahun
              </label>
              <select
                value={selectedTahun ?? ""}
                onChange={(e) =>
                  setSelectedTahun(
                    e.target.value === "" ? null : String(e.target.value)
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Tahun</option>
                {tahunOptions.map((tahun) => (
                  <option key={tahun} value={tahun}>
                    {tahun}
                  </option>
                ))}
              </select>
            </div>
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
