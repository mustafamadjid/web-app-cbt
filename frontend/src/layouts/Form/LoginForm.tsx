import { useState } from "react";

// Components
import { LoginInputField } from "../../components/common/Input/LoginInputField";

export const LoginForm = ()=> {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");


    return (
      <>
        <main>
          {/* Logo */}
          <div></div>
          {/* Welcome Text */}
          <div></div>

          {/* Form */}
          <form className="flex flex-col gap-10">
            <LoginInputField
              id="username"
              label="Username"
              type="text"
              value={username}
              onChange={setUsername}
              
              autoComplete="username"
            />
            <LoginInputField
              id="password"
              label="Password"
              type="password"
              value={password}
              onChange={setPassword}
             
              autoComplete="current-password"
            />
          </form>
        </main>
      </>
    );
}