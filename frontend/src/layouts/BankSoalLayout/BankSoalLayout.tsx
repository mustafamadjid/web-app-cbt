import BoxBankSoal from "@/components/features/BankSoal/BoxBankSoal";

import type { BankSoalItem } from "@/types/BankSoal/BankSoal";

type BankSoalLayoutProps = {
  items: BankSoalItem[];
  onPreview?: (item: BankSoalItem) => void;
  onUpload?: (item: BankSoalItem) => void;
  onKelola?: (item: BankSoalItem) => void;
  onHapus?: (item: BankSoalItem) => void;
  resolveGuruLabel?: (item: BankSoalItem) => string;
  resolveMapelLabel?: (item: BankSoalItem) => string;
  resolveKelasLabel?: (item: BankSoalItem) => string;

  /** kalau mau mulai angka dari selain 1 */
  startIndex?: number;

  className?: string;
};

const BankSoalLayout: React.FC<BankSoalLayoutProps> = ({
  items,
  onPreview,
  onUpload,
  onKelola,
  onHapus,
  resolveGuruLabel,
  resolveMapelLabel,
  resolveKelasLabel,
  startIndex = 1,
  className = "",
}) => {
  return (
    <div className={["w-full", className].join(" ")}>
      <ol className="space-y-6">
        {items.map((item, idx) => {
          const number = startIndex + idx;
          const isLast = idx === items.length - 1;

          return (
            <li key={item.id_bank_soal} className="relative ">
              {/* Timeline rail (kiri) */}
              <div className="absolute left-0 top-0 flex h-full w-12 justify-center">
                <div className="relative h-full w-1">
                  {/* shadow layer */}
                  {idx !== 0 && (
                    <div className="absolute left-1/2 top-0 h-4 w-full -translate-x-1/2 rounded-full bg-black/10" />
                  )}
                  {!isLast && (
                    <div className="absolute left-1/2 top-6 h-[calc(100%)] w-full -translate-x-1/2 rounded-full bg-black/10" />
                  )}

                  {idx !== 0 && (
                    <div className="absolute left-1/2 top-0 h-4 w-[3px] -translate-x-1/2 rounded-full bg-emerald-900" />
                  )}
                  {!isLast && (
                    <div className="absolute left-1/2 top-6 h-[calc(100%)] w-[3px] -translate-x-1/2 rounded-full bg-emerald-900" />
                  )}
                </div>
              </div>

              {/* Node angka */}
              <div className="absolute left-0 top-0 z-10 flex w-12 justify-center">
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-emerald-900 text-sm font-semibold text-white ring-4 ring-white shadow-sm">
                  {number}
                </div>
              </div>

              {/* Content card (geser kanan, sejajar node) */}
              <div className="pl-16">
                <BoxBankSoal
                  idBankSoal={item.id_bank_soal}
                  guruLabel={resolveGuruLabel ? resolveGuruLabel(item) : "-"}
                  namaBankSoal={item.nama_bank_soal}
                  kelasLabel={resolveKelasLabel ? resolveKelasLabel(item) : "-"}
                  mapelLabel={resolveMapelLabel ? resolveMapelLabel(item) : "-"}
                  materi={item.materi}
                  tglBuat={undefined}
                  onPreview={onPreview ? () => onPreview(item) : undefined}
                  onUpload={onUpload ? () => onUpload(item) : undefined}
                  onKelola={onKelola ? () => onKelola(item) : undefined}
                  onHapus={onHapus ? () => onHapus(item) : undefined}
                />
              </div>
            </li>
          );
        })}
      </ol>
    </div>
  );
};

export default BankSoalLayout;
