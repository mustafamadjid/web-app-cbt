import React from "react";
import { AlertCircle } from "lucide-react";

type ErrorFloatingProps = {
  message: string;
};

const ErrorFloating: React.FC<ErrorFloatingProps> = ({ message }) => {
  if (!message) return null;

  return (
    <div className="pointer-events-none fixed right-4 top-4 z-50 max-w-sm animate-in slide-in-from-top-2 duration-300">
      <div
        role="alert"
        className="flex items-start gap-3 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-rose-700 shadow-lg"
      >
        <AlertCircle className="mt-0.5 h-5 w-5 shrink-0" />
        <p className="text-sm font-medium leading-relaxed">{message}</p>
      </div>
    </div>
  );
};

export default ErrorFloating;
