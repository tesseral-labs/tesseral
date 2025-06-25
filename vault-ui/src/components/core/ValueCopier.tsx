import { Copy } from "lucide-react";
import React from "react";
import { toast } from "sonner";

import { cn } from "@/lib/utils";

export function ValueCopier({
  label,
  maxLength = 255,
  value,
  className = "",
}: {
  value: string;
  label?: string;
  maxLength?: number;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "inline-flex items-center bg-muted text-muted-foreground px-2 py-1 rounded text-xs font-mono cursor-pointer pr-6 relative max-w-full hover:text-foreground",
        className,
      )}
      onClick={() => {
        navigator.clipboard.writeText(value);
        toast.success(`${label ? label : "Value"} copied to clipboard`);
      }}
    >
      <span className="flex-shrink overflow-hidden max-w-full flex-grow-0">
        {value.length > maxLength ? value.substring(0, maxLength) : value}
        {value.length > maxLength ? "..." : ""}
      </span>
      <div className="bg-muted flex items-center justify-center absolute right-0 top-0.5 p-1">
        <Copy className="w-3 h-3" />
      </div>
    </div>
  );
}
