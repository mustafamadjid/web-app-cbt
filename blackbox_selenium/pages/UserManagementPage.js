import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class UserManagementPage extends BasePage {
  static guru(driver) {
    return new UserManagementPage(driver, routes.adminGuru);
  }

  static siswa(driver) {
    return new UserManagementPage(driver, routes.adminSiswa);
  }
}
