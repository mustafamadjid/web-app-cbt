// @vitest-environment jsdom

import { cleanup, render, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { useExamBrowserExitCleanup } from "./useExamBrowserExitCleanup";

function TestExamBrowserExitCleanup({
  hasActiveExamSession,
  onViolation,
}: {
  hasActiveExamSession: boolean;
  onViolation: Parameters<typeof useExamBrowserExitCleanup>[0]["onViolation"];
}) {
  useExamBrowserExitCleanup({
    hasActiveExamSession,
    onViolation,
  });

  return <a href="/ujian-lain">Ujian lain</a>;
}

describe("useExamBrowserExitCleanup", () => {
  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  function setDocumentVisibilityState(visibilityState: DocumentVisibilityState) {
    Object.defineProperty(document, "visibilityState", {
      configurable: true,
      value: visibilityState,
    });
  }

  it("mencegah shortcut tab baru saat sesi ujian aktif", async () => {
    const onViolation = vi.fn().mockResolvedValue(true);
    render(
      <TestExamBrowserExitCleanup
        hasActiveExamSession
        onViolation={onViolation}
      />,
    );

    const event = new KeyboardEvent("keydown", {
      key: "t",
      ctrlKey: true,
      bubbles: true,
      cancelable: true,
    });
    const stopPropagation = vi.spyOn(event, "stopPropagation");

    document.dispatchEvent(event);

    expect(event.defaultPrevented).toBe(true);
    expect(stopPropagation).toHaveBeenCalled();
    await waitFor(() =>
      expect(onViolation).toHaveBeenCalledWith("new-window-shortcut"),
    );
  });

  it("mencegah klik link yang membuka tab baru saat sesi ujian aktif", async () => {
    const onViolation = vi.fn().mockResolvedValue(true);
    const { getByRole } = render(
      <TestExamBrowserExitCleanup
        hasActiveExamSession
        onViolation={onViolation}
      />,
    );

    const link = getByRole("link", { name: "Ujian lain" });
    const event = new MouseEvent("click", {
      button: 0,
      ctrlKey: true,
      bubbles: true,
      cancelable: true,
    });
    const stopPropagation = vi.spyOn(event, "stopPropagation");

    link.dispatchEvent(event);

    expect(event.defaultPrevented).toBe(true);
    expect(stopPropagation).toHaveBeenCalled();
    await waitFor(() =>
      expect(onViolation).toHaveBeenCalledWith("new-tab-click"),
    );
  });

  it("meminta konfirmasi sebelum menandai tab-hidden sebagai pelanggaran", async () => {
    const onViolation = vi.fn().mockResolvedValue(true);
    const confirm = vi.spyOn(window, "confirm").mockReturnValue(false);
    render(
      <TestExamBrowserExitCleanup
        hasActiveExamSession
        onViolation={onViolation}
      />,
    );

    setDocumentVisibilityState("hidden");
    document.dispatchEvent(new Event("visibilitychange"));

    expect(confirm).toHaveBeenCalled();
    expect(onViolation).not.toHaveBeenCalled();
  });

  it("menandai tab-hidden sebagai pelanggaran setelah konfirmasi disetujui", async () => {
    const onViolation = vi.fn().mockResolvedValue(true);
    vi.spyOn(window, "confirm").mockReturnValue(true);
    render(
      <TestExamBrowserExitCleanup
        hasActiveExamSession
        onViolation={onViolation}
      />,
    );

    setDocumentVisibilityState("hidden");
    document.dispatchEvent(new Event("visibilitychange"));

    await waitFor(() => expect(onViolation).toHaveBeenCalledWith("tab-hidden"));
  });
});
