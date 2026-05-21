import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class AnnouncementPage extends BasePage {
  constructor(driver, role = "admin") {
    const paths = {
      admin: routes.pengumumanAdmin,
      guru: routes.pengumumanGuru,
      siswa: routes.pengumumanSiswa,
    };
    super(driver, paths[role] || routes.pengumumanAdmin);
  }
}
