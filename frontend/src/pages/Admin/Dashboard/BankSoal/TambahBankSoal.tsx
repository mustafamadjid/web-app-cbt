import BoxUpload from "@/components/features/Upload/BoxUpload";
import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  fileDocxOnly,
  fileMaxSize,
} from "@/helper/validate/validateForm";
import {
  getImportSoalJob,
  uploadImportSoal,
} from "@/services/Api/features-api/BankSoal/importSoal.service";
import { ApiError } from "@/services/Api/api";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import { useRef, useState, type FormEvent } from "react";
import { useParams } from "react-router";

type FileValues = {
  file: File | null;
};

const initialValues: FileValues = {
  file: null,
};

const TambahBankSoal = () => {
  const { idBankSoal: idBankSoalParam } = useParams();
  const idBankSoal = Number(idBankSoalParam) || 0;

  const [values, setValues] = useState<FileValues>(initialValues);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [isUploading, setIsUploading] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const setField = createSetField(setValues);

  const validate = createValidator<typeof initialValues>({
    file: [
      fileMaxSize(20 * 1024 * 1024, "Ukuran file maksimal 20MB."),
      fileDocxOnly("File harus berformat .docx."),
    ],
  });
  const errors = validate(values);

  const hasError = (name: keyof typeof initialValues) =>
    !!errors[name] && !!values[name];

  const clearFile = () => {
    setField("file", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSuccessMessage(null);
    setErrorMessage(null);

    // Validasi: file harus ada
    if (!values.file) {
      setErrorMessage("Pilih file .docx terlebih dahulu.");
      return;
    }

    // Validasi: tidak ada error validasi
    if (Object.keys(errors).length > 0) {
      setErrorMessage(Object.values(errors)[0] as string);
      return;
    }

    // Validasi: idBankSoal harus valid
    if (!idBankSoal || idBankSoal <= 0) {
      setErrorMessage("Data bank soal tidak valid. Silakan kembali ke halaman Bank Soal.");
      return;
    }

    setIsUploading(true);
    try {
      const { id_job } = await uploadImportSoal(idBankSoal, values.file);
      setSuccessMessage("File berhasil di-upload, memproses import soal...");

      const maxAttempts = 20;
      for (let attempt = 0; attempt < maxAttempts; attempt += 1) {
        await new Promise((resolve) => window.setTimeout(resolve, 1500));
        const job = await getImportSoalJob(id_job);

        if (job.status === "failed") {
          throw new ApiError(
            job.error_msg || "Import soal gagal diproses.",
            400,
            job,
          );
        }

        if (job.status === "completed") {
          const warningText = job.warning_msg?.trim();
          setSuccessMessage(
            warningText
              ? `Import selesai dengan peringatan: ${warningText}`
              : `Import selesai. Total soal: ${job.total_soal}.`,
          );
          clearFile();
          return;
        }
      }

      setSuccessMessage(
        "File berhasil di-upload. Import masih diproses, silakan cek kembali beberapa saat lagi.",
      );
      clearFile();
    } catch (error) {
      if (error instanceof ApiError) {
        setErrorMessage(
          getUserFriendlyErrorMessage(error, {
            action: "submit",
            entity: "import soal",
            fallbackMessage: "Terjadi kesalahan saat mengupload file.",
          }),
        );
      } else {
        setErrorMessage("Terjadi kesalahan saat mengupload file.");
      }
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <>
      <form onSubmit={handleSubmit}>
        <div className="flex flex-col items-center justify-center gap-10 px-8 py-8">
          <BoxUpload
            ref={fileInputRef}
            fileName={values.file?.name}
            size={values.file?.size}
            onChange={(e) => {
              const file = e.target.files?.[0] ?? null;
              setField("file", file);
              // Reset pesan saat user memilih file baru
              setSuccessMessage(null);
              setErrorMessage(null);
            }}
            onClear={clearFile}
          />

          {hasError("file") && (
            <p className="mt-2 text-2xl text-rose-600">{errors.file}</p>
          )}

          {/* Pesan sukses */}
          {successMessage && (
            <div className="w-full md:w-auto rounded-lg border border-emerald-200 bg-emerald-50 px-6 py-3 text-sm font-medium text-emerald-700 text-center">
              {successMessage}
            </div>
          )}

          {/* Pesan error */}
          {errorMessage && (
            <div className="w-full md:w-auto rounded-lg border border-rose-200 bg-rose-50 px-6 py-3 text-sm font-medium text-rose-600 text-center">
              {errorMessage}
            </div>
          )}

          <button
            className={`
                w-full md:w-auto px-8 py-3
                ${isUploading ? "cursor-not-allowed bg-gray-400" : "cursor-pointer bg-[#397e50] hover:bg-green-800"}
                text-white font-semibold uppercase tracking-wide text-sm
                rounded-xl shadow-md shadow-green-200/70
                transition-shadow duration-150 ease-out
                active:scale-95 active:shadow-sm
                focus:outline-none focus:ring-2 focus:ring-green-400 focus:ring-offset-2
                flex items-center justify-center gap-2
            `}
            type="submit"
            disabled={isUploading}
          >
            <span>{isUploading ? "Mengupload..." : "Submit Soal"}</span>
          </button>
        </div>
      </form>
    </>
  );
};

export default TambahBankSoal;
