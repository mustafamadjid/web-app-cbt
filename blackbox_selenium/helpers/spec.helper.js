import { takeScreenshot } from "./screenshot.helper.js";

export function installScreenshotOnFailure(getDriver) {
  afterEach(async function () {
    if (this.currentTest?.state === "failed") {
      await takeScreenshot(getDriver(), this.currentTest.title);
    }
  });
}

export function registerSkippedScenarios(scenarios, reason) {
  for (const title of scenarios) {
    it.skip(`${title} - ${reason}`, async function () {});
  }
}
