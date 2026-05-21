import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class ExamStudentPage extends BasePage {
  constructor(driver) {
    super(driver, routes.ujianSiswa);
  }
}
