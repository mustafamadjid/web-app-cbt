import React from "react";

type UseExamBrowserExitCleanupParams = {
  hasActiveExamSession: boolean;
};

export function useExamBrowserExitCleanup({
  hasActiveExamSession,
}: UseExamBrowserExitCleanupParams) {
  const pendingBrowserExitRef = React.useRef(false);

  const clearBrowserStorageCache = React.useCallback(() => {
    if (typeof window === "undefined") {
      return;
    }

    try {
      window.localStorage.clear();
    } catch {
      // Ignore storage cleanup errors.
    }
  }, []);

  React.useEffect(() => {
    if (!hasActiveExamSession || typeof window === "undefined") {
      pendingBrowserExitRef.current = false;
      return;
    }

    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      pendingBrowserExitRef.current = true;
      event.preventDefault();
      event.returnValue = "";
    };

    const handlePageHide = () => {
      if (!pendingBrowserExitRef.current) {
        return;
      }

      clearBrowserStorageCache();
    };

    const resetPendingBrowserExit = () => {
      pendingBrowserExitRef.current = false;
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    window.addEventListener("pagehide", handlePageHide);
    window.addEventListener("focus", resetPendingBrowserExit);
    window.addEventListener("pageshow", resetPendingBrowserExit);

    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
      window.removeEventListener("pagehide", handlePageHide);
      window.removeEventListener("focus", resetPendingBrowserExit);
      window.removeEventListener("pageshow", resetPendingBrowserExit);
      pendingBrowserExitRef.current = false;
    };
  }, [clearBrowserStorageCache, hasActiveExamSession]);
}
