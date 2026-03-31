import React from "react";
import ReactDOM from "react-dom/client";

import App from "../App.jsx";
import "../index.css";
import "../styles/signin.css";

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App initialScreen="signin" />
  </React.StrictMode>,
);
