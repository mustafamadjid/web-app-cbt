type StatistikHasilUjianProps = {
    title:string,
    value:number
}

const StatistikHasilUjianWidget = ({value,title = "Nilai Tertinggi"}:StatistikHasilUjianProps) => {
    return (
      <>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">
            {title}
          </p>
          <p className="mt-3 text-2xl font-bold text-slate-800">
            {value}
          </p>
        </div>
      </>
    );
}

export default StatistikHasilUjianWidget;