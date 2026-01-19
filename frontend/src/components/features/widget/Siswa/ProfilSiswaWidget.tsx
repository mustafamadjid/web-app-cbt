import React from "react";
import { GraduationCap, Users } from "lucide-react";
import type { SiswaProfile } from "@/types/Widget/SiswaDashboard";

type ProfilSiswaWidgetProps = {
  profil: SiswaProfile;
  className?: string;
};

const ProfilSiswaWidget: React.FC<ProfilSiswaWidgetProps> = ({
  profil,
  className,
}) => {
  return (
    <section
      className={[
        "flex h-full flex-col gap-4 rounded-xl border border-gray-200 bg-white p-5 shadow-sm",
        "transition-all duration-300 hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      <header className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
          <GraduationCap className="h-5 w-5" />
        </div>
        <div>
          <p className="text-xs font-bold uppercase tracking-wider text-gray-500">
            Profil Siswa
          </p>
          <h3 className="text-lg font-bold text-[#37513d]">{profil.nama}</h3>
        </div>
      </header>

      <div className="grid gap-3 text-sm text-gray-600">
        <div className="flex items-center gap-2 rounded-lg border border-gray-100 bg-gray-50 px-3 py-2">
          <Users className="h-4 w-4 text-[#397e50]" />
          <div>
            <p className="text-xs text-gray-500">Kelas</p>
            <p className="font-semibold text-[#37513d]">{profil.kelas}</p>
          </div>
        </div>
        <div className="rounded-lg border border-gray-100 bg-gray-50 px-3 py-2">
          <p className="text-xs text-gray-500">Target Pekan Ini</p>
          <p className="font-semibold text-[#37513d]">
            Selesaikan latihan sebelum ujian terdekat
          </p>
        </div>
      </div>
    </section>
  );
};

export default ProfilSiswaWidget;
