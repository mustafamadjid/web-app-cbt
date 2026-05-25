import React from "react";

export type ExamBrowserViolationReason =
  | "new-tab-click"
  | "new-window-shortcut"
  | "tab-hidden"
  | "window-blur";

type UseExamBrowserExitCleanupParams = {
  hasActiveExamSession: boolean;
  onViolation: (reason: ExamBrowserViolationReason) => Promise<boolean>;
};

export function useExamBrowserExitCleanup({
  hasActiveExamSession,
  onViolation,
}: UseExamBrowserExitCleanupParams) {
  const pendingBrowserExitRef = React.useRef(false);
  const violationTriggeredRef = React.useRef(false);
  const onViolationRef = React.useRef(onViolation);

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
    onViolationRef.current = onViolation;
  }, [onViolation]);

  React.useEffect(() => {
    if (
      !hasActiveExamSession ||
      typeof window === "undefined" ||
      typeof document === "undefined"
    ) {
      pendingBrowserExitRef.current = false;
      violationTriggeredRef.current = false;
      return;
    }

    const triggerViolation = (reason: ExamBrowserViolationReason) => {
      if (violationTriggeredRef.current) {
        return;
      }

      violationTriggeredRef.current = true;

      void (async () => {
        const handled = await onViolationRef.current(reason);
        if (!handled) {
          violationTriggeredRef.current = false;
        }
      })();
    };

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

    const isNewTabMouseEvent = (event: MouseEvent) => {
      const target = event.target;
      if (!(target instanceof Element)) {
        return false;
      }

      const link = target.closest("a[href]");
      if (!link) {
        return false;
      }

      const isLeftClick = event.button === 0;
      const isMiddleClick = event.button === 1;
      const isNewTab = (event.ctrlKey || event.metaKey) && isLeftClick;
      const isNewWindow = event.shiftKey && isLeftClick;

      return isNewTab || isNewWindow || isMiddleClick;
    };

    const handleLinkClick = (event: MouseEvent) => {
      if (!isNewTabMouseEvent(event)) {
        return;
      }

      event.preventDefault();
      event.stopPropagation();
      triggerViolation("new-tab-click");
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      const key = event.key.toLowerCase();
      const isBrowserNewWindowShortcut =
        (event.ctrlKey || event.metaKey) && (key === "t" || key === "n");

      if (!isBrowserNewWindowShortcut) {
        return;
      }

      event.preventDefault();
      event.stopPropagation();
      triggerViolation("new-window-shortcut");
    };

    const handleVisibilityChange = () => {
      if (document.visibilityState !== "hidden") {
        return;
      }

      triggerViolation("tab-hidden");
    };

    const handleWindowBlur = () => {
      triggerViolation("window-blur");
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    window.addEventListener("pagehide", handlePageHide);
    window.addEventListener("focus", resetPendingBrowserExit);
    window.addEventListener("pageshow", resetPendingBrowserExit);
    window.addEventListener("blur", handleWindowBlur);
    document.addEventListener("click", handleLinkClick, true);
    document.addEventListener("auxclick", handleLinkClick, true);
    document.addEventListener("keydown", handleKeyDown, true);
    document.addEventListener("visibilitychange", handleVisibilityChange);

    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
      window.removeEventListener("pagehide", handlePageHide);
      window.removeEventListener("focus", resetPendingBrowserExit);
      window.removeEventListener("pageshow", resetPendingBrowserExit);
      window.removeEventListener("blur", handleWindowBlur);
      document.removeEventListener("click", handleLinkClick, true);
      document.removeEventListener("auxclick", handleLinkClick, true);
      document.removeEventListener("keydown", handleKeyDown, true);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      pendingBrowserExitRef.current = false;
      violationTriggeredRef.current = false;
    };
  }, [clearBrowserStorageCache, hasActiveExamSession]);
}
