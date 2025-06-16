import React from "react";

import { cn } from "@/lib/utils";

import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

export function TableCardSkeleton({
  columns = 5,
  noAction = false,
  rows = 3,
}: {
  columns?: number;
  noAction?: boolean;
  rows?: number;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="h-6 w-32 bg-gray-200 animate-pulse rounded" />
        <CardDescription className="h-4 w-md bg-gray-200 animate-pulse rounded" />
        {!noAction && (
          <CardAction>
            <div className="h-10 w-40 bg-gray-200 animate-pulse rounded-md" />
          </CardAction>
        )}
      </CardHeader>
      <CardContent>
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
      </CardContent>
    </Card>
  );
}
