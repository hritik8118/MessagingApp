import React, { useState } from "react";
import Register from "./components/Register";
import Login from "./components/Login";
import Chat from "./components/Chat";

function App() {
  const [step, setStep] = useState("register");

  if (step === "register") return <Register onRegister={() => setStep("login")} />;
  if (step === "login") return <Login onLogin={() => setStep("chat")} />;
  return <Chat />;
}

export default App;
