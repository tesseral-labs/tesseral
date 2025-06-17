import React from "react";

import { cn } from "@/lib/utils";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

export function TableSkeleton({
  columns = 5,
  rows = 3,
}: {
  columns?: number;
  rows?: number;
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          {Array.from({ length: columns }).map((_, index) => (
            <TableHead
              className={cn(index === columns - 1 ? "text-right" : "")}
              key={index}
            >
              <div
                className={cn(
                  "h-4 w-24 bg-gray-200 animate-pulse rounded",
                  index === columns - 1 ? "ml-auto" : "",
                )}
              />
            </TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {Array.from({ length: rows }).map((_, rowIndex) => (
          <TableRow key={rowIndex}>
            {Array.from({ length: columns }).map((_, colIndex) => (
              <TableCell
                className={cn(colIndex === columns - 1 ? "text-right" : "")}
                key={colIndex}
              >
                {colIndex === columns - 1 ? (
                  <div className="h-8 w-24 bg-gray-200 animate-pulse rounded-md ml-auto" />
                ) : (
                  <div className="h-4 w-32 bg-gray-200 animate-pulse rounded" />
                )}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
