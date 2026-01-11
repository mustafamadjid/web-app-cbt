
type StatistikDonutChartProps = {
  size: number;
  stroke: number;
  trackColor: string;
  gradientFrom: string;
  gradientTo: string;
  percent: number;
  id: string;
  className?: string;
};

export const StatistikDonutChart = ({
  size,
  stroke,
  trackColor,
  gradientFrom,
  gradientTo,
  percent,
  id,
  className,
}: StatistikDonutChartProps) => {
  const r = (size - stroke) / 2;
  const c = 2 * Math.PI * r;
  const dash = (percent / 100) * c;

  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      className={className}
    >
      <defs>
        <linearGradient id={id} x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor={gradientFrom} />
          <stop offset="100%" stopColor={gradientTo} />
        </linearGradient>
      </defs>

      <circle
        cx={size / 2}
        cy={size / 2}
        r={r}
        fill="none"
        stroke={trackColor}
        strokeWidth={stroke}
        className="opacity-80"
      />

      <circle
        cx={size / 2}
        cy={size / 2}
        r={r}
        fill="none"
        stroke={`url(#${id})`}
        strokeWidth={stroke}
        strokeLinecap="round"
        strokeDasharray={`${dash} ${c - dash}`}
        className="transition-all duration-1000 ease-out"
      />
    </svg>
  );
};
