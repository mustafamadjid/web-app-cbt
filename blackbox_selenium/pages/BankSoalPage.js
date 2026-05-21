import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class BankSoalPage extends BasePage {
  constructor(driver, role = "admin") {
    super(driver, role === "guru" ? routes.bankSoalGuru : routes.bankSoalAdmin);
  }
}
