import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";

export class DashboardPage extends BasePage {
  constructor(driver, role = "admin") {
    const paths = {
      admin: routes.adminDashboard,
      guru: routes.guruDashboard,
      siswa: routes.siswaDashboard,
    };
    super(driver, paths[role] || routes.adminDashboard);
  }
}
