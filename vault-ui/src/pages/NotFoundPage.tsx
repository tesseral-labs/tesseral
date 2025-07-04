import React from "react";
import { Link } from "react-router";

import { Title } from "@/components/core/Title";

export function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Title title="Not Found" />

      <div className="space-y-4 text-center">
        <h1 className="text-6xl font-bold">404</h1>
        <p className="text-muted-foreground text-xl">Page not found</p>
        <Link
          to="/"
          className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
        >
          Go back home
        </Link>
      </div>
    </div>
  );
}
