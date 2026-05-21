import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class PrintPage extends BasePage {
  constructor(driver, role = "admin") {
    super(driver, role === "guru" ? routes.cetakGuru : routes.cetakAdmin);
  }
}
