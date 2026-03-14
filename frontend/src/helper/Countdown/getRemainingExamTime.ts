const padTimePart = (value: number) => String(value).padStart(2, "0");

const formatRemainingTime = (remainingMs: number) => {
  const totalSeconds = Math.max(0, Math.floor(remainingMs / 1000));
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  return [
    padTimePart(hours),
    padTimePart(minutes),
    padTimePart(seconds),
  ].join(":");
};

export type RemainingExamTime = {
  remainingMs: number;
  formattedTime: string;
  isExpired: boolean;
};

export const getRemainingExamTime = (
  waktuSelesai: string,
  nowMs = Date.now(),
): RemainingExamTime | null => {
  if (!waktuSelesai) return null;

  const targetMs = new Date(waktuSelesai).getTime();
  if (Number.isNaN(targetMs)) return null;

  const remainingMs = Math.max(0, targetMs - nowMs);

  return {
    remainingMs,
    formattedTime: formatRemainingTime(remainingMs),
    isExpired: remainingMs === 0,
  };
};
