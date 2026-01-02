import React, { forwardRef } from "react";

type ImageUploadProps = {
  sectionTitle?: string;
  helperText?: string;
  formatText?: string;
  optionalText?: string;
  imgSrc?: string;
  imgAlt?: string;


  type?: string;
  accept?: string;
  fileName?: string;
  size?: number;
  imageFileCheck?: boolean;

  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onClick?: () => void;
};

export const ImageUpload = forwardRef<HTMLInputElement, ImageUploadProps>(
  (
    {
      sectionTitle = "Gambar",
      helperText = "Unggah Gambar (maks. 2MB).",
      formatText = "Format: JPG/PNG",
      optionalText = "Opsional",
      imgSrc,
      imgAlt = "Preview Gambar",
      type = "file",
      accept = "image/*",
      fileName,
      size,
      imageFileCheck = false,
      onChange,
      onClick,
    },
    ref
  ) => {
  return (
    <>
      <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <div className="mb-4 flex items-start justify-between gap-4">
          <div>
            <h2 className="text-sm font-semibold text-slate-800">
              {sectionTitle}
            </h2>
            <p className="text-xs text-slate-500">{helperText}</p>
          </div>

          <div className="text-right">
            <p className="text-xs text-slate-500">{formatText}</p>
          </div>
        </div>

        <div className="flex flex-col gap-4 md:flex-row md:items-center">
          <div className="h-24 w-24 overflow-hidden rounded-md border border-slate-200 bg-slate-100">
            {imgSrc ? (
              <img
                src={imgSrc}
                alt={imgAlt}
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-xs text-slate-400">
                No Photo
              </div>
            )}

            <div>
              <p className="text-sm font-medium text-slate-700">{helperText}</p>
              <p className="text-xs text-slate-500">{optionalText}</p>
            </div>
          </div>
          {/* Upload + Chip */}
          <div className="flex flex-1 flex-col gap-2 md:flex-row md:items-center md:justify-end">
            <label className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm transition hover:bg-slate-50">
              Pilih Gambar
              <input
                ref={ref}
                type={type}
                accept={accept}
                className="hidden"
                onChange={onChange}
              />
            </label>

            {imageFileCheck? (
              <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-slate-50 px-3 py-2">
                <div className="min-w-0">
                  <p className="truncate text-xs font-medium text-slate-700">
                    {fileName}
                  </p>
                  <p className="text-[11px] text-slate-500">
                    {size} MB
                  </p>
                </div>
                <button
                  type="button"
                  onClick={onClick}
                  className="rounded-md px-2 py-1 text-xs font-medium text-slate-600 hover:bg-white"
                  aria-label="Hapus foto"
                  title="Hapus foto"
                >
                  X
                </button>
              </div>
            ) : (
              <p className="text-xs text-slate-500 md:text-right">
                Belum ada file dipilih.
              </p>
            )}
          </div>
        </div>
      </div>
    </>
  );
});
