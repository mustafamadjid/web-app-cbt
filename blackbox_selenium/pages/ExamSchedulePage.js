import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class ExamSchedulePage extends BasePage {
  constructor(driver, role = "admin") {
    super(driver, role === "guru" ? routes.jadwalUjianGuru : routes.jadwalUjianAdmin);
  }
}
