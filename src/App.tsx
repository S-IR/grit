import React from "react";
import logo from "./logo.svg";
import "./App.css";
import ErrorBoundary from "./_components/ErrorBoundary";

const ws = new WebSocket("ws://localhost:3000/ws");
ws.addEventListener("message", (e) => {
  const start = performance.now();

  if (e.data instanceof Blob) return handleRepalceBlob(e);
  if (e.data === "string") return handleJSONResponse(e);
  else if (typeof e.data === "string") {
  }
  const end = performance.now();
  console.log(`Hot reloading took ${end - start} milliseconds`);
});

function App() {
  return (
    <ErrorBoundary ws={ws}>
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            <p>
            Edit dwadwadwadwadwadaw<code>src/App.tsx</code> and save to reload.
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
          <button
            className="w-12 h-12 bg-black rounded-sm text-white"
            onClick={() => {
              throw new Error("ERROR");
            }}
          >
            THROW
          </button>
        </header>
        xxxxxx``
      </div>
    </ErrorBoundary>
  );
}

export default App;

const handleRepalceBlob = (e: MessageEvent<any>) => {
  const reader = new FileReader();
  reader.onload = function () {
    const text = reader.result;
    if (typeof text !== "string") throw new Error("type of text is not string");
    const headerEndIndex = text.indexOf(":");
    const header = text.substring(0, headerEndIndex);
    const assetData = text.substring(headerEndIndex + 1);

    if (header === "js") {
      // Remove existing script and add a new one
      const oldScript = document.querySelector("#react-bundle");
      if (oldScript === null || oldScript.parentNode === null) return;
      const newScript = document.createElement("script");
      newScript.src = "bundle.js?time=" + new Date().getTime();
      newScript.id = "react-bundle";
      oldScript.parentNode.replaceChild(newScript, oldScript);
    } else if (header === "css") {
      const styleSheet = document.querySelector(
        'link[rel="stylesheet"]'
      ) as HTMLLinkElement;
      if (styleSheet === null) return;
      styleSheet.href = "styles.css?time=" + new Date().getTime();
    }
  };
  reader.readAsText(e.data);
};

type updateRes = { assetPaths: string[] };
const handleJSONResponse = (e: MessageEvent<any>) => {
  const data = JSON.parse(e.data);
  if (data.type === "update") return handleUpdate(data as updateRes);
  if (data.type === "error") return handleError(data);
};

const handleUpdate = (data: { assetPaths: string[] }) => {
  {
    for (let i = 0; i < data.assetPaths.length; i++) {
      const path = data.assetPaths[i];
      const element = document.querySelector(
        `[src="${path}"], [href="${path}"]`
      );
      if (element) {
        const newUrl = `${path}?${new Date().getTime()}`;
        if (element.tagName === "SCRIPT" || element.tagName === "IMG") {
          (element as HTMLImageElement | HTMLScriptElement).src = newUrl;
        } else if (element.tagName === "LINK") {
          (element as HTMLLinkElement).href = newUrl;
        }
      }
    }
  }
};

const handleError = (data: { errors: string[] }) => {
  for (let i = 0; i < data.errors.length; i++) {
    const error = data.errors[i];
  }
};
