import { BoxJadwalUjian } from "@/components/features/Ujian/BoxJadwalUjian"

import { getJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service"
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { useEffect, useRef, useState } from "react";


export const JadwalUjian = () => {
    const [loading, setLoading] = useState(false);
    const [errorMsg, setErrorMsg] = useState("");
    
    // Daftar ujian terjadwal
    const [daftarJadwalUjian, setDaftarJadwalUjian] = useState<JadwalUjianItem[]>([]);

   const requestSeq = useRef(0);

    // Fetch data jadwal ujian
    useEffect(() => {
        const seq = ++requestSeq.current;
        (async () => {
              try {
                setLoading(true);
                setErrorMsg("");
        
                const data = await getJadwalUjian();
        
                if (seq !== requestSeq.current) return;
        
                setDaftarJadwalUjian(data);
                
              } catch (e) {
                if (seq !== requestSeq.current) return;
                setErrorMsg("Gagal memuat data siswa.");
                setDaftarJadwalUjian([]);
              } finally {
                if (seq !== requestSeq.current) return;
                setLoading(false);
              }
            })();
          }, []);


    return (
      <>
        <div className="px-8 py-10">
          <div className="flex flex-col gap-5">
            {daftarJadwalUjian.map((ujian) => (
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
        </div>
      </>
    );
}