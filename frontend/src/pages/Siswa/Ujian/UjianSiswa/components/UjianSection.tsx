import React from "react";

type UjianSectionProps = {
  title: string;
  description: string;
  children: React.ReactNode;
  className?: string;
};

const UjianSection: React.FC<UjianSectionProps> = ({
  title,
  description,
  children,
  className,
}) => {
  return (
    <section className={["space-y-4", className ?? ""].join(" ")}>
      <header className="flex flex-wrap items-start justify-between gap-2">
        <div>
          <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
          <p className="text-sm text-gray-500">{description}</p>
        </div>
      </header>
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {children}
      </div>
    </section>
  );
};

export default UjianSection;
