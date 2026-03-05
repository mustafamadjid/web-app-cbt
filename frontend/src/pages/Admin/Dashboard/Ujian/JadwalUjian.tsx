import { useEffect, useMemo, useState } from "react";
import { Calendar, CalendarRange, Layers, MapPin, Search } from "lucide-react";

import BoxJadwalUjian from "@/components/features/Ujian/BoxJadwalUjian";
import { useAuth } from "@/contexts/AuthContext";
import { tahunOption } from "@/helper/TahunOption/TahunOption";
import { paths } from "@/routes/paths";
import { useGetDataKelasFull } from "@/services/Api/features-api/DataMaster/kelas.service";
import { useGetRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import {
  updateStatusUjian,
  useGetJadwalUjian,
} from "@/services/Api/features-api/Ujian/jadwalujian.service";
import type { TingkatKelas } from "@/types/DataMaster/Kelas";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { JadwalUjianItem, JadwalUjianStatusClient } from "@/types/Ujian/jadwalUjian";

const STATUS_SECTIONS: Array<{
  key: JadwalUjianStatusClient;
  label: string;
}> = [
  { key: "berlangsung", label: "Berlangsung" },
  { key: "belum_dimulai", label: "Belum Mulai" },
  { key: "selesai", label: "Selesai" },
  { key: "dibatalkan", label: "Dibatalkan" },
];

function useDebouncedValue<T>(value: T, delayMs: number) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const timer = window.setTimeout(() => setDebounced(value), delayMs);
    return () => window.clearTimeout(timer);
  }, [value, delayMs]);
  return debounced;
}

