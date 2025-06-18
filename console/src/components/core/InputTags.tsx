import { XIcon } from "lucide-react";
import React, { Dispatch, SetStateAction, forwardRef, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type InputTagsProps = React.ComponentProps<"input"> & {
  value: string[];
  onChange: Dispatch<SetStateAction<string[]>>;
};

export const InputTags = forwardRef<HTMLInputElement, InputTagsProps>(
  ({ value, onChange, onBlur, ...props }, ref) => {
    const [pendingDataPoint, setPendingDataPoint] = useState("");

    function addPendingDataPoint() {
      if (pendingDataPoint) {
        // trim() because a copy-pasted input may still contain leading/trailing whitespace
        const newDataPoints = new Set([...value, pendingDataPoint.trim()]);
        onChange(Array.from(newDataPoints));
        setPendingDataPoint("");
      }
    }

    return (
      <>
        <div className="flex gap-x-2">
          <Input
            value={pendingDataPoint}
            onChange={(e) => setPendingDataPoint(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.preventDefault();
                addPendingDataPoint();
              } else if (e.key === "," || e.key === " ") {
                e.preventDefault();
                addPendingDataPoint();
              }
            }}
            onBlur={(e) => {
              if (pendingDataPoint !== "") {
                addPendingDataPoint();
              }
              if (onBlur) {
                onBlur(e);
              }
            }}
            {...props}
            ref={ref}
          />
          <Button
            type="button"
            variant="secondary"
            onClick={addPendingDataPoint}
          >
            Add
          </Button>
        </div>
        {value.length > 0 && (
          <div className="flex flex-wrap items-center gap-2">
            {value.map((item, idx) => (
              <Badge key={idx} variant="secondary" className="px-2 py-1">
                {item}
                <button
                  type="button"
                  className="ml-2"
                  onClick={() => {
                    onChange(value.filter((i) => i !== item));
                  }}
                >
                  <XIcon className="text-muted-foreground size-4" />
                </button>
              </Badge>
            ))}
          </div>
        )}
      </>
    );
  },
);
InputTags.displayName = "InputTags";
