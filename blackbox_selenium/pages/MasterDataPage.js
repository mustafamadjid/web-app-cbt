import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class MasterDataPage extends BasePage {
  static mapel(driver) {
    return new MasterDataPage(driver, routes.mapel);
  }

  static kelas(driver) {
    return new MasterDataPage(driver, routes.kelas);
  }

  static ruang(driver) {
    return new MasterDataPage(driver, routes.ruang);
  }

  static sesi(driver) {
    return new MasterDataPage(driver, routes.sesi);
  }
}
