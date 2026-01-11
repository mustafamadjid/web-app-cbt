import { BoxCetakUjian } from "@/components/features/Ujian/BoxCetakUjian";
import { PrintButton } from "@/components/common/Input/PrintButton";
import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { TingkatKelasOption } from "@/types/DataMaster/Kelas";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { Calendar, Layers } from "lucide-react";
import { useEffect, useRef, useState } from "react";



type PrintJenis = "daftar-hadir" | "berita-acara" | "kartu-peserta";

export const Cetak = () => {
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [daftarUjian, setDaftarUjian] = useState<JadwalUjianItem[]>([]);
  const [tingkatKelasOptions, setTingkatKelasOptions] = useState<
    TingkatKelasOption[]
  >([]);

  const [selectedDate, setSelectedDate] = useState("");
  const [selectedTingkatId, setSelectedTingkatId] = useState<number | null>(
    null
  );

  const requestSeq = useRef(0);

  useEffect(() => {
    (async () => {
      try {
        const options = await getTingkatKelasOptions();
        setTingkatKelasOptions(options);
      } catch {
        setTingkatKelasOptions([]);
      }
    })();
  }, []);

  useEffect(() => {
    const seq = ++requestSeq.current;
    const tingkatKelas =
      selectedTingkatId != null ? selectedTingkatId : undefined;

    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getJadwalUjian({
          tanggal: selectedDate || undefined,
          tingkatKelasId: tingkatKelas,
        });
        if (seq !== requestSeq.current) return;
        setDaftarUjian(data);
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data ujian untuk cetak.");
        setDaftarUjian([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, [selectedDate, selectedTingkatId]);

  const handlePrint = (jenis: PrintJenis, ujian: JadwalUjianItem) => {
    console.info("Print", jenis, ujian);
    if (typeof window !== "undefined") {
      window.print();
    }
  };

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="flex flex-col gap-8">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl font-bold text-slate-800">
                Cetak Dokumen Ujian
              </h1>
              <p className="text-sm text-slate-500">
                Pilih ujian untuk mencetak daftar hadir, berita acara, dan kartu
                peserta.
              </p>
            </div>
          </div>

          <div className="mt-6 grid gap-6 md:grid-cols-2">
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Calendar size={16} /> Tanggal
              </label>
              <input
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              />
            </div>
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Layers size={16} /> Tingkat Kelas
              </label>
              <select
                value={selectedTingkatId ?? ""}
                onChange={(e) =>{
                  setSelectedTingkatId(
                    e.target.value === "" ? null : Number(e.target.value)
                  )
                }
                  
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Tingkat</option>
                {tingkatKelasOptions.map((option) => (
                  <option
                    key={option.id_tingkat_kelas}
                    value={option.id_tingkat_kelas}
                  >
                    Kelas {option.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>

        <div className="min-h-[300px]">
          {errorMsg && (
            <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
              {errorMsg}
            </div>
          )}

          {loading ? (
            <div className="flex h-40 items-center justify-center rounded-2xl border border-dashed border-slate-300 text-slate-500">
              <div className="flex flex-col items-center gap-2">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-[#397e50] border-t-transparent" />
                <span>Memuat data ujian...</span>
              </div>
            </div>
          ) : (
            <div className="flex flex-col gap-5">
              {daftarUjian.length > 0 ? (
                daftarUjian.map((ujian) => (
                  <BoxCetakUjian
                    key={ujian.id}
                    {...ujian}
                    actions={
                      <>
                        <PrintButton
                          label="Daftar Hadir"
                          onClick={() => handlePrint("daftar-hadir", ujian)}
                        />
                        <PrintButton
                          label="Berita Acara"
                          variant="outline"
                          onClick={() => handlePrint("berita-acara", ujian)}
                        />
                        <PrintButton
                          label="Kartu Peserta"
                          variant="outline"
                          onClick={() => handlePrint("kartu-peserta", ujian)}
                        />
                      </>
                    }
                  />
                ))
              ) : (
                <div className="flex h-40 flex-col items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white text-slate-500">
                  <p className="font-medium">Belum ada data ujian</p>
                  <p className="text-xs">
                    Ubah filter untuk menemukan ujian yang ingin dicetak.
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
