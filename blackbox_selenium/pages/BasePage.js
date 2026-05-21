import { By } from "selenium-webdriver";
import { env } from "../config/env.js";
import { openRoute } from "../helpers/navigation.helper.js";
import { clickByText, clearAndType } from "../helpers/form.helper.js";
import { waitForVisible } from "../helpers/wait.helper.js";

export class BasePage {
  constructor(driver, path = "/") {
    this.driver = driver;
    this.path = path;
  }

  async open() {
    await openRoute(this.driver, this.path);
  }

  async bodyText() {
    return this.driver.findElement(By.css("body")).getText();
  }

  async title() {
    const candidates = [
      By.css('[data-testid="page-title"]'),
      By.css("h1"),
      By.css("h2"),
    ];
    for (const locator of candidates) {
      const elements = await this.driver.findElements(locator);
      if (elements.length > 0) return elements[0].getText();
    }
    return "";
  }

  async clickText(text) {
    await clickByText(this.driver, text);
  }

  async type(locator, value) {
    await clearAndType(this.driver, locator, value);
  }

  async waitForBody() {
    await waitForVisible(this.driver, By.css("body"), env.timeout);
  }
}
