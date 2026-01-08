import { BoxJadwalUjian } from "@/components/features/Ujian/BoxJadwalUjian";

import { getTingkatKelasOptions } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { getJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { useEffect, useMemo, useRef, useState } from "react";

const STATUS_SECTIONS = [
  { key: "belum_dimulai", label: "Belum Mulai" },
  { key: "berlangsung", label: "Berlangsung" },
  { key: "selesai", label: "Selesai" },
];

const monthMap: Record<string, string> = {
  januari: "01",
  februari: "02",
  maret: "03",
  april: "04",
  mei: "05",
  juni: "06",
  juli: "07",
  agustus: "08",
  september: "09",
  oktober: "10",
  november: "11",
  desember: "12",
};

const normalize = (value: string) => value.toLowerCase().trim();

const formatTanggalToIso = (value?: string) => {
  if (!value) return null;
  const parts = value.split(",").map((item) => item.trim());
  if (parts.length < 2) return null;

  const dateParts = parts[1].split(" ").filter(Boolean);
  if (dateParts.length < 3) return null;

  const [day, monthName, year] = dateParts;
  const month = monthMap[normalize(monthName)];
  if (!month) return null;

  const paddedDay = day.padStart(2, "0");
  return `${year}-${month}-${paddedDay}`;
};

export const JadwalUjian = () => {
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  // Daftar ujian terjadwal
  const [daftarJadwalUjian, setDaftarJadwalUjian] = useState<JadwalUjianItem[]>([]);
  const [tingkatKelasOptions, setTingkatKelasOptions] = useState<number[]>([]);
  const [ruangOptions, setRuangOptions] = useState<RuangUjianRow[]>([]);

  const [searchTerm, setSearchTerm] = useState("");
  const [selectedDate, setSelectedDate] = useState("");
  const [selectedTingkat, setSelectedTingkat] = useState("");
  const [selectedRuang, setSelectedRuang] = useState("");

  const requestSeq = useRef(0);

  // Fetch data jadwal ujian
  useEffect(() => {
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");

        const [data, tingkatOptions, ruangUjianOptions] = await Promise.all([
          getJadwalUjian(),
          getTingkatKelasOptions(),
          getRuangUjian(),
        ]);

        if (seq !== requestSeq.current) return;

        setDaftarJadwalUjian(data);
        setTingkatKelasOptions(tingkatOptions);
        setRuangOptions(ruangUjianOptions);
      } catch (e) {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat data jadwal ujian.");
        setDaftarJadwalUjian([]);
        setTingkatKelasOptions([]);
        setRuangOptions([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, []);

  const filteredJadwal = useMemo(() => {
    const query = normalize(searchTerm);

    return daftarJadwalUjian.filter((ujian) => {
      if (selectedTingkat && ujian.tingkat_kelas !== Number(selectedTingkat)) {
        return false;
      }

      if (selectedRuang && ujian.ruang_ujian !== selectedRuang) {
        return false;
      }

      if (selectedDate) {
        const isoDate = formatTanggalToIso(ujian.tgl_ujian);
        if (!isoDate || isoDate !== selectedDate) return false;
      }

      if (!query) return true;

      const searchable = normalize(
        [
          ujian.nama_ujian,
          ujian.pengawas_ujian,
          ujian.tgl_ujian,
          ujian.waktu_mulai,
          ujian.sesi_ujian ? String(ujian.sesi_ujian) : "",
          ujian.ruang_ujian ?? "",
          ujian.status_ujian ?? "",
          ujian.tingkat_kelas ? String(ujian.tingkat_kelas) : "",
          ujian.nama_kelas ?? "",
        ]
          .filter(Boolean)
          .join(" ")
      );

      return searchable.includes(query);
    });
  }, [daftarJadwalUjian, searchTerm, selectedDate, selectedRuang, selectedTingkat]);

  const groupedJadwal = useMemo(() => {
    const grouped: Record<string, JadwalUjianItem[]> = {
      belum_dimulai: [],
      berlangsung: [],
      selesai: [],
    };

    filteredJadwal.forEach((ujian) => {
      const status = ujian.status_ujian ?? "";
      if (grouped[status]) {
        grouped[status].push(ujian);
      }
    });

    return grouped;
  }, [filteredJadwal]);

  return (
    <>
      <div className="px-8 py-10">
        <div className="flex flex-col gap-8">
          <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
              <div className="flex flex-col gap-2">
                <label className="text-sm font-semibold text-slate-600">
                  Cari Jadwal
                </label>
                <input
                  type="text"
                  placeholder="Cari ujian, pengawas, kelas, ruangan..."
                  value={searchTerm}
                  onChange={(event) => setSearchTerm(event.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-4 py-2 text-sm text-slate-700 shadow-sm focus:border-[#397e50] focus:outline-none focus:ring-2 focus:ring-[#397e50]/20"
                />
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-sm font-semibold text-slate-600">
                  Tanggal Ujian
                </label>
                <input
                  type="date"
                  value={selectedDate}
                  onChange={(event) => setSelectedDate(event.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-4 py-2 text-sm text-slate-700 shadow-sm focus:border-[#397e50] focus:outline-none focus:ring-2 focus:ring-[#397e50]/20"
                />
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-sm font-semibold text-slate-600">
                  Tingkat Kelas
                </label>
                <select
                  value={selectedTingkat}
                  onChange={(event) => setSelectedTingkat(event.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-4 py-2 text-sm text-slate-700 shadow-sm focus:border-[#397e50] focus:outline-none focus:ring-2 focus:ring-[#397e50]/20"
                >
                  <option value="">Semua Tingkat</option>
                  {tingkatKelasOptions.map((tingkat) => (
                    <option key={tingkat} value={tingkat}>
                      Kelas {tingkat}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-sm font-semibold text-slate-600">
                  Ruang Ujian
                </label>
                <select
                  value={selectedRuang}
                  onChange={(event) => setSelectedRuang(event.target.value)}
                  className="w-full rounded-xl border border-slate-200 px-4 py-2 text-sm text-slate-700 shadow-sm focus:border-[#397e50] focus:outline-none focus:ring-2 focus:ring-[#397e50]/20"
                >
                  <option value="">Semua Ruangan</option>
                  {ruangOptions.map((ruang) => (
                    <option key={ruang.id} value={ruang.namaRuangan}>
                      {ruang.namaRuangan}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>

          {errorMsg ? (
            <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
              {errorMsg}
            </div>
          ) : null}

          {loading ? (
            <div className="rounded-2xl border border-slate-200 bg-white px-4 py-6 text-center text-sm text-slate-500">
              Memuat jadwal ujian...
            </div>
          ) : null}

          {!loading && !errorMsg && filteredJadwal.length === 0 ? (
            <div className="rounded-2xl border border-slate-200 bg-white px-4 py-6 text-center text-sm text-slate-500">
              Tidak ada jadwal ujian yang sesuai dengan filter.
            </div>
          ) : null}

          {!loading && !errorMsg
            ? STATUS_SECTIONS.map((section) => (
                <div key={section.key} className="flex flex-col gap-4">
                  <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-slate-800">
                      {section.label}
                    </h2>
                    <span className="text-sm text-slate-500">
                      {groupedJadwal[section.key]?.length ?? 0} ujian
                    </span>
                  </div>
                  <div className="flex flex-col gap-5">
                    {(groupedJadwal[section.key] ?? []).map((ujian) => (
                      <BoxJadwalUjian
                        key={ujian.id}
                        id={ujian.id}
                        nama_ujian={ujian.nama_ujian}
                        nama_kelas={ujian.nama_kelas}
                        tingkat_kelas={ujian.tingkat_kelas}
                        pengawas_ujian={ujian.pengawas_ujian}
                        tgl_ujian={ujian.tgl_ujian}
                        waktu_mulai={ujian.waktu_mulai}
                        sesi_ujian={ujian.sesi_ujian}
                        ruang_ujian={ujian.ruang_ujian}
                        status_ujian={ujian.status_ujian}
                      />
                    ))}
                  </div>
                  {(groupedJadwal[section.key] ?? []).length === 0 ? (
                    <div className="rounded-2xl border border-dashed border-slate-200 bg-slate-50 px-4 py-6 text-center text-sm text-slate-500">
                      Tidak ada ujian untuk status ini.
                    </div>
                  ) : null}
                </div>
              ))
            : null}
        </div>
      </div>
    </>
  );
};
