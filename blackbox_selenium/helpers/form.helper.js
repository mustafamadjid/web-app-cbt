import { By, Key } from "selenium-webdriver";
import { expectToastSuccess, expectValidationMessage } from "./assertion.helper.js";
import { waitForClickable, waitForVisible } from "./wait.helper.js";

export const byTestId = (testId) => By.css(`[data-testid="${testId}"]`);

export async function clearAndType(driver, locator, value) {
  const element = await waitForVisible(driver, locator);
  await element.clear();
  await element.sendKeys(value);
  return element;
}

export async function typeByTestId(driver, testId, value) {
  return clearAndType(driver, byTestId(testId), value);
}

export async function clickByTestId(driver, testId) {
  const element = await waitForClickable(driver, byTestId(testId));
  await element.click();
}

export async function clickByText(driver, text) {
  const escaped = text.replace(/"/g, '\\"');
  const element = await waitForClickable(
    driver,
    By.xpath(`//*[self::button or self::a][contains(normalize-space(.), "${escaped}")]`)
  );
  await element.click();
}

export async function selectByText(driver, testId, visibleText) {
  const element = await waitForClickable(driver, byTestId(testId));
  await element.click();
  await element.sendKeys(visibleText, Key.ENTER);
}

export async function submitAndExpectSuccess(driver) {
  await clickByText(driver, "Simpan");
  await expectToastSuccess(driver);
}

export async function submitAndExpectValidation(driver, expectedText) {
  await clickByText(driver, "Simpan");
  await expectValidationMessage(driver, expectedText);
}
