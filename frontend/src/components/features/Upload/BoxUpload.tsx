import React, { forwardRef } from "react";
import { CloudUpload, FileText, X } from "lucide-react";

type BoxUploadProps = {
  helperText?: string;
  formatText?: string;
  type?: string;
  accept?: string;
  fileName?: string;
  size?: number;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onClear?: () => void;
};

const BoxUpload = forwardRef<HTMLInputElement, BoxUploadProps>(
  (
    {
      helperText = "Unggah File (maks. 10MB).",
      formatText = "Format: .DOCX",
      type = "file",
      accept = "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      fileName,
      size,
      onChange,
      onClear,
    },
    ref
  ) => {
    const formatFileSize = (bytes?: number) => {
      if (!bytes) return "";
      const mb = bytes / (1024 * 1024);
      return mb.toFixed(2) + " MB";
    };

    return (
      <div className="w-full max-w-2xl mx-auto group">
        <label
          className={`
            relative cursor-pointer w-full min-h-[220px] md:min-h-[280px]
            flex flex-col items-center justify-center rounded-2xl
            border-2 border-dashed transition-all duration-300
            ${
              fileName
                ? "border-[#397e50] bg-[#397e50]/5"
                : "border-gray-300 hover:border-[#397e50] bg-gray-50 hover:bg-[#397e50]/5"
            }
          `}
        >
          <input
            ref={ref}
            type={type}
            accept={accept}
            className="hidden"
            onChange={onChange}
          />

          <div className="flex flex-col items-center justify-center p-6 text-center">
            {fileName ? (
              /* Tampilan Saat File Terpilih */
              <div className="flex flex-col items-center animate-in fade-in zoom-in duration-300">
                <div className="p-4 bg-[#397e50]/10 rounded-2xl mb-4">
                  <FileText className="w-12 h-12 text-[#397e50]" />
                </div>
                <div className="space-y-1">
                  <h3 className="text-gray-800 font-bold text-lg break-all px-4">
                    {fileName}
                  </h3>
                  {size && (
                    <p className="text-sm text-gray-500 font-medium">
                      {formatFileSize(size)}
                    </p>
                  )}
                </div>
                <div className="mt-6 flex items-center gap-2 text-[#397e50] font-semibold text-sm bg-[#397e50]/10 px-4 py-2 rounded-full">
                  <span>Ganti File</span>
                </div>
              </div>
            ) : (
              /* Tampilan Standar (Kosong) */
              <div className="flex flex-col items-center">
                <div className="p-5 bg-white rounded-full shadow-sm mb-5 group-hover:scale-110 transition-transform duration-300">
                  <CloudUpload className="w-12 h-12 text-[#397e50]" />
                </div>
                <div className="space-y-2">
                  <h2 className="text-gray-700 font-bold text-lg md:text-xl">
                    {helperText}
                  </h2>
                  <p className="text-gray-400 text-sm font-medium">
                    {formatText}
                  </p>
                </div>
                <div className="mt-8 px-8 py-2.5 bg-[#397e50] text-white text-sm font-bold rounded-xl shadow-lg shadow-[#397e50]/30 hover:bg-[#2d633f] transition-colors">
                  Pilih Dokumen
                </div>
              </div>
            )}
          </div>
        </label>

        {/* Action button untuk menghapus selection */}
        {fileName && onClear && (
          <button
            onClick={(e) => {
              e.preventDefault();
              onClear();
            }}
            className="mt-4 cursor-pointer flex items-center gap-2 text-red-500 text-sm font-bold hover:text-red-700 mx-auto bg-red-50 px-4 py-2 rounded-lg transition-colors"
          >
            <X size={16} />
            Batalkan Pilihan
          </button>
        )}
      </div>
    );
  }
);

BoxUpload.displayName = "BoxUpload";

export default BoxUpload;
