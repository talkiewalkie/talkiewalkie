import React from "react";
// import * as Sentry from "@sentry/browser";

// import { Button } from "components";

export class ErrorBoundary extends React.Component {
  state: { error?: Error } = {};

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  componentDidCatch(error: Error) {
    // Sentry.captureException(error);
  }

  render() {
    const { children } = this.props;
    const { error } = this.state;
    if (error) {
      return (
        <div className="m-auto p-16 flex-center flex-col text-center">
          <span className="font-semibold">Une erreur est survenue</span>
          {error.message.includes(
            "Failed to execute 'insertBefore' on 'Node'"
          ) && (
            <span>
              Le plugin Google traduction peut être à l'origine de cette erreur.
              S'il est activé, utilisez l'option "Ne jamais traduire ce site".
            </span>
          )}
          {process.env.NODE_ENV === "development" && (
            <div className="text-danger" style={{ whiteSpace: "pre-wrap" }}>
              {error.toString()}
            </div>
          )}
          <button className="mt-16" onClick={() => window.location.reload()}>
            Recharger
          </button>
        </div>
      );
    }

    return children;
  }
}
