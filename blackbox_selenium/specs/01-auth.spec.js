import { expect } from "chai";
import { By } from "selenium-webdriver";
import { createDriver } from "../config/browser.js";
import { routes } from "../config/routes.js";
import { accountFor } from "../data/accounts.js";
import { loginAs, loginWithCredentials, logout } from "../helpers/auth.helper.js";
import { expectAccessDenied, expectCurrentUrlContains, expectTextVisible } from "../helpers/assertion.helper.js";
import { openRoute } from "../helpers/navigation.helper.js";
import { installScreenshotOnFailure } from "../helpers/spec.helper.js";

describe("F-04 Autentikasi dan Otorisasi", function () {
  let driver;

  installScreenshotOnFailure(() => driver);

  beforeEach(async function () {
    driver = await createDriver();
  });

  afterEach(async function () {
    if (driver) await driver.quit();
    driver = null;
  });

  it("[F-04-S01] Login berhasil sesuai role", async function () {
    await loginAs(driver, "admin");
    await logout(driver);

    await loginAs(driver, "guru");
    await logout(driver);

    await loginAs(driver, "siswa");
  });

  it("[F-04-S02] Login gagal dengan kredensial salah", async function () {
    const admin = accountFor("admin");
    await loginWithCredentials(driver, admin.username, "PasswordSalah123!");
    await expectTextVisible(driver, "Login");
    const body = await driver.findElement(By.css("body")).getText();
    expect(body.toLowerCase()).to.match(/gagal|password|kredensial|username/);
  });

  it("[F-04-S03] Akses menu dibatasi oleh role", async function () {
    await loginAs(driver, "siswa");
    await openRoute(driver, routes.adminDashboard);
    await expectAccessDenied(driver);
  });
});

describe("F-16 Pembatasan Login Satu Perangkat Siswa", function () {
  let firstDriver;
  let secondDriver;

  installScreenshotOnFailure(() => secondDriver || firstDriver);

  afterEach(async function () {
    if (secondDriver) await secondDriver.quit();
    if (firstDriver) await firstDriver.quit();
    secondDriver = null;
    firstDriver = null;
  });

  it("[F-16-S01] Login siswa pertama berhasil", async function () {
    firstDriver = await createDriver();
    await loginAs(firstDriver, "siswa");
    await expectCurrentUrlContains(firstDriver, routes.siswaDashboard);
  });

  it("[F-16-S02] Login siswa kedua ditolak saat sesi masih aktif", async function () {
    firstDriver = await createDriver();
    secondDriver = await createDriver();
    const siswa = accountFor("siswa");

    await loginAs(firstDriver, "siswa");
    await loginWithCredentials(secondDriver, siswa.username, siswa.password);

    const body = (await secondDriver.findElement(By.css("body")).getText()).toLowerCase();
    expect(body).to.match(/perangkat|sesi|login gagal|logout/);
  });

  it("[F-16-S03] Login kembali berhasil setelah sesi lama logout/reset", async function () {
    firstDriver = await createDriver();
    await loginAs(firstDriver, "siswa");
    await logout(firstDriver);

    secondDriver = await createDriver();
    await loginAs(secondDriver, "siswa");
    await expectCurrentUrlContains(secondDriver, routes.siswaDashboard);
  });
});
