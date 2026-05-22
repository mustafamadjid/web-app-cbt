import type { ReactNode } from "react";

import { helperText, sectionTitle } from "./constants";

type FormSectionProps = {
  title: string;
  description?: ReactNode;
  action?: ReactNode;
  children: ReactNode;
};

const FormSection = ({
  title,
  description,
  action,
  children,
}: FormSectionProps) => (
  <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
    <div className="mb-4 flex items-center justify-between">
      <div>
        <h2 className={sectionTitle}>{title}</h2>
        {description !== undefined && <p className={helperText}>{description}</p>}
      </div>
      {action}
    </div>
    {children}
  </div>
);

export default FormSection;
