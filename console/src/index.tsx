import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import * as Sentry from "@sentry/react";

Sentry.init({
  dsn: "https://9634699510defb519b8158327f7f46b4@o4505847296557056.ingest.us.sentry.io/4509096499412992",
});

const root = createRoot(document.getElementById('react-root') as HTMLElement);
root.render(<App />);
