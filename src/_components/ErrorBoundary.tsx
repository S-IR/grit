import React, { ErrorInfo, ReactNode } from "react";
import { v4 as uuidv4 } from "uuid";

interface ErrorDisplayed {
  ID: string;
  Name: string;
  Text: string;
  Line?: number;
  Column?: number;
  File?: string;
  LineText?: string;
  Notes?: string[];
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: ErrorDisplayed | null | ErrorDisplayed[];
}

interface ErrorBoundaryProps {
  children: ReactNode;
  ws: WebSocket;
}

export default class ErrorBoundary extends React.Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  ws: WebSocket;

  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
    this.ws = props.ws;
  }

  componentDidMount(): void {
    this.ws.onmessage = (e) => {
      if (typeof e.data !== "string") {
        this.setState({ hasError: false, error: null });
        return;
      }
      const data = JSON.parse(e.data) as {
        type: "error";
        errors: ErrorDisplayed[];
      };
      console.log("data", data);

      this.setState({ hasError: true, error: data.errors });
    };
  }
  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    const ID = uuidv4();
    return {
      hasError: true,
      error: { Text: error.message, Name: error.name, ID },
    };
  }

  componentDidCatch(error: Error, _: ErrorInfo): void {
    const ID = uuidv4();

    this.setState({
      hasError: true,
      error: { Text: error.message, Name: error.name, ID },
    });
  }

  render(): ReactNode {
    if (this.state.hasError) {
      return (
        <div>
          {Array.isArray(this.state.error) ? (
            (this.state.error as ErrorDisplayed[]).map((err) => (
              <ErrorRow key={err.ID} err={err} />
            ))
          ) : (
            <ErrorRow err={this.state.error as ErrorDisplayed} />
          )}
        </div>
      );
    }
    return this.props.children;
  }
}

const ErrorRow = ({ err }: { err: ErrorDisplayed }) => (
  <div>
    <h3>{err.Name}</h3>
    <p>Text: {err.Text}</p>
    {err.Line && <p>Line: {err.Line}</p>}
    {err.Column && <p>Column: {err.Column}</p>}
    {err.File && <p>File: {err.File}</p>}
    {err.LineText && <p>Line Text: {err.LineText}</p>}
    {err.Notes && <p>Notes: {err.Notes.join(", ")}</p>}
  </div>
);
