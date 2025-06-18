import React, { PropsWithChildren } from "react";

export function SideBySideLayout({ children }: PropsWithChildren) {
  return (
    <div className="bg-background w-full grid grid-cols-2 gap-0">
      <div className="bg-primary p-8 rounded-l" />
      <div className="p-8 py-16 flex flex-row justify-center items-center">
        {children}
      </div>
    </div>
  );
}

export function SideBySideLayoutWireframePreview() {
  return (
    <div className="w-full grid grid-cols-2 gap-0 rounded">
      <div className="bg-muted rounded-l-md" />
      <div className="space-y-2 py-6">
        <div className="w-full">
          <div className="w-10 h-3 bg-muted rounded mx-auto" />
        </div>

        <div className="px-2 bg-white space-y-2 mx-auto">
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
    </div>
  );
}
