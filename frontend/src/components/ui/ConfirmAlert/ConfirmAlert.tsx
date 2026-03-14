import React from "react";

type ConfirmAlertProps = {
  isOpen: boolean;
  message: string;
  onClose: () => void;
  onConfirm: () => void;
  isLoading?: boolean;
  title?: string;
  confirmLabel?: string;
  loadingLabel?: string;
  cancelLabel?: string;
  hideCancel?: boolean;
  dismissible?: boolean;
  confirmClassName?: string;
};

const ConfirmAlert: React.FC<ConfirmAlertProps> = ({
  isOpen,
  message,
  onClose,
  onConfirm,
  isLoading = false,
  title = "Konfirmasi Hapus Akun",
  confirmLabel = "Ya, Hapus",
  loadingLabel = "Menghapus...",
  cancelLabel = "Batal",
  hideCancel = false,
  dismissible = true,
  confirmClassName = "",
}) => {
  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/55 px-4"
      role="dialog"
      aria-modal="true"
    >
      <div className="w-full max-w-xl rounded-xl border border-slate-200 bg-white p-6 shadow-2xl">
        <h3 className="text-lg font-semibold text-slate-900">{title}</h3>
        <p className="mt-2 text-sm leading-relaxed text-slate-600">{message}</p>

        <div className="mt-6 flex justify-end gap-3">
          {!hideCancel && (
            <button
              type="button"
              onClick={onClose}
              disabled={isLoading || !dismissible}
              className="cursor-pointer rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {cancelLabel}
            </button>
          )}
          <button
            type="button"
            onClick={onConfirm}
            disabled={isLoading}
            className={[
              "cursor-pointer rounded-lg bg-rose-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-rose-700 disabled:cursor-not-allowed disabled:opacity-60",
              confirmClassName,
            ].join(" ")}
          >
            {isLoading ? loadingLabel : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConfirmAlert;
