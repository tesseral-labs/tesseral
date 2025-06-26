import { ChevronLeft, ChevronRight } from "lucide-react";
import React from "react";

import { cn } from "@/lib/utils";

import { Button } from "../ui/button";
import { Separator } from "../ui/separator";

export function Pagination({
  count = 0,
  className = "",
  hasNextPage = false,
  hasPreviousPage = false,
  fetchNextPage = () => {},
  fetchPreviousPage = () => {},
}: {
  className?: string;
  count?: number;
  hasNextPage?: boolean;
  hasPreviousPage?: boolean;
  fetchNextPage?: () => void;
  fetchPreviousPage?: () => void;
}) {
  return (
    <>
      {(hasNextPage || hasPreviousPage) && (
        <div
          className={cn(
            "flex items-center justify-end gap-2 w-full &:first:mb-4",
            className,
          )}
        >
          <div className="h-9 flex items-center">
            <span className="px-4 text-sm text-muted-foreground">
              Showing <span className="font-semibold">{count}</span>{" "}
              {count > 1 ? "results" : "result"}
            </span>
          </div>
          <div className="h-9 flex items-center border rounded-md">
            <Button
              className="rounded-r-none"
              variant="ghost"
              size="sm"
              onClick={fetchPreviousPage}
              disabled={!hasPreviousPage}
            >
              <ChevronLeft />
            </Button>
            <Separator orientation="vertical" />
            <Button
              className="rounded-l-none"
              variant="ghost"
              size="sm"
              onClick={fetchNextPage}
              disabled={!hasNextPage}
            >
              <ChevronRight />
            </Button>
          </div>
        </div>
      )}
    </>
  );
}
