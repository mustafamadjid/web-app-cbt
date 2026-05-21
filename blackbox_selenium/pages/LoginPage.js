import { By } from "selenium-webdriver";
import { BasePage } from "./BasePage.js";
import { routes } from "../config/routes.js";
import { loginWithCredentials } from "../helpers/auth.helper.js";

export class LoginPage extends BasePage {
  constructor(driver) {
    super(driver, routes.login);
  }

  usernameInput = By.id("username");
  passwordInput = By.id("password");

  async login(username, password) {
    await loginWithCredentials(this.driver, username, password);
  }
}
