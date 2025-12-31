import React from "react";
import { ChevronDown, Megaphone, Paperclip } from "lucide-react";

export type AnnouncementDoc = {
  id?: string;
  name: string; 
  url: string; 
  mimeType?: string;
  sizeLabel?: string; 
};


// Nanti buat mappping agar field dari response API bisa sesuai dengan type di bawah ini
// Buat mapper di utils dan panggil mapper nanti di service ketika getPEngumuman
export type PengumumanItem = {
  id: string;
  judul: string;
  isi_pengumuman: string; 
  tanggal_rilis_pengumuman: string; 
  dokumen?: AnnouncementDoc | AnnouncementDoc[] | null;
};

type PengumumanWidgetProps = {
  title?: string; // default: "Pengumuman"
  items: PengumumanItem[];

  /** default: buka item pertama */
  defaultOpenId?: string;

  /** default: false (accordion) */
  allowMultipleOpen?: boolean;

  className?: string;
};

export const PengumumanWidget: React.FC<PengumumanWidgetProps> = ({
  title = "Pengumuman",
  items,
  defaultOpenId,
  allowMultipleOpen = false,
  className,
}) => {
  const [openIds, setOpenIds] = React.useState<Set<string>>(() => {
    const s = new Set<string>();
    if (defaultOpenId) s.add(defaultOpenId);
    else if (items[0]?.id) s.add(items[0].id);
    return s;
  });

  const toggle = (id: string) => {
    setOpenIds((prev) => {
      const next = new Set(prev);
      const isOpen = next.has(id);

      if (allowMultipleOpen) {
        if (isOpen) next.delete(id);
        else next.add(id);
        return next;
      }

      next.clear();
      if (!isOpen) next.add(id);
      return next;
    });
  };

  return (
    <section
      className={[
        "rounded-2xl border border-slate-100 bg-white",
        "shadow-[0_10px_24px_rgba(15,23,42,0.06)] h-full",
        "p-4 sm:p-5",
        className ?? "",
      ].join(" ")}
    >
      {/* Header */}
      <header className="flex items-center gap-2">
        <span className="grid h-7 w-7 place-items-center rounded-full bg-slate-50 text-slate-500">
          <Megaphone size={16} strokeWidth={2} />
        </span>
        <h2 className="text-sm font-semibold text-slate-800">{title}</h2>
      </header>

      <div className="mt-3 border-t border-slate-100 pt-3">
        {items.length === 0 ? (
          <div className="rounded-xl border border-dashed border-slate-200 p-4 text-sm text-slate-500">
            Tidak ada pengumuman.
          </div>
        ) : (
          <div className="space-y-2">
            {items.map((item) => {
              const isOpen = openIds.has(item.id);
              return (
                <AnnouncementRow
                  key={item.id}
                  item={item}
                  isOpen={isOpen}
                  onToggle={() => toggle(item.id)}
                />
              );
            })}
          </div>
        )}
      </div>
    </section>
  );
};

function AnnouncementRow({
  item,
  isOpen,
  onToggle,
}: {
  item: PengumumanItem;
  isOpen: boolean;
  onToggle: () => void;
}) {
  const docs = normalizeDocs(item.dokumen);

  return (
    <article
      className={[
        "rounded-xl border bg-white cursor-pointer",
        isOpen ? "border-[#397e50]" : "border-slate-100",
      ].join(" ")}
    >
      {/* Row header (clickable) */}
      <button
        type="button"
        onClick={onToggle}
        className="flex w-full items-start justify-between gap-3 px-4 py-3 text-left cursor-pointer"
        aria-expanded={isOpen}
      >
        <div className="min-w-0 ">
          <div className="truncate text-sm font-semibold text-slate-800 cursor-pointer">
            {item.judul}
          </div>
          <div className="mt-1 text-xs text-slate-500">
            {item.tanggal_rilis_pengumuman}
          </div>
        </div>

        <span
          className={[
            "mt-0.5 shrink-0 text-slate-500 transition-transform",
            isOpen ? "rotate-180" : "rotate-0",
          ].join(" ")}
          aria-hidden="true"
        >
          <ChevronDown size={18} strokeWidth={2.5} />
        </span>
      </button>

      {/* Dropdown content */}
      <div
        className={[
          "grid transition-[grid-template-rows]  duration-200 ease-out",
          isOpen ? "grid-rows-[1fr]" : "grid-rows-[0fr]",
        ].join(" ")}
      >
        <div className="overflow-hidden">
          <div className="border-t border-slate-100 px-4 py-3">
            {/* isi_pengumuman */}
            <p className="whitespace-pre-wrap text-sm leading-relaxed text-slate-600">
              {item.isi_pengumuman}
            </p>

            {/* dokumen */}
            {docs.length > 0 && (
              <div className="mt-4">
                <div className="mb-2 text-xs font-semibold text-slate-700">
                  Dokumen
                </div>

                <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap">
                  {docs.map((d, idx) => (
                    <a
                      key={d.id ?? `${d.url}-${idx}`}
                      href={d.url}
                      target="_blank"
                      rel="noreferrer"
                      className={[
                        "inline-flex items-center justify-center gap-2",
                        "rounded-lg border border-slate-200 bg-white",
                        "px-3 py-2 text-xs font-semibold text-slate-700",
                        "hover:bg-slate-50",
                        "w-full sm:w-auto",
                      ].join(" ")}
                      title={d.name}
                    >
                      <Paperclip size={16} strokeWidth={2} />
                      <span className="truncate max-w-[240px] sm:max-w-[260px]">
                        {d.name}
                      </span>
                      {d.sizeLabel ? (
                        <span className="text-slate-400">{d.sizeLabel}</span>
                      ) : null}
                    </a>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </article>
  );
}

function normalizeDocs(dokumen: PengumumanItem["dokumen"]): AnnouncementDoc[] {
  if (!dokumen) return [];
  return Array.isArray(dokumen) ? dokumen : [dokumen];
}
