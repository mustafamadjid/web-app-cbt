// import TambahBankSoalForm from "@/layouts/Form/Admin/BankSoal/TambahBankSoal";

import BoxUpload from "@/components/features/Upload/BoxUpload";
import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  fileDocxOnly,
  fileMaxSize,
} from "@/helper/validate/validateForm";
import { useEffect, useRef, useState } from "react";

type FileValues = {
  file: File | null;
};

const initialValues: FileValues = {
  file: null,
};
const TambahBankSoal = () => {
  const [values, setValues] = useState<FileValues>(initialValues);
  const [fileUrl, setFileUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (!values.file) {
      setFileUrl("");
      return;
    }

    const url = URL.createObjectURL(values.file);
    setFileUrl(url);

    return () => URL.revokeObjectURL(url);
  }, [values.file]);

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

  // TODO : Buat onsubmit

  const clearFile = () => {
    setField("file", null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <>
      <form action="">
        <div className="flex flex-col items-center justify-center gap-10 px-8 py-8">
          <BoxUpload
            ref={fileInputRef}
            fileName={values.file?.name}
            size={values.file?.size}
            onChange={(e) => {
              const file = e.target.files?.[0] ?? null;
              setField("file", file);
            }}
            onClear={clearFile}
          />

          {hasError("file") && (
            <p className="mt-2 text-2xl text-rose-600">{errors.file}</p>
          )}
          <button
            className="
                w-full md:w-auto px-8 py-3 cursor-pointer
                bg-[#397e50] hover:bg-green-800
                text-white font-semibold uppercase tracking-wide text-sm
                rounded-xl shadow-md shadow-green-200/70
                transition-shadow duration-150 ease-out
                active:scale-95 active:shadow-sm
                focus:outline-none focus:ring-2 focus:ring-green-400 focus:ring-offset-2
                flex items-center justify-center gap-2
            "
            type="submit"
          >
            <span>Submit Soal</span>
          </button>
        </div>
      </form>
    </>
  );
};

export default TambahBankSoal;
