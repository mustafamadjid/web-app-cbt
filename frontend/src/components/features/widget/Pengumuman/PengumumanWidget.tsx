import React from "react";
import {
  ChevronDown,
  Megaphone,
  Paperclip,
  CalendarDays,
  FileText,
} from "lucide-react";

export type AnnouncementDoc = {
  id?: string;
  name: string;
  url: string;
  mimeType?: string;
  sizeLabel?: string;
};

export type PengumumanItem = {
  id: string;
  judul: string;
  isi_pengumuman: string;
  tanggal_rilis_pengumuman: string;
  dokumen?: AnnouncementDoc | AnnouncementDoc[] | null;
};

type PengumumanWidgetProps = {
  title?: string;
  items: PengumumanItem[];
  defaultOpenId?: string;
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
        "relative flex flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:shadow-lg hover:shadow-[#397e50]/5",
        "h-full", // Pastikan widget mengisi height parent jika diperlukan
        className ?? "",
      ].join(" ")}
    >
      {/* Top Accent Line */}
      <div className="h-1.5 w-full bg-[#37513d]" />

      {/* Header */}
      <header className="flex items-center gap-3 px-5 pt-5 pb-2">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
          <Megaphone className="h-5 w-5" />
        </div>
        <div>
          <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
          <p className="text-xs font-medium text-gray-500">
            {items.length} informasi terbaru
          </p>
        </div>
      </header>

      {/* List Container */}
      <div className="mt-2 flex-1 overflow-y-auto px-5 pb-5">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-gray-300 bg-gray-50 py-10 text-center mt-2">
            <Megaphone className="h-8 w-8 text-gray-300" />
            <p className="text-sm text-gray-500">
              Tidak ada pengumuman saat ini.
            </p>
          </div>
        ) : (
          <div className="space-y-3 pt-2">
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
        "group overflow-hidden rounded-xl border transition-all duration-300",
        isOpen
          ? "border-[#397e50] bg-white shadow-md shadow-[#397e50]/5 ring-1 ring-[#397e50]/10"
          : "border-gray-200 bg-white hover:border-[#397e50]/40 hover:bg-gray-50/50",
      ].join(" ")}
    >
      {/* Header Row (Clickable) */}
      <button
        type="button"
        onClick={onToggle}
        className="flex w-full cursor-pointer items-start justify-between gap-4 px-4 py-3.5 text-left focus:outline-none"
        aria-expanded={isOpen}
      >
        <div className="flex-1 space-y-1.5">
          {/* Tanggal */}
          <div className="flex items-center gap-1.5 text-xs font-medium text-gray-500">
            <CalendarDays className="h-3.5 w-3.5" />
            <span>{item.tanggal_rilis_pengumuman}</span>
          </div>

          {/* Judul */}
          <h3
            className={[
              "text-sm font-bold leading-snug transition-colors",
              isOpen
                ? "text-black"
                : "text-gray-800 group-hover:text-[#37513d]",
            ].join(" ")}
          >
            {item.judul}
          </h3>
        </div>

        {/* Chevron Icon */}
        <span
          className={[
            "mt-1 flex h-6 w-6 shrink-0 items-center justify-center rounded-full transition-all duration-300",
            isOpen
              ? "bg-[#397e50] text-white rotate-180"
              : "bg-gray-100 text-gray-500 group-hover:bg-[#397e50]/10 group-hover:text-[#397e50]",
          ].join(" ")}
          aria-hidden="true"
        >
          <ChevronDown className="h-4 w-4" strokeWidth={2.5} />
        </span>
      </button>

      {/* Accordion Content */}
      <div
        className={[
          "grid transition-[grid-template-rows] duration-300 ease-out",
          isOpen ? "grid-rows-[1fr]" : "grid-rows-[0fr]",
        ].join(" ")}
      >
        <div className="overflow-hidden">
          <div className="border-t border-gray-100 px-4 py-4">
            {/* Isi Pengumuman */}
            <div className="prose-sm text-sm leading-relaxed text-gray-600 whitespace-pre-wrap">
              {item.isi_pengumuman}
            </div>

            {/* Dokumen Section */}
            {docs.length > 0 && (
              <div className="mt-5">
                <div className="mb-3 flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-gray-400">
                  <Paperclip className="h-3.5 w-3.5" />
                  Lampiran Dokumen
                </div>

                <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                  {docs.map((d, idx) => (
                    <a
                      key={d.id ?? `${d.url}-${idx}`}
                      href={d.url}
                      target="_blank"
                      rel="noreferrer"
                      className={[
                        "group/doc flex items-center gap-3 rounded-lg border border-gray-200 bg-gray-50/50 p-2.5 transition-all",
                        "hover:border-[#397e50]/50 hover:bg-[#397e50]/5 hover:shadow-sm",
                      ].join(" ")}
                      title={d.name}
                    >
                      {/* Icon Container */}
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded bg-white text-[#397e50] shadow-sm ring-1 ring-gray-100 group-hover/doc:ring-[#397e50]/20">
                        <FileText className="h-4 w-4" />
                      </div>

                      {/* File Info */}
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-xs font-bold text-gray-700 group-hover/doc:text-[#397e50]">
                          {d.name}
                        </p>
                        {d.sizeLabel && (
                          <p className="mt-0.5 text-2xs text-gray-500">
                            {d.sizeLabel}
                          </p>
                        )}
                      </div>
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