const JadwalUjian = () => {
  const { user } = useAuth();

  const [activeTab, setActiveTab] = useState<JadwalUjianStatusClient>("berlangsung");
  const [searchTerm, setSearchTerm] = useState("");
  const debouncedSearchTerm = useDebouncedValue(searchTerm, 300);
  const [selectedDate, setSelectedDate] = useState("");
  const [selectedTingkatId, setSelectedTingkatId] = useState<number | null>(null);
  const [selectedRuang, setSelectedRuang] = useState<number | null>(null);
  const [selectedTahun, setSelectedTahun] = useState<string | null>(null);
  const [updatingIdUjian, setUpdatingIdUjian] = useState<number | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  const { data: kelasData } = useGetDataKelasFull();
  const { data: ruangData } = useGetRuangUjian();

  const tingkatKelasOptions: TingkatKelas[] = kelasData?.item_tingkat_kelas ?? [];
  const ruangOptions: RuangUjianRow[] = ruangData ?? [];
  const tahunOptions = useMemo(() => tahunOption().map((item) => String(item)), []);

  const {
    data: jadwalDataRaw,
    loading,
    error,
    refetch,
  } = useGetJadwalUjian({
    search: debouncedSearchTerm.trim() || undefined,
    tanggal: selectedDate || undefined,
    ruangUjianId: selectedRuang ?? undefined,
    tingkatKelasId: selectedTingkatId ?? undefined,
    tahun: selectedTahun || undefined,
  });

  const jadwalData: JadwalUjianItem[] = jadwalDataRaw ?? [];

  const groupedJadwal = useMemo(() => {
    const grouped: Record<JadwalUjianStatusClient, JadwalUjianItem[]> = {
      belum_dimulai: [],
      berlangsung: [],
      selesai: [],
      dibatalkan: [],
    };
    jadwalData.forEach((item) => {
      grouped[item.status_ujian].push(item);
    });
    return grouped;
  }, [jadwalData]);

  const canControlUjian = (ujian: JadwalUjianItem) => {
    if (!user) return false;
    if (user.role === "ADMIN") return true;
    if (user.role !== "GURU") return false;
    return ujian.id_guru === user.id_pengguna || ujian.id_pengawas === user.id_pengguna;
  };

  const handleStatusUpdate = async (
    idUjian: number,
    nextStatus: JadwalUjianStatusClient,
  ) => {
    setActionError(null);
    setUpdatingIdUjian(idUjian);
    try {
      await updateStatusUjian(idUjian, nextStatus);
      await refetch();
    } catch (updateError) {
      if (updateError instanceof Error) {
        setActionError(updateError.message);
      } else {
        setActionError("Gagal memperbarui status ujian.");
      }
    } finally {
      setUpdatingIdUjian(null);
    }
  };

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-8">
      <div className="flex flex-col gap-8">
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <h1 className="mb-6 text-2xl font-bold text-slate-800">Jadwal Ujian</h1>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Search size={16} /> Cari Jadwal
              </label>
              <input
                type="text"
                placeholder="Cari ujian, pengawas..."
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              />
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <CalendarRange size={16} /> Tanggal
              </label>
              <input
                type="date"
                value={selectedDate}
                onChange={(event) => setSelectedDate(event.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              />
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Calendar size={16} /> Tahun
              </label>
              <select
                value={selectedTahun ?? ""}
                onChange={(event) =>
                  setSelectedTahun(event.target.value === "" ? null : event.target.value)
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

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <Layers size={16} /> Kelas
              </label>
              <select
                value={selectedTingkatId ?? ""}
                onChange={(event) =>
                  setSelectedTingkatId(
                    event.target.value === "" ? null : Number(event.target.value),
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Kelas</option>
                {tingkatKelasOptions.map((kelas) => (
                  <option key={kelas.id_tingkat_kelas} value={kelas.id_tingkat_kelas}>
                    Kelas {kelas.tingkat_kelas}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-600">
                <MapPin size={16} /> Ruangan
              </label>
              <select
                value={selectedRuang ?? ""}
                onChange={(event) =>
                  setSelectedRuang(
                    event.target.value === "" ? null : Number(event.target.value),
                  )
                }
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm transition-all focus:border-[#397e50] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#397e50]/10"
              >
                <option value="">Semua Ruangan</option>
                {ruangOptions.map((ruang) => (
                  <option key={ruang.id_ruangan} value={ruang.id_ruangan}>
                    {ruang.nama_ruangan}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2 rounded-2xl bg-slate-100 p-1.5 sm:w-fit">
          {STATUS_SECTIONS.map((section) => {
            const isActive = activeTab === section.key;
            const count = groupedJadwal[section.key].length;
            return (
              <button
                key={section.key}
                onClick={() => setActiveTab(section.key)}
                className={`flex cursor-pointer items-center gap-2 rounded-xl px-6 py-2.5 text-sm font-bold transition-all ${
                  isActive
                    ? "bg-white text-[#397e50] shadow-sm"
                    : "text-slate-500 hover:bg-slate-200 hover:text-slate-700"
                }`}
              >
                {section.label}
                <span
                  className={`ml-1 rounded-md px-1.5 py-0.5 text-2xs ${
                    isActive ? "bg-[#397e50] text-white" : "bg-slate-200 text-slate-600"
                  }`}
                >
                  {count}
                </span>
              </button>
            );
          })}
        </div>

        <div className="min-h-[300px]">
          {(error || actionError) && (
            <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
              {actionError || error}
            </div>
          )}

          {loading ? (
            <div className="flex h-40 items-center justify-center rounded-2xl border border-dashed border-slate-300 text-slate-500">
              <div className="flex flex-col items-center gap-2">
                <div className="h-6 w-6 animate-spin rounded-full border-2 border-[#397e50] border-t-transparent" />
                <span>Memuat data...</span>
              </div>
            </div>
          ) : (
            <div className="flex flex-col gap-5">
              {groupedJadwal[activeTab].length > 0 ? (
                groupedJadwal[activeTab].map((item) => (
                  <BoxJadwalUjian
                    key={item.id}
                    {...item}
                    onStart={(idUjian) => handleStatusUpdate(idUjian, "berlangsung")}
                    onCancel={(idUjian) => handleStatusUpdate(idUjian, "dibatalkan")}
                    canControl={canControlUjian(item)}
                    updating={item.id_ujian != null && updatingIdUjian === item.id_ujian}
                    linkJadwal={paths.dashboard.detail_ujian.replace(":id", String(item.id))}
                  />
                ))
              ) : (
                <div className="flex h-40 flex-col items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white text-slate-500">
                  <p className="font-medium">Tidak ada jadwal ujian</p>
                  <p className="text-xs">
                    Ujian dengan status "{activeTab.replace("_", " ")}" tidak ditemukan.
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

export default JadwalUjian;
