import { LoaderCircle } from "lucide-react";
import React from "react";

export function PageLoading() {
  return (
    <div className="w-full h-64 flex items-center justify-center">
      <LoaderCircle className="animate-spin text-muted" />
    </div>
  );
}
