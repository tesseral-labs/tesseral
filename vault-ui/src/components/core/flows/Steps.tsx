import { Circle, CircleCheck, CircleDashed } from "lucide-react";
import React, { PropsWithChildren } from "react";

import { cn } from "@/lib/utils";

export function Steps({ children }: PropsWithChildren) {
  return (
    <div className="flex justify-between items-center w-full gap-2">
      {children}
    </div>
  );
}

interface StepProps {
  label: string;
  status: "active" | "completed" | "pending";
}

export function Step({ label, status }: StepProps) {
  return (
    <div className="flex flex-col items-center justify-center space-y-2">
      {status === "active" && <Circle className="h-4 w-4" />}
      {status === "completed" && (
        <CircleCheck className="h-4 w-4 text-muted-foreground/50" />
      )}
      {status === "pending" && (
        <CircleDashed className="h-4 w-4 text-muted-foreground/50" />
      )}

      <div
        className={cn(
          "text-xs text-center",
          ["completed", "pending"].includes(status) &&
            "text-muted-foreground/50",
          status === "active" && "font-medium",
          status === "completed" && "line-through",
        )}
      >
        {label}
      </div>
    </div>
  );
}
