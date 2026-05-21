import { expect } from "chai";
import { By, until } from "selenium-webdriver";
import { env } from "../config/env.js";
import { expectedText } from "../data/expected-text.js";
import { waitForText, waitForUrlContains, waitForVisible } from "./wait.helper.js";

export async function expectCurrentUrlContains(driver, path) {
  await waitForUrlContains(driver, path);
  expect(await driver.getCurrentUrl()).to.include(path);
}

export async function expectTextVisible(driver, text) {
  await waitForText(driver, text);
}

export async function expectTextNotVisible(driver, text) {
  const body = await driver.findElement(By.css("body")).getText();
  expect(body.toLowerCase()).to.not.include(text.toLowerCase());
}

async function expectAnyText(driver, candidates) {
  await driver.wait(async () => {
    const body = (await driver.findElement(By.css("body")).getText()).toLowerCase();
    return candidates.some((candidate) => body.includes(candidate.toLowerCase()));
  }, env.timeout);
}

export async function expectToastSuccess(driver) {
  await driver.wait(until.elementLocated(By.css('[role="status"], [aria-live], .go2072408551')), env.timeout);
}

export async function expectToastError(driver) {
  await expectAnyText(driver, ["gagal", "error", "ditolak", "tidak valid"]);
}

export async function expectValidationMessage(driver, text) {
  if (text) {
    await expectTextVisible(driver, text);
    return;
  }
  await expectAnyText(driver, expectedText.validation);
}

export async function expectEmptyState(driver) {
  await expectAnyText(driver, expectedText.emptyState);
}

export async function expectAccessDenied(driver) {
  await expectAnyText(driver, expectedText.accessDenied);
}

export async function expectElementVisible(driver, locator) {
  return waitForVisible(driver, locator);
}
