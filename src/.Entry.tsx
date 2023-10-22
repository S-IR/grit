import React from "react";
import ReactDOM from "react-dom";
import App from "./App";

const ws = new WebSocket("ws://localhost:3000/ws");

ws.addEventListener("message", (event) => {
  if (event.data === "file-changed") {
    // Perform hot module replacement
  }
});

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);
