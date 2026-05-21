import { expect } from "chai";
import { By } from "selenium-webdriver";
import { clickByText, clearAndType } from "./form.helper.js";
import { waitForVisible } from "./wait.helper.js";

const tableLocator = By.css('[data-testid="data-table"], table');

export async function expectTableContains(driver, text) {
  const table = await waitForVisible(driver, tableLocator);
  expect(await table.getText()).to.include(text);
}

export async function expectTableNotContains(driver, text) {
  const table = await waitForVisible(driver, tableLocator);
  expect(await table.getText()).to.not.include(text);
}

export async function getTableRowCount(driver) {
  const rows = await driver.findElements(By.css('[data-testid="table-row"], tbody tr'));
  return rows.length;
}

export async function searchTable(driver, keyword) {
  const locators = [
    By.css('[data-testid="search-input"]'),
    By.css('input[type="search"]'),
    By.css('input[placeholder*="Cari" i]'),
    By.css('input[placeholder*="Search" i]'),
  ];
  for (const locator of locators) {
    const elements = await driver.findElements(locator);
    if (elements.length > 0) {
      await clearAndType(driver, locator, keyword);
      return;
    }
  }
  throw new Error("Input pencarian tabel tidak ditemukan.");
}

export async function openFirstRowDetail(driver) {
  const action = await waitForVisible(driver, By.css('tbody tr a, tbody tr button, [data-testid="table-row"] a, [data-testid="table-row"] button'));
  await action.click();
}

export async function deleteRowByText(driver, text) {
  await expectTableContains(driver, text);
  await clickByText(driver, "Hapus");
}
