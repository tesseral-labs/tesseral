import * as Sentry from "@sentry/react";
import React from "react";
import { createRoot } from "react-dom/client";

import { App } from "./App";

Sentry.init({
  dsn: "https://79f15f3f4f93544077e1e70509ebbdd7@o4505847296557056.ingest.us.sentry.io/4509096519073792",
});

const root = createRoot(document.getElementById("react-root") as HTMLElement);
root.render(<App />);
