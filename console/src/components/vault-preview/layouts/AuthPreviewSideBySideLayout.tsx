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
