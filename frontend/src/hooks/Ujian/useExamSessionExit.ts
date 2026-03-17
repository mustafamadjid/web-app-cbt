import React from "react";
import { useBlocker } from "react-router";
import type { BlockerFunction } from "react-router";

import { ApiError } from "@/services/Api/api";
import { useExpireAttemptUjianSiswa } from "@/services/Api/features-api/Ujian/ujian.service";

const EXPIRE_ATTEMPT_ERROR_MESSAGE =
  "Gagal mengakhiri sesi ujian. Silakan coba lagi.";

type UseExamSessionExitParams = {
  attemptId: number | null;
  hasActiveExamSession: boolean;
  isTimeExpired: boolean;
  resetKey: number;
  clearAttemptCache: () => void;
  clearSoalCache: () => void;
  onFallbackLeave: () => void;
};

type UseExamSessionExitResult = {
  isLeaveConfirmOpen: boolean;
  sessionExitError: string | null;
  expiringAttempt: boolean;
  allowNavigation: () => void;
  clearSessionExitError: () => void;
  handleExpiredSubmit: () => Promise<void>;
  handleLeaveConfirm: () => Promise<void>;
  handleLeaveCancel: () => void;
};

const mapExpireAttemptErrorMessage = (error: unknown): string => {
  if (error instanceof ApiError && error.message) {
    return error.message;
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return EXPIRE_ATTEMPT_ERROR_MESSAGE;
};

export function useExamSessionExit({
  attemptId,
  hasActiveExamSession,
  isTimeExpired,
  resetKey,
  clearAttemptCache,
  clearSoalCache,
  onFallbackLeave,
}: UseExamSessionExitParams): UseExamSessionExitResult {
  const [isLeaveConfirmOpen, setIsLeaveConfirmOpen] = React.useState(false);
  const [sessionExitError, setSessionExitError] = React.useState<string | null>(
    null,
  );
  const allowNavigationRef = React.useRef(false);
  const expireTriggeredRef = React.useRef(false);
  const { execute: executeExpireAttempt, loading: expiringAttempt } =
    useExpireAttemptUjianSiswa();

  React.useEffect(() => {
    setIsLeaveConfirmOpen(false);
    setSessionExitError(null);
    allowNavigationRef.current = false;
    expireTriggeredRef.current = false;
  }, [resetKey]);

  const shouldBlockNavigation = React.useCallback<BlockerFunction>(
    ({ currentLocation, nextLocation }) => {
      if (allowNavigationRef.current || !hasActiveExamSession) {
        return false;
      }

      return (
        currentLocation.pathname !== nextLocation.pathname ||
        currentLocation.search !== nextLocation.search ||
        currentLocation.hash !== nextLocation.hash
      );
    },
    [hasActiveExamSession],
  );
  const navigationBlocker = useBlocker(shouldBlockNavigation);

  const expireAttemptBeforeLeave = React.useCallback(async () => {
    if (attemptId === null || expireTriggeredRef.current) {
      return true;
    }

    setSessionExitError(null);
    expireTriggeredRef.current = true;

    try {
      await executeExpireAttempt(attemptId);
      clearAttemptCache();
      clearSoalCache();
      return true;
    } catch (error) {
      expireTriggeredRef.current = false;
      setSessionExitError(mapExpireAttemptErrorMessage(error));
      return false;
    }
  }, [
    attemptId,
    clearAttemptCache,
    clearSoalCache,
    executeExpireAttempt,
  ]);

  const allowNavigation = React.useCallback(() => {
    allowNavigationRef.current = true;
  }, []);

  const clearSessionExitError = React.useCallback(() => {
    setSessionExitError(null);
  }, []);

  React.useEffect(() => {
    if (isTimeExpired && navigationBlocker.state === "blocked") {
      navigationBlocker.reset();
    }
  }, [isTimeExpired, navigationBlocker]);

  React.useEffect(() => {
    if (navigationBlocker.state === "blocked") {
      setIsLeaveConfirmOpen(true);
      return;
    }

    setIsLeaveConfirmOpen(false);
  }, [navigationBlocker.state]);

  const handleExpiredSubmit = React.useCallback(async () => {
    const expired = await expireAttemptBeforeLeave();
    if (!expired) {
      return;
    }

    allowNavigationRef.current = true;
    onFallbackLeave();
  }, [expireAttemptBeforeLeave, onFallbackLeave]);

  const handleLeaveConfirm = React.useCallback(async () => {
    const expired = await expireAttemptBeforeLeave();
    if (!expired) {
      return;
    }

    allowNavigationRef.current = true;
    setIsLeaveConfirmOpen(false);

    if (navigationBlocker.state === "blocked") {
      navigationBlocker.proceed();
      return;
    }

    onFallbackLeave();
  }, [expireAttemptBeforeLeave, navigationBlocker, onFallbackLeave]);

  const handleLeaveCancel = React.useCallback(() => {
    setIsLeaveConfirmOpen(false);

    if (navigationBlocker.state === "blocked") {
      navigationBlocker.reset();
    }
  }, [navigationBlocker]);

  return {
    isLeaveConfirmOpen,
    sessionExitError,
    expiringAttempt,
    allowNavigation,
    clearSessionExitError,
    handleExpiredSubmit,
    handleLeaveConfirm,
    handleLeaveCancel,
  };
}
