import { By, until } from "selenium-webdriver";
import { env } from "../config/env.js";

export async function waitForLocated(driver, locator, timeout = env.timeout) {
  return driver.wait(until.elementLocated(locator), timeout);
}

export async function waitForVisible(driver, locator, timeout = env.timeout) {
  const element = await waitForLocated(driver, locator, timeout);
  await driver.wait(until.elementIsVisible(element), timeout);
  return element;
}

export async function waitForClickable(driver, locator, timeout = env.timeout) {
  const element = await waitForVisible(driver, locator, timeout);
  await driver.wait(until.elementIsEnabled(element), timeout);
  return element;
}

export async function waitForUrlContains(driver, path, timeout = env.timeout) {
  await driver.wait(until.urlContains(path), timeout);
}

export async function waitForText(driver, text, timeout = env.timeout) {
  const lower = text.toLowerCase();
  await driver.wait(async () => {
    const body = await driver.findElement(By.css("body")).getText();
    return body.toLowerCase().includes(lower);
  }, timeout);
}

export async function sleep(ms) {
  await new Promise((resolve) => setTimeout(resolve, ms));
}
