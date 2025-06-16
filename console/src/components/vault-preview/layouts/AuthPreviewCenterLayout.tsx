import React, { PropsWithChildren } from "react";

export function AuthPreviewCenterLayout({ children }: PropsWithChildren) {
  return (
    <div className="bg-background w-full flex flex-row items-center justify-center p-8 py-16">
      {children}
    </div>
  );
}
