import React, { PropsWithChildren } from "react";

export function AuthPreviewCenterLayout({ children }: PropsWithChildren) {
  return (
    <div className="bg-background w-full flex flex-row items-center justify-center p-8 py-16">
      {children}
    </div>
  );
}

export function CenterLayoutWireframePreview() {
  return (
    <div className="w-full py-6 flex flex-col flex-wrap items-center justify-center bg-muted rounded-md gap-4">
      <div className="w-full">
        <div className="w-10 h-3 bg-muted rounded mx-auto" />
      </div>

      <div className="px-2 py-4 rounded-md border bg-white space-y-2">
        <div className="h-5 w-32 rounded border flex items-center justify-center gap-2 mx-auto">
          <div className="h-3 w-3 bg-muted rounded-sm" />
          <div className="h-1 w-18 bg-muted" />
        </div>
        <div className="h-5 w-32 rounded border flex items-center justify-center gap-2 mx-auto">
          <div className="h-3 w-3 bg-muted rounded-sm" />
          <div className="h-1 w-18 bg-muted" />
        </div>
        <div className="border-t mx-auto w-32 my-3" />
        <div className="h-4 w-32 rounded border flex items-center pl-2 mx-auto">
          <div className="h-1 w-24 bg-muted rounded-sm" />
        </div>
        <div className="h-4 w-32 rounded border flex items-center pl-2 mx-auto">
          <div className="h-1 w-24 bg-muted rounded-sm" />
        </div>
        <div className="h-5 w-32 rounded bg-muted flex items-center pl-2 mx-auto" />
      </div>
    </div>
  );
}
