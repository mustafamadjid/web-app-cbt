import React from "react";
import { ArrowRight, CalendarDays, Award, CheckCircle2, XCircle } from "lucide-react";
import type { UjianSiswaResultItem } from "@/types/Ujian/ujianSiswa";
import { Link } from "react-router";
import { paths } from "@/routes/paths";

const UjianResultCard: React.FC<{ item: UjianSiswaResultItem }> = ({ item }) => {
  const detailPath = paths.dashboard.hasil_ujian_detail_siswa.replace(
    ":id",
    String(item.id)
  );

  return (
    <Link
      to={detailPath}
      className="group flex h-full flex-col justify-between rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition hover:border-[#397e50] hover:shadow-lg hover:shadow-[#397e50]/5"
    >
      <div className="space-y-4">
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-gray-400">
            Hasil Ujian
          </p>
          <h3 className="text-lg font-bold text-[#37513d]">
            {item.nama_ujian}
          </h3>
          <p className="text-sm font-medium text-[#397e50]">{item.mapel}</p>
        </div>

        <div className="space-y-2 text-sm text-gray-500">
          <div className="flex items-center gap-2">
            <CalendarDays className="h-4 w-4 text-[#397e50]" />
            <span>{item.tgl_ujian}</span>
          </div>
          <div className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 text-[#397e50]" />
            <span>{item.jumlah_benar} benar</span>
          </div>
          <div className="flex items-center gap-2">
            <XCircle className="h-4 w-4 text-[#397e50]" />
            <span>{item.jumlah_salah} salah</span>
          </div>
          <div className="flex items-center gap-2">
            <Award className="h-4 w-4 text-[#397e50]" />
            <span>Total nilai: {item.nilai}</span>
          </div>
        </div>
      </div>

      <div className="mt-5 flex items-center justify-between text-sm font-semibold text-[#397e50]">
        <span>Lihat detail</span>
        <ArrowRight className="h-4 w-4 transition group-hover:translate-x-1" />
      </div>
    </Link>
  );
};

export default UjianResultCard;
