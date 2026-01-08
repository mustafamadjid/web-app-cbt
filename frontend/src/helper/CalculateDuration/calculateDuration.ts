import { parseTime } from "../Parse/parseTime";

export const calculateDuration = (start: string, end: string) => {
  const startMin = parseTime(start);
  const endMin = parseTime(end);
  if (startMin == null || endMin == null) return 0;
  const diff = endMin - startMin;
  return diff > 0 ? diff : 0;
};