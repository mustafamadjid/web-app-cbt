import React from "react";

import { getRemainingExamTime } from "@/helper/Countdown/getRemainingExamTime";

export const DEFAULT_TIMER_LABEL = "00:00:00";

type UseExamCountdownResult = {
  sisaWaktu: string;
  isTimeExpired: boolean;
  hasValidWaktuSelesai: boolean;
};

export function useExamCountdown(waktuSelesai: string): UseExamCountdownResult {
  const initialRemainingTime = React.useMemo(() => {
    if (!waktuSelesai) {
      return null;
    }

    return getRemainingExamTime(waktuSelesai);
  }, [waktuSelesai]);
  const [sisaWaktu, setSisaWaktu] = React.useState(DEFAULT_TIMER_LABEL);
  const [isTimeExpired, setIsTimeExpired] = React.useState(false);

  React.useEffect(() => {
    if (!waktuSelesai) {
      setSisaWaktu(DEFAULT_TIMER_LABEL);
      setIsTimeExpired(false);
      return;
    }

    const currentRemainingTime = getRemainingExamTime(waktuSelesai);
    if (!currentRemainingTime) {
      setSisaWaktu(DEFAULT_TIMER_LABEL);
      setIsTimeExpired(false);
      return;
    }

    setSisaWaktu(currentRemainingTime.formattedTime);
    setIsTimeExpired(currentRemainingTime.isExpired);

    if (currentRemainingTime.isExpired) {
      return;
    }

    const timerId = window.setInterval(() => {
      const nextRemainingTime = getRemainingExamTime(waktuSelesai);
      if (!nextRemainingTime) {
        window.clearInterval(timerId);
        setSisaWaktu(DEFAULT_TIMER_LABEL);
        setIsTimeExpired(false);
        return;
      }

      setSisaWaktu(nextRemainingTime.formattedTime);

      if (nextRemainingTime.isExpired) {
        setIsTimeExpired(true);
        window.clearInterval(timerId);
      }
    }, 1000);

    return () => {
      window.clearInterval(timerId);
    };
  }, [waktuSelesai]);

  return {
    sisaWaktu,
    isTimeExpired,
    hasValidWaktuSelesai: Boolean(waktuSelesai && initialRemainingTime),
  };
}
