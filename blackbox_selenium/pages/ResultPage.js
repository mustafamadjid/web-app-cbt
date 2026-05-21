import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class ResultPage extends BasePage {
  constructor(driver, role = "admin") {
    const paths = {
      admin: routes.hasilUjianAdmin,
      guru: routes.hasilUjianGuru,
      siswa: routes.hasilUjianSiswa,
    };
    super(driver, paths[role] || routes.hasilUjianAdmin);
  }
}
