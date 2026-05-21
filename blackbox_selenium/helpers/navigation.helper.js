import { env } from "../config/env.js";

export async function openRoute(driver, path) {
  await driver.get(`${env.baseUrl}${path}`);
}

export async function refresh(driver) {
  await driver.navigate().refresh();
}

export async function back(driver) {
  await driver.navigate().back();
}
