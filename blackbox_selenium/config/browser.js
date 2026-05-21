import chrome from "selenium-webdriver/chrome.js";
import { Builder } from "selenium-webdriver";
import path from "node:path";
import { env } from "./env.js";

export function downloadsDir() {
  return path.join(env.rootDir, "artifacts", "downloads");
}

export async function createDriver(options = {}) {
  if (env.browser !== "chrome") {
    throw new Error(`Browser "${env.browser}" belum didukung. Gunakan E2E_BROWSER=chrome.`);
  }

  const chromeOptions = new chrome.Options();
  chromeOptions.addArguments("--window-size=1366,768");
  chromeOptions.addArguments("--disable-dev-shm-usage");
  chromeOptions.addArguments("--no-sandbox");
  chromeOptions.setUserPreferences({
    "download.default_directory": downloadsDir(),
    "download.prompt_for_download": false,
    "download.directory_upgrade": true,
    "safebrowsing.enabled": true,
  });

  if (options.headless ?? env.headless) {
    chromeOptions.addArguments("--headless=new");
  }

  const driver = await new Builder()
    .forBrowser("chrome")
    .setChromeOptions(chromeOptions)
    .build();

  await driver.manage().setTimeouts({
    implicit: 0,
    pageLoad: env.timeout * 2,
    script: env.timeout,
  });

  return driver;
}
