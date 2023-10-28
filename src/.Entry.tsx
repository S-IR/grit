import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import { createRoot } from "react-dom/client";
import ErrorBoundary from "./_components/ErrorBoundary";

const container = document.getElementById("app");
if (!container) throw new Error("No container");
const root = createRoot(container);
root.render(<App />);
