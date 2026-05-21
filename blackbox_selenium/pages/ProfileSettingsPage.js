import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class ProfileSettingsPage extends BasePage {
  constructor(driver) {
    super(driver, routes.pengaturan);
  }
}
