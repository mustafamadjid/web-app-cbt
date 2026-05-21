import { BasePage } from "./BasePage.js";

export class ActivityLogPage extends BasePage {
  constructor(driver) {
    super(driver, "/dashboard/administrator/log-aktivitas");
  }
}
