import { By } from "selenium-webdriver";
import { accountFor } from "../data/accounts.js";
import { routes } from "../config/routes.js";
import { openRoute } from "./navigation.helper.js";
import { clearAndType, clickByText } from "./form.helper.js";
import { expectAccessDenied as assertAccessDenied, expectCurrentUrlContains } from "./assertion.helper.js";
import { waitForVisible } from "./wait.helper.js";

export async function loginWithCredentials(driver, username, password) {
  await openRoute(driver, routes.login);
  await clearAndType(driver, By.id("username"), username);
  await clearAndType(driver, By.id("password"), password);
  await clickByText(driver, "Masuk");
}

export async function loginAs(driver, role) {
  const account = accountFor(role);
  await loginWithCredentials(driver, account.username, account.password);
  await expectCurrentUrlContains(driver, account.dashboardPath);
}

export async function logout(driver) {
  const candidates = [
    By.css('[data-testid="logout-button"]'),
    By.xpath('//*[self::button or self::a][contains(normalize-space(.), "Keluar")]'),
    By.xpath('//*[self::button or self::a][contains(normalize-space(.), "Logout")]'),
  ];

  for (const locator of candidates) {
    const elements = await driver.findElements(locator);
    if (elements.length > 0) {
      await elements[0].click();
      await expectCurrentUrlContains(driver, routes.login);
      return;
    }
  }

  await driver.manage().deleteAllCookies();
  await driver.executeScript("window.localStorage.clear(); window.sessionStorage.clear();");
  await openRoute(driver, routes.login);
}

export async function expectLoggedInAsRole(driver, role) {
  const account = accountFor(role);
  await expectCurrentUrlContains(driver, account.dashboardPath);
  await waitForVisible(driver, By.css("body"));
}

export async function expectAccessDenied(driver) {
  await assertAccessDenied(driver);
}
